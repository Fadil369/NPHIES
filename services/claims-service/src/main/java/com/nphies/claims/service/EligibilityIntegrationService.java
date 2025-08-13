package com.nphies.claims.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

/**
 * Service for integrating with the Eligibility Service
 */
@Service
public class EligibilityIntegrationService {

    private final RestTemplate restTemplate;
    
    @Value("${nphies.services.eligibility.url:http://localhost:8090}")
    private String eligibilityServiceUrl;

    public EligibilityIntegrationService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /**
     * Check member eligibility for benefits
     */
    public boolean checkEligibility(String memberId, String payerId) {
        try {
            // Call eligibility service
            String url = eligibilityServiceUrl + "/api/v1/eligibility/check/" + memberId + "/" + payerId;
            
            EligibilityResponse response = restTemplate.getForObject(url, EligibilityResponse.class);
            
            return response != null && response.isEligible();
        } catch (Exception e) {
            // Log error and return false for safety
            System.err.println("Failed to check eligibility: " + e.getMessage());
            return false;
        }
    }

    /**
     * DTO for eligibility response
     */
    public static class EligibilityResponse {
        private boolean eligible;
        private String status;
        private String message;

        // Constructors
        public EligibilityResponse() {}

        // Getters and Setters
        public boolean isEligible() { return eligible; }
        public void setEligible(boolean eligible) { this.eligible = eligible; }

        public String getStatus() { return status; }
        public void setStatus(String status) { this.status = status; }

        public String getMessage() { return message; }
        public void setMessage(String message) { this.message = message; }
    }
}