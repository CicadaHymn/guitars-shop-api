CREATE TYPE order_status AS ENUM(
  'pending',
  'processing',
  'delivered',
  'completed',
  'returned',
  'canceled'
);

CREATE TABLE IF NOT EXISTS orders (
  id INT GENERATED ALWAYS AS IDENTITY
  (START WITH 1000 INCREMENT BY 1)
  PRIMARY KEY,
  phone VARCHAR(20) NOT NULL,
  address TEXT NOT NULL,
  product_id INTEGER,
  status order_status DEFAULT 'pending' NOT NULL
);
 