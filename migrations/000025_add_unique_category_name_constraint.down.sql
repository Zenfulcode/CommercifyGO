-- Remove unique constraint on category name and parent_id

-- Drop the general performance index
DROP INDEX IF EXISTS idx_categories_name_parent;

-- Drop the unique partial indexes
DROP INDEX IF EXISTS unique_child_category_name_parent;
DROP INDEX IF EXISTS unique_root_category_name;