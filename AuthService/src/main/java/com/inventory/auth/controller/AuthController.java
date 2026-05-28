package com.inventory.auth.controller;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CookieValue;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.inventory.auth.config.ApiPaths;
import com.inventory.auth.dto.GuestResponse;
import com.inventory.auth.dto.LoginRequest;
import com.inventory.auth.dto.RegisterRequest;
import com.inventory.auth.dto.TokenResponse;
import com.inventory.auth.dto.UserResponse;
import com.inventory.auth.model.User;
import com.inventory.auth.service.Auth.AuthService;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.validation.Valid;

@RestController
@RequestMapping(ApiPaths.AUTH_BASE)
public class AuthController {

    private static final String REFRESH_TOKEN_COOKIE = "refresh_token";

    private final AuthService authService;
    private final boolean refreshCookieSecure;

    public AuthController(
            AuthService authService,
            @Value("${app.auth.refresh-cookie-secure:false}") boolean refreshCookieSecure) {
        this.authService = authService;
        this.refreshCookieSecure = refreshCookieSecure;
    }

    @PostMapping("/register")
    public ResponseEntity<?> register(
            @Valid @RequestBody(required = false) RegisterRequest request,
            @RequestParam(required = false) String type,
            HttpServletResponse response) {

        if ("guest".equals(type)) {
            String email = authService.generateGuestEmail();
            String password = authService.generateRandomPassword();

            User user = authService.createGuestUser(email, password);

            String refreshToken = authService.generateRefreshToken();
            TokenResponse tokenResponse = authService.authenticate(email, password, refreshToken);

            setRefreshTokenCookie(response, refreshToken);

            GuestResponse guestResponse = new GuestResponse(
                    email,
                    password,
                    tokenResponse.accessToken(),
                    tokenResponse.tokenType(),
                    tokenResponse.expiresIn(),
                    tokenResponse.permissions());

            return ResponseEntity.status(HttpStatus.CREATED).body(guestResponse);
        }

        if (request == null) {
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).build();
        }

        User user = authService.createUser(request.email(), request.password());
        UserResponse userResponse = new UserResponse(user.getId(), user.getEmail(), user.getCreatedAt());
        return ResponseEntity.status(HttpStatus.CREATED).body(userResponse);
    }

    @PostMapping("/login")
    public ResponseEntity<TokenResponse> login(
            @Valid @RequestBody LoginRequest request,
            HttpServletResponse response) {
        String refreshToken = authService.generateRefreshToken();
        TokenResponse tokenResponse = authService.authenticate(request.email(), request.password(), refreshToken);

        setRefreshTokenCookie(response, refreshToken);

        return ResponseEntity.ok(tokenResponse);
    }

    @PostMapping("/refresh")
    public ResponseEntity<TokenResponse> refresh(
            @CookieValue(name = REFRESH_TOKEN_COOKIE, required = false) String refreshToken,
            HttpServletResponse response) {
        if (refreshToken == null || refreshToken.isBlank()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
        }

        String newRefreshToken = authService.generateRefreshToken();
        TokenResponse tokenResponse = authService.refresh(refreshToken, newRefreshToken);
        setRefreshTokenCookie(response, newRefreshToken);
        return ResponseEntity.ok(tokenResponse);
    }

    @PostMapping("/logout")
    public ResponseEntity<Void> logout(
            @CookieValue(name = REFRESH_TOKEN_COOKIE, required = false) String refreshToken,
            HttpServletResponse response) {
        if (refreshToken != null && !refreshToken.isBlank()) {
            authService.revokeRefreshToken(refreshToken);
        }

        clearRefreshTokenCookie(response);

        return ResponseEntity.noContent().build();
    }

    private void setRefreshTokenCookie(HttpServletResponse response, String token) {
        Cookie cookie = new Cookie(REFRESH_TOKEN_COOKIE, token);
        cookie.setHttpOnly(true);
        cookie.setPath("/");
        cookie.setMaxAge((int) (authService.getRefreshTokenExpirationMs() / 1000));
        cookie.setSecure(refreshCookieSecure);
        cookie.setAttribute("SameSite", "Lax");
        response.addCookie(cookie);
    }

    private void clearRefreshTokenCookie(HttpServletResponse response) {
        Cookie cookie = new Cookie(REFRESH_TOKEN_COOKIE, "");
        cookie.setHttpOnly(true);
        cookie.setPath("/");
        cookie.setMaxAge(0);
        cookie.setSecure(refreshCookieSecure);
        cookie.setAttribute("SameSite", "Lax");
        response.addCookie(cookie);
    }
}
