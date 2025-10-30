CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id),
    token TEXT UNIQUE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_auth_token ON auth (token);
CREATE INDEX idx_auth_user_id ON auth (user_id);


 -- Create User Role
CREATE TYPE user_roles AS ENUM ('user', 'seller', 'admin');

ALTER TABLE users ADD COLUMN role user_roles DEFAULT 'user';
