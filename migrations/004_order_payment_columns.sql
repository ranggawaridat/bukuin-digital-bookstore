ALTER TABLE orders ADD COLUMN payment_token TEXT;
ALTER TABLE orders ADD COLUMN payment_url TEXT;
ALTER TABLE orders ADD COLUMN payment_method TEXT;
ALTER TABLE orders ADD COLUMN transaction_id TEXT;
ALTER TABLE orders ADD COLUMN paid_at DATETIME;
