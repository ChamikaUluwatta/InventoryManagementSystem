package com.inventory.auth.exception;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(BadCredentialsException.class)
    public ResponseEntity<Map<String, Object>> handleBadCredentials(BadCredentialsException ex) {
        log.warn("Authentication failed: {}", ex.getMessage());
        return ResponseEntity.status(401).body(Map.of(
                "status", 401,
                "error", "Unauthorized",
                "message", "Invalid email or password"));
    }

    @ExceptionHandler(UsernameNotFoundException.class)
    public ResponseEntity<Map<String, Object>> handleUserNotFound(UsernameNotFoundException ex) {
        log.warn("User not found: {}", ex.getMessage());
        return ResponseEntity.status(401).body(Map.of(
                "status", 401,
                "error", "Unauthorized",
                "message", "Invalid email or password"));
    }

    @ExceptionHandler(CustomAuthException.class)
    public ResponseEntity<Map<String, Object>> handleCustomAuth(CustomAuthException ex) {
        log.warn("Auth error: {}", ex.getMessage());
        return ResponseEntity.status(ex.getStatusCode()).body(Map.of(
                "status", ex.getStatusCode(),
                "error", ex.getStatusCode() == 409 ? "Conflict" : "Unauthorized",
                "message", ex.getMessage()));
    }

    @ExceptionHandler(DuplicateKeyException.class)
    public ResponseEntity<Map<String, Object>> handleDuplicateKey(DuplicateKeyException ex) {
        log.warn("Duplicate key error: {}", ex.getMessage());
        return ResponseEntity.status(409).body(Map.of(
                "status", 409,
                "error", "Conflict",
                "message", "An account with this email already exists"));
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, Object>> handleValidation(MethodArgumentNotValidException ex) {
        String message = ex.getBindingResult().getFieldErrors().stream()
                .map(error -> error.getField() + ": " + error.getDefaultMessage())
                .findFirst()
                .orElse("Please check your input");
        log.warn("Validation failed: {}", message);
        return ResponseEntity.status(400).body(Map.of(
                "status", 400,
                "error", "Bad Request",
                "message", message));
    }

    @ExceptionHandler(DataIntegrityViolationException.class)
    public ResponseEntity<Map<String, Object>> handleDataIntegrity(DataIntegrityViolationException ex) {
        String msg = ex.getMostSpecificCause().getMessage();
        log.warn("Data integrity violation: {}", msg);

        if (msg.contains("fk_user_roles_user")) {
            return ResponseEntity.status(404).body(Map.of(
                    "status", 404, "error", "Not Found", "message", "User not found"));
        }
        if (msg.contains("fk_user_roles_role") || msg.contains("fk_role_permissions_role")) {
            return ResponseEntity.status(404).body(Map.of(
                    "status", 404, "error", "Not Found", "message", "Role not found"));
        }
        if (msg.contains("fk_role_permissions_permission")) {
            return ResponseEntity.status(404).body(Map.of(
                    "status", 404, "error", "Not Found", "message", "Permission not found"));
        }

        throw ex;
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<Map<String, Object>> handleGeneral(Exception ex) {
        log.error("Unexpected error", ex);
        return ResponseEntity.status(500).body(Map.of(
                "status", 500,
                "error", "Internal Server Error",
                "message", "Something went wrong, please try again later"));
    }
}
