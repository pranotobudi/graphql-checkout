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
