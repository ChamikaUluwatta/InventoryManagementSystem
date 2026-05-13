package com.inventory.auth.service.Health;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;

@Component
public class InventoryHealthService {
    private final RestClient restClient;
    public InventoryHealthService(@Value("${IMS_URL}") String inventoryServiceUrl) {
        this.restClient = RestClient.builder().baseUrl(inventoryServiceUrl).build();
    }
    public boolean checkInventoryHealth() {
        try {
            String response = restClient.get()
                .uri("api/v1/health")
                .retrieve()
                .body(String.class);
            return response != null && response.contains("\"status\":\"healthy\"");
        } catch (Exception e) {
            return false;
        }
    }
}
