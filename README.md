# 🚀 Advanced Go Load Balancer

<div align="center">

[![Version](https://img.shields.io/badge/version-0.1.0--beta-blue.svg)](https://github.com/rixtrayker/go-loadbalancer)
[![Go Version](https://img.shields.io/badge/go-1.18+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

**A high-performance, feature-rich HTTP/S load balancer built from scratch in Go**

*Educational • Modular • Production-Ready*

</div>

## 📑 Table of Contents

- [✨ Overview](#-overview)
- [🎯 Features](#-features)
- [📁 Project Structure](#-project-structure)
- [🚀 Quick Start](#-quick-start)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running](#running)
- [⚙️ Configuration](#️-configuration)
- [🏗️ Architecture](#️-architecture)
- [🎓 Lessons Learned & Skills](#-lessons-learned--skills-demonstrated)
- [🤝 Contributing](#-contributing)
- [📊 Roadmap](#-roadmap)

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
│   ├── serverpool/              # Backend pools & algorithms
│   ├── routing/                 # Request routing engine
│   ├── policy/                  # Policy enforcement
│   └── handler/                 # Core request handlers
├── 📦 pkg/                      # Reusable utilities
│   ├── logging/                 # Structured logging
│   ├── metrics/                 # Performance metrics
│   └── tracer/                  # Distributed tracing
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

- **Go 1.18+** installed on your system
- Backend services to load balance (for testing)

### Installation

```bash
# Clone the repository
git clone https://github.com/rixtrayker/go-loadbalancer.git

# Navigate to project directory
cd go-loadbalancer

# Build the project
go build -o go-lb .
```

### Configuration

Create a `config.yaml` file with your configuration:

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
```

### Running

```bash
# Start with default configuration
./go-lb

# Or specify a config file
./go-lb --config config.yaml
```

---

## ⚙️ Configuration

The load balancer uses a flexible configuration system supporting:

- **🎯 Listener Settings**: Address, port, TLS configuration
- **🏊 Backend Pools**: Server groups with load balancing algorithms
- **💓 Health Checks**: Monitoring intervals, timeouts, and probe types
- **🛣️ Routing Rules**: Complex request matching and forwarding
- **📋 Policies**: Rate limiting, transformations, and access control

### Configuration Sources

- YAML files
- Environment variables
- Command-line flags
- Runtime API updates

For detailed configuration options, see the [Configuration Guide](docs/configuration.md).

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

## 🎓 Lessons Learned & Skills Demonstrated

This project serves as an excellent learning resource, demonstrating several important concepts and best practices in Go development:

### 🏗️ Code Organization & Structure
- Standard Go project layout with clear separation of concerns
- Modular architecture with well-defined package boundaries
- Clean code principles and SOLID design patterns
- Effective use of Go's package system

### 🔧 Configuration Management
- YAML-based configuration with environment variable overrides
- Structured data types for configuration options
- Sensible default values and validation
- Runtime configuration updates

### ⚡ Load Balancing Implementation
- Multiple load balancing algorithms (Round Robin, Least Connections, Weighted)
- Thread-safe counter manipulation using atomic operations
- Efficient backend health checking and management
- Connection tracking and management

### 🌐 HTTP Handling
- Efficient request proxying with `httputil.NewSingleHostReverseProxy`
- Proper request/response header management
- Context-based request cancellation and timeout handling
- Graceful server shutdown implementation

### 🔒 Concurrency & Thread Safety
- Effective use of `sync.RWMutex` for concurrent operations
- Atomic operations for counter updates
- Proper locking mechanisms for shared resources
- Thread-safe backend selection algorithms

### 🛡️ Policy Implementation
- Rate limiting using token bucket algorithm
- IP-based access control (ACL)
- Header transformation policies
- Policy chain implementation

### 🏥 Health Checking
- HTTP and TCP health check implementations
- Context-based health check cancellation
- Periodic health check scheduling
- Health check result management

### 🛣️ Routing
- Regular expression-based pattern matching
- Host, path, and header-based routing
- Efficient route lookup using maps
- Dynamic route configuration

### 📊 Logging & Monitoring
- Structured logging with logrus
- Multiple log levels implementation
- Prometheus metrics integration
- Grafana dashboard setup

### 🐳 Docker & Deployment
- Multi-stage Docker builds
- Docker health check implementation
- Security-focused container configuration
- Environment variable management

### 🛠️ Development Tools
- Makefile for common development tasks
- Build and test automation
- Linting with golangci-lint
- Proper cleanup procedures

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

### Getting Started
1. 🍴 Fork the repository
2. 🌿 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. 💾 Commit your changes (`git commit -m 'Add amazing feature'`)
4. 📤 Push to the branch (`git push origin feature/amazing-feature`)
5. 🔄 Open a Pull Request

### Areas for Contribution
- 🐛 **Bug Reports**: Help us identify and fix issues
- ✨ **Feature Requests**: Suggest new functionality
- 📝 **Documentation**: Improve guides and examples
- 🧪 **Testing**: Add test cases and scenarios
- 🎨 **UI/UX**: Enhance the admin interface

---

## 📊 Roadmap

- [ ] **gRPC Load Balancing**: Support for gRPC protocols
- [ ] **Docker Integration**: Containerized deployment
- [ ] **Kubernetes Support**: Native K8s integration
- [ ] **WebSocket Proxying**: Real-time connection support
- [ ] **Advanced Metrics**: Prometheus integration
- [ ] **Circuit Breaker**: Fault tolerance patterns

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### What does this mean?

The MIT License is a permissive license that is short and to the point. It lets people do anything they want with your code as long as they provide attribution back to you and don't hold you liable.

**You are free to:**
- ✅ Use this code commercially
- ✅ Modify the code
- ✅ Distribute the code
- ✅ Use it privately
- ✅ Sublicense it

**Under the following conditions:**
- ℹ️ Include the original copyright notice
- ℹ️ Include the license text

**No liability:**
- 🛡️ The software is provided "as is", without warranty of any kind

For more information, please refer to the [LICENSE](LICENSE) file in this repository.

---

## 🙏 Acknowledgments

- **Go Community**: For excellent standard libraries and ecosystem
- **net/http/httputil**: Core reverse proxy functionality
- **Open Source Contributors**: For inspiration and best practices

---

<div align="center">

**Made with ❤️ by [rixtrayker](https://github.com/rixtrayker)**

⭐ **Star this repo if you find it helpful!** ⭐

</div>