package main

import (
	"fmt"
	"github/elangreza/edot-commerce/pkg/dbsql"

	"log"

	"github.com/elangreza/edot-commerce/order/internal/client"
	"github.com/elangreza/edot-commerce/order/internal/server"
	"github.com/elangreza/edot-commerce/order/internal/service"
	"github.com/elangreza/edot-commerce/order/internal/sqlitedb"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("order.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	cartRepo := sqlitedb.NewCartRepository(db)
	orderRepo := sqlitedb.NewOrderRepository(db)
	stockClient, err := client.NewWarehouseClient()
	errChecker(err)
	productClient, err := client.NewProductClient()
	errChecker(err)

	orderService := service.NewOrderService(
		orderRepo,
		cartRepo,
		stockClient,
		productClient)

	srv := server.New(orderService)
	address := fmt.Sprintf(":%v", 50054)
	if err := srv.Start(address); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
	// grpcServer := grpc.NewServer(
	// 	grpc.ChainUnaryInterceptor(
	// 		interceptor.UserIDParser(),
	// 	),
	// )
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
