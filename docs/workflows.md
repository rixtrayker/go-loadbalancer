# Use Cases & Workflows

This document describes the key use cases and workflows of the Go Load Balancer, illustrating how the system handles different scenarios.

## Core Workflows

### Request Processing Workflow

The primary workflow of the load balancer is processing and routing HTTP requests to backend servers.

```mermaid
sequenceDiagram
    participant Client
    participant LoadBalancer
    participant Router
    participant PolicyEngine
    participant ServerPool
    participant Backend
    
    Client->>LoadBalancer: HTTP Request
    LoadBalancer->>Router: Find matching route
    Router-->>LoadBalancer: Return matching route & pool
    
    LoadBalancer->>PolicyEngine: Apply policies
    alt Policy Violation
        PolicyEngine-->>LoadBalancer: Reject request
        LoadBalancer-->>Client: 403 Forbidden
    else Policies Passed
        PolicyEngine-->>LoadBalancer: Allow request
        LoadBalancer->>ServerPool: Get next backend
        ServerPool-->>LoadBalancer: Selected backend
        
        LoadBalancer->>Backend: Forward request
        Backend-->>LoadBalancer: HTTP Response
        LoadBalancer-->>Client: HTTP Response
    end
```

### Health Check Workflow

The health checker continuously monitors backend servers to ensure they are operational.

```mermaid
sequenceDiagram
    participant HealthChecker
    participant Probe
    participant Backend
    participant ServerPool
    
    loop Every interval
        HealthChecker->>Probe: Check backend health
        Probe->>Backend: Send health check request
        
        alt Backend is healthy
            Backend-->>Probe: 200 OK
            Probe-->>HealthChecker: Healthy
            HealthChecker->>ServerPool: Mark backend as healthy
        else Backend is unhealthy
            Backend-->>Probe: Error or timeout
            Probe-->>HealthChecker: Unhealthy
            HealthChecker->>ServerPool: Mark backend as unhealthy
        end
    end
```

### Configuration Reload Workflow

The load balancer supports runtime configuration updates through the admin API.

```mermaid
sequenceDiagram
    participant Admin
    participant AdminAPI
    participant ConfigManager
    participant Router
    participant ServerPool
    participant HealthChecker
    
    Admin->>AdminAPI: Update configuration
    AdminAPI->>ConfigManager: Parse and validate config
    
    alt Valid Configuration
        ConfigManager->>Router: Update routing rules
        ConfigManager->>ServerPool: Update backend pools
        ConfigManager->>HealthChecker: Update health check config
        ConfigManager-->>AdminAPI: Configuration applied
        AdminAPI-->>Admin: 200 OK
    else Invalid Configuration
        ConfigManager-->>AdminAPI: Validation error
        AdminAPI-->>Admin: 400 Bad Request
    end
```

## Use Cases

### Use Case 1: API Gateway

Using the load balancer as an API gateway to route requests to different microservices.

```mermaid
graph TD
    Client[Client] -->|HTTP Request| LB[Load Balancer]
    
    LB -->|/auth/*| AuthService[Auth Service]
    LB -->|/users/*| UserService[User Service]
    LB -->|/products/*| ProductService[Product Service]
    LB -->|/orders/*| OrderService[Order Service]
```

Configuration:

```yaml
backend_pools:
  - name: "auth-service"
    algorithm: "round_robin"
    backends:
      - url: "http://auth-service:8000"
  - name: "user-service"
    algorithm: "round_robin"
    backends:
      - url: "http://user-service:8000"
  - name: "product-service"
    algorithm: "round_robin"
    backends:
      - url: "http://product-service:8000"
  - name: "order-service"
    algorithm: "round_robin"
    backends:
      - url: "http://order-service:8000"

routing_rules:
  - match:
      path: "/auth/*"
    target_pool: "auth-service"
  - match:
      path: "/users/*"
    target_pool: "user-service"
  - match:
      path: "/products/*"
    target_pool: "product-service"
  - match:
      path: "/orders/*"
    target_pool: "order-service"
```

### Use Case 2: High Availability Web Service

Using the load balancer to distribute traffic across multiple web servers for high availability.

```mermaid
graph TD
    Client[Client] -->|HTTP Request| LB[Load Balancer]
    
    subgraph "Web Server Pool"
        Web1[Web Server 1]
        Web2[Web Server 2]
        Web3[Web Server 3]
    end
    
    LB -->|Round Robin| Web1
    LB -->|Round Robin| Web2
    LB -->|Round Robin| Web3
    
    HealthChecker[Health Checker] -.->|Monitor| Web1
    HealthChecker -.->|Monitor| Web2
    HealthChecker -.->|Monitor| Web3
```

Configuration:

```yaml
backend_pools:
  - name: "web-servers"
    algorithm: "round_robin"
    backends:
      - url: "http://web1:80"
      - url: "http://web2:80"
      - url: "http://web3:80"
    health_check:
      path: "/health"
      interval: "10s"
      timeout: "2s"

routing_rules:
  - match:
      path: "/*"
    target_pool: "web-servers"
```

### Use Case 3: A/B Testing

Using the load balancer to split traffic between different versions of a service for A/B testing.

```mermaid
graph TD
    Client[Client] -->|HTTP Request| LB[Load Balancer]
    
    subgraph "Version A (80%)"
        ServiceA1[Service A - Instance 1]
        ServiceA2[Service A - Instance 2]
    end
    
    subgraph "Version B (20%)"
        ServiceB[Service B]
    end
    
    LB -->|80% of traffic| ServiceA1
    LB -->|80% of traffic| ServiceA2
    LB -->|20% of traffic| ServiceB
```

Configuration:

```yaml
backend_pools:
  - name: "service-a"
    algorithm: "round_robin"
    backends:
      - url: "http://service-a-1:8000"
      - url: "http://service-a-2:8000"
  - name: "service-b"
    algorithm: "round_robin"
    backends:
      - url: "http://service-b:8000"

routing_rules:
  - match:
      headers:
        X-Test-Group: "B"
    target_pool: "service-b"
  - match:
      path: "/*"
    target_pool: "service-a"
```

### Use Case 4: Rate Limiting and Security

Using the load balancer to apply rate limiting and security policies to protect backend services.

```mermaid
sequenceDiagram
    participant Client
    participant LoadBalancer
    participant RateLimiter
    participant ACL
    participant Backend
    
    Client->>LoadBalancer: HTTP Request
    LoadBalancer->>RateLimiter: Check rate limit
    
    alt Rate limit exceeded
        RateLimiter-->>LoadBalancer: Reject request
        LoadBalancer-->>Client: 429 Too Many Requests
    else Rate limit ok
        RateLimiter-->>LoadBalancer: Allow request
        LoadBalancer->>ACL: Check access control
        
        alt Access denied
            ACL-->>LoadBalancer: Reject request
            LoadBalancer-->>Client: 403 Forbidden
        else Access allowed
            ACL-->>LoadBalancer: Allow request
            LoadBalancer->>Backend: Forward request
            Backend-->>LoadBalancer: HTTP Response
            LoadBalancer-->>Client: HTTP Response
        end
    end
```

Configuration:

```yaml
backend_pools:
  - name: "api-servers"
    algorithm: "least_conn"
    backends:
      - url: "http://api1:8000"
      - url: "http://api2:8000"

routing_rules:
  - match:
      path: "/api/*"
    target_pool: "api-servers"
    policies:
      - rate_limit: "100/minute"
      - acl: "allow:192.168.1.0/24,deny:10.0.0.1"
```

### Use Case 5: Header Transformation

Using the load balancer to transform headers for backend compatibility.

```mermaid
sequenceDiagram
    participant Client
    participant LoadBalancer
    participant HeaderTransformer
    participant Backend
    
    Client->>LoadBalancer: HTTP Request
    LoadBalancer->>HeaderTransformer: Transform headers
    HeaderTransformer-->>LoadBalancer: Modified request
    LoadBalancer->>Backend: Forward modified request
    Backend-->>LoadBalancer: HTTP Response
    LoadBalancer-->>Client: HTTP Response
```

Configuration:

```yaml
routing_rules:
  - match:
      path: "/api/*"
    target_pool: "api-servers"
    policies:
      - transform: "add-header:X-Forwarded-Host:api.example.com,remove-header:Referer"
```

## Advanced Workflows

### Blue-Green Deployment

Using the load balancer to implement blue-green deployments for zero-downtime updates.

```mermaid
graph TD
    Client[Client] -->|HTTP Request| LB[Load Balancer]
    
    subgraph "Blue Environment (Active)"
        Blue1[Blue Server 1]
        Blue2[Blue Server 2]
    end
    
    subgraph "Green Environment (Staging)"
        Green1[Green Server 1]
        Green2[Green Server 2]
    end
    
    LB -->|Production Traffic| Blue1
    LB -->|Production Traffic| Blue2
    LB -->|Test Traffic| Green1
    LB -->|Test Traffic| Green2
    
    Admin[Admin] -.->|Switch Traffic| LB
```

### Canary Deployment

Using the load balancer to implement canary deployments for gradual rollouts.

```mermaid
graph TD
    Client[Client] -->|HTTP Request| LB[Load Balancer]
    
    subgraph "Stable Version (90%)"
        Stable1[Stable Server 1]
        Stable2[Stable Server 2]
    end
    
    subgraph "Canary Version (10%)"
        Canary[Canary Server]
    end
    
    LB -->|90% of traffic| Stable1
    LB -->|90% of traffic| Stable2
    LB -->|10% of traffic| Canary
    
    Metrics[Metrics Collector] -.->|Monitor| Stable1
    Metrics -.->|Monitor| Stable2
    Metrics -.->|Monitor| Canary
    
    Admin[Admin] -.->|Adjust Traffic %| LB
```

### Circuit Breaking

Using the load balancer to implement circuit breaking for fault tolerance.

```mermaid
stateDiagram-v2
    [*] --> Closed
    Closed --> Open: Error threshold exceeded
    Open --> HalfOpen: After timeout
    HalfOpen --> Closed: Successful requests
    HalfOpen --> Open: Continued errors
```

```mermaid
sequenceDiagram
    participant Client
    participant LoadBalancer
    participant CircuitBreaker
    participant Backend
    
    Client->>LoadBalancer: HTTP Request
    LoadBalancer->>CircuitBreaker: Check circuit state
    
    alt Circuit Open
        CircuitBreaker-->>LoadBalancer: Circuit open
        LoadBalancer-->>Client: 503 Service Unavailable
    else Circuit Closed or Half-Open
        CircuitBreaker-->>LoadBalancer: Allow request
        LoadBalancer->>Backend: Forward request
        
        alt Backend Success
            Backend-->>LoadBalancer: HTTP Response
            LoadBalancer->>CircuitBreaker: Record success
            LoadBalancer-->>Client: HTTP Response
        else Backend Failure
            Backend-->>LoadBalancer: Error
            LoadBalancer->>CircuitBreaker: Record failure
            LoadBalancer-->>Client: 502 Bad Gateway
        end
    end
```
