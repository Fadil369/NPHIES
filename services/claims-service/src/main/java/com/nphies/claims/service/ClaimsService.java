package com.nphies.claims.service;

import com.nphies.claims.dto.ClaimSubmissionRequest;
import com.nphies.claims.dto.ClaimSubmissionResponse;
import com.nphies.claims.model.Claim;
import com.nphies.claims.model.ClaimLine;
import com.nphies.claims.model.DiagnosisCode;
import com.nphies.claims.repository.ClaimRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Core service for claims processing, validation, and audit trails
 */
@Service
@Transactional
public class ClaimsService {

    private final ClaimRepository claimRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;
    private final EligibilityIntegrationService eligibilityService;
    private final ValidationService validationService;

    @Autowired
    public ClaimsService(ClaimRepository claimRepository, 
                        KafkaTemplate<String, Object> kafkaTemplate,
                        EligibilityIntegrationService eligibilityService,
                        ValidationService validationService) {
        this.claimRepository = claimRepository;
        this.kafkaTemplate = kafkaTemplate;
        this.eligibilityService = eligibilityService;
        this.validationService = validationService;
    }

    /**
     * Submit a new claim for processing
     */
    public ClaimSubmissionResponse submitClaim(ClaimSubmissionRequest request) {
        // Check for duplicate submission using idempotency key
        if (request.getIdempotencyKey() != null) {
            Optional<Claim> existingClaim = claimRepository.findByIdempotencyKey(request.getIdempotencyKey());
            if (existingClaim.isPresent()) {
                return createSubmissionResponse(existingClaim.get(), "DUPLICATE");
            }
        }

        // Validate claim data
        List<ClaimSubmissionResponse.ValidationMessage> validationMessages = 
            validationService.validateClaim(request);

        // Check eligibility
        boolean isEligible = eligibilityService.checkEligibility(
            request.getMemberId(), request.getPayerId());

        if (!isEligible) {
            ClaimSubmissionResponse response = new ClaimSubmissionResponse();
            response.setStatus("REJECTED");
            List<ClaimSubmissionResponse.ValidationMessage> messages = new ArrayList<>();
            messages.add(new ClaimSubmissionResponse.ValidationMessage(
                "ERROR", "ELIGIBILITY_FAILED", "Member is not eligible for benefits"));
            response.setMessages(messages);
            return response;
        }

        // Create claim entity
        Claim claim = createClaimFromRequest(request);
        claim.setStatus("SUBMITTED");
        claim.setTrackingNumber(generateTrackingNumber());

        // Save claim
        Claim savedClaim = claimRepository.save(claim);

        // Publish claim submission event
        publishClaimEvent("claim.submitted", savedClaim);

        // Create response
        ClaimSubmissionResponse response = createSubmissionResponse(savedClaim, "SUBMITTED");
        response.setMessages(validationMessages);

        return response;
    }

    /**
     * Get claim status by claim ID
     */
    @Transactional(readOnly = true)
    public Optional<Claim> getClaimStatus(String claimId) {
        return claimRepository.findByClaimId(claimId);
    }

    /**
     * Reprocess a claim
     */
    public ClaimSubmissionResponse reprocessClaim(String claimId) {
        Optional<Claim> claimOpt = claimRepository.findByClaimId(claimId);
        
        if (claimOpt.isEmpty()) {
            throw new IllegalArgumentException("Claim not found: " + claimId);
        }

        Claim claim = claimOpt.get();
        
        // Update status and timestamp
        claim.setStatus("PROCESSING");
        claim.setUpdatedAt(LocalDateTime.now());

        Claim savedClaim = claimRepository.save(claim);

        // Publish reprocessing event
        publishClaimEvent("claim.reprocessing", savedClaim);

        return createSubmissionResponse(savedClaim, "REPROCESSING");
    }

    /**
     * Get claims by provider
     */
    @Transactional(readOnly = true)
    public Page<Claim> getClaimsByProvider(String providerId, Pageable pageable) {
        return claimRepository.findByProviderId(providerId, pageable);
    }

    /**
     * Get claims by member
     */
    @Transactional(readOnly = true)
    public Page<Claim> getClaimsByMember(String memberId, Pageable pageable) {
        return claimRepository.findByMemberId(memberId, pageable);
    }

    private Claim createClaimFromRequest(ClaimSubmissionRequest request) {
        Claim claim = new Claim();
        claim.setClaimId(generateClaimId());
        claim.setProviderId(request.getProviderId());
        claim.setMemberId(request.getMemberId());
        claim.setPayerId(request.getPayerId());
        claim.setServiceDate(request.getServiceDate());
        claim.setTotalAmount(request.getTotalAmount());
        claim.setType(request.getType().toUpperCase());
        claim.setIdempotencyKey(request.getIdempotencyKey());
        claim.setCreatedBy("system"); // TODO: Get from security context

        // Add claim lines
        for (ClaimSubmissionRequest.ClaimLineRequest lineRequest : request.getClaimLines()) {
            ClaimLine claimLine = new ClaimLine();
            claimLine.setServiceCode(lineRequest.getServiceCode());
            claimLine.setServiceDate(lineRequest.getServiceDate());
            claimLine.setUnits(lineRequest.getUnits());
            claimLine.setChargedAmount(lineRequest.getChargedAmount());
            claimLine.setPlaceOfService(lineRequest.getPlaceOfService());
            claimLine.setDescription(lineRequest.getDescription());
            
            claim.addClaimLine(claimLine);
        }

        // Add diagnosis codes
        for (ClaimSubmissionRequest.DiagnosisCodeRequest diagRequest : request.getDiagnosisCodes()) {
            DiagnosisCode diagnosisCode = new DiagnosisCode();
            diagnosisCode.setCode(diagRequest.getCode());
            diagnosisCode.setCodeType(diagRequest.getCodeType());
            diagnosisCode.setDescription(diagRequest.getDescription());
            diagnosisCode.setIsPrimary(diagRequest.getIsPrimary());
            diagnosisCode.setSequenceNumber(diagRequest.getSequenceNumber());
            
            claim.addDiagnosisCode(diagnosisCode);
        }

        return claim;
    }

    private ClaimSubmissionResponse createSubmissionResponse(Claim claim, String status) {
        ClaimSubmissionResponse response = new ClaimSubmissionResponse();
        response.setClaimId(claim.getClaimId());
        response.setStatus(status);
        response.setTrackingNumber(claim.getTrackingNumber());
        response.setSubmissionDate(claim.getCreatedAt());
        return response;
    }

    private String generateClaimId() {
        return "CLM" + System.currentTimeMillis();
    }

    private String generateTrackingNumber() {
        return "TRK" + UUID.randomUUID().toString().replace("-", "").substring(0, 12).toUpperCase();
    }

    private void publishClaimEvent(String eventType, Claim claim) {
        try {
            ClaimEvent event = new ClaimEvent(eventType, claim.getClaimId(), 
                claim.getStatus(), LocalDateTime.now());
            kafkaTemplate.send("claims.events.v1", claim.getClaimId(), event);
        } catch (Exception e) {
            // Log error but don't fail the transaction
            System.err.println("Failed to publish claim event: " + e.getMessage());
        }
    }

    // Inner class for Kafka events
    private static class ClaimEvent {
        private String eventType;
        private String claimId;
        private String status;
        private LocalDateTime timestamp;

        public ClaimEvent(String eventType, String claimId, String status, LocalDateTime timestamp) {
            this.eventType = eventType;
            this.claimId = claimId;
            this.status = status;
            this.timestamp = timestamp;
        }

        // Getters for JSON serialization
        public String getEventType() { return eventType; }
        public String getClaimId() { return claimId; }
        public String getStatus() { return status; }
        public LocalDateTime getTimestamp() { return timestamp; }
    }
}