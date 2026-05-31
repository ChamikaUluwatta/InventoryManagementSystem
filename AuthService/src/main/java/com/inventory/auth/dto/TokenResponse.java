package com.inventory.auth.dto;

import java.util.Set;

public record TokenResponse(
    String accessToken,
    String tokenType,
    long expiresIn,
    String email,
    Set<String> permissions
) {}
