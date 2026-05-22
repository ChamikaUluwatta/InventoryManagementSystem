ALTER TABLE user_roles
    DROP CONSTRAINT IF EXISTS user_roles_user_id_fkey,
    DROP CONSTRAINT IF EXISTS user_roles_role_id_fkey,
    ADD CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE role_permissions
    DROP CONSTRAINT IF EXISTS role_permissions_role_id_fkey,
    DROP CONSTRAINT IF EXISTS role_permissions_permission_id_fkey,
    ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;
