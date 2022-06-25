package store

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"sync"

// 	"github.com/asvvvad/exchange"
// 	"github.com/pranotobudi/graphql-checkout/database"
// 	"github.com/pranotobudi/graphql-checkout/graph/model"
// )

// type StoreGetter interface {
// 	GetStoreService() *StoreService
// }
// type ProductGetter interface {
// 	GetProducts() ([]*model.Product, error)
// }
// type IStoreService interface {
// 	StoreGetter
// 	ProductGetter
// }

// type StoreService struct {
// 	Repo  IStoreRepository
// 	Store Store
// }

// type Store struct {
// 	Carts map[string]*Cart
// }

// var StoreInstance *Store
// var once sync.Once

// func NewStore() *Store {
// 	once.Do(func() {
// 		StoreInstance = &Store{
// 			Carts: make(map[string]*Cart),
// 		}
// 	})

// 	return StoreInstance
// }
// func NewStoreService(store Store, repo IStoreRepository) *StoreService {
// 	return &StoreService{
// 		Repo:  repo,
// 		Store: store,
// 	}
// }

// func (s *StoreService) addCart(userID string) {
// 	s.Store.Carts[userID] = NewCart(userID)
// }

// func (s *StoreService) removeCart(userID string) {
// 	delete(s.Store.Carts, userID)
// }

// // Get all products from database
// func (s *StoreService) GetProducts() ([]*model.Product, error) {
// 	// Get products from database
// 	postgres := database.GetDB()
// 	products, err := postgres.GetAllProducts()
// 	if err != nil {
// 		return nil, err
// 	}
// 	// return all products
// 	return products, nil
// }

// type CartAdder interface {
// 	AddToCart(sku string, qty int) (string, error)
// }
// type CheckoutGetter interface {
// 	GetCheckout(cur string) (*model.CheckoutReport, error)
// }
// type SummaryGetter interface {
// 	GetCartSummary() ([]*model.CartProduct, error)
// }
// type ICartService interface {
// 	CartAdder
// 	CheckoutGetter
// 	SummaryGetter
// }

// type Cart struct {
// 	UserID         string
// 	CartProducts   []*model.CartProduct
// 	AddedProducts  []model.Product
// 	BonusProducts  []model.Product
// 	CartSummary    []*model.CartProduct
// 	CheckoutReport model.CheckoutReport
// }

// func NewCart(userID string) *Cart {
// 	cart := Cart{
// 		UserID:        userID,
// 		CartProducts:  []*model.CartProduct{},
// 		AddedProducts: []model.Product{},
// 		BonusProducts: []model.Product{},
// 		CartSummary:   []*model.CartProduct{},
// 		CheckoutReport: model.CheckoutReport{
// 			Items: []*model.ProductName{},
// 			Total: 0,
// 		},
// 	}
// 	return &cart
// }

// // Add product to Cart
// func (c *StoreService) AddToCart(userID string, sku string, qty int) (string, error) {
// 	log.Println("Store - AddToCart")

// 	// check cart availability
// 	_, found := c.Store.Carts[userID]
// 	if !found {
// 		log.Println("AddToCart Resolver NewCart")
// 		c.Store.Carts[userID] = NewCart(userID)
// 	}

// 	// get product from db
// 	product, err := c.Repo.GetProduct(sku)
// 	if err != nil {
// 		return "", err
// 	}
// 	log.Println("product: ", product)

// 	// check Inventory
// 	if (product.InventoryQty - qty) < 0 {
// 		return "", fmt.Errorf("current inventory only %d", product.InventoryQty)
// 	}

// 	// add to global StoreInstance variable as much as quantity
// 	for i := 0; i < qty; i++ {
// 		c.Store.Carts[userID].AddedProducts = append(c.Store.Carts[userID].AddedProducts, *product)
// 	}
// 	log.Println("AddedProducts: ", c.Store.Carts[userID].AddedProducts)
// 	// return success message
// 	return "Product successfully added to cart", nil
// }

// // Get all products from database
// func (c *StoreService) GetCheckout(userID string, cur string) (*model.CheckoutReport, error) {
// 	log.Println("GetCheckout")
// 	// clear CheckoutReport - fresh CheckoutReport for current transaction
// 	c.Store.Carts[userID].CheckoutReport = model.CheckoutReport{
// 		Items: []*model.ProductName{},
// 		Total: 0,
// 	}

// 	// if CartSummary empty, load CartSummary (user directly to checkout page without visiting cartSummary page)
// 	if len(c.Store.Carts[userID].CartSummary) == 0 {
// 		c.LoadCartSummary(userID)
// 	}

// 	fmt.Println("cartSummary: ", c.Store.Carts[userID].CartSummary)

// 	// checkout bought product
// 	for _, cartProduct := range c.Store.Carts[userID].CartSummary {
// 		log.Println("cartProduct: ", cartProduct)

// 		c.CountProductPrice(userID, *cartProduct)

// 	}
// 	// checkout bonus product, adjust price based on Rasberry bonus
// 	c.AdjustToBonusProduct(userID)

// 	// }

// 	// clear CartSummary, BonusProducts, AddedProduct - fresh for future transaction
// 	c.Store.Carts[userID].CartSummary = []*model.CartProduct{}
// 	c.Store.Carts[userID].AddedProducts = []model.Product{}
// 	c.Store.Carts[userID].BonusProducts = []model.Product{}

// 	// checkout formatting
// 	c.CheckoutReportFormatting(userID, cur)
// 	// return report
// 	return &c.Store.Carts[userID].CheckoutReport, nil
// }

// // Get all products summary in cart
// func (c *StoreService) GetCartSummary(userID string) ([]*model.CartProduct, error) {
// 	// clear previous CartSummary, we need to count from the latest update of c.cart.AddedProducts
// 	c.Store.Carts[userID].CartSummary = []*model.CartProduct{}

// 	log.Println("Store - GetCartSummary")
// 	// iterate all addedProduct and add to cartSummary
// 	for _, product := range c.Store.Carts[userID].AddedProducts {
// 		err := c.AddToProductsInCartSummary(userID, product)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return c.Store.Carts[userID].CartSummary, nil
// }

// // LoadCartSummary - Load all products summary in cart internally, no return value
// func (c *StoreService) LoadCartSummary(userID string) error {
// 	log.Println("LoadCartSummary")
// 	// clear previous CartSummary, we need to count from the latest update of c.AddedProducts
// 	c.Store.Carts[userID].CartSummary = []*model.CartProduct{}

// 	// iterate all addedProduct and add to cartSummary
// 	for _, product := range c.Store.Carts[userID].AddedProducts {
// 		err := c.AddToProductsInCartSummary(userID, product)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // AddToCartSummary will modify cartSummary values.
// func (c *StoreService) AddToProductsInCartSummary(userID string, product model.Product) error {
// 	// base
// 	if len(c.Store.Carts[userID].CartSummary) == 0 {
// 		cartProduct := model.CartProduct{
// 			Product:   &product,
// 			TotalItem: 1,
// 		}
// 		c.Store.Carts[userID].CartSummary = append(c.Store.Carts[userID].CartSummary, &cartProduct)
// 		return nil
// 	}
// 	// if available increment TotalItem, if not available add as the last elmnt
// 	idx := c.FindCartSummarySkuIdx(userID, product.Sku)
// 	// idx not found
// 	if idx == -1 {
// 		cartProduct := model.CartProduct{
// 			Product:   &product,
// 			TotalItem: 1,
// 		}
// 		c.Store.Carts[userID].CartSummary = append(c.Store.Carts[userID].CartSummary, &cartProduct)
// 		return nil
// 	}

// 	// idx found, increment TotalItem
// 	c.Store.Carts[userID].CartSummary[idx].TotalItem++
// 	return nil
// }

// // FindCartSummarySkuIdx will return index of the element with the same sku or return -1.
// func (c *StoreService) FindCartSummarySkuIdx(userID string, sku string) int {
// 	for idx, cartProduct := range c.Store.Carts[userID].CartSummary {
// 		if cartProduct.Product.Sku == sku {
// 			return idx
// 		}
// 	}
// 	// not found
// 	return -1
// }

// type PriceCounter interface {
// 	CountProductPrice(userID string, cartProduct model.CartProduct) error
// }
// type PromoType struct{}

// type PromoType1 struct{}
// type PromoType2 struct{}
// type PromoType3 struct{}
// type PromoType4 struct{}
// type PromoType5 struct{}

// func (c *PromoType1) CountProductPrice(userID string, cartProduct model.CartProduct) error {
// 	ProcessPromoType1(userID, cartProduct)
// }

// // reference: https://www.codementor.io/@uditrastogi/replace-conditional-statements-if-else-or-switch-with-polymorphism-ryl8mx4ns
// // https://golangbot.com/polymorphism/
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

// // ProcessNoPromo
// func (c *StoreService) ProcessNoPromo(userID string, cartProduct model.CartProduct) error {
// 	log.Println("ProcessNoPromo")
// 	// count price
// 	price := float64(cartProduct.TotalItem) * cartProduct.Product.Price
// 	log.Println("ProcessNoPromo price: ", price)

// 	// add to CheckoutReport
// 	c.Store.Carts[userID].CheckoutReport.Total += math.Round(price*100) / 100 //it will round it to 2 digit decimal
// 	productName := model.ProductName{
// 		Name: cartProduct.Product.Name,
// 	}
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 	}
// 	return nil
// }

// // ProcessPromoType1: Process product with this rule:
// // Each sale of a MacBook Pro comes with a free Raspberry Pi B
// func (c *StoreService) ProcessPromoType1(userID string, cartProduct model.CartProduct) error {
// 	log.Println("ProcessPromoType1")
// 	// count price
// 	discountPrice := float64(cartProduct.TotalItem) * cartProduct.Product.Price
// 	log.Println("ProcessPromoType2 discountPrice: ", discountPrice)

// 	// add bonus product
// 	postgres := database.GetDB()

// 	// 1.1. get bonus product sku related with this product
// 	promoType1, err := postgres.GetPromoType1(cartProduct.Product.Sku)
// 	if err != nil {
// 		return err
// 	}

// 	// 1.2. get bonus product from db
// 	product, err := postgres.GetProduct(promoType1.PromoSku)
// 	if err != nil {
// 		return err
// 	}

// 	// 1.3. add bonus product to store BonusProducts
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].BonusProducts = append(c.Store.Carts[userID].BonusProducts, *product)
// 	}
// 	log.Println("bonus product: ", c.Store.Carts[userID].BonusProducts)

// 	// add to CheckoutReport
// 	c.Store.Carts[userID].CheckoutReport.Total += math.Round(discountPrice*100) / 100 //it will round it to 2 digit decimal
// 	productName := model.ProductName{
// 		Name: cartProduct.Product.Name,
// 	}
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 	}
// 	return nil
// }

// // ProcessPromoType2: Process product with this rule:
// // Buy 3 Google Homes for the price of 2
// func (c *StoreService) ProcessPromoType2(userID string, cartProduct model.CartProduct) error {
// 	log.Println("ProcessPromoType2")
// 	// count price
// 	var price float64
// 	if cartProduct.TotalItem >= 3 {
// 		price = float64(cartProduct.TotalItem)*cartProduct.Product.Price - float64(int(cartProduct.TotalItem/3))*cartProduct.Product.Price
// 		log.Println("ProcessPromoType2 discountPrice: ", price)

// 	} else {
// 		price = float64(cartProduct.TotalItem) * cartProduct.Product.Price
// 	}
// 	// add to CheckoutReport
// 	c.Store.Carts[userID].CheckoutReport.Total += math.Round(price*100) / 100 //it will round it to 2 digit decimal
// 	productName := model.ProductName{
// 		Name: cartProduct.Product.Name,
// 	}
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 	}
// 	return nil
// }

// // ProcessPromoType3: Process product with this rule:
// // Buying more than 3 Alexa Speakers will have a 10% discount on all Alexa speakers
// func (c *StoreService) ProcessPromoType3(userID string, cartProduct model.CartProduct) error {
// 	log.Println("ProcessPromoType3")
// 	// count price
// 	var discountPrice float64
// 	log.Println("cartProduct", cartProduct)
// 	if cartProduct.TotalItem >= 3 {
// 		//discount 10% for all
// 		discountPrice = float64(cartProduct.TotalItem) * cartProduct.Product.Price * (0.9)
// 	} else {
// 		discountPrice = float64(cartProduct.TotalItem) * cartProduct.Product.Price
// 	}

// 	log.Println("ProcessPromoType3 discountPrice: ", discountPrice)

// 	// add to CheckoutReport
// 	c.Store.Carts[userID].CheckoutReport.Total += math.Round(discountPrice*100) / 100 //it will round it to 2 digit decimal
// 	productName := model.ProductName{
// 		Name: cartProduct.Product.Name,
// 	}
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 	}

// 	return nil
// }

// // ProcessPromoType4: Process product with this rule:
// // Each sale of MacBook Pro will be discounted flat $200
// func (c *StoreService) ProcessPromoType4(userID string, cartProduct model.CartProduct) error {
// 	log.Println("ProcessPromoType4")
// 	// count price
// 	var discountPrice float64
// 	log.Println("cartProduct", cartProduct)
// 	discountPrice = cartProduct.Product.Price - 200

// 	log.Println("ProcessPromoType4 discountPrice: ", discountPrice)

// 	// add to CheckoutReport
// 	c.Store.Carts[userID].CheckoutReport.Total += math.Round(discountPrice*100) / 100 //it will round it to 2 digit decimal
// 	productName := model.ProductName{
// 		Name: cartProduct.Product.Name,
// 	}
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 	}

// 	return nil
// }

// // ProcessPromoType5: Process product with this rule:
// // Each purchase of 2 MacBook, the customer will get free one item on the cart that
// // is less than $50, if there are two items that are below $50, the customer will get
// // free of one the biggest of two
// func (c *StoreService) ProcessPromoType5(userID string, cartProduct model.CartProduct) error {
// 	log.Println("ProcessPromoType5")
// 	// find items < $50
// 	lessFiftyProduct := []model.Product{}
// 	for _, item := range c.Store.Carts[userID].AddedProducts {
// 		if item.Price < 50 {
// 			lessFiftyProduct = append(lessFiftyProduct, item)
// 		}
// 	}
// 	// find the biggest price & add as bonus
// 	if len(lessFiftyProduct) > 0 {
// 		maxIdx := 0
// 		for i, item := range lessFiftyProduct {
// 			if item.Price > lessFiftyProduct[maxIdx].Price {
// 				maxIdx = i
// 			}
// 		}
// 		// add as bonus
// 		totalBonus := int(cartProduct.TotalItem / 2)
// 		for i := 0; i < totalBonus; i++ {
// 			c.Store.Carts[userID].BonusProducts = append(c.Store.Carts[userID].BonusProducts, lessFiftyProduct[maxIdx])
// 		}
// 	}

// 	// count price
// 	var price float64
// 	log.Println("cartProduct", cartProduct)
// 	price = cartProduct.Product.Price * float64(cartProduct.TotalItem)

// 	log.Println("ProcessPromoType3 discountPrice: ", price)

// 	// add to CheckoutReport
// 	c.Store.Carts[userID].CheckoutReport.Total += math.Round(price*100) / 100 //it will round it to 2 digit decimal
// 	productName := model.ProductName{
// 		Name: cartProduct.Product.Name,
// 	}
// 	for i := 0; i < cartProduct.TotalItem; i++ {
// 		c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 	}

// 	return nil
// }

// func (c *StoreService) AdjustToBonusProduct(userID string) error {
// 	log.Println("CountPriceBonusProduct")
// 	// assumption: bonus only Rasberry (sku==234234), so all slice element is the same product.
// 	// check product availability on CartSummary
// 	if len(c.Store.Carts[userID].BonusProducts) > 0 {
// 		for idx, cartProduct := range c.Store.Carts[userID].CartSummary {
// 			if "234234" == cartProduct.Product.Sku {
// 				// Rasberry is added to cart
// 				if len(c.Store.Carts[userID].BonusProducts) > c.Store.Carts[userID].CartSummary[idx].TotalItem {
// 					c.Store.Carts[userID].CheckoutReport.Total -= float64(c.Store.Carts[userID].CartSummary[idx].TotalItem) * c.Store.Carts[userID].CartSummary[idx].Product.Price
// 					c.Store.Carts[userID].CartSummary[idx].TotalItem = len(c.Store.Carts[userID].BonusProducts)
// 				} else {
// 					c.Store.Carts[userID].CheckoutReport.Total -= float64(len(c.Store.Carts[userID].BonusProducts)) * c.Store.Carts[userID].CartSummary[idx].Product.Price
// 				}

// 				return nil
// 			}
// 		}

// 		// product not available on CartSummary
// 		// add to CartSummary for free
// 		productName := model.ProductName{
// 			Name: c.Store.Carts[userID].BonusProducts[0].Name,
// 		}
// 		for i := 0; i < len(c.Store.Carts[userID].BonusProducts); i++ {
// 			c.Store.Carts[userID].CheckoutReport.Items = append(c.Store.Carts[userID].CheckoutReport.Items, &productName)
// 		}
// 	}

// 	return nil
// }

// // CheckoutReportFormatting will format Total as 2 digit decimal
// func (c *StoreService) CheckoutReportFormatting(userID string, cur string) error {
// 	if cur == "IDR" {
// 		// covert from $ to IDR
// 		ex := exchange.New("USD")
// 		bigFloatConvTimes100, _ := ex.ConvertTo("IDR", int(c.Store.Carts[userID].CheckoutReport.Total*100))
// 		float64ConvTimes100, _ := bigFloatConvTimes100.Float64()
// 		conv := float64ConvTimes100 / float64(100)
// 		fmt.Printf("Total:%d conv:%.2f", int(c.Store.Carts[userID].CheckoutReport.Total), conv)

// 		c.Store.Carts[userID].CheckoutReport.Total = conv
// 	}
// 	c.Store.Carts[userID].CheckoutReport.Total = math.Round(c.Store.Carts[userID].CheckoutReport.Total*100) / 100

// 	return nil
// }
