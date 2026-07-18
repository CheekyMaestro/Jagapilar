package database

import (
	"log"
)

// RunMigrations creates all required tables if they don't exist
// Schema follows the PRD: schools, children, informants, assessment_items,
// assessment_sessions, item_responses, pillar_scores, composite_scores,
// referrals, audit_log, literature_repository
func RunMigrations() error {
	migrations := []string{
		// Schools table
		`CREATE TABLE IF NOT EXISTS schools (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			city TEXT NOT NULL,
			grade_level TEXT,
			principal_name TEXT,
			total_classes INT DEFAULT 0,
			principal_consent_doc_url TEXT,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Children table (anonymized)
		`CREATE TABLE IF NOT EXISTS children (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			anon_code TEXT UNIQUE NOT NULL,
			school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
			grade SMALLINT,
			birth_year SMALLINT,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Informants table (parent/teacher/student linked to a child)
		`CREATE TABLE IF NOT EXISTS informants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			child_id UUID REFERENCES children(id) ON DELETE CASCADE,
			role TEXT NOT NULL CHECK (role IN ('parent', 'teacher', 'student')),
			contact_hash TEXT,
			access_token TEXT UNIQUE,
			token_expires_at TIMESTAMPTZ,
			consent_signed_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Assessment items master table (18 items, seeded once)
		`CREATE TABLE IF NOT EXISTS assessment_items (
			code TEXT PRIMARY KEY,
			pillar TEXT NOT NULL CHECK (pillar IN ('parent', 'teacher', 'student')),
			question_text TEXT NOT NULL,
			is_reverse BOOLEAN DEFAULT false,
			construct TEXT,
			reference_source TEXT
		)`,

		// Assessment sessions
		`CREATE TABLE IF NOT EXISTS assessment_sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			child_id UUID REFERENCES children(id) ON DELETE CASCADE,
			informant_id UUID REFERENCES informants(id),
			pillar TEXT NOT NULL,
			submitted_at TIMESTAMPTZ,
			status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'submitted')),
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Item responses (progressive save per item)
		`CREATE TABLE IF NOT EXISTS item_responses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			session_id UUID REFERENCES assessment_sessions(id) ON DELETE CASCADE,
			item_code TEXT REFERENCES assessment_items(code),
			raw_score SMALLINT NOT NULL CHECK (raw_score BETWEEN 1 AND 5),
			saved_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Pillar scores (calculated server-side)
		`CREATE TABLE IF NOT EXISTS pillar_scores (
			child_id UUID REFERENCES children(id) ON DELETE CASCADE,
			pillar TEXT NOT NULL,
			subtotal NUMERIC NOT NULL,
			zone TEXT NOT NULL CHECK (zone IN ('hijau', 'kuning', 'merah')),
			calculated_at TIMESTAMPTZ DEFAULT now(),
			PRIMARY KEY (child_id, pillar)
		)`,

		// Composite scores (360° aggregation)
		`CREATE TABLE IF NOT EXISTS composite_scores (
			child_id UUID PRIMARY KEY REFERENCES children(id) ON DELETE CASCADE,
			composite NUMERIC NOT NULL,
			total_raw NUMERIC,
			composite_zone TEXT NOT NULL,
			needs_manual_review BOOLEAN DEFAULT false,
			calculated_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Referrals for high-risk cases
		`CREATE TABLE IF NOT EXISTS referrals (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			child_id UUID REFERENCES children(id) ON DELETE CASCADE,
			referred_at TIMESTAMPTZ DEFAULT now(),
			referred_to TEXT,
			handbook_sent BOOLEAN DEFAULT false,
			notes JSONB
		)`,

		// Audit log for data access tracking
		`CREATE TABLE IF NOT EXISTS audit_log (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			actor_role TEXT,
			action TEXT NOT NULL,
			target_table TEXT,
			target_id UUID,
			occurred_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Literature repository (JSONB instead of separate NoSQL)
		`CREATE TABLE IF NOT EXISTS literature_repository (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			category TEXT,
			citation_harvard TEXT,
			metadata JSONB,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_children_school ON children(school_id)`,
		`CREATE INDEX IF NOT EXISTS idx_informants_child ON informants(child_id)`,
		`CREATE INDEX IF NOT EXISTS idx_informants_token ON informants(access_token)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_child ON assessment_sessions(child_id)`,
		`CREATE INDEX IF NOT EXISTS idx_responses_session ON item_responses(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_occurred ON audit_log(occurred_at)`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			return err
		}
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}
