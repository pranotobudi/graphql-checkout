package store

import "github.com/pranotobudi/graphql-checkout/graph/model"

type IStoreRepository interface {
	GetProduct(sku string) (*model.Product, error)
	GetPromoType1(sku string) (*model.Product, error)
}
