package store

import (
	"log"
	"sync"

	"github.com/pranotobudi/graphql-checkout/database"
	"github.com/pranotobudi/graphql-checkout/graph/model"
)

type Store struct {
	AddedProducts []model.Product
	CartSummary   []*model.CartProduct
}

var StoreInstance *Store
var once sync.Once

func GetStore() *Store {
	once.Do(func() {
		StoreInstance = &Store{
			AddedProducts: []model.Product{},
			CartSummary:   []*model.CartProduct{},
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
	log.Println("Store - AddToCart")
	// get product from db
	postgres := database.GetDB()
	product, err := postgres.GetProduct(sku)
	if err != nil {
		return "", err
	}
	log.Println("product: ", product)

	// add to global StoreInstance variable as much as quantity
	for i := 0; i < qty; i++ {
		s.AddedProducts = append(s.AddedProducts, *product)
	}
	log.Println("AddedProducts: ", s.AddedProducts)
	// return success message
	return "Product successfully added to cart", nil
}

// Get all products from database
func (s *Store) GetCartSummary() ([]*model.CartProduct, error) {
	// clear previous CartSummary, we need to count from the latest update of s.AddedProducts
	s.CartSummary = []*model.CartProduct{}

	log.Println("Store - GetCartSummary")
	// iterate all addedProduct and add to cartSummary
	for _, product := range s.AddedProducts {
		err := s.AddToCartSummary(product)
		if err != nil {
			return nil, err
		}
	}
	return s.CartSummary, nil
}

// AddToCartSummary will modify cartSummary values.
func (s *Store) AddToCartSummary(product model.Product) error {
	// base
	if len(s.CartSummary) == 0 {
		cartProduct := model.CartProduct{
			Sku:          product.Sku,
			Name:         product.Name,
			Price:        product.Price,
			InventoryQty: product.InventoryQty,
			PromoType:    product.PromoType,
			TotalItem:    1,
		}
		s.CartSummary = append(s.CartSummary, &cartProduct)
		return nil
	}
	// if available increment TotalItem, if not available add as the last elmnt
	idx := s.FindCartSummarySkuIdx(product.Sku)
	// idx not found
	if idx == -1 {
		cartProduct := model.CartProduct{
			Sku:          product.Sku,
			Name:         product.Name,
			Price:        product.Price,
			InventoryQty: product.InventoryQty,
			PromoType:    product.PromoType,
			TotalItem:    1,
		}
		s.CartSummary = append(s.CartSummary, &cartProduct)
		return nil
	}

	// idx found, increment TotalItem
	s.CartSummary[idx].TotalItem++
	return nil
}

// FindCartSummarySkuIdx will return index of the element with the same sku or return -1.
func (s *Store) FindCartSummarySkuIdx(sku string) int {
	for idx, cartProduct := range s.CartSummary {
		if cartProduct.Sku == sku {
			return idx
		}
	}
	// not found
	return -1
}
