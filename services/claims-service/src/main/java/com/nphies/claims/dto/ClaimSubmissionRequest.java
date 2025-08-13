package com.nphies.claims.dto;

import com.fasterxml.jackson.annotation.JsonFormat;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;

/**
 * DTO for claim submission requests
 */
public class ClaimSubmissionRequest {

    @NotBlank(message = "Provider ID is required")
    private String providerId;

    @NotBlank(message = "Member ID is required")
    private String memberId;

    @NotBlank(message = "Payer ID is required")
    private String payerId;

    @NotNull(message = "Service date is required")
    @JsonFormat(pattern = "yyyy-MM-dd'T'HH:mm:ss")
    private LocalDateTime serviceDate;

    @NotNull(message = "Total amount is required")
    private BigDecimal totalAmount;

    @NotBlank(message = "Claim type is required")
    private String type;

    private String idempotencyKey;

    @NotEmpty(message = "At least one claim line is required")
    @Valid
    private List<ClaimLineRequest> claimLines;

    @NotEmpty(message = "At least one diagnosis code is required")
    @Valid
    private List<DiagnosisCodeRequest> diagnosisCodes;

    // Constructors
    public ClaimSubmissionRequest() {}

    // Getters and Setters
    public String getProviderId() { return providerId; }
    public void setProviderId(String providerId) { this.providerId = providerId; }

    public String getMemberId() { return memberId; }
    public void setMemberId(String memberId) { this.memberId = memberId; }

    public String getPayerId() { return payerId; }
    public void setPayerId(String payerId) { this.payerId = payerId; }

    public LocalDateTime getServiceDate() { return serviceDate; }
    public void setServiceDate(LocalDateTime serviceDate) { this.serviceDate = serviceDate; }

    public BigDecimal getTotalAmount() { return totalAmount; }
    public void setTotalAmount(BigDecimal totalAmount) { this.totalAmount = totalAmount; }

    public String getType() { return type; }
    public void setType(String type) { this.type = type; }

    public String getIdempotencyKey() { return idempotencyKey; }
    public void setIdempotencyKey(String idempotencyKey) { this.idempotencyKey = idempotencyKey; }

    public List<ClaimLineRequest> getClaimLines() { return claimLines; }
    public void setClaimLines(List<ClaimLineRequest> claimLines) { this.claimLines = claimLines; }

    public List<DiagnosisCodeRequest> getDiagnosisCodes() { return diagnosisCodes; }
    public void setDiagnosisCodes(List<DiagnosisCodeRequest> diagnosisCodes) { this.diagnosisCodes = diagnosisCodes; }

    public static class ClaimLineRequest {
        @NotBlank(message = "Service code is required")
        private String serviceCode;

        @NotNull(message = "Service date is required")
        @JsonFormat(pattern = "yyyy-MM-dd'T'HH:mm:ss")
        private LocalDateTime serviceDate;

        @NotNull(message = "Units is required")
        private Integer units;

        @NotNull(message = "Charged amount is required")
        private BigDecimal chargedAmount;

        private String placeOfService;
        private List<String> modifiers;
        private String description;

        // Getters and Setters
        public String getServiceCode() { return serviceCode; }
        public void setServiceCode(String serviceCode) { this.serviceCode = serviceCode; }

        public LocalDateTime getServiceDate() { return serviceDate; }
        public void setServiceDate(LocalDateTime serviceDate) { this.serviceDate = serviceDate; }

        public Integer getUnits() { return units; }
        public void setUnits(Integer units) { this.units = units; }

        public BigDecimal getChargedAmount() { return chargedAmount; }
        public void setChargedAmount(BigDecimal chargedAmount) { this.chargedAmount = chargedAmount; }

        public String getPlaceOfService() { return placeOfService; }
        public void setPlaceOfService(String placeOfService) { this.placeOfService = placeOfService; }

        public List<String> getModifiers() { return modifiers; }
        public void setModifiers(List<String> modifiers) { this.modifiers = modifiers; }

        public String getDescription() { return description; }
        public void setDescription(String description) { this.description = description; }
    }

    public static class DiagnosisCodeRequest {
        @NotBlank(message = "Diagnosis code is required")
        private String code;

        @NotBlank(message = "Code type is required")
        private String codeType;

        private String description;
        private Boolean isPrimary = false;
        private Integer sequenceNumber;

        // Getters and Setters
        public String getCode() { return code; }
        public void setCode(String code) { this.code = code; }

        public String getCodeType() { return codeType; }
        public void setCodeType(String codeType) { this.codeType = codeType; }

        public String getDescription() { return description; }
        public void setDescription(String description) { this.description = description; }

        public Boolean getIsPrimary() { return isPrimary; }
        public void setIsPrimary(Boolean isPrimary) { this.isPrimary = isPrimary; }

        public Integer getSequenceNumber() { return sequenceNumber; }
        public void setSequenceNumber(Integer sequenceNumber) { this.sequenceNumber = sequenceNumber; }
    }
}