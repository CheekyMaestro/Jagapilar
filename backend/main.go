package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"jagapilar-backend/config"
	"jagapilar-backend/database"
	"jagapilar-backend/handlers"
	"jagapilar-backend/middleware"
	"jagapilar-backend/services"
)

func main() {
	// Load .env file if present (ignored if missing, e.g. in production)
	_ = godotenv.Load()

	// Banner
	fmt.Println(`
     ██╗ █████╗  ██████╗  █████╗ ██████╗ ██╗██╗      █████╗ ██████╗ 
     ██║██╔══██╗██╔════╝ ██╔══██╗██╔══██╗██║██║     ██╔══██╗██╔══██╗
     ██║███████║██║  ███╗███████║██████╔╝██║██║     ███████║██████╔╝
██   ██║██╔══██║██║   ██║██╔══██║██╔═══╝ ██║██║     ██╔══██║██╔══██╗
╚█████╔╝██║  ██║╚██████╔╝██║  ██║██║     ██║███████╗██║  ██║██║  ██║
 ╚════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
    Platform Asesmen 360° Deteksi Brain-Rotting
    `)

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Println("📌 Running in offline mode — API will work but data won't be persisted")
		log.Println("📌 Set environment variables: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	} else {
		defer database.Close()

		// Run migrations
		if err := database.RunMigrations(); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}

		// Seed assessment items
		if err := services.SeedAssessmentItems(); err != nil {
			log.Printf("⚠️  Seed failed: %v", err)
		}
	}

	// Setup router
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/schools", handlers.SchoolsHandler)
	mux.HandleFunc("/api/schools/", handlers.SchoolsHandler)
	mux.HandleFunc("/api/children", handlers.ChildrenHandler)
	mux.HandleFunc("/api/children/", handlers.ChildrenHandler)
	mux.HandleFunc("/api/assessment/", handlers.AssessmentHandler)
	mux.HandleFunc("/api/auth/", handlers.AuthHandler)

	// Health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"jagapilar","version":"1.0.0"}`)
	})

	// Serve static frontend files
	frontendDir := "../frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "./frontend"
	}
	fs := http.FileServer(http.Dir(frontendDir))
	mux.Handle("/", fs)

	// Apply middleware
	handler := middleware.CORS(middleware.Auth(mux))

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 JAGAPILAR server starting on http://localhost%s", addr)
	log.Printf("📁 Serving frontend from: %s", frontendDir)
	log.Println("─────────────────────────────────────────────")
	log.Println("  API Endpoints:")
	log.Println("  GET  /api/health                    → Health check")
	log.Println("  POST /api/schools                   → Registrasi sekolah")
	log.Println("  GET  /api/schools                   → List sekolah")
	log.Println("  POST /api/children                  → Daftarkan anak")
	log.Println("  GET  /api/children/{id}             → Data anak")
	log.Println("  POST /api/children/{id}/informants  → Buat informan + token")
	log.Println("  GET  /api/assessment/items           → Item kuesioner")
	log.Println("  POST /api/assessment/sessions        → Mulai sesi")
	log.Println("  POST /api/assessment/responses       → Simpan jawaban")
	log.Println("  POST /api/assessment/submit          → Submit + skoring")
	log.Println("  GET  /api/assessment/results/{id}    → Hasil asesmen")
	log.Println("  POST /api/auth/validate-token        → Validasi token")
	log.Println("─────────────────────────────────────────────")

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
