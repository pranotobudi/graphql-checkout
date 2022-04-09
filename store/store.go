package store

import (
	"sync"

	"github.com/pranotobudi/graphql-checkout/database"
	"github.com/pranotobudi/graphql-checkout/graph/model"
)

type Store struct {
	AddedProducts []model.Product
}

var StoreInstance *Store
var once sync.Once

func GetStore() *Store {
	once.Do(func() {
		StoreInstance = &Store{
			AddedProducts: []model.Product{},
		}
	})

	return StoreInstance
}

// Get all products from database
func (s *Store) GetProducts() ([]*model.Product, error) {
	// Get products from database
	postgres := database.GetDB()
	products, err := postgres.GetAllProducts()
	if err != nil {
		return nil, err
	}
	// return all products
	return products, nil
}

// Add product to Cart
func (s *Store) AddToCart(sku string, qty int) (string, error) {
	// get product from db
	postgres := database.GetDB()
	product, err := postgres.GetProduct(sku)
	if err != nil {
		return "", err
	}

	// add to global StoreInstance variable as much as quantity
	for i := 0; i <= qty; i++ {
		s.AddedProducts = append(s.AddedProducts, *product)
	}

	// return success message
	return "Product successfully added to cart", nil
}
