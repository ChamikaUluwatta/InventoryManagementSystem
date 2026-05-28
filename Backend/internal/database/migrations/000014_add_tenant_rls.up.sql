ALTER TABLE categories
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

ALTER TABLE companies
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

ALTER TABLE products
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

ALTER TABLE locations
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

ALTER TABLE inventories
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

ALTER TABLE supplier_returns
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

ALTER TABLE supplier_return_items
    ADD COLUMN tenant_id UUID NOT NULL,
    ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON categories
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation ON companies
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation ON products
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation ON locations
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation ON inventories
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation ON supplier_returns
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation ON supplier_return_items
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);
