package com.inventory.auth.dto;

import java.time.Instant;
import java.util.UUID;

public record UserResponse(
    UUID id,
    String email,
    Instant createdAt
) {}
