# 🚀 Advanced Go Load Balancer

<div align="center">

[![Version](https://img.shields.io/badge/version-0.2.0-blue.svg)](https://github.com/rixtrayker/go-loadbalancer)
[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

**A high-performance, feature-rich HTTP/S load balancer built from scratch in Go**

*Educational • Modular • Production-Ready*

</div>

## 📑 Table of Contents

- [🚀 Advanced Go Load Balancer](#-advanced-go-load-balancer)
  - [📑 Table of Contents](#-table-of-contents)
  - [✨ Overview](#-overview)
  - [🎯 Features](#-features)
    - [🔄 **Load Balancing**](#-load-balancing)
    - [🏥 **Health Monitoring**](#-health-monitoring)
    - [🛣️ **Advanced Routing**](#️-advanced-routing)
    - [🛡️ **Policy Enforcement**](#️-policy-enforcement)
    - [🔧 **Additional Features**](#-additional-features)
  - [📁 Project Structure](#-project-structure)
  - [🚀 Quick Start](#-quick-start)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Running](#running)
  - [⚙️ Configuration](#️-configuration)
    - [Configuration Example](#configuration-example)
  - [🏗️ Architecture](#️-architecture)
    - [Design Principles](#design-principles)
    - [Core Components](#core-components)
  - [🔍 Observability](#-observability)
    - [Logging](#logging)
    - [Metrics](#metrics)
    - [Tracing](#tracing)
  - [🎓 Lessons Learned \& Skills Demonstrated](#-lessons-learned--skills-demonstrated)
  - [🤝 Contributing](#-contributing)
  - [📊 Roadmap](#-roadmap)
  - [📄 License](#-license)

---

## ✨ Overview

> **⚠️ Beta Release:** This project is currently in beta and actively under development. Features and APIs may evolve.

This project showcases an advanced HTTP/S load balancer implementation in Go, designed as both an educational resource and a practical networking infrastructure component. Built with modularity and extensibility in mind, it demonstrates core concepts of distributed systems, load balancing, and network programming.

## 🎯 Features

<table>
<tr>
<td width="50%">

### 🔄 **Load Balancing**
- **Multiple Algorithms**: Round Robin, Least Connections, Weighted
- **Smart Distribution**: Intelligent traffic routing
- **Backend Pool Management**: Organized backend grouping

### 🏥 **Health Monitoring**
- **Active Health Checks**: Continuous backend monitoring
- **Pluggable Probes**: HTTP, TCP health check support
- **Automatic Failover**: Seamless backend switching

</td>
<td width="50%">

### 🛣️ **Advanced Routing**
- **Multi-Factor Routing**: Host, path, method, headers
- **Flexible Rules**: Complex routing configurations
- **Dynamic Updates**: Runtime route modifications

### 🛡️ **Policy Enforcement**
- **Rate Limiting**: Request throttling and control
- **Header Transformation**: Request/response modification
- **IP Access Control**: Whitelist/blacklist support

</td>
</tr>
</table>

### 🔧 **Additional Features**
- **HTTP/S Reverse Proxying** with high performance
- **Graceful Shutdown** for clean terminations
- **Observability Integration** with logging, metrics, and tracing
- **Modular Architecture** for easy extension

---

## 📁 Project Structure

The project follows the standard Go project layout:

```
go-loadbalancer/
├── 📄 cmd/go-lb/                # Application entry point
│   └── main.go                  # Main application file
├── 🌐 api/                      # API definitions
│   └── http/v1/                 # HTTP API version 1
├── ⚙️  configs/                  # Configuration files and templates
│   ├── config.go                # Configuration structures
│   └── loader.go                # Configuration loading logic
├── 🔒 internal/                  # Internal business logic
│   ├── app/                     # Application initialization
│   ├── admin/                   # Admin interface
│   ├── backend/                 # Backend server management
│   ├── healthcheck/             # Health checking system
│   ├── logging/                 # Structured logging
│   ├── middleware/              # HTTP middleware
│   ├── monitoring/              # Metrics and monitoring
│   ├── serverpool/              # Backend pools & algorithms
│   ├── routing/                 # Request routing engine
│   ├── policy/                  # Policy enforcement
│   ├── tracing/                 # Distributed tracing
│   └── handler/                 # Core request handlers
├── 📦 pkg/                      # Reusable utilities
├── 🧪 test/                     # Additional test applications
├── 🔧 scripts/                  # Scripts for various tasks
├── 🚢 deployments/              # Deployment configurations
│   └── docker/                  # Docker-related files
├── 📊 tools/                    # Tools and utilities
│   └── k6/                      # K6 load testing scripts
└── 📚 docs/                     # Documentation files
```

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.23+** installed on your system
- Backend services to load balance (for testing)

### Installation

```bash
# Clone the repository
git clone https://github.com/rixtrayker/go-loadbalancer.git

# Navigate to project directory
cd go-loadbalancer

# Build the project
go build -o build/go-lb ./cmd/go-lb
```

### Running

```bash
# Start with default configuration
./build/go-lb

# Or specify a config file
./build/go-lb --config configs/config.yml
```

---

## ⚙️ Configuration

The load balancer uses a flexible configuration system supporting:

- **🎯 Listener Settings**: Address, port, TLS configuration
- **🏊 Backend Pools**: Server groups with load balancing algorithms
- **💓 Health Checks**: Monitoring intervals, timeouts, and probe types
- **🛣️ Routing Rules**: Complex request matching and forwarding
- **📋 Policies**: Rate limiting, transformations, and access control
- **📊 Monitoring**: Logging, metrics, and tracing configuration

### Configuration Example

```yaml
server:
  address: ":8080"
  
backend_pools:
  - name: "web-servers"
    algorithm: "round_robin"
    backends:
      - url: "http://localhost:3001"
        weight: 1
      - url: "http://localhost:3002"
        weight: 2
    health_check:
      path: "/health"
      interval: "30s"
      timeout: "5s"

routing_rules:
  - match:
      host: "example.com"
      path: "/api/*"
    target_pool: "web-servers"
    policies:
      - rate_limit: "100/minute"

monitoring:
  prometheus:
    enabled: true
    path: "/metrics"
    port: 9090
  tracing:
    enabled: true
    service_name: "go-loadbalancer"
    endpoint: "localhost:4317"
```

---

## 🏗️ Architecture

### Design Principles

- **🔧 Modularity**: Clean separation of concerns
- **🔌 Extensibility**: Plugin-based architecture
- **⚡ Performance**: Optimized for high throughput
- **🛡️ Reliability**: Robust error handling and recovery

### Core Components

| Component | Purpose |
|-----------|---------|
| **Router** | Intelligent request routing and matching |
| **Server Pool** | Backend management and load balancing |
| **Health Checker** | Continuous backend monitoring |
| **Policy Engine** | Request/response transformation and control |
| **Admin Interface** | Runtime configuration and monitoring |

---

## 🔍 Observability

The load balancer includes comprehensive observability features:

### Logging

- Structured JSON logging with zap
- Configurable log levels and formats
- Trace context correlation

### Metrics

- Prometheus metrics for all components
- Request/response metrics
- Backend health and performance metrics
- System resource usage

### Tracing

- OpenTelemetry integration
- Distributed tracing support
- Trace context propagation
- Span attributes for detailed analysis

---

## 🎓 Lessons Learned & Skills Demonstrated

This project serves as an excellent learning resource, demonstrating several important concepts and best practices in Go development:

- Standard Go project layout
- Modular architecture with clean interfaces
- Effective use of middleware patterns
- Comprehensive observability implementation
- Advanced configuration management
- Multiple load balancing algorithms
- Policy-based request handling

---

## 🤝 Contributing

We welcome contributions! See the [Contributing Guide](docs/contributing.md) for more information.

---

## 📊 Roadmap

- [ ] **gRPC Load Balancing**: Support for gRPC protocols
- [ ] **Docker Integration**: Containerized deployment
- [ ] **Kubernetes Support**: Native K8s integration
- [ ] **WebSocket Proxying**: Real-time connection support
- [ ] **Circuit Breaker**: Fault tolerance patterns

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ by [rixtrayker](https://github.com/rixtrayker)**

⭐ **Star this repo if you find it helpful!** ⭐

</div>
