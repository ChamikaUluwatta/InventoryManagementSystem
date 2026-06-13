package com.inventory.auth.service.Health;

import org.springframework.stereotype.Service;

@Service
public class HealthService {
    private final DatabaseHealthService databaseHealthService;
    private final InventoryHealthService inventoryHealthService;
    
    public HealthService(DatabaseHealthService databaseHealthService, InventoryHealthService inventoryHealthService) {
        this.databaseHealthService = databaseHealthService;
        this.inventoryHealthService = inventoryHealthService;
    }

    public String checkHealth() {
        boolean isDatabaseHealthy = databaseHealthService.checkDatabaseHealth();
        boolean isInventoryHealthy = inventoryHealthService.checkInventoryHealth();

        if (isDatabaseHealthy && isInventoryHealthy) {
            return "Service is healthy!";
        } else {
            if (!isDatabaseHealthy) {
                System.err.println("Database connection is unhealthy!");
            }
            if (!isInventoryHealthy) {
                System.err.println("Inventory service is unhealthy!");
            }
            return "Service is experiencing issues!";
        }
    }

}