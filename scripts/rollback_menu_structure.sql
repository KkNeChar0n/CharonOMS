-- Menu Structure Rollback Script
-- Rollback menu changes to restore from backup
-- Date: 2026-01-27

USE charonoms;

-- Check if backup exists
SELECT 'Checking backup table...' AS status;
SELECT COUNT(*) AS backup_rows FROM menu_backup_20260127;

-- Temporarily disable foreign key checks
SET FOREIGN_KEY_CHECKS = 0;

-- Clear current menu data
TRUNCATE TABLE menu;

-- Restore from backup
INSERT INTO menu SELECT * FROM menu_backup_20260127;

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;

-- Verification
SELECT 'Rollback completed!' AS status;
SELECT 'Total menus restored:' AS label, COUNT(*) AS count FROM menu;
SELECT * FROM menu ORDER BY parent_id, sort_order;
