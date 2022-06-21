package graph

import (
	"github.com/pranotobudi/graphql-checkout/store"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	StoreService store.StoreService
}
