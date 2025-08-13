# NPHIES Platform - Project Structure

## Recommended Directory Structure

```
NPHIES/
├── README.md
├── .gitignore
├── LICENSE
├── docker-compose.yml
├── Makefile
│
├── docs/                           # Documentation
│   ├── architecture/               # Architecture diagrams and docs
│   ├── api/                       # API documentation
│   ├── deployment/                # Deployment guides
│   └── Unified_Digital_Healthcare_Insurance_PRD.md
│
├── services/                      # Microservices
│   ├── api-gateway/               # FHIR API Gateway
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── handlers/              # FHIR endpoint handlers
│   │   ├── middleware/            # Auth, rate limiting
│   │   └── config/
│   │
│   ├── eligibility-service/       # Eligibility checking
│   │   ├── Dockerfile
│   │   ├── src/
│   │   ├── tests/
│   │   └── config/
│   │
│   ├── claims-service/            # Claims intake and processing
│   │   ├── Dockerfile
│   │   ├── src/
│   │   ├── tests/
│   │   └── config/
│   │
│   ├── terminology-service/       # Code systems management
│   │   ├── Dockerfile
│   │   ├── src/
│   │   ├── data/                  # Terminology data
│   │   └── config/
│   │
│   └── shared/                    # Shared libraries
│       ├── fhir/                  # FHIR models and utilities
│       ├── kafka/                 # Kafka producers/consumers
│       ├── auth/                  # Authentication utilities
│       └── logging/               # Structured logging
│
├── infrastructure/                # Infrastructure as Code
│   ├── terraform/                 # Terraform modules
│   │   ├── aws/                   # AWS resources
│   │   ├── azure/                 # Azure resources
│   │   ├── gcp/                   # GCP resources
│   │   └── modules/               # Reusable modules
│   │
│   ├── kubernetes/                # K8s manifests
│   │   ├── namespaces/
│   │   ├── services/
│   │   ├── ingress/
│   │   ├── monitoring/
│   │   └── security/
│   │
│   └── helm/                      # Helm charts
│       ├── nphies-platform/       # Main platform chart
│       ├── kafka/                 # Kafka chart
│       ├── redis/                 # Redis chart
│       └── monitoring/            # Observability stack
│
├── scripts/                       # Automation scripts
│   ├── setup/                     # Environment setup
│   ├── deployment/                # Deployment automation
│   ├── testing/                   # Test automation
│   └── data-migration/            # Data migration tools
│
├── tests/                         # Integration and E2E tests
│   ├── integration/               # Service integration tests
│   ├── e2e/                      # End-to-end tests
│   ├── load/                     # Performance tests
│   └── fixtures/                 # Test data
│
├── configs/                       # Configuration files
│   ├── development/               # Dev environment configs
│   ├── staging/                   # Staging environment configs
│   ├── production/                # Production environment configs
│   └── schemas/                   # Kafka schemas
│
├── monitoring/                    # Observability configurations
│   ├── prometheus/                # Prometheus configs
│   ├── grafana/                   # Grafana dashboards
│   ├── jaeger/                    # Tracing configs
│   └── alerts/                    # Alerting rules
│
└── .github/                       # GitHub workflows
    ├── workflows/                 # CI/CD pipelines
    ├── ISSUE_TEMPLATE/            # Issue templates
    └── PULL_REQUEST_TEMPLATE.md
```

## Service Architecture

### Core Services

1. **API Gateway**

   - Technology: Go with Gin/Echo framework
   - Responsibilities: FHIR R4 endpoints, authentication, rate limiting
   - Dependencies: Redis (caching), Auth service

2. **Eligibility Service**

   - Technology: Go or Java with Spring Boot
   - Responsibilities: Coverage validation, member status checks
   - Dependencies: PostgreSQL, Redis, Kafka

3. **Claims Service**

   - Technology: Java with Spring Boot
   - Responsibilities: Claims intake, validation, basic processing
   - Dependencies: PostgreSQL, Kafka, Terminology service

4. **Terminology Service**
   - Technology: Go or Python with FastAPI
   - Responsibilities: Code system management, lookup, mapping
   - Dependencies: PostgreSQL, Redis

### Infrastructure Components

1. **Apache Kafka**

   - Event streaming backbone
   - Topics for claims, eligibility, audit events
   - Schema registry for event validation

2. **PostgreSQL**

   - Primary database for transactional data
   - Separate databases per service
   - Connection pooling and read replicas

3. **Redis**

   - Caching layer for eligibility and terminology
   - Session management
   - Rate limiting counters

4. **Observability Stack**
   - Prometheus + Grafana for metrics
   - Jaeger for distributed tracing
   - ELK/OpenSearch for logs

## Development Guidelines

### Language and Framework Choices

- **Go**: API Gateway, Terminology Service (performance critical)
- **Java/Spring Boot**: Claims Service, complex business logic
- **Python/FastAPI**: Data processing, ML services (future)
- **React**: Patient-facing applications (future phases)

### Code Standards

- Follow language-specific best practices
- Implement comprehensive error handling
- Use structured logging (JSON format)
- Include health check endpoints (/health, /ready)
- Implement graceful shutdown patterns

### Security Considerations

- Never commit secrets to version control
- Use environment variables for configuration
- Implement mTLS between services
- Encrypt sensitive data at rest
- Audit all data access

### Testing Strategy

- Unit tests for business logic (>80% coverage)
- Integration tests for service interactions
- Contract tests for API compatibility
- Load tests for performance validation
- Chaos engineering for resilience testing

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Kubernetes cluster (local or cloud)
- Helm 3.x
- Terraform (for infrastructure)
- Go 1.19+ and Java 11+ (for development)

### Quick Start Commands

```bash
# Clone repository
git clone https://github.com/Fadil369/NPHIES.git
cd NPHIES

# Start local development environment
make dev-up

# Build all services
make build

# Run tests
make test

# Deploy to staging
make deploy-staging
```

### Development Workflow

1. Create feature branch from main
2. Implement changes with tests
3. Run local validation (lint, test, security scan)
4. Create pull request
5. Automated CI/CD pipeline runs
6. Code review and approval
7. Merge and deploy

This structure provides a solid foundation for implementing the NPHIES platform while maintaining scalability, security, and operational excellence.
