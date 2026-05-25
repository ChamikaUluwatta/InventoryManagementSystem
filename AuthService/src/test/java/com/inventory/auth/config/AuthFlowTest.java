package com.inventory.auth.config;

import com.inventory.auth.dto.LoginRequest;
import com.inventory.auth.model.User;
import com.inventory.auth.service.Auth.AuthService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.http.MediaType;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.context.TestPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.cookie;

@SpringBootTest
@AutoConfigureMockMvc
@TestPropertySource(properties = {
    "app.jwt.access-token-expiration-ms=900000",
    "app.jwt.refresh-token-expiration-ms=2592000000",
    "app.inventory.service-url=http://localhost:9999"
})
@Testcontainers
class AuthFlowTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine");

    @Container
    static GenericContainer<?> redis = new GenericContainer<>("redis:7-alpine")
            .withExposedPorts(6379);

    @DynamicPropertySource
    static void redisProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.data.redis.host", redis::getHost);
        registry.add("spring.data.redis.port", () -> redis.getMappedPort(6379));
    }

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private AuthService authService;

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Autowired
    private StringRedisTemplate redisTemplate;

    private final ObjectMapper objectMapper = new ObjectMapper();

    private User testUser;

    @BeforeEach
    void seedData() {
        cleanUp();
        jdbcTemplate.update("INSERT INTO roles (name) VALUES (?)", "ADMIN");
        jdbcTemplate.update("INSERT INTO roles (name) VALUES (?)", "VIEWER");
        jdbcTemplate.update("INSERT INTO permissions (name) VALUES (?)", "products:read");
        jdbcTemplate.update("INSERT INTO permissions (name) VALUES (?)", "products:write");

        jdbcTemplate.update("INSERT INTO role_permissions (role_id, permission_id) " +
                "SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'ADMIN'");

        testUser = authService.createUser("test@example.com", "password123");
    }

    @AfterEach
    void cleanUp() {
        jdbcTemplate.update("DELETE FROM user_roles");
        jdbcTemplate.update("DELETE FROM role_permissions");
        jdbcTemplate.update("DELETE FROM users");
        jdbcTemplate.update("DELETE FROM permissions");
        jdbcTemplate.update("DELETE FROM roles");
    }

    @Test
    void loginWithValidCredentialsReturnsToken() throws Exception {
        LoginRequest request = new LoginRequest("test@example.com", "password123");

        mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.accessToken").exists())
                .andExpect(jsonPath("$.tokenType").value("Bearer"))
                .andExpect(jsonPath("$.expiresIn").value(900))
                .andExpect(jsonPath("$.email").value("test@example.com"))
                .andExpect(jsonPath("$.permissions").isArray())
                .andExpect(cookie().exists("refresh_token"));
    }

    @Test
    void loginWithInvalidPasswordReturns401() throws Exception {
        LoginRequest request = new LoginRequest("test@example.com", "wrongpassword");

        mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.message").value("Invalid email or password"));
    }

    @Test
    void loginWithNonExistentUserReturns401() throws Exception {
        LoginRequest request = new LoginRequest("nonexistent@example.com", "password123");

        mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.message").value("Invalid email or password"));
    }

    @Test
    void loginWithMissingEmailReturns400() throws Exception {
        LoginRequest request = new LoginRequest("", "password123");

        mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }

    @Test
    void refreshWithValidCookieReturnsNewToken() throws Exception {
        LoginRequest loginRequest = new LoginRequest("test@example.com", "password123");

        String refreshToken = mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andReturn().getResponse().getCookie("refresh_token").getValue();

        mockMvc.perform(post("/api/v1/auth/refresh")
                        .cookie(new jakarta.servlet.http.Cookie("refresh_token", refreshToken)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.accessToken").exists())
                .andExpect(jsonPath("$.tokenType").value("Bearer"));
    }

    @Test
    void refreshRotatesTokenAndNewTokenWorksForSubsequentRefresh() throws Exception {
        LoginRequest loginRequest = new LoginRequest("test@example.com", "password123");

        String firstToken = mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isOk())
                .andExpect(cookie().exists("refresh_token"))
                .andReturn().getResponse().getCookie("refresh_token").getValue();

        String secondToken = mockMvc.perform(post("/api/v1/auth/refresh")
                        .cookie(new jakarta.servlet.http.Cookie("refresh_token", firstToken)))
                .andExpect(status().isOk())
                .andExpect(cookie().exists("refresh_token"))
                .andReturn().getResponse().getCookie("refresh_token").getValue();

        assertThat(redisTemplate.hasKey("refresh:active:" + secondToken)).isTrue();

        mockMvc.perform(post("/api/v1/auth/refresh")
                        .cookie(new jakarta.servlet.http.Cookie("refresh_token", secondToken)))
                .andExpect(status().isOk())
                .andExpect(cookie().exists("refresh_token"));
    }

    @Test
    void refreshWithoutCookieReturns401() throws Exception {
        mockMvc.perform(post("/api/v1/auth/refresh"))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void logoutRevokesRefreshToken() throws Exception {
        LoginRequest loginRequest = new LoginRequest("test@example.com", "password123");

        String refreshToken = mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andReturn().getResponse().getCookie("refresh_token").getValue();

        assertThat(redisTemplate.opsForValue().get("auth:user:jwt:version:" + testUser.getId()))
                .isEqualTo("1");

        mockMvc.perform(post("/api/v1/auth/logout")
                        .cookie(new jakarta.servlet.http.Cookie("refresh_token", refreshToken)))
                .andExpect(status().isNoContent());

        // Verify JWT version was incremented on logout
        assertThat(redisTemplate.opsForValue().get("auth:user:jwt:version:" + testUser.getId()))
                .isEqualTo("2");

        // Try to refresh with the revoked token
        mockMvc.perform(post("/api/v1/auth/refresh")
                        .cookie(new jakarta.servlet.http.Cookie("refresh_token", refreshToken)))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void logoutWithoutCookieReturns204() throws Exception {
        mockMvc.perform(post("/api/v1/auth/logout"))
                .andExpect(status().isNoContent());
    }

    @Test
    void jwksEndpointReturnsPublicKey() throws Exception {
        mockMvc.perform(get("/api/v1/auth/.well-known/jwks.json"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.keys").isArray())
                .andExpect(jsonPath("$.keys[0].kty").value("RSA"))
                .andExpect(jsonPath("$.keys[0].kid").exists())
                .andExpect(jsonPath("$.keys[0].use").value("sig"))
                .andExpect(jsonPath("$.keys[0].alg").value("RS256"))
                .andExpect(jsonPath("$.keys[0].n").exists())
                .andExpect(jsonPath("$.keys[0].e").exists());
    }

    @Test
    void accessTokenContainsJwtVersion() throws Exception {
        LoginRequest request = new LoginRequest("test@example.com", "password123");

        String response = mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        String accessToken = objectMapper.readTree(response).get("accessToken").asText();

        // Parse the JWT to verify jwt_version claim
        String[] parts = accessToken.split("\\.");
        String payload = new String(java.util.Base64.getUrlDecoder().decode(parts[1]));
        assertThat(payload).contains("jwt_version");
        assertThat(payload).contains("\"jwt_version\":1");
    }

    @Test
    void assignRoleIncrementsJwtVersion() throws Exception {

        // Get VIEWER role ID
        Integer viewerRoleId = jdbcTemplate.queryForObject(
                "SELECT id FROM roles WHERE name = 'VIEWER'",
                Integer.class);

        // Assign role
        authService.assignRoleToUser(testUser.getId(), viewerRoleId);

        // Verify version incremented in Redis
        String redisVersion = redisTemplate.opsForValue().get("auth:user:jwt:version:" + testUser.getId());
        assertThat(redisVersion).isEqualTo("2");
    }

    @Test
    void removeRoleIncrementsJwtVersion() throws Exception {
        // Get VIEWER role ID
        Integer viewerRoleId = jdbcTemplate.queryForObject(
                "SELECT id FROM roles WHERE name = 'VIEWER'",
                Integer.class);

        // Assign then remove role
        authService.assignRoleToUser(testUser.getId(), viewerRoleId);
        authService.removeRoleFromUser(testUser.getId(), viewerRoleId);

        // Verify version incremented in Redis
        String redisVersion = redisTemplate.opsForValue().get("auth:user:jwt:version:" + testUser.getId());
        assertThat(redisVersion).isEqualTo("3");
    }

    @Test
    void addPermissionToRoleIncrementsVersionForAffectedUsers() throws Exception {
        // Create a new permission
        jdbcTemplate.update("INSERT INTO permissions (name) VALUES (?)", "inventory:read");
        Integer permissionId = jdbcTemplate.queryForObject(
                "SELECT id FROM permissions WHERE name = 'inventory:read'",
                Integer.class);

        // Get ADMIN role ID
        Integer adminRoleId = jdbcTemplate.queryForObject(
                "SELECT id FROM roles WHERE name = 'ADMIN'",
                Integer.class);

        // Add permission to role
        authService.addPermissionToRole(adminRoleId, permissionId);

        // Verify version incremented for user with ADMIN role
        String redisVersion = redisTemplate.opsForValue().get("auth:user:jwt:version:" + testUser.getId());
        assertThat(redisVersion).isEqualTo("2");
    }

    @Test
    void refreshReturnsUpdatedJwtVersion() throws Exception {
        LoginRequest loginRequest = new LoginRequest("test@example.com", "password123");

        // Login to get initial token
        String response = mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        String initialToken = objectMapper.readTree(response).get("accessToken").asText();

        // Get VIEWER role ID and assign it (increments version)
        Integer viewerRoleId = jdbcTemplate.queryForObject(
                "SELECT id FROM roles WHERE name = 'VIEWER'",
                Integer.class);
        authService.assignRoleToUser(testUser.getId(), viewerRoleId);

        // Refresh token
        String refreshToken = mockMvc.perform(post("/api/v1/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andReturn().getResponse().getCookie("refresh_token").getValue();

        String refreshResponse = mockMvc.perform(post("/api/v1/auth/refresh")
                        .cookie(new jakarta.servlet.http.Cookie("refresh_token", refreshToken)))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        String newToken = objectMapper.readTree(refreshResponse).get("accessToken").asText();

        // Verify new token has updated version
        String[] parts = newToken.split("\\.");
        String payload = new String(java.util.Base64.getUrlDecoder().decode(parts[1]));
        assertThat(payload).contains("\"jwt_version\":2");
    }
}
