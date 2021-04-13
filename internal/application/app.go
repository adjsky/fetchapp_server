package application

import (
	"database/sql"
	"log"
	"net/http"
	"server/config"
	"server/internal/services/auth"
	"server/internal/services/ege"
	"server/pkg/handlers"
	"server/pkg/middlewares"
	"time"

	"github.com/gorilla/mux"

	// initialize a driver
	_ "github.com/lib/pq"
)

const migrationScheme string = "CREATE TABLE IF NOT EXISTS Users (" +
	"ID SERIAL PRIMARY KEY," +
	"email TEXT NOT NULL UNIQUE," +
	"password TEXT NOT NULL);"

type app struct {
	Config   *config.Config
	Database *sql.DB
	Router   *mux.Router
}

// New creates the application instance
func New() *app {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	migrateTable(db)
	return &app{
		Config:   cfg,
		Database: db,
		Router:   mux.NewRouter(),
	}
}

// Close does cleaning operations on the application
func (app *app) Close() {
	_ = app.Database.Close()
}

func (app *app) initializeServices() {
	app.Router.Use(middlewares.Log)
	app.Router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	authRouter := app.Router.PathPrefix("/auth").Subrouter()
	authService := auth.NewService(app.Config, app.Database)
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			authService.CheckExpire()
		}
	}()
	authService.Register(authRouter)

	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	apiRouter.Use(authService.AuthMiddleware)

	egeRouter := apiRouter.PathPrefix("/ege").Subrouter()
	egeService := ege.NewService(app.Config)
	egeService.Register(egeRouter)
}

func migrateTable(db *sql.DB) {
	_, err := db.Exec(migrationScheme)
	if err != nil {
		log.Fatal("failed to migrate the database scheme")
	}
}

// Start the server
func (app *app) Start() {
	app.initializeServices()
	log.Println("Starting server on port: " + app.Config.Port)
	log.Fatal(http.ListenAndServe(":"+app.Config.Port, app.Router))
	//log.Fatal(http.ListenAndServeTLS(cfg.Realm+":"+cfg.Port, cfg.CertFile, cfg.KeyFile, r))
}
