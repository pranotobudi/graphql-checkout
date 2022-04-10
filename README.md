# GRAPHQL CHECKOUT

GraphQL Checkout is an application to manage checkout process in a online store with GraphQL endpoint. 

About This App:
>- This app using Postgres SQL as the database. 
>- Automatic build and run at Heroku platform when code pushed to Github.
>- Automatic build and push docker image to docker hub using Github Action when code pushed to Github. 

### How to Access this app: 
* #1 On the Cloud: https://graphql-checkout.herokuapp.com/
> This app is deployed at Heroku. The main url could be used as GraphQL Query editor.

* #2 Locally: 2 easy step to run GraphQL Checkout app:
>- clone this repo: 
>>- git clone https://github.com/pranotobudi/graphql-checkout.git
>- Run this command from root folder:
>>- docker-compose up


## GraphQL API endpoints

### https://graphql-checkout.herokuapp.com/

> Add Product to Cart (Success)

_Mutation_
```
mutation {
  addToCart(input: {sku: "43N23P", qty: 1})
}

```

_Response (200)_
```
{
    "data": {
        "addToCart": "Product successfully added to cart"
    }
}
```

> Add Product to Cart (Failed, out of stock)

_Mutation_
```
mutation {
  addToCart(input: {sku: "43N23P", qty: 100})
}

```

_Response_
```
    "errors": [
        {
            "message": "current inventory only 5",
            "path": [
                "addToCart"
            ]
        }
    ],
    "data": null
}
```


> Checkout

_Query_
```
query {
	checkout {
        items {
            name
        }
        total
    }
}

```

_Response (200)_
```
{
    "data": {
        "checkout": {
            "items": [
                {
                    "name": "MacBook Pro"
                },
                {
                    "name": "Rasberry Pi B"
                }
            ],
            "total": 5399.99
        }
    }
}
```

### GraphQL API endpoints (Additional)

> Get All Products

_Query_
```
query {
	products {
        sku
        name
        price
        inventory_qty
        promo_type
    }
}

```

_Response (200)_
```
{
    "data": {
        "products": [
            {
                "sku": "120P90",
                "name": "Google Home",
                "price": 49.99,
                "inventory_qty": 10,
                "promo_type": 2
            },
            {
                "sku": "43N23P",
                "name": "MacBook Pro",
                "price": 5399.99,
                "inventory_qty": 5,
                "promo_type": 1
            },
            {
                "sku": "A304SD",
                "name": "Alexa Speaker",
                "price": 109.5,
                "inventory_qty": 10,
                "promo_type": 3
            },
            {
                "sku": "234234",
                "name": "Rasberry Pi B",
                "price": 30,
                "inventory_qty": 2,
                "promo_type": 0
            }
        ]
    }
}
```

> CartSummary 
>- Summary of All Products in the Cart. In the frontend This data can be used for the summary of added to cart products before the checkout button. 

_Query_
```
query {
	cartSummary {
        sku
        name
        price
        inventory_qty
        promo_type
        total_item
    }
}
```

_Response (200)_
```
{
    "data": {
        "cartSummary": [
            {
                "sku": "43N23P",
                "name": "MacBook Pro",
                "price": 5399.99,
                "inventory_qty": 5,
                "promo_type": 1,
                "total_item": 3
            }
        ]
    }
}
```


## Assumptions

* Each product only handle 1 type of promotion
* Buying More than three interpreted as: buying 3 or 4 or 5 ... and so on. (3 is included)
* User is already logged in
* Products is already in the the database (loaded from database/init.sql when the app is run)
