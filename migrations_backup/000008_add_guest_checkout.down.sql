-- Revert all guest checkout related changes

-- 1. First, delete all guest carts to prevent constraint violations
DELETE FROM carts WHERE user_id IS NULL;

-- 2. Delete all guest orders to prevent constraint violations
DELETE FROM orders WHERE user_id IS NULL;

-- 3. Drop the modified foreign key constraint for orders
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_user_id_fkey;

-- 4. Restore the original foreign key constraint for orders without the ON DELETE SET NULL
ALTER TABLE orders ADD CONSTRAINT orders_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

-- 5. Drop the modified foreign key constraint for carts
ALTER TABLE carts DROP CONSTRAINT IF EXISTS carts_user_id_fkey;

-- 6. Restore the original foreign key constraint for carts without the ON DELETE SET NULL
ALTER TABLE carts ADD CONSTRAINT carts_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

-- 7. Make user_id required again in orders table
ALTER TABLE orders ALTER COLUMN user_id SET NOT NULL;

-- 8. Remove guest information columns from orders table
ALTER TABLE orders DROP COLUMN IF EXISTS guest_email;
ALTER TABLE orders DROP COLUMN IF EXISTS guest_phone;
ALTER TABLE orders DROP COLUMN IF EXISTS guest_full_name;
ALTER TABLE orders DROP COLUMN IF EXISTS is_guest_order;

-- 9. Make user_id required again in carts table
ALTER TABLE carts ALTER COLUMN user_id SET NOT NULL;

-- 10. Drop the session_id index
DROP INDEX IF EXISTS idx_carts_session_id;

-- 11. Remove session_id column from carts table
ALTER TABLE carts DROP COLUMN IF EXISTS session_id;