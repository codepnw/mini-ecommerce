-- ลบ Index ทั้งหมด
DROP INDEX IF EXISTS idx_unique_guest_cart_session;
DROP INDEX IF EXISTS idx_unique_active_cart_user;
DROP INDEX IF EXISTS idx_carts_session_id;
DROP INDEX IF EXISTS idx_carts_user_id;

-- ลบ Foreign Key ของตารางลูก (cart_items) ก่อนเปลี่ยน Type กลับ
ALTER TABLE cart_items
DROP CONSTRAINT IF EXISTS cart_items_cart_id_fkey;

-- ลบ Foreign Key ของตารางแม่ (carts)
ALTER TABLE carts
DROP CONSTRAINT IF EXISTS cart_user_id_fkey, -- Drop FK ON DELETE SET NULL
ADD CONSTRAINT cart_user_id_fkey             -- Add FK เดิม (สมมติว่าเป็น NO ACTION)
FOREIGN KEY (user_id) REFERENCES users(id);


-- ย้อนกลับชนิดข้อมูลในตารางลูก (cart_items)
ALTER TABLE cart_items
ALTER COLUMN cart_id TYPE BIGINT USING NULL;
-- เราใช้ NULL เพื่อล้างข้อมูล UUID ใหม่ในคอลัมน์นี้ให้เป็นค่าว่าง เพราะแปลงกลับไม่ได้

-- ย้อนกลับชนิดข้อมูลในตารางแม่ (carts)
ALTER TABLE carts
ALTER COLUMN cart_id DROP DEFAULT;
ALTER TABLE carts
ALTER COLUMN cart_id TYPE BIGINT USING CAST(SUBSTRING(cart_id::text FROM 1 FOR 15) AS BIGINT);
-- พยายามแปลง UUID กลับเป็น BIGINT ด้วยการตัดสตริง (อาจทำให้เกิดค่าซ้ำซ้อนหรือผิดพลาด)

-- สร้าง Foreign Key ของตารางลูก (cart_items) กลับคืนมา
ALTER TABLE cart_items
ADD CONSTRAINT cart_items_cart_id_fkey
FOREIGN KEY (cart_id) REFERENCES carts (cart_id)
ON DELETE CASCADE;


-- ลบคอลัมน์ที่เพิ่มเข้ามา (ต้องแยกคำสั่ง)
ALTER TABLE carts DROP COLUMN IF EXISTS session_id;
ALTER TABLE carts DROP COLUMN IF EXISTS status;

-- เปลี่ยนชื่อคอลัมน์กลับเป็น id
ALTER TABLE carts RENAME COLUMN cart_id TO id;

-- ลบ Type และ Extension
DROP TYPE cart_status CASCADE;
DROP EXTENSION "uuid-ossp" CASCADE;
