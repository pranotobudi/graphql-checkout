package store_test

import (
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/pranotobudi/graphql-checkout/graph"
	"github.com/pranotobudi/graphql-checkout/graph/generated"
	"github.com/pranotobudi/graphql-checkout/graph/model"
	"github.com/pranotobudi/graphql-checkout/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedPostgres struct {
	mock.Mock
}

func (p *MockedPostgres) GetProduct(sku string) (*model.Product, error) {
	args := p.Called()
	return args.Get(0).(*model.Product), args.Error(1)
}

// type MockedStore struct {
// 	mock.Mock
// }

// func (s *MockedStore) GetStore() *store.Store {
// 	args := s.Called()
// 	return args.Get(0).(*store.Store)
// }

// func (s *MockedStore) AddToCart(sku string, qty int) (string, error) {
// 	args := s.Called()
// 	return args.String(0), args.Error(1)
// }
func TestAddToCart(t *testing.T) {

	tt := []struct {
		Name       string
		CodeWant   int
		HttpMethod string
	}{
		{
			Name:       "AddToCart Success",
			CodeWant:   http.StatusOK,
			HttpMethod: http.MethodPost,
		},
		{
			Name:       "AddToCart Out Of Stock",
			CodeWant:   http.StatusBadRequest,
			HttpMethod: http.MethodPost,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			// need to mock Store also, it is called inside AddToCart function

			testStore := store.Store{
				Carts: make(map[string]*store.Cart),
			}
			mockedPostgres := new(MockedPostgres)
			product1 := model.Product{
				Sku:          "43N23P",
				Name:         "Google IPhone",
				Price:        5,
				InventoryQty: 10,
			}
			mockedPostgres.On("GetProduct", mock.Anything).Return(&product1, nil)
			storeService := store.NewStoreService(testStore, mockedPostgres)
			// mockedCartAdder.On("AddToCart", mock.Anything).Return("Success", nil)

			resolvers := graph.Resolver{
				StoreService: *storeService,
			}

			c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
			query := `mutation {
				addToCart(added_product: {user_id:"01", sku: "43N23P", qty: 1})
				}`
			// var expectedResponse struct {
			// 	response string
			// }
			var expectedResponse map[string]interface{}
			c.MustPost(query, &expectedResponse)

			assert.Equal(t, "43N23P", storeService.Store.Carts["01"].AddedProducts[0].Sku)
			assert.Equal(t, 10, storeService.Store.Carts["01"].AddedProducts[0].InventoryQty)
		})
	}
}

type MockedProductGetter struct {
	mock.Mock
}

func (s *MockedProductGetter) GetProducts() ([]*model.Product, error) {
	args := s.Called()
	return args.Get(0).([]*model.Product), nil
}
