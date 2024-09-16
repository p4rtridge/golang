ALTER TABLE IF EXISTS orders ADD IF NOT EXISTS status order_status NOT NULL DEFAULT 'pending';
