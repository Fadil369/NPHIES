package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger *logrus.Logger
}

func NewHandler(logger *logrus.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Health check endpoint
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "terminology-service",
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
		"requests_total": 0,
		"uptime_seconds": 0,
	})
}

// GetCodeSystems returns available code systems
func (h *Handler) GetCodeSystems(c *gin.Context) {
	codeSystems := []map[string]interface{}{
		{
			"url":         "http://snomed.info/sct",
			"name":        "SNOMED CT",
			"version":     "2024-03",
			"status":      "active",
			"description": "SNOMED Clinical Terms",
		},
		{
			"url":         "http://loinc.org",
			"name":        "LOINC",
			"version":     "2.77",
			"status":      "active",
			"description": "Logical Observation Identifiers Names and Codes",
		},
		{
			"url":         "http://hl7.org/fhir/sid/icd-10",
			"name":        "ICD-10-AM",
			"version":     "11th Edition",
			"status":      "active",
			"description": "International Classification of Diseases - Australian Modification",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"total":        len(codeSystems),
		"entry":        codeSystems,
	})
}

// LookupCode looks up a specific code in a code system
func (h *Handler) LookupCode(c *gin.Context) {
	system := c.Param("system")
	code := c.Param("code")

	h.logger.WithFields(logrus.Fields{
		"system": system,
		"code":   code,
	}).Info("Looking up code")

	// Mock response for demonstration
	response := gin.H{
		"resourceType": "Parameters",
		"parameter": []gin.H{
			{
				"name":        "name",
				"valueString": "Sample Code Display",
			},
			{
				"name":        "version",
				"valueString": "2024-03",
			},
			{
				"name":        "display",
				"valueString": "Sample code for " + code,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// ValidateCode validates a code against a code system
func (h *Handler) ValidateCode(c *gin.Context) {
	system := c.Param("system")

	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"system": system,
		"request": request,
	}).Info("Validating code")

	// Mock validation response
	response := gin.H{
		"resourceType": "Parameters",
		"parameter": []gin.H{
			{
				"name":         "result",
				"valueBoolean": true,
			},
			{
				"name":        "message",
				"valueString": "Code validation successful",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// SearchCodes searches for codes in a code system
func (h *Handler) SearchCodes(c *gin.Context) {
	system := c.Param("system")
	query := c.Query("q")

	h.logger.WithFields(logrus.Fields{
		"system": system,
		"query":  query,
	}).Info("Searching codes")

	// Mock search results
	results := []gin.H{
		{
			"code":    "12345",
			"display": "Sample result 1",
			"system":  system,
		},
		{
			"code":    "67890",
			"display": "Sample result 2",
			"system":  system,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"total":        len(results),
		"entry":        results,
	})
}

// MapCodes maps codes between different terminology systems
func (h *Handler) MapCodes(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	h.logger.WithField("request", request).Info("Mapping codes")

	// Mock mapping response
	response := gin.H{
		"resourceType": "ConceptMap",
		"status":       "active",
		"group": []gin.H{
			{
				"source": request["sourceSystem"],
				"target": request["targetSystem"],
				"element": []gin.H{
					{
						"code": request["sourceCode"],
						"target": []gin.H{
							{
								"code":         "mapped-code",
								"display":      "Mapped Code Display",
								"equivalence": "equivalent",
							},
						},
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetConcept gets details about a specific concept
func (h *Handler) GetConcept(c *gin.Context) {
	concept := c.Param("concept")

	h.logger.WithField("concept", concept).Info("Getting concept details")

	// Mock concept response
	response := gin.H{
		"resourceType": "CodeSystem",
		"concept": []gin.H{
			{
				"code":       concept,
				"display":    "Concept Display Name",
				"definition": "This is a sample concept definition",
				"property": []gin.H{
					{
						"code":        "status",
						"valueCode": "active",
					},
					{
						"code":         "effectiveDate",
						"valueDateTime": "2024-01-01",
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, response)
}