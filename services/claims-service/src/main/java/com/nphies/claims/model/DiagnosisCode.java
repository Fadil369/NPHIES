package com.nphies.claims.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;

/**
 * DiagnosisCode Entity - represents diagnosis codes associated with a claim
 */
@Entity
@Table(name = "diagnosis_codes")
public class DiagnosisCode {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "claim_id", nullable = false)
    private Claim claim;

    @Column(name = "code", nullable = false)
    @NotBlank
    private String code;

    @Column(name = "code_type", nullable = false)
    @NotBlank
    private String codeType; // ICD-10, ICD-11, etc.

    @Column(name = "description")
    private String description;

    @Column(name = "is_primary", nullable = false)
    private Boolean isPrimary = false;

    @Column(name = "sequence_number")
    private Integer sequenceNumber;

    // Constructors
    public DiagnosisCode() {}

    public DiagnosisCode(String code, String codeType, Boolean isPrimary) {
        this.code = code;
        this.codeType = codeType;
        this.isPrimary = isPrimary;
    }

    // Getters and Setters
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public Claim getClaim() { return claim; }
    public void setClaim(Claim claim) { this.claim = claim; }

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