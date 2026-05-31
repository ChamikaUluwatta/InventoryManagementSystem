package com.inventory.auth.exception;

public class CustomAuthException extends RuntimeException {

    private final int statusCode;

    public CustomAuthException(int statusCode, String message) {
        super(message);
        this.statusCode = statusCode;
    }

    public CustomAuthException(int statusCode, String message, Throwable cause) {
        super(message, cause);
        this.statusCode = statusCode;
    }

    public int getStatusCode() {
        return statusCode;
    }
}
