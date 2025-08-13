package com.nphies.claims.controller;

import com.nphies.claims.dto.ClaimSubmissionRequest;
import com.nphies.claims.dto.ClaimSubmissionResponse;
import com.nphies.claims.model.Claim;
import com.nphies.claims.service.ClaimsService;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Optional;

/**
 * REST Controller for Claims Service
 */
@RestController
@RequestMapping("/api/v1/claims")
@CrossOrigin(origins = "*")
public class ClaimsController {

    private final ClaimsService claimsService;

    @Autowired
    public ClaimsController(ClaimsService claimsService) {
        this.claimsService = claimsService;
    }

    /**
     * Submit a new claim
     */
    @PostMapping("/submit")
    public ResponseEntity<ClaimSubmissionResponse> submitClaim(
            @Valid @RequestBody ClaimSubmissionRequest request) {
        try {
            ClaimSubmissionResponse response = claimsService.submitClaim(request);
            
            if ("REJECTED".equals(response.getStatus())) {
                return ResponseEntity.badRequest().body(response);
            }
            
            return ResponseEntity.status(HttpStatus.ACCEPTED).body(response);
        } catch (Exception e) {
            ClaimSubmissionResponse errorResponse = new ClaimSubmissionResponse();
            errorResponse.setStatus("ERROR");
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(errorResponse);
        }
    }

    /**
     * Get claim status
     */
    @GetMapping("/{claimId}/status")
    public ResponseEntity<?> getClaimStatus(@PathVariable String claimId) {
        Optional<Claim> claim = claimsService.getClaimStatus(claimId);
        
        if (claim.isEmpty()) {
            return ResponseEntity.notFound().build();
        }
        
        return ResponseEntity.ok(createClaimStatusResponse(claim.get()));
    }

    /**
     * Reprocess a claim
     */
    @PostMapping("/{claimId}/reprocess")
    public ResponseEntity<ClaimSubmissionResponse> reprocessClaim(@PathVariable String claimId) {
        try {
            ClaimSubmissionResponse response = claimsService.reprocessClaim(claimId);
            return ResponseEntity.accepted().body(response);
        } catch (IllegalArgumentException e) {
            return ResponseEntity.notFound().build();
        } catch (Exception e) {
            ClaimSubmissionResponse errorResponse = new ClaimSubmissionResponse();
            errorResponse.setStatus("ERROR");
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(errorResponse);
        }
    }

    /**
     * Get claims by provider
     */
    @GetMapping("/provider/{providerId}")
    public ResponseEntity<Page<Claim>> getClaimsByProvider(
            @PathVariable String providerId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        
        Pageable pageable = PageRequest.of(page, size);
        Page<Claim> claims = claimsService.getClaimsByProvider(providerId, pageable);
        return ResponseEntity.ok(claims);
    }

    /**
     * Get claims by member
     */
    @GetMapping("/member/{memberId}")
    public ResponseEntity<Page<Claim>> getClaimsByMember(
            @PathVariable String memberId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        
        Pageable pageable = PageRequest.of(page, size);
        Page<Claim> claims = claimsService.getClaimsByMember(memberId, pageable);
        return ResponseEntity.ok(claims);
    }

    /**
     * Health check endpoint
     */
    @GetMapping("/health")
    public ResponseEntity<HealthResponse> health() {
        HealthResponse health = new HealthResponse();
        health.setStatus("UP");
        health.setService("claims-service");
        health.setTimestamp(java.time.LocalDateTime.now());
        return ResponseEntity.ok(health);
    }

    private ClaimStatusResponse createClaimStatusResponse(Claim claim) {
        ClaimStatusResponse response = new ClaimStatusResponse();
        response.setClaimId(claim.getClaimId());
        response.setStatus(claim.getStatus());
        response.setTrackingNumber(claim.getTrackingNumber());
        response.setLastUpdated(claim.getUpdatedAt());
        response.setTotalAmount(claim.getTotalAmount());
        return response;
    }

    // DTOs for responses
    public static class ClaimStatusResponse {
        private String claimId;
        private String status;
        private String trackingNumber;
        private java.time.LocalDateTime lastUpdated;
        private java.math.BigDecimal totalAmount;

        // Getters and Setters
        public String getClaimId() { return claimId; }
        public void setClaimId(String claimId) { this.claimId = claimId; }

        public String getStatus() { return status; }
        public void setStatus(String status) { this.status = status; }

        public String getTrackingNumber() { return trackingNumber; }
        public void setTrackingNumber(String trackingNumber) { this.trackingNumber = trackingNumber; }

        public java.time.LocalDateTime getLastUpdated() { return lastUpdated; }
        public void setLastUpdated(java.time.LocalDateTime lastUpdated) { this.lastUpdated = lastUpdated; }

        public java.math.BigDecimal getTotalAmount() { return totalAmount; }
        public void setTotalAmount(java.math.BigDecimal totalAmount) { this.totalAmount = totalAmount; }
    }

    public static class HealthResponse {
        private String status;
        private String service;
        private java.time.LocalDateTime timestamp;

        // Getters and Setters
        public String getStatus() { return status; }
        public void setStatus(String status) { this.status = status; }

        public String getService() { return service; }
        public void setService(String service) { this.service = service; }

        public java.time.LocalDateTime getTimestamp() { return timestamp; }
        public void setTimestamp(java.time.LocalDateTime timestamp) { this.timestamp = timestamp; }
    }
}