DROP POLICY IF EXISTS tenant_isolation ON supplier_return_items;
DROP POLICY IF EXISTS tenant_isolation ON supplier_returns;
DROP POLICY IF EXISTS tenant_isolation ON inventories;
DROP POLICY IF EXISTS tenant_isolation ON products;
DROP POLICY IF EXISTS tenant_isolation ON locations;
DROP POLICY IF EXISTS tenant_isolation ON companies;
DROP POLICY IF EXISTS tenant_isolation ON categories;

ALTER TABLE supplier_return_items DISABLE ROW LEVEL SECURITY;
ALTER TABLE supplier_returns DISABLE ROW LEVEL SECURITY;
ALTER TABLE inventories DISABLE ROW LEVEL SECURITY;
ALTER TABLE products DISABLE ROW LEVEL SECURITY;
ALTER TABLE locations DISABLE ROW LEVEL SECURITY;
ALTER TABLE companies DISABLE ROW LEVEL SECURITY;
ALTER TABLE categories DISABLE ROW LEVEL SECURITY;

ALTER TABLE supplier_return_items DROP COLUMN tenant_id;
ALTER TABLE supplier_returns DROP COLUMN tenant_id;
ALTER TABLE inventories DROP COLUMN tenant_id;
ALTER TABLE products DROP COLUMN tenant_id;
ALTER TABLE locations DROP COLUMN tenant_id;
ALTER TABLE companies DROP COLUMN tenant_id;
ALTER TABLE categories DROP COLUMN tenant_id;
