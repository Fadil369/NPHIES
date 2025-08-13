# Product Requirements Document (PRD)
Unified Digital Healthcare Insurance Platform for Saudi Arabia (Aligned with NPHIES & Vision 2030)

Version: 1.0  
Date: 2025-08-13  
Owner: Product / Architecture Joint Working Group  
Status: Draft

--------------------------------------------------
1. Executive Summary
--------------------------------------------------
This platform will unify national healthcare insurance interactions (eligibility, prior authorization, claims, payments, fraud detection, patient benefits, population analytics) across payers, providers, regulators, and citizens. It builds on and extends NPHIES by adding a cloud-native, event-driven, AI-augmented services layer. The solution aligns with Vision 2030 goals of digital transformation, health system efficiency, localization of healthtech, improved patient outcomes, and sustainable cost management.

High-level Outcomes:
- Reduce average claim adjudication time from days to minutes (<15 min target; <2 min for straight-through claims)
- Lower fraud/leakage by ≥25% in 3 years
- Achieve ≥90% electronic prior authorization auto-decisioning in <60 seconds
- Enable real-time eligibility lookups <1 second P99
- Provide transparent pre-service cost estimates to ≥70% of insured population within 24 months
- Create longitudinal patient benefits & utilization dataset for population health analytics

--------------------------------------------------
2. Business Objectives & KPIs
--------------------------------------------------
Primary Objectives:
1. Operational Efficiency: Automate and standardize claims and prior auth flows.
2. Financial Integrity: Detect and prevent fraud, waste, abuse earlier in the lifecycle.
3. Patient Empowerment: Provide wallet, benefits clarity, pricing transparency.
4. Regulatory Insight: Near real-time national health economics analytics.
5. Interoperability: Enforce standards adoption (FHIR R4, ICD-10-AM, ACHI, LOINC, SNOMED CT).
6. AI Enablement: Embed governed ML for fraud, risk stratification, coding QA.

Key KPIs (Initial + Targets):
- Eligibility API P99 latency: Initial 1500 ms → Target <900 ms
- Prior auth auto-approval rate: 30% → 70% (Year 2) → 90% (Year 3)
- Average prior auth decision time: >24h → <5 min (Year 2) → <60 sec (Year 3)
- Straight-through processing (STP) claims %: 20% → 55% (Year 2) → 75% (Year 3)
- Fraud detection precision@top1% alerts: ≥0.6 → ≥0.8 (Year 3)
- Claim rejection rework rate: 18% → 8% (Year 3)
- Data completeness (FHIR compliance score): Baseline TBD → ≥95%
- Patient app monthly active rate of insured population: 0% → 40% (Year 2) → 65% (Year 3)
- Time to integrate a new provider system: 8–12 weeks → ≤3 weeks
- Cost per processed claim (avg ops overhead): Baseline (index=1.0) → 0.65 (Year 3)

--------------------------------------------------
3. Stakeholders & Personas
--------------------------------------------------
Stakeholder Groups:
- Regulators: NPHIES (NHIC), Ministry of Health, SDAIA (data governance), CST (telecom + cloud compliance)
- Payers: Insurance companies / TPA
- Providers: Hospitals, Clinics, Pharmacies, Labs
- Patients / Members
- Technology Partners: Cloud providers (stc Cloud, Oracle Cloud Jeddah, AWS Saudi Region when available), AI vendors
- Auditors / Compliance Officers
- Data Scientists & Actuaries

Personas:
1. Claims Adjudicator: Needs consolidated view, risk signals, coding suggestions.
2. Provider Billing Clerk: Needs pre-check claim validations, coding normalization.
3. Prior Auth Nurse Reviewer: Needs clinical guidelines, decision support, history.
4. Patient / Member: Wants to see benefits, remaining limits, claim status, cost estimates.
5. Fraud Analyst: Needs ranked anomaly cases, investigation workflow, explainability.
6. Data Scientist: Requires governed feature store, reproducible experimentation.
7. Regulator Analyst: Needs aggregated de-identified utilization and cost dashboards.

--------------------------------------------------
4. Scope
--------------------------------------------------
In-Scope (Phases defined below):
- Central API gateway and integration adapters (FHIR, X12/EDI, proprietary EMRs)
- Event streaming backbone (Kafka/Kinesis)
- Core microservices (Eligibility, Prior Auth, Claims, Fraud, Coding Normalization, Payment Orchestration, Notification, Terminology)
- AI/ML services (Fraud, Risk Stratification, Code Suggestion, Cost Forecasting)
- Patient digital wallet & benefits app (mobile + web)
- Data lake + warehouse + semantic layer + governed data catalog
- Blockchain sidechain for consent & immutable claims hash references
- Observability platform (traces, metrics, logs, audit)
- Security & compliance enforcement (zero-trust, PDPL alignment)
- MLOps pipeline & model governance registry

Out of Scope (initial release):
- Direct clinical EMR functionalities (focus is insurance exchange)
- Full telemedicine platform (only integration points)
- GenAI conversational clinical coding assistant (future enhancement)
- IoT device management platform (only ingestion interface placeholder in Phase 3)
- Payment rails settlement engine (integrate with existing bank/fintech systems)

--------------------------------------------------
5. Assumptions & Constraints
--------------------------------------------------
Assumptions:
- Provider systems can adopt FHIR R4 gradually; interim adapters required.
- Saudi data residency mandatory; cryptographic keys stored locally (HSM).
- National digital identity (Absher/eID) integration is accessible via OAuth2/OpenID Connect.
- Terminology services license for SNOMED CT, LOINC legally secured.
- PDPL and future health data regulations continue to evolve—architect for adaptability.

Constraints:
- Latency requirements for real-time check operations (<1 sec).
- Multi-lingual support (Arabic primary, English secondary) for patient interfaces.
- Intermittent provider connectivity in rural areas → need eventual consistency & retry.
- Regulatory audit logs retention (≥10 years immutable).

--------------------------------------------------
6. High-Level Architecture
--------------------------------------------------
Layers:
1. Edge & Access: API Gateway (REST FHIR + GraphQL) / WAF / Rate Limiter
2. Integration Adapters: EMR Connectors (HL7 v2 → FHIR), EDI Translators (X12 837/835)
3. Event Bus: Kafka clusters (ClaimsEvents, EligibilityEvents, AuthEvents, FraudAlerts, AuditTrail)
4. Microservices (each stateless + sidecar for security):
   - Eligibility Service
   - Prior Authorization Service
   - Claims Intake Service
   - Claims Adjudication Engine
   - Coding Normalization (NLP)
   - Fraud Scoring Service
   - Payment Orchestrator
   - Terminology Service
   - Consent & Wallet Service
   - Patient App Backend
   - Notification Service (SMS, email, push)
   - Analytics Extractor / Data Ingestion
5. Data Layer:
   - OLTP: PostgreSQL/CockroachDB
   - Cache: Redis for low-latency eligibility & provider config
   - Object Store: S3-compatible (medical docs, attachments)
   - Warehouse: Snowflake/BigQuery/ClickHouse (depends on residency compliance; Snowflake on regional cloud or ClickHouse self-hosted)
   - Data Lake: Raw, curated, feature store partitions (Parquet)
   - Blockchain Sidechain: Consortium nodes (hash claims summaries & consent transactions)
6. AI & MLOps:
   - Feature Store (Feast or custom) sourced from Kafka & curated lake
   - Model Registry (MLflow or SageMaker Model Registry)
   - Serving: Real-time (REST/gRPC) + Batch (scheduled)
7. Observability:
   - Metrics: Prometheus / Cloud native
   - Traces: OpenTelemetry → Jaeger/Tempo
   - Logs: Elastic/OpenSearch
   - Audit WORM storage (Immutable S3 bucket + Glacier tiering)
8. Security Fabric:
   - mTLS service mesh (Istio/Linkerd)
   - IAM / RBAC / ABAC
   - Secrets Management (Vault / KMS)
   - Tokenization of PHI fields

--------------------------------------------------
7. Detailed Functional Requirements
--------------------------------------------------
7.1 API Gateway & Integration
FR-API-01: Support FHIR R4 resources: Patient, Coverage, Claim, ClaimResponse, PriorAuthorization (PA as Claim with use=preauthorization), Practitioner, Organization, Encounter, Procedure, Medication.
FR-API-02: Provide GraphQL endpoint for aggregated multi-entity queries (e.g., claimsByMember(memberId, lastNMonths)).
FR-API-03: Enforce OAuth2 + OIDC tokens mapped to Absher/eID subject.
FR-API-04: Rate limiting configurable per client (default 500 req/min).
FR-API-05: Provide API key + mutual TLS for system-to-system provider integrations.
FR-API-06: X12 837/835 translation endpoints with acknowledgement events.

7.2 Eligibility Service
FR-ELIG-01: Validate coverage using payer rules and member status cache.
FR-ELIG-02: Return determination (Active, Inactive, Limited) <900 ms P99.
FR-ELIG-03: Cache refresh TTL default 5 minutes; proactive invalidation on coverage update events.

7.3 Prior Authorization Service
FR-PA-01: Accept structured request (FHIR Claim use=preauthorization).
FR-PA-02: Invoke rules engine (clinical guidelines, necessity rules).
FR-PA-03: AI model to predict approval likelihood; auto-approve if confidence ≥ threshold and no rule conflicts.
FR-PA-04: Provide real-time status streaming (topic: PriorAuthStatus).
FR-PA-05: Support manual review workflow with SLA timers & escalation.

7.4 Claims Intake & Adjudication
FR-CLM-01: Ingest claims from API/EDI/batch via standardized envelope.
FR-CLM-02: Validate mandatory FHIR elements & coding sets.
FR-CLM-03: Invoke coding normalization (NLP) for ambiguous or legacy codes.
FR-CLM-04: Run adjudication rules (benefit limits, coordination of benefits, bundling).
FR-CLM-05: Produce ClaimResponse with line-level decisions & explanation of benefit (EOB).
FR-CLM-06: Support straight-through processing if no manual flags raised.
FR-CLM-07: Persist claim hash to blockchain sidechain for integrity.

7.5 Coding Normalization
FR-CODE-01: Map local codes to ICD-10-AM, ACHI, SNOMED CT, LOINC using terminology index.
FR-CODE-02: Provide code suggestions with confidence values.
FR-CODE-03: Flag inconsistent code combinations (procedure w/o appropriate diagnosis).

7.6 Fraud Detection
FR-FRD-01: Real-time scoring on claim ingest (<500 ms additional latency budget).
FR-FRD-02: Hybrid models (supervised + unsupervised outlier detection).
FR-FRD-03: Explanation interface (top contributing features).
FR-FRD-04: Alert streaming (FraudAlerts topic) with prioritization.
FR-FRD-05: Feedback loop—investigator disposition labels fed back to training set.

7.7 Payment Orchestration
FR-PAY-01: Generate settlement instructions after claim finalization.
FR-PAY-02: Integrate with external financial systems through secure APIs.
FR-PAY-03: Track status (Initiated, Pending, Completed, Reconciled).

7.8 Consent & Wallet
FR-WLT-01: Maintain patient wallet ledger (benefits remaining, prior auth approvals).
FR-WLT-02: Record consents (data sharing, telemedicine) hashed on sidechain.
FR-WLT-03: Display pre-service cost estimate (range & coverage portion).
FR-WLT-04: Multi-language UI (Arabic default).

7.9 Patient App
FR-PAPP-01: Secure login via national eID.
FR-PAPP-02: View benefits, claim history, prior auth statuses.
FR-PAPP-03: Push notifications for claim milestones, expiring authorizations.
FR-PAPP-04: Provide provider search (coverage network) & indicative rates.
FR-PAPP-05: Consent management UI.

7.10 Analytics & Population Health
FR-ANL-01: Daily ingestion of curated claims + coverage + outcomes to warehouse.
FR-ANL-02: Risk stratification batch model run nightly.
FR-ANL-03: Provide API for aggregated metrics (cost per member per month, chronic disease prevalence).
FR-ANL-04: De-identification pipeline for research queries.

7.11 Terminology Service
FR-TERM-01: CRUD for code systems (versioned).
FR-TERM-02: Bulk import SNOMED CT releases.
FR-TERM-03: Fast concept lookup (<50 ms cache hit).
FR-TERM-04: Concept mapping historical traceability.

7.12 Notification Service
FR-NOT-01: Event-driven push (eligibility change, auth approval, claim paid).
FR-NOT-02: Preference management (channels: SMS, email, push).
FR-NOT-03: Localized templates (Arabic/English).

--------------------------------------------------
8. Non-Functional Requirements (NFR)
--------------------------------------------------
Performance:
- Eligibility P99 <0.9s; Prior auth decision (auto) P95 <1.5s
- Throughput: Sustained 300 claims/sec national peak (scalable >5x)

Scalability:
- Horizontal scaling via Kubernetes HPA & partitioned Kafka topics
- Data retention: Raw events ≥7 years (tiered)

Availability:
- Core transactional services ≥99.95% SLA
- Active-active across Riyadh & Jeddah (RPO ≤5 min, RTO ≤30 min)

Security:
- Zero-trust: mTLS + SPIFFE identities
- AES-256 at rest (KMS), TLS 1.3 in transit
- Role-based & attribute-based access; fine-grained scoping to patient-level resources

Compliance:
- PDPL, Saudi Health Data Regulations, alignment with HIPAA principles
- Audit immutability (WORM S3 + blockchain hash references)

Observability:
- 100% trace sampling for regulated flows; 10% for low-risk flows
- Real-time anomaly detection in log pipeline

Data Quality:
- Daily data completeness & conformance scoring
- Automatic quarantine of malformed events

Localization:
- Full Arabic UI & RTL support
- Encoding normalization (UTF-8) across services

Privacy:
- Differential privacy options for population analytics exports
- Pseudonymization for model training datasets

Resilience:
- Chaos testing monthly
- Backpressure strategies (circuit breakers, queue depth thresholds)

--------------------------------------------------
9. Data Model (Conceptual Highlights)
--------------------------------------------------
Core Entities:
- Member (Patient, Coverage, BenefitLimits)
- Claim (Header, LineItems, Codes, Attachments, AdjudicationResults)
- PriorAuthorization (Request, Decision, RulesTrace)
- Provider (Organization, Practitioner, Facility)
- CodeSystem / TerminologyMap
- FraudCase (Scores, Features, InvestigatorNotes)
- Consent (Type, Scope, Status, HashRef)
- WalletTransaction (Type: BenefitDebit/Credit, Timestamp, Amount, BalanceImpact)

Key Relationships:
- Member 1..* Coverage
- Coverage 1..* Claim
- Claim 0..1 PriorAuthorization reference
- Claim 1..* LineItem with codes
- Claim 0..1 FraudCase
- Consent 1..1 Member; hashed to Sidechain
- WalletTransaction aggregates by Member

Partitioning Strategy:
- Kafka topics partitioned by memberId / providerId depending on domain
- OLTP sharding by memberId hash modulo N

--------------------------------------------------
10. Event Streaming Design
--------------------------------------------------
Primary Topics (example config):
- claims.intake.v1 (Partitions: 48; Key: claimId)
- claims.adjudicated.v1
- eligibility.requests.v1
- eligibility.responses.v1
- priorauth.requests.v1
- priorauth.status.v1
- fraud.alerts.v1
- coding.suggestions.v1
- wallet.transactions.v1
- audit.trail.v1 (compacted + archived)
- model.inference.events.v1
- terminology.updates.v1

Schema Registry:
- Avro/JSON Schema with semantic versioning (MAJOR incompatible, MINOR additive)
- Compatibility: BACKWARD by default

Retention:
- Operational topics: 7–30 days
- Compacted state topics (e.g., coverage snapshot)
- Long-term archival to lake (batch export every hour)

Event Envelope:
{
  "eventId": "...",
  "eventType": "ClaimAdjudicated",
  "version": "1.0",
  "timestamp": "...",
  "correlationId": "...",
  "producer": "...",
  "payload": {...},
  "hash": "sha256(payload)"
}

--------------------------------------------------
11. API Design Guidelines
--------------------------------------------------
- FHIR REST endpoints: /fhir/Claim, /fhir/Patient, /fhir/Coverage
- GraphQL: Single endpoint /graphql with persisted queries
- Pagination: Cursor-based for large lists
- Error Format: RFC 7807 Problem+JSON with FHIR OperationOutcome mapping
- Idempotency: Idempotency-Key header for POST claim submissions
- Versioning: /v1 prefix for non-FHIR endpoints; FHIR version in metadata

--------------------------------------------------
12. AI / ML Components & Governance
--------------------------------------------------
Models:
- Fraud Detection (Supervised Gradient Boosting + Unsupervised Isolation Forest)
- Coding Suggestion (NLP: fine-tuned transformer on coding pairs)
- Prior Auth Auto-Approval Classifier
- Risk Stratification (Chronic disease model)
- Cost Forecasting (Time-series + gradient boosting hybrid)

Feature Store:
- Real-time features (last N claim counts, provider anomaly scores)
- Batch features (12-month cost aggregates, comorbidity indices)

MLOps Lifecycle:
- Data lineage tracked via metadata catalog
- Model registry with version & performance metrics
- Shadow deployment for new model versions
- Bias & fairness checks (age, gender, region)
- Drift detection (population stability index thresholds)

Explainability:
- SHAP value computation stored for each high-risk fraud alert
- Retain inference logs 24 months

Retraining Cadence:
- Fraud model: Monthly incremental
- Coding model: Quarterly
- Cost forecast: Monthly rolling window
- Prior auth model: Bi-weekly if data volume sufficient

--------------------------------------------------
13. Blockchain Sidechain Scope
--------------------------------------------------
Purpose:
- Immutable referencing for claim summaries & consent records (hash anchoring)
- Not for full PHI storage; only hashed pointers & metadata

Technology:
- Consortium permissioned chain (e.g., Hyperledger Fabric or Quorum)
- Nodes: Regulator, 2 major payers, 2 major providers
Data Structure:
- Transaction { refType: CLAIM|CONSENT, refId, hash, timestamp, signer }
SLAs:
- Hash submission latency <5 seconds

--------------------------------------------------
14. Security & Compliance Details
--------------------------------------------------
Identity:
- Service identities via SPIFFE IDs
- User identities via OIDC (Absher integration)

Authorization:
- Fine-grained scopes: claim.read, claim.write, priorauth.approve, fraud.review
- Attribute policies (providerId, memberId scoping)

Data Protection:
- PHI fields tokenized in logs
- Encryption in DB (column-level for sensitive fields)
- Key rotation every 90 days

Audit:
- Immutable logs: Event + Action + Actor + IP + Before/After (masked)
- Real-time audit anomaly detection (e.g., mass record access)

Privacy:
- Right-to-access & right-to-erasure workflows (PDPL) with restricted exceptions (regulator override)
- Data minimization enforcement in feature extraction

--------------------------------------------------
15. Patient-Facing Experience (High-Level User Journey)
--------------------------------------------------
Journey: Viewing Claim Status
1. Login via eID
2. Dashboard fetches wallet balance & open claims
3. Select claim → timeline (Submitted → In Review → Adjudicated → Paid)
4. Tap for Explanation of Benefits (cost breakdown & coverage)
5. Option to dispute (opens ticket workflow)

Journey: Pre-Service Cost Estimate
1. Select provider/service
2. System queries eligibility + typical coding package
3. Runs cost estimator with historical provider-specific allowed amounts
4. Displays estimate: Provider price range, covered amount, patient share
5. Option to request prior authorization if required

--------------------------------------------------
16. Deployment & DevOps
--------------------------------------------------
Kubernetes:
- Namespaces: core, data, ml, integration, security
- Service mesh injected sidecars
CI/CD:
- GitHub Actions triggers; security scans (SAST, DAST), IaC validation (Terraform)
- ArgoCD for GitOps sync to clusters
Environments:
- Dev → Test → Staging → Production (promotion gates)
- Data anonymization in non-prod
Infrastructure as Code:
- Terraform modules for network, clusters, databases
- Policy-as-code (OPA) for security guardrails
Disaster Recovery:
- Cross-region replication (Riyadh ↔ Jeddah)
- Quarterly DR simulation tests

--------------------------------------------------
17. Phased Implementation Roadmap
--------------------------------------------------
Phase 1 (0–6 months):
- Deliver: API Gateway (FHIR baseline), Eligibility Service, Kafka foundation (3 clusters: core, analytics, audit), Basic Claims Intake, Terminology Service MVP, Observability baseline
- KPIs: Eligibility latency baseline, FHIR compliance 70%
- Risks Mitigated: Integration standardization early

Phase 2 (6–18 months):
- Add: Prior Auth Engine (rules + ML pilot), Full Claims Adjudication rules, Coding Normalization NLP, Fraud Detection v1, Patient App MVP (wallet & claim tracking), Model Registry & Feature Store, Payment Orchestrator
- KPIs: Prior auth auto-approval 50%, Fraud precision 0.65
- Scaling: Kafka partitions expansion, multi-region active-active

Phase 3 (18–36 months):
- Add: National Digital Health Wallet (blockchain anchoring), Predictive Risk Stratification, Cost Forecasting, IoT ingestion adapters, Advanced analytics APIs, Differential privacy exports
- KPIs: STP claims 75%, Fraud precision 0.8, Coverage of patient app 65%
- Hardening: Advanced privacy, dynamic consent, continuous learning loops

--------------------------------------------------
18. Risk Register (Selected)
--------------------------------------------------
| ID | Risk | Impact | Likelihood | Mitigation |
| R1 | Provider slow FHIR adoption | High | Medium | Adapters + phased mandates |
| R2 | Model bias leading to unfair denials | High | Medium | Fairness audits + human override |
| R3 | Latency spikes during peaks | Medium | High | Autoscaling + performance load tests |
| R4 | Regulatory changes (PDPL updates) | Medium | Medium | Modular policy layer |
| R5 | Blockchain consortium governance delays | Low | Medium | Start with limited nodes MVP |
| R6 | Data quality from legacy EMRs | High | High | Validation pipelines + quarantine |
| R7 | Security breach attempt | High | Medium | Zero-trust, continuous monitoring |
| R8 | Talent scarcity (ML & health informatics) | Medium | Medium | Training program + vendor augmentation |

--------------------------------------------------
19. Open Questions
--------------------------------------------------
1. Final decision on warehouse technology (Snowflake vs ClickHouse vs BigQuery regional availability)?
2. Which terminology hosting strategy (self-host vs managed) aligned with licensing constraints?
3. Blockchain sidechain governance charter—who operates validating nodes?
4. Differential privacy parameters acceptable to regulators?
5. Standard SLA for manual prior auth review (target <4 hours?)—regulator alignment needed.
6. EMR vendor prioritized list for adapter development sequencing.
7. Data retention exact durations per data class (claims, audit, model features)?

--------------------------------------------------
20. Success Metrics & Post-Launch Monitoring
--------------------------------------------------
- Time-to-integrate new provider (weeks → days)
- Reduction in re-submitted claims
- Member NPS for digital experience
- Model drift frequency vs expected
- Fraud recovery amounts vs baseline
- Infrastructure cost per processed claim trending downward
- Regulatory audit pass rate (≥98% requirements satisfied)

--------------------------------------------------
21. Appendix: SLA & SLO Targets
--------------------------------------------------
Service | SLO Latency (P95) | Error Budget / Month | Availability Target
Eligibility | 800 ms | 0.5% | 99.95%
Prior Auth Auto Decision | 1500 ms | 1% | 99.9%
Claims Intake | 1200 ms | 1% | 99.9%
Fraud Scoring | +500 ms additive | 1% | 99.9%
Patient App API | 1000 ms | 1% | 99.9%

--------------------------------------------------
22. Appendix: Technology Choices (Preliminary)
--------------------------------------------------
- Gateway: Kong / Apigee / AWS API Gateway (depending on residency)
- Event Streaming: Apache Kafka (self-managed or MSK-like in-region)
- Databases: CockroachDB for horizontal scaling, PostgreSQL for smaller services
- Cache: Redis (clustered)
- NLP: HuggingFace transformers fine-tuned on coding corpora
- Rules Engine: Drools / OpenCDS for clinical necessity + benefit rules
- ML Platform: SageMaker (if regionally compliant) or on-prem Kubeflow
- Blockchain: Hyperledger Fabric (permissioned)
- IaC: Terraform + OPA policies
- Observability: OpenTelemetry + Prometheus + Grafana + OpenSearch

--------------------------------------------------
23. Change Control
--------------------------------------------------
All modifications to this PRD require Architecture + Product governance board approval with documented impact assessment.

--------------------------------------------------
End of Document
--------------------------------------------------