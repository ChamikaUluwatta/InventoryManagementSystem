ALTER TABLE locations            ALTER COLUMN tenant_id DROP DEFAULT;
ALTER TABLE categories           ALTER COLUMN tenant_id DROP DEFAULT;
ALTER TABLE companies            ALTER COLUMN tenant_id DROP DEFAULT;
ALTER TABLE products             ALTER COLUMN tenant_id DROP DEFAULT;
ALTER TABLE inventories          ALTER COLUMN tenant_id DROP DEFAULT;
ALTER TABLE supplier_returns      ALTER COLUMN tenant_id DROP DEFAULT;
ALTER TABLE supplier_return_items ALTER COLUMN tenant_id DROP DEFAULT;
