-- Drop triggers
DROP TRIGGER IF EXISTS update_applications_updated_at ON applications;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_applications_search_vector ON applications;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS applications_search_vector_update();

-- Drop indexes
DROP INDEX IF EXISTS applications_search_idx;
DROP INDEX IF EXISTS applications_developer_id_idx;
DROP INDEX IF EXISTS applications_category_idx;
DROP INDEX IF EXISTS applications_price_idx;
DROP INDEX IF EXISTS applications_rating_idx;
DROP INDEX IF EXISTS applications_downloads_idx;
DROP INDEX IF EXISTS downloads_user_id_idx;
DROP INDEX IF EXISTS downloads_application_id_idx;
DROP INDEX IF EXISTS downloads_timestamp_idx;

-- Drop tables
DROP TABLE IF EXISTS downloads;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS categories; 