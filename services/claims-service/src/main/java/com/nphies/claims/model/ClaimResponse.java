package com.nphies.claims.model;

import java.time.LocalDateTime;
import java.util.List;

public class ClaimResponse {
    
    private String claimId;
    private String memberId;
    private String providerId;
    private String status;
    private LocalDateTime submissionDate;
    private LocalDateTime processedDate;
    private Double totalAmount;
    private Double approvedAmount;
    private Double rejectedAmount;
    private String rejectReason;
    private List<ClaimLineItemResponse> lineItemResponses;
    private String explanationOfBenefits;

    // Constructors
    public ClaimResponse() {}

    // Getters and Setters
    public String getClaimId() {
        return claimId;
    }

    public void setClaimId(String claimId) {
        this.claimId = claimId;
    }

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

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public LocalDateTime getSubmissionDate() {
        return submissionDate;
    }

    public void setSubmissionDate(LocalDateTime submissionDate) {
        this.submissionDate = submissionDate;
    }

    public LocalDateTime getProcessedDate() {
        return processedDate;
    }

    public void setProcessedDate(LocalDateTime processedDate) {
        this.processedDate = processedDate;
    }

    public Double getTotalAmount() {
        return totalAmount;
    }

    public void setTotalAmount(Double totalAmount) {
        this.totalAmount = totalAmount;
    }

    public Double getApprovedAmount() {
        return approvedAmount;
    }

    public void setApprovedAmount(Double approvedAmount) {
        this.approvedAmount = approvedAmount;
    }

    public Double getRejectedAmount() {
        return rejectedAmount;
    }

    public void setRejectedAmount(Double rejectedAmount) {
        this.rejectedAmount = rejectedAmount;
    }

    public String getRejectReason() {
        return rejectReason;
    }

    public void setRejectReason(String rejectReason) {
        this.rejectReason = rejectReason;
    }

    public List<ClaimLineItemResponse> getLineItemResponses() {
        return lineItemResponses;
    }

    public void setLineItemResponses(List<ClaimLineItemResponse> lineItemResponses) {
        this.lineItemResponses = lineItemResponses;
    }

    public String getExplanationOfBenefits() {
        return explanationOfBenefits;
    }

    public void setExplanationOfBenefits(String explanationOfBenefits) {
        this.explanationOfBenefits = explanationOfBenefits;
    }
}