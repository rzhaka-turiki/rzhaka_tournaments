package app

import (
	"context"
	"log"

	"github.com/a1uka/rzhaka_tournaments/internal/config"
	"github.com/a1uka/rzhaka_tournaments/internal/database"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/handlers"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/middleware"
	"github.com/a1uka/rzhaka_tournaments/internal/transport/rest/routes"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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
	router.Use(middleware.Logger(), middleware.Recovery())

	// container w/ handlers
	container := NewContainer(db)
	restHandlers := rest.Handlers{
		Health: handlers.NewHealthHandler(db),
		User:   handlers.NewUserHandler(container.Services.User),
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
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5,
	)
	defer cancel()

	if a.db != nil {
		a.db.Close()
	}

	log.Println("application stopped")

	_ = ctx
}
