package com.inventory.auth.repository;

import com.inventory.auth.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;
import java.util.Set;
import java.util.UUID;

@Repository
public interface UserRepository extends JpaRepository<User, UUID> {

	Optional<User> findByEmail(String email);

	boolean existsByEmail(String email);

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
		value = "UPDATE users SET permissions_version = permissions_version + 1 WHERE id = :userId",
		nativeQuery = true
	)
	void incrementPermissionsVersion(@Param("userId") UUID userId);

	@Query(
		value = "SELECT permissions_version FROM users WHERE id = :userId",
		nativeQuery = true
	)
	Integer getPermissionsVersion(@Param("userId") UUID userId);
}
