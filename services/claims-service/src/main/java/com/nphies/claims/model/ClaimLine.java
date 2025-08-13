package com.nphies.claims.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;

import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * ClaimLine Entity - represents individual service lines within a claim
 */
@Entity
@Table(name = "claim_lines")
public class ClaimLine {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "claim_id", nullable = false)
    private Claim claim;

    @Column(name = "service_code", nullable = false)
    @NotBlank
    private String serviceCode;

    @Column(name = "service_date", nullable = false)
    @NotNull
    private LocalDateTime serviceDate;

    @Column(name = "units", nullable = false)
    @Positive
    private Integer units;

    @Column(name = "charged_amount", nullable = false, precision = 10, scale = 2)
    @NotNull
    private BigDecimal chargedAmount;

    @Column(name = "approved_amount", precision = 10, scale = 2)
    private BigDecimal approvedAmount;

    @Column(name = "place_of_service")
    private String placeOfService;

    @Column(name = "modifiers")
    private String modifiers; // JSON string for list of modifiers

    @Column(name = "description")
    private String description;

    // Constructors
    public ClaimLine() {}

    public ClaimLine(String serviceCode, LocalDateTime serviceDate, 
                     Integer units, BigDecimal chargedAmount) {
        this.serviceCode = serviceCode;
        this.serviceDate = serviceDate;
        this.units = units;
        this.chargedAmount = chargedAmount;
    }

    // Getters and Setters
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public Claim getClaim() { return claim; }
    public void setClaim(Claim claim) { this.claim = claim; }

    public String getServiceCode() { return serviceCode; }
    public void setServiceCode(String serviceCode) { this.serviceCode = serviceCode; }

    public LocalDateTime getServiceDate() { return serviceDate; }
    public void setServiceDate(LocalDateTime serviceDate) { this.serviceDate = serviceDate; }

    public Integer getUnits() { return units; }
    public void setUnits(Integer units) { this.units = units; }

    public BigDecimal getChargedAmount() { return chargedAmount; }
    public void setChargedAmount(BigDecimal chargedAmount) { this.chargedAmount = chargedAmount; }

    public BigDecimal getApprovedAmount() { return approvedAmount; }
    public void setApprovedAmount(BigDecimal approvedAmount) { this.approvedAmount = approvedAmount; }

    public String getPlaceOfService() { return placeOfService; }
    public void setPlaceOfService(String placeOfService) { this.placeOfService = placeOfService; }

    public String getModifiers() { return modifiers; }
    public void setModifiers(String modifiers) { this.modifiers = modifiers; }

    public String getDescription() { return description; }
    public void setDescription(String description) { this.description = description; }
}