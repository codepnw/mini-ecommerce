CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE cart_status AS ENUM ('active', 'guest', 'saved', 'ordered');

-- 1. DROP ตารางเดิมทิ้ง (ต้องใช้ CASCADE ถ้ามี FK ชี้มา)
DROP TABLE IF EXISTS carts CASCADE;

-- 2. CREATE ตารางใหม่ทั้งหมดตาม Schema ที่ต้องการ
CREATE TABLE carts (
    cart_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      BIGINT, -- สมมติ user_id เป็น BIGINT
    session_id   UUID,
    status       cart_status NOT NULL DEFAULT 'guest',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign Keys
    CONSTRAINT cart_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 3. สร้าง Index ทั้งหมด
CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_session_id ON carts(session_id);

-- 4. สร้าง Partial Unique Index (เพื่อบังคับให้มีตะกร้าเดียว)
CREATE UNIQUE INDEX idx_unique_active_cart_user
ON carts(user_id, status)
WHERE (status = 'active');

CREATE UNIQUE INDEX idx_unique_guest_cart_session
ON carts(session_id, status)
WHERE (status = 'guest');
