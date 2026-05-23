package com.inventory.auth.service.Auth;

import com.inventory.auth.exception.CustomAuthException;
import com.inventory.auth.model.User;
import com.inventory.auth.repository.RoleRepository;
import com.inventory.auth.repository.UserRepository;
import org.junit.jupiter.api.Test;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.security.crypto.password.PasswordEncoder;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

public class AuthServiceTests {

    @Test
    void createUser_encodesPasswordAndSaves() {
        PasswordEncoder passwordEncoder = mock(PasswordEncoder.class);
        UserRepository userRepository = mock(UserRepository.class);
        RoleRepository roleRepository = mock(RoleRepository.class);
        JwtService jwtService = mock(JwtService.class);
        TokenRevocationService tokenRevocationService = mock(TokenRevocationService.class);

        AuthService authService = new AuthService(
                passwordEncoder, userRepository, roleRepository,
                jwtService, tokenRevocationService);

        String email = "new@example.com";
        String rawPassword = "plain-password";
        String encodedPassword = "encoded-password";

        when(passwordEncoder.encode(rawPassword)).thenReturn(encodedPassword);

        User savedUser = new User();
        savedUser.setId(UUID.randomUUID());
        savedUser.setEmail(email);
        savedUser.setPassword(encodedPassword);
        when(userRepository.saveAndFlush(any(User.class))).thenReturn(savedUser);

        User result = authService.createUser(email, rawPassword);

        assertEquals(email, result.getEmail());
        assertEquals(encodedPassword, result.getPassword());
        verify(passwordEncoder).encode(rawPassword);
        verify(userRepository).saveAndFlush(any(User.class));
    }

    @Test
    void createUser_duplicateEmail_throwsException() {
        PasswordEncoder passwordEncoder = mock(PasswordEncoder.class);
        UserRepository userRepository = mock(UserRepository.class);
        RoleRepository roleRepository = mock(RoleRepository.class);
        JwtService jwtService = mock(JwtService.class);
        TokenRevocationService tokenRevocationService = mock(TokenRevocationService.class);

        AuthService authService = new AuthService(
                passwordEncoder, userRepository, roleRepository,
                jwtService, tokenRevocationService);

        when(userRepository.saveAndFlush(any())).thenThrow(new DataIntegrityViolationException("duplicate"));

        CustomAuthException ex = assertThrows(CustomAuthException.class, () ->
                authService.createUser("existing@example.com", "password"));

        assertEquals(409, ex.getStatusCode());
        assertEquals("An account with this email already exists", ex.getMessage());
    }
}
