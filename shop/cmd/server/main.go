package main

import (
	"context"
	"github/elangreza/edot-commerce/pkg/dbsql"
	"github/elangreza/edot-commerce/pkg/gracefulshutdown"
	"github/elangreza/edot-commerce/shop/internal/client"
	"github/elangreza/edot-commerce/shop/internal/server"
	"github/elangreza/edot-commerce/shop/internal/service"
	"github/elangreza/edot-commerce/shop/internal/sqlitedb"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("shop.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	warehouseClient, err := client.NewWarehouseClient()
	errChecker(err)

	shopRepo := sqlitedb.NewShopRepo(db)
	shopService := service.NewShopService(shopRepo, warehouseClient)

	address := "localhost:50055"

	srv := server.New(shopService)
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
