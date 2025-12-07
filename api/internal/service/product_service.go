package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"

	params "github.com/elangreza/edot-commerce/api/internal/params"
	"github.com/elangreza/edot-commerce/gen"
)

func NewProductService(
	pClient gen.ProductServiceClient,
	sClient gen.ShopServiceClient,
) *productService {
	return &productService{
		productServiceClient: pClient,
		shopServiceClient:    sClient,
	}
}

type productService struct {
	productServiceClient gen.ProductServiceClient
	shopServiceClient    gen.ShopServiceClient
}

func (s *productService) ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error) {
	listProduct, err := s.productServiceClient.ListProducts(ctx, &gen.ListProductsRequest{
		Search: req.Search,
		Limit:  req.Limit,
		Page:   req.Page,
		SortBy: req.SortBy,
	})

	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.ListProductsResponse{
		Products:   []*params.Product{},
		Total:      listProduct.GetTotal(),
		TotalPages: listProduct.GetTotalPages(),
	}

	shopIDs := []int64{}
	for _, product := range listProduct.Products {
		shopIDs = append(shopIDs, product.ShopId)
	}

	shops, err := s.shopServiceClient.GetShops(ctx, &gen.GetShopsRequest{
		Ids:            shopIDs,
		WithWarehouses: false,
	})
	if err != nil {
		return nil, convertErrGrpc(err)
	}

	shopsMap := make(map[int64]string)
	for _, shop := range shops.Shops {
		shopsMap[shop.GetId()] = shop.Name
	}

	for _, product := range listProduct.Products {
		p := &params.Product{
			Id:          product.GetId(),
			Name:        product.GetName(),
			Description: product.GetDescription(),
			ImageUrl:    product.GetImageUrl(),
			Stock:       product.GetStock(),
			Price: &params.Money{
				Units:        product.Price.GetUnits(),
				CurrencyCode: product.Price.GetCurrencyCode(),
			},
			ShopID: product.GetShopId(),
		}
		shopName, ok := shopsMap[product.GetShopId()]
		if ok {
			p.ShopName = shopName
		}
		res.Products = append(res.Products, p)
	}

	return res, nil
}
