package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// BlockchainClient interface for blockchain operations
type BlockchainClient interface {
	SubmitTransaction(refType, refID, data string) (*BlockchainResponse, error)
	VerifyHash(hash string) (*VerificationResult, error)
	GetTransaction(txID string) (*TransactionDetails, error)
}

// HyperledgerClient implements blockchain operations for Hyperledger Fabric
type HyperledgerClient struct {
	nodeURL    string
	privateKey string
	chainID    int
	logger     *logrus.Logger
}

// BlockchainResponse represents the response from blockchain submission
type BlockchainResponse struct {
	TransactionID string    `json:"transaction_id"`
	Hash          string    `json:"hash"`
	BlockNumber   *int64    `json:"block_number,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Status        string    `json:"status"`
}

// VerificationResult represents the result of hash verification
type VerificationResult struct {
	Valid         bool      `json:"valid"`
	Hash          string    `json:"hash"`
	TransactionID string    `json:"transaction_id"`
	BlockNumber   *int64    `json:"block_number,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// TransactionDetails represents detailed information about a blockchain transaction
type TransactionDetails struct {
	TransactionID string    `json:"transaction_id"`
	Hash          string    `json:"hash"`
	RefType       string    `json:"ref_type"`
	RefID         string    `json:"ref_id"`
	BlockNumber   int64     `json:"block_number"`
	Timestamp     time.Time `json:"timestamp"`
	Signer        string    `json:"signer"`
	Status        string    `json:"status"`
}

// NewHyperledgerClient creates a new Hyperledger Fabric client
func NewHyperledgerClient(nodeURL, privateKey string, chainID int, logger *logrus.Logger) *HyperledgerClient {
	return &HyperledgerClient{
		nodeURL:    nodeURL,
		privateKey: privateKey,
		chainID:    chainID,
		logger:     logger,
	}
}

// SubmitTransaction submits a transaction to the blockchain
func (h *HyperledgerClient) SubmitTransaction(refType, refID, data string) (*BlockchainResponse, error) {
	h.logger.WithFields(logrus.Fields{
		"ref_type": refType,
		"ref_id":   refID,
	}).Info("Submitting transaction to blockchain")

	// Generate hash of the data
	hash := h.generateHash(data)

	// Create transaction payload
	payload := map[string]interface{}{
		"refType":   refType,
		"refId":     refID,
		"hash":      hash,
		"timestamp": time.Now(),
	}

	// In a real implementation, this would interact with Hyperledger Fabric SDK
	// For now, we'll simulate the blockchain interaction
	txID := h.generateTransactionID(refType, refID)

	h.logger.WithField("payload", payload).Debug("Transaction payload created")

	response := &BlockchainResponse{
		TransactionID: txID,
		Hash:          hash,
		Timestamp:     time.Now(),
		Status:        "PENDING",
	}

	// Simulate blockchain confirmation delay
	go h.simulateConfirmation(txID, hash)

	h.logger.WithFields(logrus.Fields{
		"transaction_id": txID,
		"hash":           hash,
	}).Info("Transaction submitted to blockchain")

	return response, nil
}

// VerifyHash verifies if a hash exists on the blockchain
func (h *HyperledgerClient) VerifyHash(hash string) (*VerificationResult, error) {
	h.logger.WithField("hash", hash).Info("Verifying hash on blockchain")

	// In a real implementation, this would query the blockchain
	// For now, we'll simulate verification
	result := &VerificationResult{
		Valid:         true,
		Hash:          hash,
		TransactionID: "tx-" + hash[:8],
		Timestamp:     time.Now(),
	}

	return result, nil
}

// GetTransaction retrieves transaction details from blockchain
func (h *HyperledgerClient) GetTransaction(txID string) (*TransactionDetails, error) {
	h.logger.WithField("transaction_id", txID).Info("Retrieving transaction from blockchain")

	// In a real implementation, this would query the blockchain
	// For now, we'll simulate transaction retrieval
	details := &TransactionDetails{
		TransactionID: txID,
		Hash:          "0x" + txID[3:],
		RefType:       "CLAIM",
		RefID:         "sample-ref-id",
		BlockNumber:   12345,
		Timestamp:     time.Now(),
		Signer:        "nphies-platform",
		Status:        "CONFIRMED",
	}

	return details, nil
}

// generateHash creates a SHA-256 hash of the input data
func (h *HyperledgerClient) generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// generateTransactionID creates a unique transaction ID
func (h *HyperledgerClient) generateTransactionID(refType, refID string) string {
	payload := fmt.Sprintf("%s-%s-%d", refType, refID, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(payload))
	return "tx-" + hex.EncodeToString(hash[:8])
}

// simulateConfirmation simulates blockchain confirmation process
func (h *HyperledgerClient) simulateConfirmation(txID, hash string) {
	// Simulate network delay
	time.Sleep(5 * time.Second)

	h.logger.WithFields(logrus.Fields{
		"transaction_id": txID,
		"hash":           hash,
	}).Info("Transaction confirmed on blockchain")
}

// CreateDataHash creates a hash for anchoring complex data structures
func CreateDataHash(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:]), nil
}

// ValidateDataIntegrity validates that data matches its stored hash
func ValidateDataIntegrity(data interface{}, expectedHash string) (bool, error) {
	actualHash, err := CreateDataHash(data)
	if err != nil {
		return false, err
	}

	return actualHash == expectedHash, nil
}