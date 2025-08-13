package com.nphies.claims.dto;

import com.fasterxml.jackson.annotation.JsonFormat;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;

/**
 * DTO for claim submission responses
 */
public class ClaimSubmissionResponse {

    private String claimId;
    private String status;
    private String trackingNumber;

    @JsonFormat(pattern = "yyyy-MM-dd'T'HH:mm:ss")
    private LocalDateTime submissionDate;

    private List<ValidationMessage> messages;

    // Constructors
    public ClaimSubmissionResponse() {}

    public ClaimSubmissionResponse(String claimId, String status, String trackingNumber, LocalDateTime submissionDate) {
        this.claimId = claimId;
        this.status = status;
        this.trackingNumber = trackingNumber;
        this.submissionDate = submissionDate;
    }

    // Getters and Setters
    public String getClaimId() { return claimId; }
    public void setClaimId(String claimId) { this.claimId = claimId; }

    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }

    public String getTrackingNumber() { return trackingNumber; }
    public void setTrackingNumber(String trackingNumber) { this.trackingNumber = trackingNumber; }

    public LocalDateTime getSubmissionDate() { return submissionDate; }
    public void setSubmissionDate(LocalDateTime submissionDate) { this.submissionDate = submissionDate; }

    public List<ValidationMessage> getMessages() { return messages; }
    public void setMessages(List<ValidationMessage> messages) { this.messages = messages; }

    public static class ValidationMessage {
        private String level; // INFO, WARNING, ERROR
        private String code;
        private String message;
        private String field;

        public ValidationMessage() {}

        public ValidationMessage(String level, String code, String message) {
            this.level = level;
            this.code = code;
            this.message = message;
        }

        // Getters and Setters
        public String getLevel() { return level; }
        public void setLevel(String level) { this.level = level; }

        public String getCode() { return code; }
        public void setCode(String code) { this.code = code; }

        public String getMessage() { return message; }
        public void setMessage(String message) { this.message = message; }

        public String getField() { return field; }
        public void setField(String field) { this.field = field; }
    }
}

/**
 * DTO for claim status responses
 */
class ClaimStatusResponse {
    private String claimId;
    private String status;
    private String trackingNumber;

    @JsonFormat(pattern = "yyyy-MM-dd'T'HH:mm:ss")
    private LocalDateTime lastUpdated;

    private BigDecimal totalAmount;
    private BigDecimal approvedAmount;
    private List<ClaimSubmissionResponse.ValidationMessage> processingMessages;

    // Constructors
    public ClaimStatusResponse() {}

    // Getters and Setters
    public String getClaimId() { return claimId; }
    public void setClaimId(String claimId) { this.claimId = claimId; }

    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }

    public String getTrackingNumber() { return trackingNumber; }
    public void setTrackingNumber(String trackingNumber) { this.trackingNumber = trackingNumber; }

    public LocalDateTime getLastUpdated() { return lastUpdated; }
    public void setLastUpdated(LocalDateTime lastUpdated) { this.lastUpdated = lastUpdated; }

    public BigDecimal getTotalAmount() { return totalAmount; }
    public void setTotalAmount(BigDecimal totalAmount) { this.totalAmount = totalAmount; }

    public BigDecimal getApprovedAmount() { return approvedAmount; }
    public void setApprovedAmount(BigDecimal approvedAmount) { this.approvedAmount = approvedAmount; }

    public List<ClaimSubmissionResponse.ValidationMessage> getProcessingMessages() { return processingMessages; }
    public void setProcessingMessages(List<ClaimSubmissionResponse.ValidationMessage> processingMessages) { this.processingMessages = processingMessages; }
}