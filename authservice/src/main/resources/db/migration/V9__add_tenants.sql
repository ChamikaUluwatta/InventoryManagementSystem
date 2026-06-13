CREATE TABLE tenants (
    tenant_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_name VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE users ADD COLUMN tenant_id UUID NOT NULL
    REFERENCES tenants(tenant_id);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
