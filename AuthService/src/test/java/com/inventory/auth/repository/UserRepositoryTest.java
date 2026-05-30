package com.inventory.auth.repository;

import com.inventory.auth.model.Tenant;
import com.inventory.auth.model.User;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest;
import org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
@Testcontainers
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class UserRepositoryTest {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private TenantRepository tenantRepository;

    private Tenant testTenant;

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine");

    @BeforeEach
    void setUp() {
        testTenant = new Tenant();
        testTenant.setTenantName("test-tenant");
        testTenant = tenantRepository.save(testTenant);
    }

    private User createUser(String email) {
        User user = new User();
        user.setEmail(email);
        user.setPassword("hashed-password");
        user.setTenant(testTenant);
        return user;
    }

    @Test
    void savePersistsUser() {
        User user = createUser("save@example.com");
        User saved = userRepository.save(user);

        assertThat(saved.getId()).isNotNull();
        assertThat(saved.getEmail()).isEqualTo("save@example.com");
        assertThat(saved.getPassword()).isEqualTo("hashed-password");
    }

    @Test
    void findByIdReturnsSavedUser() {
        User saved = userRepository.save(createUser("findbyid@example.com"));

        User found = userRepository.findById(saved.getId()).orElse(null);

        assertThat(found).isNotNull();
        assertThat(found.getEmail()).isEqualTo("findbyid@example.com");
    }

    @Test
    void findByEmailReturnsSavedUser() {
        userRepository.save(createUser("findbyemail@example.com"));

        User found = userRepository.findByEmail("findbyemail@example.com").orElse(null);

        assertThat(found).isNotNull();
        assertThat(found.getEmail()).isEqualTo("findbyemail@example.com");
    }

    @Test
    void findByEmailReturnsEmptyWhenNotFound() {
        var result = userRepository.findByEmail("nonexistent@example.com");

        assertThat(result).isEmpty();
    }
}
