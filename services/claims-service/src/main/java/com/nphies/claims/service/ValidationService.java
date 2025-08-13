package com.nphies.claims.service;

import com.nphies.claims.dto.ClaimSubmissionRequest;
import com.nphies.claims.dto.ClaimSubmissionResponse;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

/**
 * Service for validating claims data and business rules
 */
@Service
public class ValidationService {

    /**
     * Validate claim submission request
     */
    public List<ClaimSubmissionResponse.ValidationMessage> validateClaim(ClaimSubmissionRequest request) {
        List<ClaimSubmissionResponse.ValidationMessage> messages = new ArrayList<>();

        // Basic validation
        if (request.getClaimLines().isEmpty()) {
            messages.add(new ClaimSubmissionResponse.ValidationMessage(
                "ERROR", "MISSING_CLAIM_LINES", "At least one claim line is required"));
        }

        if (request.getDiagnosisCodes().isEmpty()) {
            messages.add(new ClaimSubmissionResponse.ValidationMessage(
                "ERROR", "MISSING_DIAGNOSIS", "At least one diagnosis code is required"));
        }

        // Validate total amount matches sum of claim lines
        var calculatedTotal = request.getClaimLines().stream()
            .map(line -> line.getChargedAmount().multiply(java.math.BigDecimal.valueOf(line.getUnits())))
            .reduce(java.math.BigDecimal.ZERO, java.math.BigDecimal::add);

        if (calculatedTotal.compareTo(request.getTotalAmount()) != 0) {
            messages.add(new ClaimSubmissionResponse.ValidationMessage(
                "WARNING", "AMOUNT_MISMATCH", 
                "Total amount does not match sum of claim line amounts"));
        }

        // Validate service codes (placeholder - would integrate with terminology service)
        for (var claimLine : request.getClaimLines()) {
            if (!isValidServiceCode(claimLine.getServiceCode())) {
                messages.add(new ClaimSubmissionResponse.ValidationMessage(
                    "ERROR", "INVALID_SERVICE_CODE", 
                    "Invalid service code: " + claimLine.getServiceCode()));
            }
        }

        // Validate diagnosis codes (placeholder - would integrate with terminology service)
        for (var diagCode : request.getDiagnosisCodes()) {
            if (!isValidDiagnosisCode(diagCode.getCode(), diagCode.getCodeType())) {
                messages.add(new ClaimSubmissionResponse.ValidationMessage(
                    "ERROR", "INVALID_DIAGNOSIS_CODE", 
                    "Invalid diagnosis code: " + diagCode.getCode()));
            }
        }

        return messages;
    }

    private boolean isValidServiceCode(String serviceCode) {
        // Placeholder implementation - would call terminology service
        // For now, just check basic format
        return serviceCode != null && serviceCode.matches("\\d{4,5}[A-Z]?");
    }

    private boolean isValidDiagnosisCode(String code, String codeType) {
        // Placeholder implementation - would call terminology service
        // For now, just check basic format
        if ("ICD-10".equals(codeType)) {
            return code != null && code.matches("[A-Z]\\d{2}(\\.[0-9A-Z]{1,4})?");
        }
        return code != null && !code.trim().isEmpty();
    }
}