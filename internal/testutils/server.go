package testutils

import (
	"context"
	"log"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"

	"github.com/ananthvk/gochat/internal"
	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/config"
	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/group"
	"github.com/ananthvk/gochat/internal/message"
	"github.com/ananthvk/gochat/internal/middleware"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/ananthvk/gochat/internal/token"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestDatabase struct {
	Container testcontainers.Container
	ConnStr   string
}

func NewTestDatabase(t *testing.T, ctx context.Context) *TestDatabase {
	t.Helper()

	// Create PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	return &TestDatabase{
		Container: postgresContainer,
		ConnStr:   connStr,
	}
}

func (td *TestDatabase) Close(ctx context.Context) error {
	if td.Container != nil {
		return td.Container.Terminate(ctx)
	}
	return nil
}

func newTestServerWithDatabase(t *testing.T, ctx context.Context) (*app.App, *chi.Mux, *TestDatabase) {
	t.Helper()

	// Setup database
	testDB := NewTestDatabase(t, ctx)

	router := chi.NewRouter()
	cfg := &config.Config{
		Env:                       "test",
		EnableDetailedHealthCheck: true,
		DbDSN:                     testDB.ConnStr,
		DbPingTimeout:             5 * time.Second,
		DbQueryTimeout:            5 * time.Second,
	}

	rtService := realtime.NewRealtimeService(ctx)
	dbService, err := database.NewDatabaseService(ctx, cfg)
	if err != nil {
		log.Fatalf("could not create database service %s", err)
	}
	groupService := group.NewGroupService(dbService)
	mesageService := message.NewMessageService(dbService)
	tokenService := token.NewTokenService(dbService)
	authService := auth.NewAuthService(dbService, tokenService)

	time.Sleep(50 * time.Millisecond)
	// Run migrations
	// Runs the migrate CLI tool, TODO: Later integrate the library
	if err := exec.CommandContext(ctx, "migrate", "-database", cfg.DbDSN, "-path", "../../internal/database/migrations", "up").Run(); err != nil {
		log.Fatalf("Could not migrate db: %s", err)
	}

	app := &app.App{
		Ctx:             ctx,
		RealtimeService: rtService,
		DatabaseService: dbService,
		GroupService:    groupService,
		MessageService:  mesageService,
		AuthService:     authService,
		TokenService:    tokenService,
		Config:          cfg,
	}
	app.RealtimeService.StartHubEventLoop()
	middlewares := middleware.Middlewares{
		Authenticate: auth.AuthMiddleware(tokenService),
	}
	router.Mount("/api/v1/", internal.Routes(app, middlewares))
	return app, router, testDB
}

func NewTestServerWithDatabaseAndCancel(t *testing.T) (*app.App, *httptest.Server, *TestDatabase, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	app, router, testDB := newTestServerWithDatabase(t, ctx)
	srv := httptest.NewServer(router)

	// Cleanup function that also closes database
	originalCancel := cancel
	cancel = func() {
		app.DatabaseService.Pool.Close()
		testDB.Close(ctx)
		originalCancel()
	}

	return app, srv, testDB, cancel
}
