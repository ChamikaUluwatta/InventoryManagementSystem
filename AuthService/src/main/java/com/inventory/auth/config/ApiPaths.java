package com.inventory.auth.config;

public final class ApiPaths {
    public static final String AUTH_BASE = "/api/v1/auth";
    public static final String REGISTER = AUTH_BASE + "/register";
    public static final String LOGIN = AUTH_BASE + "/login";
    public static final String REFRESH = AUTH_BASE + "/refresh";
    public static final String LOGOUT = AUTH_BASE + "/logout";
    public static final String JWKS = AUTH_BASE + "/.well-known/jwks.json";

    private ApiPaths() {}
}
