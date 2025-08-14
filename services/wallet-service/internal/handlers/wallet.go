package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Fadil369/NPHIES/services/wallet-service/internal/blockchain"
	"github.com/Fadil369/NPHIES/services/wallet-service/internal/config"
	"github.com/Fadil369/NPHIES/services/wallet-service/internal/models"
)

type Handler struct {
	logger          *logrus.Logger
	config          *config.Config
	blockchainClient blockchain.BlockchainClient
}

func NewHandler(logger *logrus.Logger, cfg *config.Config) (*Handler, error) {
	// Initialize blockchain client
	blockchainClient := blockchain.NewHyperledgerClient(
		cfg.Blockchain.NodeURL,
		cfg.Blockchain.PrivateKey,
		cfg.Blockchain.ChainID,
		logger,
	)

	return &Handler{
		logger:          logger,
		config:          cfg,
		blockchainClient: blockchainClient,
	}, nil
}

// Health check endpoint
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "wallet-service",
		"version": "1.0.0",
	})
}

// Readiness check endpoint
func (h *Handler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Metrics endpoint
func (h *Handler) Metrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"requests_total":     0,
		"uptime_seconds":     0,
		"blockchain_txs":     0,
		"wallet_balances":    0,
		"active_consents":    0,
	})
}

// GetWallet retrieves wallet information for a member
func (h *Handler) GetWallet(c *gin.Context) {
	memberID := c.Param("memberId")

	h.logger.WithField("member_id", memberID).Info("Retrieving wallet information")

	// Mock wallet data
	wallet := gin.H{
		"member_id": memberID,
		"balance": gin.H{
			"total_balance":     1500.00,
			"available_balance": 1200.00,
			"reserved_balance":  300.00,
			"currency":          "SAR",
			"last_updated":      time.Now(),
		},
		"benefits": []gin.H{
			{
				"type":              "MEDICAL",
				"limit_amount":      5000.00,
				"used_amount":       1200.00,
				"remaining_amount":  3800.00,
				"period_start":      "2024-01-01",
				"period_end":        "2024-12-31",
			},
			{
				"type":              "DENTAL",
				"limit_amount":      1000.00,
				"used_amount":       200.00,
				"remaining_amount":  800.00,
				"period_start":      "2024-01-01",
				"period_end":        "2024-12-31",
			},
		},
		"recent_transactions": []gin.H{
			{
				"id":          "tx-001",
				"type":        "BENEFIT_USAGE",
				"amount":      -250.00,
				"description": "Doctor visit copay",
				"timestamp":   time.Now().Add(-24 * time.Hour),
				"status":      "COMPLETED",
			},
		},
	}

	c.JSON(http.StatusOK, wallet)
}

// CreateTransaction creates a new wallet transaction
func (h *Handler) CreateTransaction(c *gin.Context) {
	memberID := c.Param("memberId")
	
	var request models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"member_id": memberID,
		"type":      request.Type,
		"amount":    request.Amount,
	}).Info("Creating wallet transaction")

	// Create transaction
	transaction := models.WalletTransaction{
		ID:            uuid.New().String(),
		MemberID:      memberID,
		TransactionID: "TXN-" + uuid.New().String()[:8],
		Type:          request.Type,
		Amount:        request.Amount,
		Currency:      "SAR",
		Description:   request.Description,
		Status:        "COMPLETED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ClaimID:       request.ClaimID,
		ProviderID:    request.ProviderID,
	}

	// Anchor transaction to blockchain
	if blockchainTx, err := h.anchorTransactionToBlockchain(transaction); err == nil {
		transaction.BlockchainTx = &blockchainTx
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions retrieves transaction history for a member
func (h *Handler) GetTransactions(c *gin.Context) {
	memberID := c.Param("memberId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	h.logger.WithFields(logrus.Fields{
		"member_id": memberID,
		"page":      page,
		"size":      size,
	}).Info("Retrieving transaction history")

	// Mock transaction data
	transactions := []gin.H{
		{
			"id":             "tx-001",
			"transaction_id": "TXN-ABC123",
			"type":           "BENEFIT_USAGE",
			"amount":         -250.00,
			"currency":       "SAR",
			"description":    "Doctor visit copay",
			"status":         "COMPLETED",
			"created_at":     time.Now().Add(-24 * time.Hour),
			"blockchain_tx":  "0x1234567890abcdef",
		},
		{
			"id":             "tx-002",
			"transaction_id": "TXN-DEF456",
			"type":           "CREDIT",
			"amount":         1000.00,
			"currency":       "SAR",
			"description":    "Insurance reimbursement",
			"status":         "COMPLETED",
			"created_at":     time.Now().Add(-72 * time.Hour),
			"blockchain_tx":  "0xabcdef1234567890",
		},
	}

	response := gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"page":        page,
			"size":        size,
			"total":       len(transactions),
			"total_pages": 1,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetBalance retrieves current balance for a member
func (h *Handler) GetBalance(c *gin.Context) {
	memberID := c.Param("memberId")

	h.logger.WithField("member_id", memberID).Info("Retrieving wallet balance")

	balance := models.WalletBalance{
		MemberID:         memberID,
		TotalBalance:     1500.00,
		AvailableBalance: 1200.00,
		ReservedBalance:  300.00,
		Currency:         "SAR",
		LastUpdated:      time.Now(),
	}

	c.JSON(http.StatusOK, balance)
}

// AnchorToBlockchain anchors data to blockchain
func (h *Handler) AnchorToBlockchain(c *gin.Context) {
	var request models.BlockchainAnchorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"ref_type": request.RefType,
		"ref_id":   request.RefID,
	}).Info("Anchoring data to blockchain")

	response, err := h.blockchainClient.SubmitTransaction(request.RefType, request.RefID, request.Data)
	if err != nil {
		h.logger.WithError(err).Error("Failed to submit transaction to blockchain")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to anchor to blockchain"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// VerifyHash verifies a hash on the blockchain
func (h *Handler) VerifyHash(c *gin.Context) {
	hash := c.Param("hash")

	h.logger.WithField("hash", hash).Info("Verifying hash on blockchain")

	result, err := h.blockchainClient.VerifyHash(hash)
	if err != nil {
		h.logger.WithError(err).Error("Failed to verify hash")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify hash"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetBlockchainTransaction retrieves blockchain transaction details
func (h *Handler) GetBlockchainTransaction(c *gin.Context) {
	txID := c.Param("txId")

	h.logger.WithField("transaction_id", txID).Info("Retrieving blockchain transaction")

	details, err := h.blockchainClient.GetTransaction(txID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// anchorTransactionToBlockchain anchors a transaction to blockchain
func (h *Handler) anchorTransactionToBlockchain(transaction models.WalletTransaction) (string, error) {
	data, err := blockchain.CreateDataHash(transaction)
	if err != nil {
		return "", err
	}

	response, err := h.blockchainClient.SubmitTransaction("TRANSACTION", transaction.ID, data)
	if err != nil {
		return "", err
	}

	return response.TransactionID, nil
}