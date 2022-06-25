package store

import (
	"log"
	"sort"

	"github.com/pranotobudi/graphql-checkout/graph/model"
)

type PriceCounter interface {
}

type PromoChecker interface {
	CheckPromotion(cartProduct model.CartProduct) ([]model.CartProduct, error)
}
type PromotionProcessor struct {
	promoChecker PromoChecker
}

func NewPromotionProcessor(promoType int, cart Cart, repo IStoreRepository) PromotionProcessor {
	var checker PromoChecker
	switch promoType {
	case 1:
		checker = &PromoType1{
			Repo: repo,
		}
	case 2:
		checker = &PromoType2{}
	case 3:
		checker = &PromoType3{}
	case 4:
		checker = &PromoType4{}
	case 5:
		checker = &PromoType5{
			cart: cart,
		}
	default:
		checker = &PromoTypeDefault{}
	}
	return PromotionProcessor{
		promoChecker: checker,
	}
}

type PromoTypeDefault struct{}

// ProcessNoPromo
func (c *PromoTypeDefault) CheckPromotion(product model.CartProduct) ([]model.CartProduct, error) {
	log.Println("ProcessNoPromo")
	var result []model.CartProduct
	result = append(result, product)
	return result, nil
}

type PromoType1 struct {
	Repo IStoreRepository
}

// ProcessPromoType1: Process product with this rule:
// Each sale of a MacBook Pro comes with a free Raspberry Pi B
func (c *PromoType1) CheckPromotion(cartProduct model.CartProduct) ([]model.CartProduct, error) {
	log.Println("ProcessPromoType1")

	c.Repo.GetPromoType1(cartProduct.Product.Sku)
	promoProductType1, err := c.Repo.GetPromoType1(cartProduct.Product.Sku)
	if err != nil {
		return nil, err
	}
	promoProduct := model.CartProduct{
		Product:   promoProductType1,
		TotalItem: 1,
	}

	var result []model.CartProduct
	result = append(result, cartProduct, promoProduct)
	return result, nil
}

type PromoType2 struct{}

// ProcessPromoType2: Process product with this rule:
// Buy 3 Google Homes for the price of 2
func (c *PromoType2) CheckPromotion(cartProduct model.CartProduct) ([]model.CartProduct, error) {
	log.Println("ProcessPromoType2")
	var price float64
	if cartProduct.TotalItem >= 3 {
		price = float64(cartProduct.TotalItem)*cartProduct.Product.Price - float64(int(cartProduct.TotalItem/3))*cartProduct.Product.Price
		log.Println("ProcessPromoType2 discountPrice: ", cartProduct.TotalPricePerProduct)

	} else {
		price = float64(cartProduct.TotalItem) * cartProduct.Product.Price
	}
	cartProduct.TotalPricePerProduct = price

	var result []model.CartProduct
	result = append(result, cartProduct)
	return result, nil
}

type PromoType3 struct{}

// ProcessPromoType3: Process product with this rule:
// Buying more than 3 Alexa Speakers will have a 10% discount on all Alexa speakers
func (c *PromoType3) CheckPromotion(cartProduct model.CartProduct) ([]model.CartProduct, error) {
	log.Println("ProcessPromoType3")
	// count price
	var discountPrice float64
	log.Println("cartProduct", cartProduct)
	if cartProduct.TotalItem >= 3 {
		discountPrice = float64(cartProduct.TotalItem) * cartProduct.Product.Price * (0.9)
	} else {
		discountPrice = float64(cartProduct.TotalItem) * cartProduct.Product.Price
	}
	cartProduct.TotalPricePerProduct = discountPrice
	log.Println("ProcessPromoType3 discountPrice: ", discountPrice)

	var result []model.CartProduct
	result = append(result, cartProduct)
	return result, nil
}

type PromoType4 struct{}

// ProcessPromoType4: Process product with this rule:
// Each sale of MacBook Pro will be discounted flat $200
func (c *PromoType4) CheckPromotion(cartProduct model.CartProduct) ([]model.CartProduct, error) {
	log.Println("ProcessPromoType4")
	// count price
	cartProduct.Product.Price = cartProduct.Product.Price - 200

	var result []model.CartProduct
	result = append(result, cartProduct)
	return result, nil
}

type PromoType5 struct {
	cart Cart
}

// ProcessPromoType5: Process product with this rule:
// Each purchase of 2 MacBook, the customer will get free one item on the cart that
// is less than $50, if there are two items that are below $50, the customer will get
// free of one the biggest of two
func (c *PromoType5) CheckPromotion(cartProduct model.CartProduct) ([]model.CartProduct, error) {
	log.Println("ProcessPromoType5")
	// find items < $50
	lessFiftyProducts := []model.CartProduct{}
	for _, cartProduct := range c.cart.CartProducts {
		if cartProduct.Product.Price < 50 {
			lessFiftyProducts = append(lessFiftyProducts, *cartProduct)
		}
	}
	lessFiftyProducts = sortLessFiftyProductsByPrice(lessFiftyProducts)

	totalBonus := int(cartProduct.TotalItem / 2)
	if totalBonus > len(lessFiftyProducts) {
		totalBonus = len(lessFiftyProducts)
	}

	var result []model.CartProduct
	result = append(result, cartProduct)
	for i := 0; i < totalBonus; i++ {
		result = append(result, lessFiftyProducts[i])
	}
	return result, nil
}

type LessFiftyProductsByPrice []model.CartProduct

func (a LessFiftyProductsByPrice) Len() int           { return len(a) }
func (a LessFiftyProductsByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LessFiftyProductsByPrice) Less(i, j int) bool { return a[i].Product.Price > a[j].Product.Price }

func sortLessFiftyProductsByPrice(cartProducts []model.CartProduct) []model.CartProduct {
	sort.Sort(LessFiftyProductsByPrice(cartProducts))
	return cartProducts
}
