package main

import (
	"context"
	"log"
	"net/http"

	"github.com/M0F3/docstore-api/internal/auth"
	"github.com/M0F3/docstore-api/internal/config"
	"github.com/M0F3/docstore-api/internal/database"
	"github.com/M0F3/docstore-api/internal/handlers"
	"github.com/M0F3/docstore-api/internal/middleware"
	"github.com/M0F3/docstore-api/internal/repositories"
	"github.com/M0F3/docstore-api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	config := config.LoadAppConfig()
	db := database.Connect(config.AppDatabaseConnectionUrl)
	adminDb := database.ConnectAdmin(config.AppAdminDatabaseConnectionUrl)
	if err := db.Ping(ctx); err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
	defer db.Close()

	if err := adminDb.Ping(ctx); err != nil {
        log.Fatalf("failed to connect to admin database: %v", err)
    }
	defer adminDb.Close()
	router := chi.NewRouter()
	buildApp(router, config, db, adminDb)

	log.Println("Start server at :8080")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func buildApp(router *chi.Mux, config *config.AppConfig, db *pgxpool.Pool, adminDb *pgxpool.Pool) {

	container := newContainer(db, adminDb)

	auth.Init(config.JWTSecret)

	router.Use(middleware.Logging)
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AttachDatabaseConnection)
		r.Get("/healthz", handlers.Healthz)
		r.Post("/register", container.UserHandler.Register)
		r.Post("/login", container.UserHandler.Login)
		r.Route("/", func(ra chi.Router) {
			ra.Use(jwtauth.Verifier(auth.TokenAuth))
			ra.Use(jwtauth.Authenticator) // populate context
			ra.Use(middleware.AttachUserToContext)
			ra.Use(middleware.AttachDatabaseConnectionWithSession)
			ra.Get("/users", container.UserHandler.ListUsers)
		})
	})

}

type Container struct {
	UserHandler handlers.UserHandler
}

func newContainer(db *pgxpool.Pool, adminDb *pgxpool.Pool) *Container {
	userRepo := repositories.NewUserRepository(db,adminDb)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)


	return &Container{
		UserHandler: *userHandler,
	}
}
