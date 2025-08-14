package com.nphies.claims.model;

import java.time.LocalDateTime;
import java.util.List;

public class ClaimStatus {
    
    private String claimId;
    private String status;
    private LocalDateTime lastUpdated;
    private List<String> statusHistory;
    private String currentStep;
    private String nextStepDescription;
    private Integer estimatedProcessingDays;

    // Constructors
    public ClaimStatus() {}

    // Getters and Setters
    public String getClaimId() {
        return claimId;
    }

    public void setClaimId(String claimId) {
        this.claimId = claimId;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public LocalDateTime getLastUpdated() {
        return lastUpdated;
    }

    public void setLastUpdated(LocalDateTime lastUpdated) {
        this.lastUpdated = lastUpdated;
    }

    public List<String> getStatusHistory() {
        return statusHistory;
    }

    public void setStatusHistory(List<String> statusHistory) {
        this.statusHistory = statusHistory;
    }

    public String getCurrentStep() {
        return currentStep;
    }

    public void setCurrentStep(String currentStep) {
        this.currentStep = currentStep;
    }

    public String getNextStepDescription() {
        return nextStepDescription;
    }

    public void setNextStepDescription(String nextStepDescription) {
        this.nextStepDescription = nextStepDescription;
    }

    public Integer getEstimatedProcessingDays() {
        return estimatedProcessingDays;
    }

    public void setEstimatedProcessingDays(Integer estimatedProcessingDays) {
        this.estimatedProcessingDays = estimatedProcessingDays;
    }
}