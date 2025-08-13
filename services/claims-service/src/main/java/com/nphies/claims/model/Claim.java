package com.nphies.claims.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

/**
 * Claim Entity - represents a healthcare claim in the NPHIES system
 */
@Entity
@Table(name = "claims")
@EntityListeners(AuditingEntityListener.class)
public class Claim {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "claim_id", unique = true, nullable = false)
    @NotBlank
    private String claimId;

    @Column(name = "provider_id", nullable = false)
    @NotBlank
    private String providerId;

    @Column(name = "member_id", nullable = false)
    @NotBlank
    private String memberId;

    @Column(name = "payer_id", nullable = false)
    @NotBlank
    private String payerId;

    @Column(name = "service_date", nullable = false)
    @NotNull
    private LocalDateTime serviceDate;

    @Column(name = "total_amount", nullable = false, precision = 10, scale = 2)
    @NotNull
    private BigDecimal totalAmount;

    @Column(name = "status", nullable = false)
    private String status;

    @Column(name = "type", nullable = false)
    private String type;

    @Column(name = "idempotency_key")
    private String idempotencyKey;

    @Column(name = "tracking_number")
    private String trackingNumber;

    @OneToMany(mappedBy = "claim", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private List<ClaimLine> claimLines = new ArrayList<>();

    @OneToMany(mappedBy = "claim", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private List<DiagnosisCode> diagnosisCodes = new ArrayList<>();

    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(name = "updated_at", nullable = false)
    private LocalDateTime updatedAt;

    @Column(name = "created_by")
    private String createdBy;

    // Constructors
    public Claim() {}

    public Claim(String claimId, String providerId, String memberId, String payerId, 
                 LocalDateTime serviceDate, BigDecimal totalAmount, String type) {
        this.claimId = claimId;
        this.providerId = providerId;
        this.memberId = memberId;
        this.payerId = payerId;
        this.serviceDate = serviceDate;
        this.totalAmount = totalAmount;
        this.type = type;
        this.status = "SUBMITTED";
    }

    // Getters and Setters
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public String getClaimId() { return claimId; }
    public void setClaimId(String claimId) { this.claimId = claimId; }

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

    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }

    public String getType() { return type; }
    public void setType(String type) { this.type = type; }

    public String getIdempotencyKey() { return idempotencyKey; }
    public void setIdempotencyKey(String idempotencyKey) { this.idempotencyKey = idempotencyKey; }

    public String getTrackingNumber() { return trackingNumber; }
    public void setTrackingNumber(String trackingNumber) { this.trackingNumber = trackingNumber; }

    public List<ClaimLine> getClaimLines() { return claimLines; }
    public void setClaimLines(List<ClaimLine> claimLines) { this.claimLines = claimLines; }

    public List<DiagnosisCode> getDiagnosisCodes() { return diagnosisCodes; }
    public void setDiagnosisCodes(List<DiagnosisCode> diagnosisCodes) { this.diagnosisCodes = diagnosisCodes; }

    public LocalDateTime getCreatedAt() { return createdAt; }
    public void setCreatedAt(LocalDateTime createdAt) { this.createdAt = createdAt; }

    public LocalDateTime getUpdatedAt() { return updatedAt; }
    public void setUpdatedAt(LocalDateTime updatedAt) { this.updatedAt = updatedAt; }

    public String getCreatedBy() { return createdBy; }
    public void setCreatedBy(String createdBy) { this.createdBy = createdBy; }

    // Helper methods
    public void addClaimLine(ClaimLine claimLine) {
        claimLines.add(claimLine);
        claimLine.setClaim(this);
    }

    public void addDiagnosisCode(DiagnosisCode diagnosisCode) {
        diagnosisCodes.add(diagnosisCode);
        diagnosisCode.setClaim(this);
    }
}