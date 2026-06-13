INSERT INTO roles (name) VALUES
    ('ADMIN'),
    ('VIEWER')
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (name) VALUES
    ('products:read'),
    ('products:write'),
    ('categories:read'),
    ('categories:write'),
    ('companies:read'),
    ('companies:write'),
    ('locations:read'),
    ('locations:write'),
    ('inventories:read'),
    ('inventories:write'),
    ('supplier_returns:read'),
    ('supplier_returns:write')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'ADMIN'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'VIEWER' AND p.name LIKE '%:read'
ON CONFLICT DO NOTHING;
