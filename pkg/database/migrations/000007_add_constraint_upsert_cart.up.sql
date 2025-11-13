CREATE UNIQUE INDEX idx_unique_product_in_cart
ON cart_items(cart_id, product_id);
