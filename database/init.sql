DROP TABLE IF EXISTS promo_type_1;
DROP TABLE IF EXISTS promo_type_2;
DROP TABLE IF EXISTS promo_type_3;
DROP TABLE IF EXISTS products;

CREATE TABLE IF NOT EXISTS products (
    sku VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR (1000),
    price FLOAT,
    inventory_qty INT,
    promo_type INT,
    PRIMARY KEY (sku) 
);

CREATE TABLE IF NOT EXISTS promo_type_1 (
    product_sku VARCHAR(255) NOT NULL UNIQUE,
    promo_sku VARCHAR (255),
    PRIMARY KEY (product_sku), 
    FOREIGN KEY (product_sku) REFERENCES products(sku)
);
CREATE TABLE IF NOT EXISTS promo_type_2 (
    product_sku VARCHAR(255) NOT NULL UNIQUE,
    reference_qty INT,
    promo_qty INT,
    PRIMARY KEY (product_sku),
    FOREIGN KEY (product_sku) REFERENCES products(sku)
);
CREATE TABLE IF NOT EXISTS promo_type_3 (
    product_sku VARCHAR(255) NOT NULL UNIQUE,
    minimum_qty INT,
    discount INT,
    PRIMARY KEY (product_sku),
    FOREIGN KEY (product_sku) REFERENCES products(sku)
);

INSERT INTO products
  ( sku, name, price, inventory_qty, promo_type )
VALUES
  ('120P90', 'Google Home', 49.99, 10, 2), 
  ('43N23P', 'MacBook Pro', 5399.99, 5, 1), 
  ('A304SD', 'Alexa Speaker', 109.99, 10, 3), 
  ('234234', 'Rasberry Pi B', 30.00, 2, 0);

INSERT INTO promo_type_1
  ( product_sku, promo_sku )
VALUES
  ('43N23P', '234234');

INSERT INTO promo_type_2
  ( product_sku, reference_qty, promo_qty )
VALUES
  ('120P90', 3, 2);

INSERT INTO promo_type_3
  ( product_sku, minimum_qty, discount)
VALUES
  ('A304SD', 3, 10);


