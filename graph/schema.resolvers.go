package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/pranotobudi/graphql-checkout/graph/generated"
	"github.com/pranotobudi/graphql-checkout/graph/model"
)

func (r *mutationResolver) AddToCart(ctx context.Context, input model.AddedProduct) (string, error) {
	result := "Success, Product added to cart"
	return result, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	dummyProduct := model.Product{
		Sku:          "ABC",
		Name:         "Product 1",
		Price:        30.00,
		InventoryQty: 10,
	}
	products = append(products, &dummyProduct)

	return products, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkout(ctx context.Context) (*model.CheckoutReport, error) {
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
