package store

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/asvvvad/exchange"
	"github.com/pranotobudi/graphql-checkout/database"
	"github.com/pranotobudi/graphql-checkout/graph/model"
)

type StoreGetter interface {
	GetStoreService() *StoreService
}
type ProductGetter interface {
	GetProducts() ([]*model.Product, error)
}
type IStoreService interface {
	StoreGetter
	ProductGetter
}

type StoreService struct {
	Repo  IStoreRepository
	Store Store
}

type Store struct {
	Carts map[string]*Cart
}

var StoreInstance *Store
var once sync.Once

func NewStore() *Store {
	once.Do(func() {
		StoreInstance = &Store{
			Carts: make(map[string]*Cart),
		}
	})

	return StoreInstance
}
func NewStoreService(store Store, repo IStoreRepository) *StoreService {
	return &StoreService{
		Repo:  repo,
		Store: store,
	}
}

// Get all products from database
func (s *StoreService) GetProducts() ([]*model.Product, error) {
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
func (s *StoreService) AddToCart(userID string, sku string, qty int) (string, error) {
	log.Println("Store - AddToCart")

	_, found := s.Store.Carts[userID]
	if !found {
		log.Println("AddToCart Resolver NewCart")
		s.Store.Carts[userID] = NewCart(userID)
	}

	// get product from db
	product, err := s.Repo.GetProduct(sku)
	if err != nil {
		return "", err
	}
	log.Println("product: ", product)

	// check Inventory
	if (product.InventoryQty - qty) < 0 {
		return "", fmt.Errorf("current inventory only %d", product.InventoryQty)
	}
	_, ok := s.Store.Carts[userID].CartProducts[sku]
	if ok {
		s.Store.Carts[userID].CartProducts[sku].TotalItem += qty
		s.Store.Carts[userID].CartProducts[sku].TotalPricePerProduct += product.Price * float64(qty)
	} else {
		cartProduct := model.CartProduct{
			Product:              product,
			TotalItem:            qty,
			TotalPricePerProduct: product.Price * float64(qty),
		}
		s.Store.Carts[userID].CartProducts[sku] = &cartProduct
	}

	return "Product successfully added to cart", nil
}

type CartAdder interface {
	AddToCart(sku string, qty int) (string, error)
}
type CheckoutGetter interface {
	GetCheckout(cur string) (*model.CheckoutReport, error)
}

type ICartService interface {
	CartAdder
	CheckoutGetter
}

type Cart struct {
	UserID       string
	CartProducts map[string]*model.CartProduct
}

func NewCart(userID string) *Cart {
	cart := Cart{
		UserID:       userID,
		CartProducts: make(map[string]*model.CartProduct),
	}
	return &cart
}

// Get all products from database
func (c *Cart) GetCheckout(cur string, repo IStoreRepository) (*model.CheckoutReport, error) {
	log.Println("GetCheckout")
	checkoutReport := model.CheckoutReport{
		Items:      []string{},
		TotalPrice: 0,
	}
	for _, cartProduct := range c.CartProducts {
		promotionProcessor := NewPromotionProcessor(cartProduct.Product.PromoType, *c, repo)
		finalProducts, err := promotionProcessor.promoChecker.CheckPromotion(*cartProduct)
		if err != nil {
			return nil, err
		}

		for _, cartProduct := range finalProducts {
			checkoutReport.Items = append(checkoutReport.Items, cartProduct.Product.Name)
			checkoutReport.TotalPrice += cartProduct.TotalPricePerProduct
		}
		// c.CountProductPrice(userID, *cartProduct)
	}

	// checkout formatting
	checkoutReport, err := c.CheckoutReportFormatting(cur, checkoutReport)
	if err != nil {
		return nil, err
	}
	// return report
	return &checkoutReport, nil
}

// CheckoutReportFormatting will format Total as 2 digit decimal
func (c *Cart) CheckoutReportFormatting(cur string, checkoutReport model.CheckoutReport) (model.CheckoutReport, error) {
	if cur == "IDR" {
		// covert from $ to IDR
		ex := exchange.New("USD")
		bigFloatConvTimes100, _ := ex.ConvertTo("IDR", int(checkoutReport.TotalPrice*100))
		float64ConvTimes100, _ := bigFloatConvTimes100.Float64()
		conv := float64ConvTimes100 / float64(100)
		fmt.Printf("Total:%d conv:%.2f", int(checkoutReport.TotalPrice), conv)

		checkoutReport.TotalPrice = conv
	}
	checkoutReport.TotalPrice = math.Round(checkoutReport.TotalPrice*100) / 100

	return checkoutReport, nil
}

// reference: https://www.codementor.io/@uditrastogi/replace-conditional-statements-if-else-or-switch-with-polymorphism-ryl8mx4ns
// https://golangbot.com/polymorphism/
// func (c *StoreService) CountProductPrice(userID string, cartProduct model.CartProduct) error {
// 	log.Println("CountPriceNormalProduct")
// 	log.Println("promoType: ", cartProduct.Product.PromoType)
// 	switch cartProduct.Product.PromoType {
// 	case 1:
// 		c.ProcessPromoType1(userID, cartProduct)
// 	case 2:
// 		c.ProcessPromoType2(userID, cartProduct)
// 	case 3:
// 		c.ProcessPromoType3(userID, cartProduct)
// 	case 4:
// 		c.ProcessPromoType4(userID, cartProduct)
// 	case 5:
// 		c.ProcessPromoType5(userID, cartProduct)
// 	default:
// 		c.ProcessNoPromo(userID, cartProduct)
// 	}
// 	return nil
// }
