ALTER TABLE cart_items DROP COLUMN price_at_add;

ALTER TABLE cart_items
ADD CONSTRAINT cart_items_cart_id_fkey
FOREIGN KEY (cart_id) REFERENCES carts (cart_id)
ON DELETE CASCADE;
