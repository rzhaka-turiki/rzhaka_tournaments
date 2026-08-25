package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/config"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/database"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/auth"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/middleware"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/transport/rest/routes"
)

type App struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	router    *gin.Engine
	container *Container
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	// db conn
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	// http router create
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		auth.Middleware(),
	)

	// container w/ handlers
	container, err := NewContainer(db, cfg)
	if err != nil {
		db.Close()
		return nil, err
	}
	restHandlers := rest.Handlers{
		Health:         handlers.NewHealthHandler(db),
		User:           handlers.NewUserHandler(container.Services.User),
		Role:           handlers.NewRoleHandler(container.Services.Role),
		Permission:     handlers.NewPermissionHandler(container.Repositories.Permission),
		RolePermission: handlers.NewRolePermissionHandler(container.Services.RolePermission),
		ApexAccounts:   handlers.NewApexAccountHandler(container.Services.ApexAccount),
	}
	routes.RegisterRoutes(router, &restHandlers)

	return &App{
		cfg:       cfg,
		db:        db,
		router:    router,
		container: container,
	}, nil
}

func (a *App) Run() error {
	return a.router.Run(":" + a.cfg.Server.Port)
}

func (a *App) Shutdown() {
	if a.container != nil && a.container.Clients.ApexVerifier != nil {
		if err := a.container.Clients.ApexVerifier.Close(); err != nil {
			log.Printf("failed to close apex verifier client: %v", err)
		}
	}

	if a.db != nil {
		a.db.Close()
	}

	log.Println("application stopped")
}
