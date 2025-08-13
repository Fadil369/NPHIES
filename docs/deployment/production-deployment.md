# NPHIES Platform - Production Deployment Guide

## Overview

This document provides comprehensive guidance for deploying the NPHIES platform in a production environment with monitoring, logging, backup, and disaster recovery strategies.

## Infrastructure Requirements

### Minimum Production Requirements

- **Kubernetes Cluster**: v1.24+
- **Nodes**: 6-10 worker nodes (8 CPU, 32GB RAM each)
- **Storage**: 500GB+ persistent storage for databases
- **Network**: Load balancer with SSL termination
- **Monitoring**: Prometheus + Grafana stack
- **Logging**: ELK/OpenSearch stack

### Recommended Production Setup

```yaml
Production Cluster:
  - 3 Control Plane nodes (4 CPU, 16GB RAM)
  - 6-12 Worker nodes (16 CPU, 64GB RAM)
  - High-availability storage (1TB+ with replication)
  - Multi-AZ deployment for disaster recovery
```

## Deployment Strategy

### 1. Infrastructure Setup

#### Prerequisites
```bash
# Install required tools
kubectl, helm, kustomize, argocd-cli

# Verify cluster access
kubectl cluster-info
kubectl get nodes
```

#### Namespace and Resource Setup
```bash
# Apply namespaces
kubectl apply -f infrastructure/kubernetes/namespaces/

# Apply ConfigMaps and Secrets
kubectl apply -f infrastructure/kubernetes/configmaps/

# Create storage classes for persistent volumes
kubectl apply -f infrastructure/kubernetes/storage/
```

### 2. Database and Infrastructure Services

#### PostgreSQL High Availability
```bash
# Deploy PostgreSQL with replication
helm install postgres-cluster bitnami/postgresql-ha \
  --namespace nphies-data \
  --set postgresql.replicaCount=3 \
  --set persistence.size=100Gi
```

#### Redis Cluster
```bash
# Deploy Redis cluster
helm install redis-cluster bitnami/redis-cluster \
  --namespace nphies-data \
  --set cluster.nodes=6 \
  --set persistence.size=20Gi
```

#### Kafka Cluster
```bash
# Deploy Kafka cluster
helm install kafka-cluster bitnami/kafka \
  --namespace nphies-data \
  --set replicaCount=3 \
  --set persistence.size=50Gi
```

### 3. Application Services Deployment

#### Blue-Green Deployment Strategy
```bash
# Deploy applications with zero-downtime
kubectl apply -f infrastructure/kubernetes/deployments/

# Verify deployment status
kubectl get deployments -n nphies-core
kubectl get pods -n nphies-core
```

#### Rolling Update Configuration
```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxUnavailable: 25%
    maxSurge: 25%
```

## Monitoring and Observability

### 1. Prometheus Setup

#### Metrics Collection
```yaml
# Prometheus configuration for NPHIES services
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "nphies_rules.yml"

scrape_configs:
  - job_name: 'nphies-services'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
```

#### Key Metrics to Monitor
- **API Gateway**: Request rate, response time, error rate
- **Claims Service**: Claims processing rate, validation errors
- **Eligibility Service**: Cache hit rate, eligibility check latency
- **Terminology Service**: Code lookup performance, cache efficiency

### 2. Grafana Dashboards

#### Essential Dashboards
1. **System Overview**: Cluster health, resource utilization
2. **Application Performance**: Service metrics, SLA compliance
3. **Business Metrics**: Claims volume, processing times
4. **Error Tracking**: Error rates, failed transactions

#### Sample Dashboard Configuration
```json
{
  "dashboard": {
    "title": "NPHIES Platform Overview",
    "panels": [
      {
        "title": "Claims Processing Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(claims_processed_total[5m])"
          }
        ]
      }
    ]
  }
}
```

### 3. Alerting Rules

#### Critical Alerts
```yaml
groups:
  - name: nphies.critical
    rules:
      - alert: ServiceDown
        expr: up{job=~"nphies-.*"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "NPHIES service is down"
      
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
```

## Logging Strategy

### 1. Centralized Logging with ELK Stack

#### Elasticsearch Configuration
```yaml
elasticsearch:
  replicas: 3
  minimumMasterNodes: 2
  resources:
    requests:
      cpu: 1000m
      memory: 2Gi
    limits:
      cpu: 2000m
      memory: 4Gi
```

#### Logstash Processing
```ruby
input {
  beats {
    port => 5044
  }
}

filter {
  if [fields][service] == "nphies" {
    json {
      source => "message"
    }
    
    # Add correlation ID tracking
    if [correlation_id] {
      mutate {
        add_tag => ["correlated"]
      }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "nphies-logs-%{+YYYY.MM.dd}"
  }
}
```

### 2. Log Retention and Archival

#### Retention Policy
- **Application Logs**: 30 days hot, 90 days warm, 1 year cold
- **Audit Logs**: 7 years (regulatory requirement)
- **System Logs**: 14 days hot, 30 days warm

#### Archival Strategy
```bash
# Automated archival using curator
curator --config curator.yml delete_indices.yml

# S3 archival for long-term storage
aws s3 sync /var/log/elasticsearch s3://nphies-logs-archive/
```

## Backup and Data Protection

### 1. Database Backup Strategy

#### PostgreSQL Backup
```bash
# Daily full backup
pg_dump -h $DB_HOST -U $DB_USER -d nphies_claims > backup_$(date +%Y%m%d).sql

# Point-in-time recovery setup
wal-e backup-push
```

#### Backup Schedule
- **Full Backup**: Daily at 2 AM
- **Incremental Backup**: Every 4 hours
- **Transaction Log Backup**: Continuous

### 2. Application State Backup

#### Kubernetes Resource Backup
```bash
# Backup all Kubernetes resources
kubectl get all --all-namespaces -o yaml > k8s-backup-$(date +%Y%m%d).yaml

# Use Velero for comprehensive backup
velero backup create nphies-backup --include-namespaces nphies-core,nphies-data
```

#### Redis Backup
```bash
# Redis snapshot backup
redis-cli --rdb backup.rdb
aws s3 cp backup.rdb s3://nphies-redis-backup/
```

## Disaster Recovery

### 1. Recovery Time Objectives (RTO) and Recovery Point Objectives (RPO)

- **RTO**: 30 minutes for critical services
- **RPO**: 5 minutes (maximum data loss)
- **Availability Target**: 99.95%

### 2. Disaster Recovery Procedures

#### Complete System Recovery
```bash
# 1. Restore infrastructure
kubectl apply -f infrastructure/kubernetes/

# 2. Restore databases
pg_restore -h $DB_HOST -U $DB_USER -d nphies_claims backup.sql

# 3. Verify service health
kubectl get pods -n nphies-core
tests/integration/run-integration-tests.sh

# 4. Update DNS to point to recovered environment
```

#### Partial Service Recovery
```bash
# Rolling restart of specific service
kubectl rollout restart deployment/claims-service -n nphies-core

# Verify recovery
kubectl rollout status deployment/claims-service -n nphies-core
```

### 3. Failover Procedures

#### Database Failover
```bash
# Promote read replica to primary
pg_promote -D /var/lib/postgresql/data

# Update service configurations
kubectl patch deployment claims-service -p '{"spec":{"template":{"spec":{"containers":[{"name":"claims-service","env":[{"name":"SPRING_DATASOURCE_URL","value":"jdbc:postgresql://postgres-replica:5432/nphies_claims"}]}]}}}}'
```

#### Service Failover
```bash
# Scale up standby region
kubectl scale deployment --replicas=3 -n nphies-core

# Update load balancer configuration
kubectl patch service api-gateway-service -p '{"spec":{"selector":{"region":"standby"}}}'
```

## Security and Compliance

### 1. Security Monitoring

#### Security Alerts
```yaml
- alert: UnauthorizedAccess
  expr: increase(http_requests_total{status="401"}[5m]) > 10
  labels:
    severity: warning

- alert: DataExfiltration
  expr: rate(egress_bytes_total[5m]) > 1000000
  labels:
    severity: critical
```

### 2. Compliance Auditing

#### Audit Log Format
```json
{
  "timestamp": "2025-08-13T10:30:00Z",
  "user_id": "user123",
  "action": "claim_submission",
  "resource": "claim_id_12345",
  "result": "success",
  "ip_address": "192.168.1.100"
}
```

## Performance Optimization

### 1. Resource Scaling

#### Horizontal Pod Autoscaling
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: claims-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: claims-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 2. Performance Tuning

#### Database Optimization
```sql
-- Index optimization for claims queries
CREATE INDEX CONCURRENTLY idx_claims_provider_date 
ON claims(provider_id, service_date);

-- Partition large tables
CREATE TABLE claims_2025 PARTITION OF claims 
FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
```

#### Cache Optimization
```yaml
redis:
  maxmemory: 4gb
  maxmemory-policy: allkeys-lru
  save: "900 1 300 10 60 10000"
```

## Troubleshooting Guide

### Common Issues

1. **Service Not Starting**
   ```bash
   kubectl logs deployment/claims-service -n nphies-core
   kubectl describe pod <pod-name> -n nphies-core
   ```

2. **Database Connection Issues**
   ```bash
   kubectl exec -it postgres-0 -n nphies-data -- psql -U postgres
   ```

3. **Performance Degradation**
   ```bash
   kubectl top pods -n nphies-core
   kubectl get hpa -n nphies-core
   ```

### Emergency Contacts

- **On-Call Engineer**: +1-XXX-XXX-XXXX
- **Database Administrator**: +1-XXX-XXX-XXXX
- **Security Team**: security@nphies.sa
- **Management Escalation**: management@nphies.sa

## Conclusion

This production deployment guide provides the foundation for a robust, scalable, and highly available NPHIES platform. Regular review and updates of these procedures ensure continued operational excellence and compliance with healthcare regulations.