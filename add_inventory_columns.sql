-- Add tenant_id and seller_id columns to inventories table
ALTER TABLE inventories ADD COLUMN IF NOT EXISTS tenant_id TEXT;
ALTER TABLE inventories ADD COLUMN IF NOT EXISTS seller_id TEXT;

-- Add sample inventory data for testing
-- This matches the data in test.csv: SKU23, HU001, quantity 10

INSERT INTO inventories (product_id, sku, location, tenant_id, seller_id, quantity, created_at, updated_at)
VALUES 
    ('PROD001', 'SKU23', 'HU001', 'TE001', 'SEL001', 15, NOW(), NOW()),
    ('PROD002', 'SKU24', 'HU001', 'TE001', 'SEL001', 20, NOW(), NOW()),
    ('PROD003', 'SKU25', 'HU002', 'TE001', 'SEL001', 8, NOW(), NOW()),
    ('PROD004', 'SKU23', 'HU002', 'TE001', 'SEL001', 12, NOW(), NOW());

-- Add sample SKU data
INSERT INTO skus (code, name, description, tenant_id, seller_id)
VALUES 
    ('SKU23', 'Laptop Model X', 'High-performance laptop', 'TE001', 'SEL001'),
    ('SKU24', 'Mouse Wireless', 'Wireless optical mouse', 'TE001', 'SEL001'),
    ('SKU25', 'Keyboard Mechanical', 'Mechanical gaming keyboard', 'TE001', 'SEL001');

-- Add sample Hub data
INSERT INTO hubs (name, location, tenant_id, seller_id)
VALUES 
    ('Warehouse A', 'HU001', 'TE001', 'SEL001'),
    ('Warehouse B', 'HU002', 'TE001', 'SEL001'); 