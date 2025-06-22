-- Create hubs table
CREATE TABLE hubs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    location TEXT,
    tenant_id TEXT,
    seller_id TEXT
);

-- Create skus table
CREATE TABLE skus (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT,
    description TEXT,
    tenant_id TEXT,
    seller_id TEXT
);

-- Create inventories table
CREATE TABLE inventories (
    id SERIAL PRIMARY KEY,
    product_id TEXT NOT NULL,
    sku TEXT NOT NULL,
    location TEXT NOT NULL,
    tenant_id TEXT,
    seller_id TEXT,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_sku_location UNIQUE (sku, location)
);
