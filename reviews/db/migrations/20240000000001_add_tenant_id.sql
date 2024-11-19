ALTER TABLE reviews ADD COLUMN tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';
CREATE INDEX idx_reviews_tenant_id ON reviews(tenant_id);

-- Update existing records (if any) with a default value
UPDATE reviews SET tenant_id = 'default' WHERE tenant_id = 'default';

-- Make tenant_id required for new records by removing the default
ALTER TABLE reviews ALTER COLUMN tenant_id DROP DEFAULT;
