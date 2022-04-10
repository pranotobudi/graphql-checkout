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
	log.Println("query: ", query)

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

func (pg *PostgresDB) GetProduct(sku string) (*model.Product, error) {
	query := fmt.Sprintf(`
		SELECT * FROM products WHERE sku='%s';
	`, sku)
	log.Println("query: ", query)

	// execute query
	row := pg.DB.QueryRow(query)

	p := model.Product{}
	err := row.Scan(&p.Sku, &p.Name, &p.Price, &p.InventoryQty, &p.PromoType)
	if err != nil {
		log.Println("SQL row scan failed, no row matched: ", err)
		return nil, err
	}
	log.Println("product: ", p)

	log.Println("Success to execute SQL query, GetProduct success: ")
	return &p, nil

}

type PromoType1 struct {
	ProductSku string `json:"product_sku"`
	PromoSku   string `json:"promo_sku"`
}

func (pg *PostgresDB) GetPromoType1(sku string) (*PromoType1, error) {
	query := fmt.Sprintf(`
		SELECT * FROM promo_type_1 WHERE product_sku='%s';
	`, sku)
	log.Println("query: ", query)

	// execute query
	row := pg.DB.QueryRow(query)

	p := PromoType1{}
	err := row.Scan(&p.ProductSku, &p.PromoSku)
	if err != nil {
		log.Println("SQL row scan failed, no row matched: ", err)
		return nil, err
	}
	log.Println("product: ", p)

	log.Println("Success to execute SQL query, GetPromoType1 success")
	return &p, nil
}
