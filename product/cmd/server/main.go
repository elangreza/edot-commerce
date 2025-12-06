package main

import (
	"context"
	"fmt"
	"github/elangreza/edot-commerce/pkg/dbsql"
	"github/elangreza/edot-commerce/pkg/gracefulshutdown"
	"log"
	"time"

	"github.com/elangreza/edot-commerce/product/internal/client"
	"github.com/elangreza/edot-commerce/product/internal/server"
	"github.com/elangreza/edot-commerce/product/internal/service"
	"github.com/elangreza/edot-commerce/product/internal/sqlitedb"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("product.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	// productRepo, err := mockjson.LoadProductJson()
	productRepo := sqlitedb.NewProductRepository(db)
	stockClient, err := client.NewWarehouseClient()
	errChecker(err)

	address := fmt.Sprintf("localhost:%v", 50050)

	productService := service.NewProductService(productRepo, stockClient)
	srv := server.New(productService)
	go func() {
		if err := srv.Start(address); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	gs := gracefulshutdown.New(context.Background(), 5*time.Second,
		gracefulshutdown.Operation{
			Name: "grpc",
			ShutdownFunc: func(ctx context.Context) error {
				srv.Close()
				return nil
			},
		},
		gracefulshutdown.Operation{
			Name: "sqlite",
			ShutdownFunc: func(ctx context.Context) error {
				return db.Close()
			},
		},
	)
	<-gs
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
