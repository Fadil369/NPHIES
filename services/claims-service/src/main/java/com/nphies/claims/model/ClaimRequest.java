package com.nphies.claims.model;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import java.util.List;

public class ClaimRequest {
    
    @NotBlank(message = "Member ID is required")
    private String memberId;
    
    @NotBlank(message = "Provider ID is required")
    private String providerId;
    
    @NotNull(message = "Total amount is required")
    @Positive(message = "Total amount must be positive")
    private Double totalAmount;
    
    private String serviceDate;
    private String diagnosisCode;
    private List<String> procedureCodes;
    private String placeOfService;
    private String referralId;
    private String priorAuthorizationId;
    private List<ClaimLineItem> lineItems;

    // Constructors
    public ClaimRequest() {}

    // Getters and Setters
    public String getMemberId() {
        return memberId;
    }

    public void setMemberId(String memberId) {
        this.memberId = memberId;
    }

    public String getProviderId() {
        return providerId;
    }

    public void setProviderId(String providerId) {
        this.providerId = providerId;
    }

    public Double getTotalAmount() {
        return totalAmount;
    }

    public void setTotalAmount(Double totalAmount) {
        this.totalAmount = totalAmount;
    }

    public String getServiceDate() {
        return serviceDate;
    }

    public void setServiceDate(String serviceDate) {
        this.serviceDate = serviceDate;
    }

    public String getDiagnosisCode() {
        return diagnosisCode;
    }

    public void setDiagnosisCode(String diagnosisCode) {
        this.diagnosisCode = diagnosisCode;
    }

    public List<String> getProcedureCodes() {
        return procedureCodes;
    }

    public void setProcedureCodes(List<String> procedureCodes) {
        this.procedureCodes = procedureCodes;
    }

    public String getPlaceOfService() {
        return placeOfService;
    }

    public void setPlaceOfService(String placeOfService) {
        this.placeOfService = placeOfService;
    }

    public String getReferralId() {
        return referralId;
    }

    public void setReferralId(String referralId) {
        this.referralId = referralId;
    }

    public String getPriorAuthorizationId() {
        return priorAuthorizationId;
    }

    public void setPriorAuthorizationId(String priorAuthorizationId) {
        this.priorAuthorizationId = priorAuthorizationId;
    }

    public List<ClaimLineItem> getLineItems() {
        return lineItems;
    }

    public void setLineItems(List<ClaimLineItem> lineItems) {
        this.lineItems = lineItems;
    }
}