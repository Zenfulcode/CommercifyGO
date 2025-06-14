-- Restore the original trigger and function
CREATE OR REPLACE FUNCTION check_product_has_variants()
RETURNS trigger AS $$
BEGIN
    -- For INSERT operations on products, we allow it but warn that variants must be added
    IF TG_OP = 'INSERT' THEN
        RETURN NEW;
    END IF;
    
    -- For DELETE operations on product_variants, ensure at least one variant remains
    IF TG_OP = 'DELETE' THEN
        IF (SELECT COUNT(*) FROM product_variants WHERE product_id = OLD.product_id) <= 1 THEN
            RAISE EXCEPTION 'Cannot delete the last variant of a product. Products must have at least one variant.';
        END IF;
        RETURN OLD;
    END IF;
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Restore the trigger
CREATE TRIGGER prevent_last_variant_deletion
    BEFORE DELETE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION check_product_has_variants();

-- Restore the comment
COMMENT ON TABLE product_variants IS 'Product variants are mandatory. Every product must have at least one variant.';
