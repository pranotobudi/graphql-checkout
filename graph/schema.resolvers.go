package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/pranotobudi/graphql-checkout/graph/generated"
	"github.com/pranotobudi/graphql-checkout/graph/model"
	"github.com/pranotobudi/graphql-checkout/store"
)

func (r *mutationResolver) AddToCart(ctx context.Context, input model.AddedProduct) (string, error) {
	log.Println("AddToCart Resolver")
	s := store.GetStore()
	result, err := s.AddToCart(input.Sku, input.Qty)
	if err != nil {
		return "", err
	}
	// result := "Success, Product added to cart"
	return result, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	log.Println("Products Resolver")
	s := store.GetStore()
	products, err := s.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *queryResolver) Checkout(ctx context.Context) (*model.CheckoutReport, error) {
	log.Println("Checkout Resolver")
	productName1 := model.ProductName{
		Name: "Product1",
	}
	dummyCheckoutReport := model.CheckoutReport{
		Items: []*model.ProductName{&productName1},
		Total: 3000.99,
	}
	return &dummyCheckoutReport, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CartSummary(ctx context.Context) ([]*model.CartProduct, error) {
	log.Println("CartSummary Resolver")
	s := store.GetStore()
	cartSummary, err := s.GetCartSummary()
	if err != nil {
		return nil, err
	}
	return cartSummary, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
