package com.nphies.claims.repository;

import com.nphies.claims.model.Claim;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

/**
 * Repository interface for Claim entities
 */
@Repository
public interface ClaimRepository extends JpaRepository<Claim, Long> {

    /**
     * Find claim by claim ID
     */
    Optional<Claim> findByClaimId(String claimId);

    /**
     * Find claim by idempotency key to prevent duplicate submissions
     */
    Optional<Claim> findByIdempotencyKey(String idempotencyKey);

    /**
     * Find claims by provider ID
     */
    Page<Claim> findByProviderId(String providerId, Pageable pageable);

    /**
     * Find claims by member ID
     */
    Page<Claim> findByMemberId(String memberId, Pageable pageable);

    /**
     * Find claims by payer ID
     */
    Page<Claim> findByPayerId(String payerId, Pageable pageable);

    /**
     * Find claims by status
     */
    Page<Claim> findByStatus(String status, Pageable pageable);

    /**
     * Find claims by service date range
     */
    @Query("SELECT c FROM Claim c WHERE c.serviceDate BETWEEN :startDate AND :endDate")
    Page<Claim> findByServiceDateBetween(
        @Param("startDate") LocalDateTime startDate,
        @Param("endDate") LocalDateTime endDate,
        Pageable pageable
    );

    /**
     * Find claims by provider and status
     */
    Page<Claim> findByProviderIdAndStatus(String providerId, String status, Pageable pageable);

    /**
     * Find claims requiring review (pending status for more than specified hours)
     */
    @Query("SELECT c FROM Claim c WHERE c.status = 'PENDING_REVIEW' AND c.updatedAt < :thresholdDate")
    List<Claim> findClaimsRequiringReview(@Param("thresholdDate") LocalDateTime thresholdDate);

    /**
     * Count claims by status for dashboard
     */
    long countByStatus(String status);

    /**
     * Count claims by provider and status
     */
    long countByProviderIdAndStatus(String providerId, String status);
}