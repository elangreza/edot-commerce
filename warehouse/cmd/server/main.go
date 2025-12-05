package main

import (
	"github/elangreza/edot-commerce/pkg/dbsql"
	"github/elangreza/edot-commerce/warehouse/internal/server"
	"github/elangreza/edot-commerce/warehouse/internal/service"
	"github/elangreza/edot-commerce/warehouse/internal/sqlitedb"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("warehouse.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	warehouseRepo := sqlitedb.NewWarehouseRepo(db)
	warehouseService := service.NewWarehouseService(warehouseRepo)

	address := "localhost:50052"

	srv := server.New(warehouseService)

	if err := srv.Start(address); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
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
