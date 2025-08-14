package com.nphies.claims.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.validation.annotation.Validated;
import com.nphies.claims.service.ClaimsService;
import com.nphies.claims.model.ClaimRequest;
import com.nphies.claims.model.ClaimResponse;
import com.nphies.claims.model.ClaimStatus;

import java.util.Map;
import java.util.List;

@RestController
@RequestMapping("/api/v1/claims")
@Validated
public class ClaimsController {

    @Autowired
    private ClaimsService claimsService;

    @PostMapping("/submit")
    public ResponseEntity<ClaimResponse> submitClaim(@RequestBody @Validated ClaimRequest claimRequest) {
        ClaimResponse response = claimsService.submitClaim(claimRequest);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{claimId}")
    public ResponseEntity<ClaimResponse> getClaim(@PathVariable String claimId) {
        ClaimResponse claim = claimsService.getClaimById(claimId);
        return claim != null ? ResponseEntity.ok(claim) : ResponseEntity.notFound().build();
    }

    @GetMapping("/{claimId}/status")
    public ResponseEntity<ClaimStatus> getClaimStatus(@PathVariable String claimId) {
        ClaimStatus status = claimsService.getClaimStatus(claimId);
        return status != null ? ResponseEntity.ok(status) : ResponseEntity.notFound().build();
    }

    @PostMapping("/{claimId}/reprocess")
    public ResponseEntity<ClaimResponse> reprocessClaim(@PathVariable String claimId) {
        ClaimResponse response = claimsService.reprocessClaim(claimId);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/search")
    public ResponseEntity<List<ClaimResponse>> searchClaims(
            @RequestParam(required = false) String memberId,
            @RequestParam(required = false) String providerId,
            @RequestParam(required = false) String status,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        
        List<ClaimResponse> claims = claimsService.searchClaims(memberId, providerId, status, page, size);
        return ResponseEntity.ok(claims);
    }

    @GetMapping("/statistics")
    public ResponseEntity<Map<String, Object>> getClaimsStatistics() {
        Map<String, Object> stats = claimsService.getClaimsStatistics();
        return ResponseEntity.ok(stats);
    }
}