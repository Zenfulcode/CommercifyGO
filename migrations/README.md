# Migration Cleanup - Legacy Schema Removal

## Overview

This directory contains a major cleanup of the Commercify database migrations. The original 34 migration files have been consolidated into a single, clean migration that represents the current state of the database schema without any legacy components.

## What Was Cleaned Up

### 1. **Deprecated Cart System**
- **Removed**: `carts` and `cart_items` tables that were deprecated in migration 021
- **Replaced with**: Modern `checkouts` and `checkout_items` system
- **Benefit**: Eliminates confusion and legacy code paths

### 2. **Consolidated Currency Support**
- **Removed**: Multiple migrations (014, 017, 027) that added and fixed currency support
- **Consolidated**: All currency features into the main schema
- **Benefit**: Clean multi-currency implementation from the start

### 3. **Fixed Data Types**
- **Removed**: Migration 011 that converted DECIMAL to BIGINT for money fields
- **Implemented**: All money fields use BIGINT (cents) from the beginning
- **Benefit**: No floating-point precision issues

### 4. **Eliminated Redundant Fixes**
- **Removed**: Migrations 018, 019, 022, 024, 025, 026, 028 that were constraint fixes
- **Implemented**: Proper constraints from the start
- **Benefit**: No need for iterative fixes

### 5. **Removed Legacy Features**
- **Removed**: Disabled cart triggers and archive tables
- **Removed**: Temporary fixes and workarounds
- **Benefit**: Clean codebase without technical debt

## New Schema Features

### Modern Architecture
- ✅ **Checkouts** instead of carts (supports guest checkout and expiration)
- ✅ **Multi-currency** support with proper exchange rates
- ✅ **Product variants** with individual pricing
- ✅ **Comprehensive shipping** system with zones and methods
- ✅ **Discount system** with various types and constraints
- ✅ **Payment transactions** logging
- ✅ **Webhooks** for external integrations
- ✅ **Friendly IDs** for products and orders

### Data Integrity
- ✅ **Proper constraints** on all monetary fields (positive values)
- ✅ **Foreign key constraints** with appropriate cascade behaviors
- ✅ **Unique constraints** where needed (session IDs, SKUs, etc.)
- ✅ **Indexes** for optimal query performance

### Developer Experience
- ✅ **Automatic timestamps** with triggers
- ✅ **Friendly ID generation** for user-facing identifiers
- ✅ **Comprehensive comments** documenting all tables and important columns
- ✅ **Consistent naming** following project conventions

## Migration Strategy

### For New Deployments
Use the new `000001_initial_schema.sql` migration. This creates a clean, modern database schema.

### For Existing Deployments
⚠️ **IMPORTANT**: Do not apply this to existing production databases. This cleanup is intended for:
1. New development environments
2. Fresh staging environments
3. New deployments

For existing production systems, continue using the existing migration sequence in the `legacy/` directory.

## Files Structure

```
migrations/
├── 000001_initial_schema.up.sql     # New consolidated schema
├── 000001_initial_schema.down.sql   # Rollback for new schema
├── legacy/                          # Original migrations (preserved)
│   ├── 000001_create_tables.up.sql
│   ├── 000001_create_tables.down.sql
│   ├── 000002_add_product_variants.up.sql
│   ├── ... (all 34 original migrations)
└── README.md                        # This file
```

## Key Improvements

1. **Performance**: Proper indexing strategy from the start
2. **Maintainability**: Single source of truth for schema
3. **Reliability**: No legacy components or workarounds
4. **Scalability**: Multi-currency and multi-tenant ready
5. **Developer Experience**: Clear documentation and naming

## Next Steps

1. Test the new schema in a development environment
2. Verify all application code works with the consolidated schema
3. Update any documentation that references old migration numbers
4. Consider this the new baseline for future migrations

---

*This cleanup was performed on June 29, 2025, consolidating 34 migration files into a single, clean schema.*
