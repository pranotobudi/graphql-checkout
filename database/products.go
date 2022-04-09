package database

import (
	"fmt"
	"log"

	"github.com/pranotobudi/graphql-checkout/graph/model"
)

func (pg *PostgresDB) GetAllProducts() ([]*model.Product, error) {
	query := fmt.Sprintf(`
		SELECT * FROM products;
	`)

	// execute query
	rows, err := pg.DB.Query(query)
	if err != nil {
		return nil, err
	}

	products := []*model.Product{}
	for rows.Next() {
		p := model.Product{}
		err = rows.Scan(&p.Sku, &p.Name, &p.Price, &p.InventoryQty, &p.PromoType)
		if err != nil {
			log.Println("SQL row scan failed, no row matched: ", err)
			return nil, err
		}
		log.Println("product: ", p)
		products = append(products, &p)
	}

	log.Println("Success to execute SQL query, GetAllProducts success: ")
	return products, nil

}
