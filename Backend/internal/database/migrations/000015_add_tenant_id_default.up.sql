ALTER TABLE locations            ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
ALTER TABLE categories           ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
ALTER TABLE companies            ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
ALTER TABLE products             ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
ALTER TABLE inventories          ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
ALTER TABLE supplier_returns      ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
ALTER TABLE supplier_return_items ALTER COLUMN tenant_id SET DEFAULT current_setting('app.current_tenant_id')::uuid;
