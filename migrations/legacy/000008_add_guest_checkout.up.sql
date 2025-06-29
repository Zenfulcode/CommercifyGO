-- Consolidate all guest checkout related migrations in one file

-- 1. Add session_id column to carts table for guest carts
ALTER TABLE carts ADD COLUMN IF NOT EXISTS session_id VARCHAR(255) NULL;

-- 2. Create index on session_id for efficient lookups
CREATE INDEX IF NOT EXISTS idx_carts_session_id ON carts(session_id);

-- 3. Make user_id optional in carts table (NULL for guest carts)
ALTER TABLE carts ALTER COLUMN user_id DROP NOT NULL;

-- 4. Add guest information to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS guest_email VARCHAR(255) NULL;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS guest_phone VARCHAR(100) NULL;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS guest_full_name VARCHAR(255) NULL;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS is_guest_order BOOLEAN DEFAULT FALSE;

-- 5. Make user_id optional in orders table (NULL for guest orders)
ALTER TABLE orders ALTER COLUMN user_id DROP NOT NULL;

-- 6. Drop the existing foreign key constraint for orders (if it exists)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.table_constraints 
    WHERE constraint_name = 'orders_user_id_fkey' AND table_name = 'orders'
  ) THEN
    ALTER TABLE orders DROP CONSTRAINT orders_user_id_fkey;
  END IF;
END $$;

-- 7. Add the constraint back with ON DELETE SET NULL option and allow nulls
ALTER TABLE orders ADD CONSTRAINT orders_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

-- 8. Drop the existing foreign key constraint for carts (if it exists)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.table_constraints 
    WHERE constraint_name = 'carts_user_id_fkey' AND table_name = 'carts'
  ) THEN
    ALTER TABLE carts DROP CONSTRAINT carts_user_id_fkey;
  END IF;
END $$;

-- 9. Add the constraint back with ON DELETE SET NULL option and allow nulls
ALTER TABLE carts ADD CONSTRAINT carts_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;