package com.inventory.auth.dto;

import java.util.Set;

public record GuestResponse(
    String email,
    String password,
    String accessToken,
    String tokenType,
    long expiresIn,
    Set<String> permissions
) {}
