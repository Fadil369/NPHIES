package com.nphies.claims.model;

public class ClaimLineItem {
    
    private String lineNumber;
    private String procedureCode;
    private String description;
    private Double unitPrice;
    private Integer quantity;
    private Double totalAmount;
    private String serviceDate;
    private String placeOfService;

    // Constructors
    public ClaimLineItem() {}

    // Getters and Setters
    public String getLineNumber() {
        return lineNumber;
    }

    public void setLineNumber(String lineNumber) {
        this.lineNumber = lineNumber;
    }

    public String getProcedureCode() {
        return procedureCode;
    }

    public void setProcedureCode(String procedureCode) {
        this.procedureCode = procedureCode;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Double getUnitPrice() {
        return unitPrice;
    }

    public void setUnitPrice(Double unitPrice) {
        this.unitPrice = unitPrice;
    }

    public Integer getQuantity() {
        return quantity;
    }

    public void setQuantity(Integer quantity) {
        this.quantity = quantity;
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

    public String getPlaceOfService() {
        return placeOfService;
    }

    public void setPlaceOfService(String placeOfService) {
        this.placeOfService = placeOfService;
    }
}

class ClaimLineItemResponse {
    
    private String lineNumber;
    private String procedureCode;
    private String status;
    private Double approvedAmount;
    private Double rejectedAmount;
    private String rejectReason;
    private String adjudicationNote;

    // Constructors
    public ClaimLineItemResponse() {}

    // Getters and Setters
    public String getLineNumber() {
        return lineNumber;
    }

    public void setLineNumber(String lineNumber) {
        this.lineNumber = lineNumber;
    }

    public String getProcedureCode() {
        return procedureCode;
    }

    public void setProcedureCode(String procedureCode) {
        this.procedureCode = procedureCode;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
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

    public String getAdjudicationNote() {
        return adjudicationNote;
    }

    public void setAdjudicationNote(String adjudicationNote) {
        this.adjudicationNote = adjudicationNote;
    }
}