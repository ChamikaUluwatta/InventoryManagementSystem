package com.inventory.auth.repository;

import com.inventory.auth.model.User;
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

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine");


    @Test
    void savePersistsUser() {
        User user = new User();
        user.setEmail("save@example.com");
        user.setPassword("hashed-password");

        User saved = userRepository.save(user);

        assertThat(saved.getId()).isNotNull();
        assertThat(saved.getEmail()).isEqualTo("save@example.com");
        assertThat(saved.getPassword()).isEqualTo("hashed-password");
    }

    @Test
    void findByIdReturnsSavedUser() {
        User user = new User();
        user.setEmail("findbyid@example.com");
        user.setPassword("hashed-password");
        User saved = userRepository.save(user);

        User found = userRepository.findById(saved.getId()).orElse(null);

        assertThat(found).isNotNull();
        assertThat(found.getEmail()).isEqualTo("findbyid@example.com");
    }

    @Test
    void findByEmailReturnsSavedUser() {
        User user = new User();
        user.setEmail("findbyemail@example.com");
        user.setPassword("hashed-password");
        userRepository.save(user);

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
