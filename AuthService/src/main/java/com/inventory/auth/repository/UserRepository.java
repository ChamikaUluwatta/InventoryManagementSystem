package com.inventory.auth.repository;

import com.inventory.auth.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.UUID;

@Repository
public interface UserRepository extends JpaRepository<User, UUID> {

	Optional<User> findByEmail(String email);

	@Query(
		value = """
			SELECT DISTINCT p.name
			FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			JOIN roles r ON r.id = ur.role_id
			JOIN role_permissions rp ON rp.role_id = r.id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE u.email = :email
		""",
  		nativeQuery = true
	)
	Set<String> findPermissionNamesByEmail(@Param("email") String email);

	@Modifying
	@Transactional
	@Query(
		value = "INSERT INTO user_roles (user_id, role_id) VALUES (:userId, :roleId)",
		nativeQuery = true
	)
	void addRoleToUser(@Param("userId") UUID userId, @Param("roleId") Integer roleId);

	@Modifying
	@Transactional
	@Query(
		value = "DELETE FROM user_roles WHERE user_id = :userId AND role_id = :roleId",
		nativeQuery = true
	)
	void removeRoleFromUser(@Param("userId") UUID userId, @Param("roleId") Integer roleId);

	@Query(
		value = """
			SELECT u.id FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			WHERE ur.role_id = :roleId
			""",
		nativeQuery = true
	)
	List<UUID> findUsersWithRole(@Param("roleId") Integer roleId);

	List<User> findByIsGuestTrueAndCreatedAtBefore(Instant cutoff);
}
