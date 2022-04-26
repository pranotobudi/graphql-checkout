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

func (r *mutationResolver) AddToCart(ctx context.Context, addedProduct model.AddedProduct) (string, error) {
	log.Println("AddToCart Resolver")
	cart := store.GetStore().Carts[addedProduct.UserID]
	result, err := cart.AddToCart(addedProduct.Sku, addedProduct.Qty)
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

func (r *queryResolver) Checkout(ctx context.Context, userID string) (*model.CheckoutReport, error) {
	log.Println("Checkout Resolver")
	cart := store.GetStore().Carts[userID]
	checkoutReport, err := cart.GetCheckout()
	if err != nil {
		return nil, err
	}
	return checkoutReport, nil
}

func (r *queryResolver) CartSummary(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	log.Println("CartSummary Resolver")
	cart := store.GetStore().Carts[userID]
	cartSummary, err := cart.GetCartSummary()
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
