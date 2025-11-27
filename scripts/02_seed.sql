-- Clear Data & Reset ID
TRUNCATE TABLE users, products, carts, cart_items, orders, order_items RESTART IDENTITY CASCADE;

-- Create Users (Password: 123456)
INSERT INTO users (email, password, first_name, last_name, role) VALUES
('admin@example.com',  '$2a$10$Y/M8QY1USL52SzgvC5mLb.OVzZWpDLgHAdLQbI53VWcyZKvNnqB0K', 'Admin', 'System', 'admin'),
('user@example.com',   '$2a$10$Y/M8QY1USL52SzgvC5mLb.OVzZWpDLgHAdLQbI53VWcyZKvNnqB0K', 'John',  'Doe',    'user'),
('seller@example.com', '$2a$10$Y/M8QY1USL52SzgvC5mLb.OVzZWpDLgHAdLQbI53VWcyZKvNnqB0K', 'Shop',  'Owner',  'seller');
-- ID 1 = Admin
-- ID 2 = User
-- ID 3 = Seller

-- Create Products
INSERT INTO products (name, description, price, stock, sku, owner_id) VALUES
-- Product Admin (ID 1)
('iPhone 15 Pro',    'Titanium design, A17 Pro chip', 41900.00, 10, 'IP15-PRO-TI', 1),
('MacBook Air M3',   'Supercharged by M3',            39900.00, 5,  'MAC-AIR-M3',  1),

-- Product Seller (ID 3)
('Mechanical Keyboard', 'Blue Switch, RGB Light',      2500.00,  20, 'KEY-MECH-RGB', 3),
('Gaming Mouse',        'Wireless, 20000 DPI',         1200.00,  15, 'MSE-GAME-WL',  3),
('4K Monitor 27"',      'IPS Panel, 144Hz',            8900.00,  8,  'MON-4K-27',    3);
