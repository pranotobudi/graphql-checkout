package store

import (
	"fmt"
	"log"
	"math"

	"github.com/pranotobudi/graphql-checkout/database"
	"github.com/pranotobudi/graphql-checkout/graph/model"
)

type Cart struct {
	UserID         string
	CartProducts   []*model.CartProduct
	AddedProducts  []model.Product
	BonusProducts  []model.Product
	CartSummary    []*model.CartProduct
	CheckoutReport model.CheckoutReport
}

func NewCart(userID string) *Cart {
	cart := Cart{
		UserID:        userID,
		CartProducts:  []*model.CartProduct{},
		AddedProducts: []model.Product{},
		BonusProducts: []model.Product{},
		CartSummary:   []*model.CartProduct{},
		CheckoutReport: model.CheckoutReport{
			Items: []*model.ProductName{},
			Total: 0,
		},
	}
	return &cart
}

// Add product to Cart
func (s *Cart) AddToCart(sku string, qty int) (string, error) {
	log.Println("Store - AddToCart")
	// get product from db
	postgres := database.GetDB()
	product, err := postgres.GetProduct(sku)
	if err != nil {
		return "", err
	}
	log.Println("product: ", product)

	// check Inventory
	if (product.InventoryQty - qty) < 0 {
		return "", fmt.Errorf("current inventory only %d", product.InventoryQty)
	}

	// add to global StoreInstance variable as much as quantity
	for i := 0; i < qty; i++ {
		s.AddedProducts = append(s.AddedProducts, *product)
	}
	log.Println("AddedProducts: ", s.AddedProducts)
	// return success message
	return "Product successfully added to cart", nil
}

// Get all products summary in cart
func (s *Cart) GetCartSummary() ([]*model.CartProduct, error) {
	// clear previous CartSummary, we need to count from the latest update of s.AddedProducts
	s.CartSummary = []*model.CartProduct{}

	log.Println("Store - GetCartSummary")
	// iterate all addedProduct and add to cartSummary
	for _, product := range s.AddedProducts {
		err := s.AddToProductsInCartSummary(product)
		if err != nil {
			return nil, err
		}
	}
	return s.CartSummary, nil
}

// LoadCartSummary - Load all products summary in cart internally, no return value
func (s *Cart) LoadCartSummary() error {
	log.Println("LoadCartSummary")
	// clear previous CartSummary, we need to count from the latest update of s.AddedProducts
	s.CartSummary = []*model.CartProduct{}

	// iterate all addedProduct and add to cartSummary
	for _, product := range s.AddedProducts {
		err := s.AddToProductsInCartSummary(product)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddToCartSummary will modify cartSummary values.
func (s *Cart) AddToProductsInCartSummary(product model.Product) error {
	// base
	if len(s.CartSummary) == 0 {
		cartProduct := model.CartProduct{
			Product:   &product,
			TotalItem: 1,
		}
		s.CartSummary = append(s.CartSummary, &cartProduct)
		return nil
	}
	// if available increment TotalItem, if not available add as the last elmnt
	idx := s.FindCartSummarySkuIdx(product.Sku)
	// idx not found
	if idx == -1 {
		cartProduct := model.CartProduct{
			Product:   &product,
			TotalItem: 1,
		}
		s.CartSummary = append(s.CartSummary, &cartProduct)
		return nil
	}

	// idx found, increment TotalItem
	s.CartSummary[idx].TotalItem++
	return nil
}

// FindCartSummarySkuIdx will return index of the element with the same sku or return -1.
func (s *Cart) FindCartSummarySkuIdx(sku string) int {
	for idx, cartProduct := range s.CartSummary {
		if cartProduct.Product.Sku == sku {
			return idx
		}
	}
	// not found
	return -1
}

// Get all products from database
func (s *Cart) GetCheckout(cur string) (*model.CheckoutReport, error) {
	log.Println("GetCheckout")
	// clear CheckoutReport - fresh CheckoutReport for current transaction
	s.CheckoutReport = model.CheckoutReport{
		Items: []*model.ProductName{},
		Total: 0,
	}

	// if CartSummary empty, load CartSummary (user directly to checkout page without visiting cartSummary page)
	if len(s.CartSummary) == 0 {
		s.LoadCartSummary()
	}

	fmt.Println("cartSummary: ", s.CartSummary)

	// checkout bought product
	for _, cartProduct := range s.CartSummary {
		log.Println("cartProduct: ", cartProduct)
		s.CountProductPrice(*cartProduct)

	}
	// checkout bonus product, adjust price based on Rasberry bonus
	s.AdjustToBonusProduct()

	// }

	// clear CartSummary, BonusProducts, AddedProduct - fresh for future transaction
	s.CartSummary = []*model.CartProduct{}
	s.AddedProducts = []model.Product{}
	s.BonusProducts = []model.Product{}

	// checkout formatting
	s.CheckoutReportFormatting(cur)
	// return report
	return &s.CheckoutReport, nil
}

func (s *Cart) CountProductPrice(cartProduct model.CartProduct) error {
	log.Println("CountPriceNormalProduct")
	log.Println("promoType: ", cartProduct.Product.PromoType)
	switch cartProduct.Product.PromoType {
	case 0:
		s.ProcessNoPromo(cartProduct)
	case 1:
		s.ProcessPromoType1(cartProduct)
	case 2:
		s.ProcessPromoType2(cartProduct)
	case 3:
		s.ProcessPromoType3(cartProduct)
	case 4:
		s.ProcessPromoType4(cartProduct)
	case 5:
		s.ProcessPromoType5(cartProduct)
	}
	return nil
}

// ProcessNoPromo
func (s *Cart) ProcessNoPromo(cartProduct model.CartProduct) error {
	log.Println("ProcessNoPromo")
	// count price
	price := float64(cartProduct.TotalItem) * cartProduct.Product.Price
	log.Println("ProcessNoPromo price: ", price)

	// add to CheckoutReport
	s.CheckoutReport.Total += math.Round(price*100) / 100 //it will round it to 2 digit decimal
	productName := model.ProductName{
		Name: cartProduct.Product.Name,
	}
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
	}
	return nil
}

// ProcessPromoType1: Process product with this rule:
// Each sale of a MacBook Pro comes with a free Raspberry Pi B
func (s *Cart) ProcessPromoType1(cartProduct model.CartProduct) error {
	log.Println("ProcessPromoType1")
	// count price
	discountPrice := float64(cartProduct.TotalItem) * cartProduct.Product.Price
	log.Println("ProcessPromoType2 discountPrice: ", discountPrice)

	// add bonus product
	postgres := database.GetDB()

	// 1.1. get bonus product sku related with this product
	promoType1, err := postgres.GetPromoType1(cartProduct.Product.Sku)
	if err != nil {
		return err
	}

	// 1.2. get bonus product from db
	product, err := postgres.GetProduct(promoType1.PromoSku)
	if err != nil {
		return err
	}

	// 1.3. add bonus product to store BonusProducts
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.BonusProducts = append(s.BonusProducts, *product)
	}
	log.Println("bonus product: ", s.BonusProducts)

	// add to CheckoutReport
	s.CheckoutReport.Total += math.Round(discountPrice*100) / 100 //it will round it to 2 digit decimal
	productName := model.ProductName{
		Name: cartProduct.Product.Name,
	}
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
	}
	return nil
}

// ProcessPromoType2: Process product with this rule:
// Buy 3 Google Homes for the price of 2
func (s *Cart) ProcessPromoType2(cartProduct model.CartProduct) error {
	log.Println("ProcessPromoType2")
	// count price
	var price float64
	if cartProduct.TotalItem >= 3 {
		price = float64(cartProduct.TotalItem)*cartProduct.Product.Price - float64(int(cartProduct.TotalItem/3))*cartProduct.Product.Price
		log.Println("ProcessPromoType2 discountPrice: ", price)

	} else {
		price = float64(cartProduct.TotalItem) * cartProduct.Product.Price
	}
	// add to CheckoutReport
	s.CheckoutReport.Total += math.Round(price*100) / 100 //it will round it to 2 digit decimal
	productName := model.ProductName{
		Name: cartProduct.Product.Name,
	}
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
	}
	return nil
}

// ProcessPromoType3: Process product with this rule:
// Buying more than 3 Alexa Speakers will have a 10% discount on all Alexa speakers
func (s *Cart) ProcessPromoType3(cartProduct model.CartProduct) error {
	log.Println("ProcessPromoType3")
	// count price
	var discountPrice float64
	log.Println("cartProduct", cartProduct)
	if cartProduct.TotalItem >= 3 {
		//discount 10% for all
		discountPrice = float64(cartProduct.TotalItem) * cartProduct.Product.Price * (0.9)
	} else {
		discountPrice = float64(cartProduct.TotalItem) * cartProduct.Product.Price
	}

	log.Println("ProcessPromoType3 discountPrice: ", discountPrice)

	// add to CheckoutReport
	s.CheckoutReport.Total += math.Round(discountPrice*100) / 100 //it will round it to 2 digit decimal
	productName := model.ProductName{
		Name: cartProduct.Product.Name,
	}
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
	}

	return nil
}

// ProcessPromoType4: Process product with this rule:
// Each sale of MacBook Pro will be discounted flat $200
func (s *Cart) ProcessPromoType4(cartProduct model.CartProduct) error {
	log.Println("ProcessPromoType4")
	// count price
	var discountPrice float64
	log.Println("cartProduct", cartProduct)
	discountPrice = cartProduct.Product.Price - 200

	log.Println("ProcessPromoType4 discountPrice: ", discountPrice)

	// add to CheckoutReport
	s.CheckoutReport.Total += math.Round(discountPrice*100) / 100 //it will round it to 2 digit decimal
	productName := model.ProductName{
		Name: cartProduct.Product.Name,
	}
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
	}

	return nil
}

// ProcessPromoType5: Process product with this rule:
// Each purchase of 2 MacBook, the customer will get free one item on the cart that
// is less than $50, if there are two items that are below $50, the customer will get
// free of one the biggest of two
func (s *Cart) ProcessPromoType5(cartProduct model.CartProduct) error {
	log.Println("ProcessPromoType5")
	// find items < $50
	lessFiftyProduct := []model.Product{}
	for _, item := range s.AddedProducts {
		if item.Price < 50 {
			lessFiftyProduct = append(lessFiftyProduct, item)
		}
	}
	// find the biggest price & add as bonus
	if len(lessFiftyProduct) > 0 {
		maxIdx := 0
		for i, item := range lessFiftyProduct {
			if item.Price > lessFiftyProduct[maxIdx].Price {
				maxIdx = i
			}
		}
		// add as bonus
		totalBonus := int(cartProduct.TotalItem / 2)
		for i := 0; i < totalBonus; i++ {
			s.BonusProducts = append(s.BonusProducts, lessFiftyProduct[maxIdx])
		}
	}

	// count price
	var price float64
	log.Println("cartProduct", cartProduct)
	price = cartProduct.Product.Price * float64(cartProduct.TotalItem)

	log.Println("ProcessPromoType3 discountPrice: ", price)

	// add to CheckoutReport
	s.CheckoutReport.Total += math.Round(price*100) / 100 //it will round it to 2 digit decimal
	productName := model.ProductName{
		Name: cartProduct.Product.Name,
	}
	for i := 0; i < cartProduct.TotalItem; i++ {
		s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
	}

	return nil
}

func (s *Cart) AdjustToBonusProduct() error {
	log.Println("CountPriceBonusProduct")
	// assumption: bonus only Rasberry (sku==234234), so all slice element is the same product.
	// check product availability on CartSummary
	if len(s.BonusProducts) > 0 {
		for idx, cartProduct := range s.CartSummary {
			if "234234" == cartProduct.Product.Sku {
				// Rasberry is added to cart
				if len(s.BonusProducts) > s.CartSummary[idx].TotalItem {
					s.CheckoutReport.Total -= float64(s.CartSummary[idx].TotalItem) * s.CartSummary[idx].Product.Price
					s.CartSummary[idx].TotalItem = len(s.BonusProducts)
				} else {
					s.CheckoutReport.Total -= float64(len(s.BonusProducts)) * s.CartSummary[idx].Product.Price
				}

				return nil
			}
		}

		// product not available on CartSummary
		// add to CartSummary for free
		productName := model.ProductName{
			Name: s.BonusProducts[0].Name,
		}
		for i := 0; i < len(s.BonusProducts); i++ {
			s.CheckoutReport.Items = append(s.CheckoutReport.Items, &productName)
		}
	}

	return nil
}

// CheckoutReportFormatting will format Total as 2 digit decimal
func (s *Cart) CheckoutReportFormatting(cur string) error {

	s.CheckoutReport.Total = math.Round(s.CheckoutReport.Total*100) / 100

	return nil
}
