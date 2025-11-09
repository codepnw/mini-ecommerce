ALTER TABLE cart_items
ALTER COLUMN cart_id TYPE UUID USING uuid_generate_v4();

ALTER TABLE cart_items
ALTER COLUMN quantity SET NOT NULL;

ALTER TABLE cart_items
ADD COLUMN price_at_add DECIMAL(10, 2) NOT NULL,
ADD COLUMN created_at TIMESTAMP DEFAULT NOW(),
ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();

ALTER TABLE cart_items
ADD CONSTRAINT cart_items_cart_id_fkey
FOREIGN KEY (cart_id) REFERENCES carts (cart_id)
ON DELETE CASCADE;
