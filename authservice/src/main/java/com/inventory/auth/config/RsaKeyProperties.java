package com.inventory.auth.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app.jwt.rsa")
public record RsaKeyProperties(String privateKeyPath) {}
