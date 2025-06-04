# Architecture Overview

The Go Load Balancer is designed with modularity and extensibility in mind. This document outlines the high-level architecture and component interactions.

## System Architecture

The load balancer follows a layered architecture pattern, with clear separation of concerns between components. This design allows for easy extension and modification of individual components without affecting the rest of the system.

```mermaid
graph TD
    Client[Client] -->|HTTP Request| LB[Load Balancer]
    LB -->|Route Request| Router[Router]
    Router -->|Apply Policies| PolicyEngine[Policy Engine]
    PolicyEngine -->|Select Backend| ServerPool[Server Pool]
    ServerPool -->|Forward Request| Backend1[Backend Server 1]
    ServerPool -->|Forward Request| Backend2[Backend Server 2]
    ServerPool -->|Forward Request| BackendN[Backend Server N]
    
    HealthChecker[Health Checker] -.->|Monitor| Backend1
    HealthChecker -.->|Monitor| Backend2
    HealthChecker -.->|Monitor| BackendN
    
    Admin[Admin API] -.->|Configure| LB
    Metrics[Metrics Collector] -.->|Monitor| LB
```

## Core Components

### HTTP Handler

The HTTP handler is the entry point for all incoming requests. It is responsible for:

- Accepting incoming HTTP connections
- Parsing HTTP requests
- Routing requests to the appropriate backend
- Handling errors and returning appropriate responses

```mermaid
classDiagram
    class Handler {
        +router *mux.Router
        +config *Config
        +logger *Logger
        +metrics *Metrics
        +ServeHTTP(w ResponseWriter, r *Request)
        +setupRoutes()
    }
    
    Handler --> Router
    Handler --> Config
    Handler --> Logger
    Handler --> Metrics
```

### Router

The router determines which backend pool should handle a request based on configurable rules. It supports matching on:

- Host
- Path
- HTTP method
- Headers

```mermaid
classDiagram
    class Router {
        +rules []*Rule
        +pools map[string]*Pool
        +logger *Logger
        +Route(req *Request) *Pool
        +AddRule(rule *Rule)
        +RemoveRule(index int)
    }
    
    class Rule {
        +HostPattern *regexp.Regexp
        +PathPattern *regexp.Regexp
        +Method string
        +HeaderRules map[string]*regexp.Regexp
        +TargetPool string
        +Policies []PolicyConfig
        +Matches(req *Request) bool
    }
    
    Router --> Rule
    Router --> Pool
```

### Server Pool

The server pool manages a group of backend servers and implements load balancing algorithms to distribute requests among them.

```mermaid
classDiagram
    class Pool {
        +Name string
        +Backends []*Backend
        +Algorithm Algorithm
        +NextBackend(r *Request) *Backend
        +MarkBackendStatus(url string, healthy bool)
    }
    
    class Backend {
        +URL *url.URL
        +Weight int
        +Healthy bool
        +ActiveConns int32
        +TotalRequests int64
        +IsHealthy() bool
        +SetHealth(healthy bool)
        +IncrementConnections()
        +DecrementConnections()
        +GetActiveConnections() int
        +IncrementRequests()
        +GetTotalRequests() int64
    }
    
    class Algorithm {
        <<interface>>
        +NextBackend(r *Request) *Backend
    }
    
    Pool --> Backend
    Pool --> Algorithm
```

### Health Checker

The health checker continuously monitors the health of backend servers and updates their status in the server pool.

```mermaid
classDiagram
    class HealthChecker {
        +pools map[string]*Pool
        +configs map[string]HealthCheckConfig
        +probes map[string]Probe
        +logger *Logger
        +Start(ctx context)
        +Stop()
        -checkBackend(ctx context, pool *Pool, backendURL string, interval time.Duration)
    }
    
    class Probe {
        <<interface>>
        +Check() bool
    }
    
    class HTTPProbe {
        +url *url.URL
        +path string
        +method string
        +timeout time.Duration
        +client *http.Client
        +Check() bool
    }
    
    class TCPProbe {
        +url *url.URL
        +timeout time.Duration
        +Check() bool
    }
    
    HealthChecker --> Probe
    Probe <|.. HTTPProbe
    Probe <|.. TCPProbe
```

### Policy Engine

The policy engine applies configurable policies to requests, such as rate limiting, access control, and header transformation.

```mermaid
classDiagram
    class Policy {
        +Apply(policy PolicyConfig, r *Request) error
    }
    
    class RateLimiter {
        +limits map[string]*limit
        +mutex sync.RWMutex
        +cleanupInt time.Duration
        +Allow(key string, rate int, per time.Duration) error
        -cleanup()
    }
    
    class ACL {
        +allowList []string
        +denyList []string
        +mutex sync.RWMutex
        +Check(ipStr string) bool
        +AddAllowRule(cidr string)
        +AddDenyRule(cidr string)
    }
    
    class HeaderTransformer {
        +AddHeaders map[string]string
        +RemoveHeaders []string
        +Apply(r *Request)
    }
    
    Policy --> RateLimiter
    Policy --> ACL
    Policy --> HeaderTransformer
```

## Data Flow

The following diagram illustrates the data flow through the system when processing a request:

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Router
    participant PolicyEngine
    participant ServerPool
    participant Backend
    
    Client->>Handler: HTTP Request
    Handler->>Router: Route Request
    Router-->>Handler: Selected Pool
    Handler->>PolicyEngine: Apply Policies
    alt Policy Violation
        PolicyEngine-->>Handler: Error
        Handler-->>Client: 403 Forbidden
    else Policies Passed
        PolicyEngine-->>Handler: Success
        Handler->>ServerPool: Get Backend
        ServerPool-->>Handler: Selected Backend
        Handler->>Backend: Forward Request
        Backend-->>Handler: Response
        Handler-->>Client: HTTP Response
    end
```

## Component Interactions

The following diagram shows how the different components interact with each other:

```mermaid
graph TD
    App[Application] -->|Initializes| Handler[HTTP Handler]
    App -->|Initializes| HealthChecker[Health Checker]
    App -->|Initializes| AdminAPI[Admin API]
    
    Handler -->|Uses| Router[Router]
    Handler -->|Uses| PolicyEngine[Policy Engine]
    Handler -->|Uses| ServerPool[Server Pool]
    
    HealthChecker -->|Updates| ServerPool
    AdminAPI -->|Configures| ServerPool
    AdminAPI -->|Configures| Router
    
    ServerPool -->|Selects| Backend[Backend Servers]
    
    Metrics[Metrics Collector] -.->|Monitors| Handler
    Metrics -.->|Monitors| ServerPool
    Metrics -.->|Monitors| Backend
    
    Logger[Logger] -.->|Logs from| App
    Logger -.->|Logs from| Handler
    Logger -.->|Logs from| HealthChecker
    Logger -.->|Logs from| ServerPool
```

## Deployment Architecture

The load balancer can be deployed in various configurations, depending on the requirements:

```mermaid
graph TD
    subgraph "Load Balancer Tier"
        LB1[Load Balancer 1]
        LB2[Load Balancer 2]
    end
    
    subgraph "Backend Tier"
        BE1[Backend 1]
        BE2[Backend 2]
        BE3[Backend 3]
    end
    
    Client1[Client] -->|HTTP| LB1
    Client2[Client] -->|HTTP| LB2
    
    LB1 -->|HTTP| BE1
    LB1 -->|HTTP| BE2
    LB1 -->|HTTP| BE3
    
    LB2 -->|HTTP| BE1
    LB2 -->|HTTP| BE2
    LB2 -->|HTTP| BE3
```

## Extension Points

The architecture is designed to be extensible in several ways:

1. **Load Balancing Algorithms**: New algorithms can be added by implementing the Algorithm interface
2. **Health Check Probes**: Additional health check types can be added
3. **Routing Rules**: The routing engine can be extended with new matching criteria
4. **Policies**: New policy types can be implemented and integrated
5. **Metrics and Monitoring**: Additional observability integrations can be added
