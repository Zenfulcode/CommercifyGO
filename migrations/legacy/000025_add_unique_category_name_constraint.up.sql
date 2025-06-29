-- Add unique constraint on category name and parent_id
-- This ensures that category names are unique within the same parent level
-- (including at the root level when parent_id is NULL)

-- First, let's handle any potential duplicate data by updating conflicting categories
-- Handle duplicates where parent_id is NULL (root categories)
WITH root_duplicates AS (
    SELECT id, name,
           ROW_NUMBER() OVER (PARTITION BY name ORDER BY id) as rn
    FROM categories
    WHERE parent_id IS NULL
)
UPDATE categories 
SET name = categories.name || '_' || root_duplicates.rn
FROM root_duplicates 
WHERE categories.id = root_duplicates.id 
AND root_duplicates.rn > 1;

-- Handle duplicates where parent_id is NOT NULL (child categories)
WITH child_duplicates AS (
    SELECT id, name, parent_id,
           ROW_NUMBER() OVER (PARTITION BY name, parent_id ORDER BY id) as rn
    FROM categories
    WHERE parent_id IS NOT NULL
)
UPDATE categories 
SET name = categories.name || '_' || child_duplicates.rn
FROM child_duplicates 
WHERE categories.id = child_duplicates.id 
AND child_duplicates.rn > 1;

-- Create a unique partial index for root categories (where parent_id IS NULL)
CREATE UNIQUE INDEX unique_root_category_name 
ON categories (name) 
WHERE parent_id IS NULL;

-- Create a unique partial index for child categories (where parent_id IS NOT NULL)
CREATE UNIQUE INDEX unique_child_category_name_parent 
ON categories (name, parent_id) 
WHERE parent_id IS NOT NULL;

-- Create a general index to improve query performance for category lookups
CREATE INDEX IF NOT EXISTS idx_categories_name_parent ON categories(name, parent_id);