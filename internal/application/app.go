package application

import (
	"database/sql"
	"log"
	"server/config"
	"server/internal/services"
	"server/internal/services/auth"
	"server/internal/services/ege"
	"server/pkg/handlers"
	"server/pkg/middlewares"
	"time"

	"github.com/gin-gonic/gin"

	// initialize the database driver
	_ "github.com/lib/pq"
)

const (
	migrationScheme = "CREATE TABLE IF NOT EXISTS Users (" +
		"ID SERIAL PRIMARY KEY," +
		"email VARCHAR(100) NOT NULL UNIQUE," +
		"password VARCHAR(100) NOT NULL," +
		"created_at TIMESTAMP NOT NULL DEFAULT NOW());"
	refreshPeriod = time.Minute * 5
)

type App struct {
	Config   *config.Config
	Database *sql.DB
	Router   *gin.Engine
	Services []services.Service
}

// New creates the application instance
func New() *App {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	migrateTable(db)

	return &App{
		Config:   cfg,
		Database: db,
		Router:   gin.New(),
		Services: make([]services.Service, 0, 2),
	}
}

// Close does cleaning operations on the application
func (app *App) Close() {
	_ = app.Database.Close()
}

func (app *App) initializeServices() {
	app.Router.Use(middlewares.Logger())
	app.Router.NoRoute(handlers.NotFound)
	app.Router.HandleMethodNotAllowed = true
	app.Router.NoMethod(handlers.NoMethod)

	authRouter := app.Router.Group("auth/")
	authService := auth.NewService(app.Config, app.Database)
	go func() {
		for {
			time.Sleep(refreshPeriod)
			authService.CheckExpire()
		}
	}()
	authService.Register(authRouter)
	app.Services = append(app.Services, authService)

	apiRouter := app.Router.Group("api/")
	apiRouter.Use(authService.AuthMiddleware())

	egeRouter := apiRouter.Group("ege/")
	egeService := ege.NewService(app.Config)
	egeService.Register(egeRouter)
}

func migrateTable(db *sql.DB) {
	_, err := db.Exec(migrationScheme)
	if err != nil {
		log.Fatal("table migration: ", err)
	}
}

// Start the server
func (app *App) Start() {
	app.initializeServices()
	if gin.Mode() != gin.DebugMode {
		log.Println("Starting server on port: " + app.Config.Port)
	}
	log.Fatal(app.Router.Run(":" + app.Config.Port))
	// log.Fatal(http.ListenAndServeTLS(cfg.Realm+":"+cfg.Port, cfg.CertFile, cfg.KeyFile, r))
}
