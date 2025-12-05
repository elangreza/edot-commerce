package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"
	"errors"
	"github/elangreza/edot-commerce/pkg/globalcontanta"

	"github.com/elangreza/edot-commerce/api/internal/constanta"
	params "github.com/elangreza/edot-commerce/api/internal/params"
	"github.com/elangreza/edot-commerce/gen"
	"google.golang.org/grpc/metadata"
)

func NewOrderService(pClient gen.OrderServiceClient) *orderService {
	return &orderService{
		orderServiceClient: pClient,
	}
}

type orderService struct {
	orderServiceClient gen.OrderServiceClient
}

func (s *orderService) ListProducts(ctx context.Context, req params.AddToCartRequest) error {

	userID, ok := ctx.Value(constanta.LocalUserID).(string)
	if !ok {
		return errors.New("error when parsing userID")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, string(globalcontanta.UserIDKey), userID)

	s.orderServiceClient.AddProductToCart(ctx, &gen.AddCartItemRequest{
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
	})

	return nil
}
