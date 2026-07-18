package database

import (
	"log"
)

// RunMigrations creates all required tables if they don't exist
// Schema follows the PRD: users, children, teacher_children, assessment_items,
// assessment_sessions, item_responses, pillar_scores, composite_scores,
func RunMigrations() error {
	migrations := []string{
		// Schools table (Kept for compatibility or if teachers want to select it)
		`CREATE TABLE IF NOT EXISTS schools (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			city TEXT,
			grade_level TEXT,
			principal_name TEXT,
			total_classes INT DEFAULT 0,
			principal_consent_doc_url TEXT,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Users table (Parents & Teachers)
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			email_contact TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL CHECK (role IN ('parent', 'teacher')),
			school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Children table (with Neuro ID)
		`CREATE TABLE IF NOT EXISTS children (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			neuro_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			birth_year SMALLINT,
			gender TEXT,
			created_by UUID REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Teacher-Child Mapping (Dashboard Guru)
		`CREATE TABLE IF NOT EXISTS teacher_children (
			teacher_id UUID REFERENCES users(id) ON DELETE CASCADE,
			child_id UUID REFERENCES children(id) ON DELETE CASCADE,
			added_at TIMESTAMPTZ DEFAULT now(),
			PRIMARY KEY (teacher_id, child_id)
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
			user_id UUID REFERENCES users(id) ON DELETE SET NULL,
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

		// Audit log for data access tracking
		`CREATE TABLE IF NOT EXISTS audit_log (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			actor_role TEXT,
			action TEXT NOT NULL,
			target_table TEXT,
			target_id UUID,
			occurred_at TIMESTAMPTZ DEFAULT now()
		)`,

		// Indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email_contact)`,
		`CREATE INDEX IF NOT EXISTS idx_children_neuro_id ON children(neuro_id)`,
		`CREATE INDEX IF NOT EXISTS idx_children_created_by ON children(created_by)`,
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
