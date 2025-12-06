package client

import (
	"context"
	"errors"
	"github/elangreza/edot-commerce/pkg/globalcontanta"

	"github.com/elangreza/edot-commerce/gen"
	"github.com/elangreza/edot-commerce/order/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type (
	warehouseServiceClient struct {
		client gen.WarehouseServiceClient
	}
)

func NewWarehouseClient() (*warehouseServiceClient, error) {
	grpcClient, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	stockClient := gen.NewWarehouseServiceClient(grpcClient)
	return &warehouseServiceClient{client: stockClient}, nil
}

func (s *warehouseServiceClient) GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error) {
	return s.client.GetStocks(ctx, &gen.GetStockRequest{ProductIds: productIds})
}

// reserve stock after order is created
func (s *warehouseServiceClient) ReserveStock(ctx context.Context, cartItem []entity.CartItem) (*gen.ReserveStockResponse, error) {

	userID, ok := ctx.Value(globalcontanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("not valid user id")
	}

	md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	stocks := []*gen.Stock{}
	for _, item := range cartItem {
		stocks = append(stocks, &gen.Stock{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// add user id in context
	return s.client.ReserveStock(newCtx, &gen.ReserveStockRequest{
		Stocks: stocks,
	})
}

// release stock when creating order is failed or order is cancelled
func (s *warehouseServiceClient) ReleaseStock(ctx context.Context, reservedStockIds []int64) (*gen.ReleaseStockResponse, error) {
	userID, ok := ctx.Value(globalcontanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("not valid user id")
	}

	md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	return s.client.ReleaseStock(newCtx, &gen.ReleaseStockRequest{
		ReservedStockIds: reservedStockIds,
	})
}
