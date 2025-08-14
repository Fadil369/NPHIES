package com.nphies.claims.service;

import org.springframework.stereotype.Service;
import com.nphies.claims.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.util.*;
import java.time.LocalDateTime;

@Service
public class ClaimsService {

    private static final Logger logger = LoggerFactory.getLogger(ClaimsService.class);
    
    // In-memory storage for demonstration (replace with actual repository)
    private final Map<String, ClaimResponse> claimsStorage = new HashMap<>();

    public ClaimResponse submitClaim(ClaimRequest claimRequest) {
        logger.info("Processing claim submission for member: {}", claimRequest.getMemberId());
        
        String claimId = "CLM-" + UUID.randomUUID().toString().substring(0, 8).toUpperCase();
        
        ClaimResponse response = new ClaimResponse();
        response.setClaimId(claimId);
        response.setMemberId(claimRequest.getMemberId());
        response.setProviderId(claimRequest.getProviderId());
        response.setStatus("SUBMITTED");
        response.setSubmissionDate(LocalDateTime.now());
        response.setTotalAmount(claimRequest.getTotalAmount());
        response.setApprovedAmount(0.0); // Will be calculated during adjudication
        response.setRejectedAmount(0.0);
        
        // Basic validation and processing logic
        if (validateClaim(claimRequest)) {
            response.setStatus("UNDER_REVIEW");
            logger.info("Claim {} accepted for processing", claimId);
        } else {
            response.setStatus("REJECTED");
            response.setRejectReason("Failed initial validation");
            logger.warn("Claim {} rejected during validation", claimId);
        }
        
        claimsStorage.put(claimId, response);
        
        // TODO: Send to Kafka for further processing
        publishClaimEvent(response, "CLAIM_SUBMITTED");
        
        return response;
    }

    public ClaimResponse getClaimById(String claimId) {
        return claimsStorage.get(claimId);
    }

    public ClaimStatus getClaimStatus(String claimId) {
        ClaimResponse claim = claimsStorage.get(claimId);
        if (claim == null) {
            return null;
        }
        
        ClaimStatus status = new ClaimStatus();
        status.setClaimId(claimId);
        status.setStatus(claim.getStatus());
        status.setLastUpdated(claim.getSubmissionDate());
        status.setStatusHistory(Arrays.asList(
            "SUBMITTED -> " + claim.getSubmissionDate(),
            claim.getStatus() + " -> " + claim.getSubmissionDate()
        ));
        
        return status;
    }

    public ClaimResponse reprocessClaim(String claimId) {
        ClaimResponse claim = claimsStorage.get(claimId);
        if (claim == null) {
            return null;
        }
        
        logger.info("Reprocessing claim: {}", claimId);
        claim.setStatus("REPROCESSING");
        
        // TODO: Implement reprocessing logic
        publishClaimEvent(claim, "CLAIM_REPROCESSED");
        
        return claim;
    }

    public List<ClaimResponse> searchClaims(String memberId, String providerId, String status, int page, int size) {
        // Simple filtering for demonstration
        return claimsStorage.values().stream()
                .filter(claim -> memberId == null || claim.getMemberId().equals(memberId))
                .filter(claim -> providerId == null || claim.getProviderId().equals(providerId))
                .filter(claim -> status == null || claim.getStatus().equals(status))
                .skip(page * size)
                .limit(size)
                .toList();
    }

    public Map<String, Object> getClaimsStatistics() {
        Map<String, Object> stats = new HashMap<>();
        stats.put("totalClaims", claimsStorage.size());
        stats.put("submittedClaims", claimsStorage.values().stream()
                .mapToLong(claim -> "SUBMITTED".equals(claim.getStatus()) ? 1 : 0).sum());
        stats.put("underReviewClaims", claimsStorage.values().stream()
                .mapToLong(claim -> "UNDER_REVIEW".equals(claim.getStatus()) ? 1 : 0).sum());
        stats.put("rejectedClaims", claimsStorage.values().stream()
                .mapToLong(claim -> "REJECTED".equals(claim.getStatus()) ? 1 : 0).sum());
        stats.put("totalAmount", claimsStorage.values().stream()
                .mapToDouble(ClaimResponse::getTotalAmount).sum());
        
        return stats;
    }

    private boolean validateClaim(ClaimRequest claimRequest) {
        // Basic validation logic
        if (claimRequest.getMemberId() == null || claimRequest.getMemberId().trim().isEmpty()) {
            logger.warn("Claim validation failed: Missing member ID");
            return false;
        }
        
        if (claimRequest.getProviderId() == null || claimRequest.getProviderId().trim().isEmpty()) {
            logger.warn("Claim validation failed: Missing provider ID");
            return false;
        }
        
        if (claimRequest.getTotalAmount() <= 0) {
            logger.warn("Claim validation failed: Invalid total amount");
            return false;
        }
        
        return true;
    }

    private void publishClaimEvent(ClaimResponse claim, String eventType) {
        // TODO: Implement Kafka event publishing
        logger.info("Publishing event {} for claim {}", eventType, claim.getClaimId());
    }
}