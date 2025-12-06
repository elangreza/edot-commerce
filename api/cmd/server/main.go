package main

import (
	"context"
	"github/elangreza/edot-commerce/pkg/dbsql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/elangreza/edot-commerce/api/internal/rest"
	"github.com/elangreza/edot-commerce/api/internal/service"
	sqlitedb "github.com/elangreza/edot-commerce/api/internal/sqlite"
	"github.com/elangreza/edot-commerce/gen"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	handler := chi.NewRouter()

	handler.Use(middleware.Recoverer)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Timeout(60 * time.Second))
	handler.Use(middleware.RequestID)
	handler.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("auth.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)

	// repositories
	userRepo := sqlitedb.NewUserRepo(db)
	tokenRepo := sqlitedb.NewTokenRepo(db)

	// service
	authService := service.NewAuthService(userRepo, tokenRepo)

	grpcClientProduct, err := grpc.NewClient("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)
	productServiceClient := service.NewProductService(gen.NewProductServiceClient(grpcClientProduct))

	grpcClientOrder, err := grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)
	orderService := service.NewOrderService(gen.NewOrderServiceClient(grpcClientOrder))

	rest.NewAuthHandler(handler, authService)
	rest.NewProductHandler(handler, productServiceClient)
	rest.NewOrderHandler(handler, authService, orderService)

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-gracefulShutdown(context.Background(), 5*time.Second,
		operation{
			name: "server",
			shutdownFunc: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			}},
		operation{
			name: "sqlite",
			shutdownFunc: func(ctx context.Context) error {
				return db.Close()
			}},
	)
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type operation struct {
	name         string
	shutdownFunc func(ctx context.Context) error
}

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops ...operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		slog.Info("shutting down")

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		go func() {
			<-ctx.Done()
			slog.Info("force quit the app")
			wait <- struct{}{}
		}()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			go func(key int, op operation) {
				defer wg.Done()

				slog.Info(op.name, "shutdown", "started")

				if err := op.shutdownFunc(ctx); err != nil {
					slog.Error(op.name, "err", err.Error())
					return
				}

				slog.Info(op.name, "shutdown", "finished")
			}(key, op)
		}

		wg.Wait()
	}()

	return wait
}
