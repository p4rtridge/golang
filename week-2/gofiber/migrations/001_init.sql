CREATE TABLE IF NOT EXISTS users (
  id        serial,
  username  text  NOT NULL UNIQUE,
  password  text  NOT NULL,
  balance   real  DEFAULT 0.0,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS products (
  id              serial,
  productName     text  NOT NULL,
  productQuantity int   NOT NULL,
  productPrice    real  NOT NULL,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS orders (
  id            serial,
  userId        int       NOT NULL,
  productID     int       NOT NULL,
  orderQuantity int       NOT NULL,
  createdAt     timestamp DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_user    FOREIGN KEY (userId)    REFERENCES users(id),
  CONSTRAINT fk_product FOREIGN KEY (productID) REFERENCES products(id)
);
