package com.inventory.auth.service.Auth;

import java.util.List;
import java.util.Set;
import java.util.UUID;

import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.inventory.auth.dto.TokenResponse;
import com.inventory.auth.exception.CustomAuthException;
import com.inventory.auth.model.User;
import com.inventory.auth.repository.RoleRepository;
import com.inventory.auth.repository.UserRepository;
import com.inventory.auth.service.Auth.TokenRevocationService.RefreshTokenData;

@Service
public class AuthService {

    private final PasswordEncoder passwordEncoder;
    private final UserRepository userRepository;
    private final RoleRepository roleRepository;
    private final JwtService jwtService;
    private final TokenRevocationService tokenRevocationService;

    public AuthService(
            PasswordEncoder passwordEncoder,
            UserRepository userRepository,
            RoleRepository roleRepository,
            JwtService jwtService,
            TokenRevocationService tokenRevocationService) {
        this.passwordEncoder = passwordEncoder;
        this.userRepository = userRepository;
        this.roleRepository = roleRepository;
        this.jwtService = jwtService;
        this.tokenRevocationService = tokenRevocationService;
    }

    public User createUser(String email, String rawPassword) {
        User user = new User();
        user.setEmail(email);
        user.setPassword(passwordEncoder.encode(rawPassword));
        try {
            User savedUser = userRepository.saveAndFlush(user);
            return savedUser;
        } catch (DataIntegrityViolationException e) {
            throw new CustomAuthException(409, "An account with this email already exists");
        }
    }

    public TokenResponse authenticate(String email, String password, String refreshToken) {
        User user = userRepository.findByEmail(email)
                .orElseThrow(() -> new BadCredentialsException("Invalid credentials"));

        if (!passwordEncoder.matches(password, user.getPassword())) {
            throw new BadCredentialsException("Invalid credentials");
        }

        Set<String> permissions = userRepository.findPermissionNamesByEmail(email);
        Integer permissionsVersion = user.getPermissionsVersion();
        String accessToken = jwtService.generateAccessToken(user.getId(), email, permissions, permissionsVersion);

        tokenRevocationService.storeRefreshToken(
                refreshToken,
                user.getId(),
                email,
                jwtService.getRefreshTokenExpirationMs());

        return new TokenResponse(
                accessToken,
                "Bearer",
                jwtService.getAccessTokenExpirationMs() / 1000,
                email,
                permissions);
    }

    public TokenResponse refresh(String refreshToken, String newRefreshToken) {
        if (tokenRevocationService.isRevoked(refreshToken)) {
            throw new CustomAuthException(401, "Your session has expired, please log in again");
        }

        RefreshTokenData data = tokenRevocationService.getRefreshTokenData(refreshToken);
        if (data == null) {
            throw new CustomAuthException(401, "Your session has expired, please log in again");
        }

        Set<String> permissions = userRepository.findPermissionNamesByEmail(data.email());
        Integer permissionsVersion = userRepository.getPermissionsVersion(data.userId());
        String accessToken = jwtService.generateAccessToken(data.userId(), data.email(), permissions, permissionsVersion);
        tokenRevocationService.storeRefreshToken(
                newRefreshToken,
                data.userId(),
                data.email(),
                jwtService.getRefreshTokenExpirationMs());
        tokenRevocationService.revokeRefreshToken(refreshToken, jwtService.getRefreshTokenExpirationMs());
        return new TokenResponse(
                accessToken,
                "Bearer",
                jwtService.getAccessTokenExpirationMs() / 1000,
                data.email(),
                permissions);
    }

    public void revokeRefreshToken(String refreshToken) {
        tokenRevocationService.revokeRefreshToken(refreshToken, jwtService.getRefreshTokenExpirationMs());
    }

    public String generateRefreshToken() {
        return jwtService.generateRefreshToken();
    }

    public long getRefreshTokenExpirationMs() {
        return jwtService.getRefreshTokenExpirationMs();
    }

    @Transactional
    public void assignRoleToUser(UUID userId, Integer roleId) {
        userRepository.addRoleToUser(userId, roleId);
        incrementAndBroadcastVersion(userId);
    }

    @Transactional
    public void removeRoleFromUser(UUID userId, Integer roleId) {
        userRepository.removeRoleFromUser(userId, roleId);
        incrementAndBroadcastVersion(userId);
    }

    @Transactional
    public void addPermissionToRole(Integer roleId, Integer permissionId) {
        roleRepository.addPermissionToRole(roleId, permissionId);
        incrementVersionForAllUsersWithRole(roleId);
    }

    @Transactional
    public void removePermissionFromRole(Integer roleId, Integer permissionId) {
        roleRepository.removePermissionFromRole(roleId, permissionId);
        incrementVersionForAllUsersWithRole(roleId);
    }

    private void incrementAndBroadcastVersion(UUID userId) {
        userRepository.incrementPermissionsVersion(userId);
        Integer newVersion = userRepository.getPermissionsVersion(userId);
        tokenRevocationService.setPermissionVersion(userId, newVersion);
    }

    //need to look if this is good pattern
    private void incrementVersionForAllUsersWithRole(Integer roleId) {
        userRepository.incrementPermissionsVersionForUsersWithRole(roleId);
        List<Object[]> users = userRepository.findUsersWithRole(roleId);
        for (Object[] row : users) {
            UUID userId = (UUID) row[0];
            Integer version = (Integer) row[1];
            tokenRevocationService.setPermissionVersion(userId, version);
        }
    }
}
