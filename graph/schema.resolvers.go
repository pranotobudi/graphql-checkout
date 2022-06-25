package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/pranotobudi/graphql-checkout/graph/generated"
	"github.com/pranotobudi/graphql-checkout/graph/model"
)

func (r *mutationResolver) AddToCart(ctx context.Context, addedProduct model.AddedProduct) (string, error) {
	log.Println("AddToCart Resolver")

	// cart := r.StoreGetter.GetStore().Carts[addedProduct.UserID]
	result, err := r.StoreService.AddToCart(addedProduct.UserID, addedProduct.Sku, addedProduct.Qty)
	if err != nil {
		return "", err
	}
	// result := "Success, Product added to cart"
	return result, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	log.Println("Products Resolver")
	products, err := r.StoreService.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *queryResolver) Checkout(ctx context.Context, userID string, cur string) (*model.CheckoutReport, error) {
	log.Println("Checkout Resolver")
	checkoutReport, err := r.StoreService.Store.Carts[userID].GetCheckout(cur, r.StoreService.Repo)
	if err != nil {
		return nil, err
	}
	return checkoutReport, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
