package com.inventory.auth.service.Health;

import javax.sql.DataSource;

import org.springframework.stereotype.Component;

@Component
public class DatabaseHealthService {
    private final DataSource dataSource;

    public DatabaseHealthService(DataSource dataSource) {
        this.dataSource = dataSource;
    }

    public boolean checkDatabaseHealth() {
        try (var connection = dataSource.getConnection()) {
            return connection.isValid(2);
        } catch (Exception e) {
            return false;
        }
    }
}
