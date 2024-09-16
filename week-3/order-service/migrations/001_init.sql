CREATE TABLE IF NOT EXISTS users (
  id                serial,
  username          text NOT NULL UNIQUE,
  password          text NOT NULL,
  balance           real DEFAULT 0.0,
  created_at        timestamp DEFAULT NOW(),
  updated_at        timestamp,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS products (
  id          serial,
  name        text  NOT NULL,
  quantity    int   NOT NULL,
  price       real  NOT NULL,
  created_at  timestamp DEFAULT NOW(),
  updated_at  timestamp,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS orders (
  id              serial,
  user_id         int       NOT NULL,
  total_price     real      DEFAULT 0.0,
  created_at      timestamp DEFAULT NOW(),
  updated_at      timestamp,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS order_items (
  order_id      int,
  product_id    int,
  product_name  text,
  product_price real,
  quantity      real,

  PRIMARY KEY (order_id, product_id),

  CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id),
  CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
);
