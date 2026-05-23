package com.inventory.auth.service.Auth;

import java.time.Duration;
import java.util.UUID;

import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class TokenRevocationService {

    private static final String ACTIVE_TOKEN_PREFIX = "refresh:active:";
    private static final String REVOKED_TOKEN_PREFIX = "refresh:revoked:";
    private static final String PERMISSION_VERSION_PREFIX = "auth:user:version:";
    /*we can use only refresh:active here instead of using both active and revoked but decision is to keep both for clarity 
    in future we can track the usage of each token*/

    private final StringRedisTemplate redisTemplate;
    private final ObjectMapper objectMapper;

    public TokenRevocationService(StringRedisTemplate redisTemplate) {
        this.redisTemplate = redisTemplate;
        this.objectMapper = new ObjectMapper();
    }

    public void storeRefreshToken(String token, UUID userId, String email, long ttlMs) {
        String key = ACTIVE_TOKEN_PREFIX + token;
        try {
            String value = objectMapper.writeValueAsString(new RefreshTokenData(userId, email));
            redisTemplate.opsForValue().set(key, value, Duration.ofMillis(ttlMs));
        } catch (Exception e) {
            throw new RuntimeException("Failed to store refresh token", e);
        }
    }

    public RefreshTokenData getRefreshTokenData(String token) {
        String key = ACTIVE_TOKEN_PREFIX + token;
        String value = redisTemplate.opsForValue().get(key);
        if (value == null) {
            return null;
        }
        try {
            return objectMapper.readValue(value, RefreshTokenData.class);
        } catch (Exception e) {
            return null;
        }
    }

    public void revokeRefreshToken(String token, long ttlMs) {
        String activeKey = ACTIVE_TOKEN_PREFIX + token;
        String revokedKey = REVOKED_TOKEN_PREFIX + token;

        redisTemplate.delete(activeKey);

        redisTemplate.opsForValue().set(revokedKey, "1", Duration.ofMillis(ttlMs));
    }

    public void deleteRefreshToken(String token) {
        String activeKey = ACTIVE_TOKEN_PREFIX + token;

        redisTemplate.delete(activeKey);
    }

    public boolean isRevoked(String token) {
        String key = REVOKED_TOKEN_PREFIX + token;
        return Boolean.TRUE.equals(redisTemplate.hasKey(key));
    }

    public void setPermissionVersion(UUID userId, Integer version) {
        String key = PERMISSION_VERSION_PREFIX + userId.toString();
        redisTemplate.opsForValue().set(key, version.toString());
    }

    public record RefreshTokenData(UUID userId, String email) {}
}
