package com.inventory.auth.repository;

import com.inventory.auth.model.Role;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;

public interface RoleRepository extends JpaRepository<Role, Integer> {
    Optional<Role> findByName(String name);

    @Modifying
    @Transactional
    @Query(
        value = "INSERT INTO role_permissions (role_id, permission_id) VALUES (:roleId, :permissionId)",
        nativeQuery = true
    )
    void addPermissionToRole(@Param("roleId") Integer roleId, @Param("permissionId") Integer permissionId);

    @Modifying
    @Transactional
    @Query(
        value = "DELETE FROM role_permissions WHERE role_id = :roleId AND permission_id = :permissionId",
        nativeQuery = true
    )
    void removePermissionFromRole(@Param("roleId") Integer roleId, @Param("permissionId") Integer permissionId);
}
