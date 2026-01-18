# Warehouse Management System – Software Product Line (SPL)

## Overview
This project is a **Software Product Line (SPL)** implementation of a **Warehouse Management System (WMS)** written in **Go**.

The goal of the project is **not** to build a full-blown WMS, but to demonstrate SPL concepts in a **clear, realistic, and technically sound way**:

- Feature modeling
- Variability (compile-time & runtime)
- Product derivation
- Cross-tree constraints
- Conscious selection of implementation techniques

The system is **CLI-first**, modular, and intentionally kept compact.

---

## SPL Core Idea

- **One code base**
- **Multiple products (variants)**
- Variability is implemented via:
  - Modular feature packages
  - **Compile-time variability** using Go build tags
  - **Runtime variability** via configuration (environment variables)
  - **Compile-time guards** for cross-tree constraints

---

## Technology Stack

- **Language:** Go
- **Database:** SQLite
- **ORM:** ent (entgo)
- **CLI:** Custom modular CLI (registry-based)
- **Build:** `go run`, `go build` with Go build tags
- **SQLite Driver:** `modernc.org/sqlite` (CGO-free, Windows compatible)

---

## Architecture Overview

### Modular Feature Structure (Component-Based Implementation)

- Each feature is implemented as an **independent Go package**
- Variability is realized through **compile-time composition** (build tags)
- The **core system is unaware of optional features**
- The primary extension point is the **CLI registry**
  - Each feature provides a `register.go`
  - Features register their own CLI commands
- No feature logic is hard-coded into the core CLI
- Additional Go interfaces are intentionally kept to a minimum
- Modularization is achieved through **package boundaries and explicit registration**, not through Clone-and-Own or runtime feature flags

```
internal/
  features/
    audit/
    auth/
    barcode/
    interfaces/
    logistics/
    multiwarehouse/
    notifications/
    picking/
    reporting/
    tracking/
  core/
    inventory/
    ordermanagement/
```

---

## Core Features (Always Active)

- **Item Management**  
  Management of items including SKU, name, and description

- **Location Management**  
  Creation and management of storage locations

- **Stock Movements**  
  Support for goods receipt (IN) and goods issue (OUT)

- **Order Handling**  
  Processing of inbound and outbound orders

- **Modular CLI**  
  Originally planned as an optional feature, the modular command-line interface is considered a core feature due to time constraints

---

## Optional Features

- **Authentication**
  - User login
  - Role-based access control
  - Access guards for CLI, API, and WebUI
- **Barcode**
  - Barcode generation
  - Barcode scanning
  - Integration with inventory and picking workflows
- **Logistics Management**
  - Definition of storage locations and zones
  - User-defined item placement
  - Inventory overview
  - Inventory planning
  - Stock warnings
  - Reservations
- **Picking**
  - Manual creation of pick lists
  - Pick list monitoring
  - Optional scanner support
- **Tracking**
  - Shipment and delivery tracking
  - External tracker integration
- **Reporting**
  - Inventory reports (current stock levels)
  - Movement reports (inbound and outbound)
  - KPI dashboards
- **Notifications**
  - Configurable notifications for users and staff
  - Rule-based notification handling
- **Automation**
  - Automatic replenishment
  - Logistics advisor
  - Robot advisor
  - Rule-based automation logic
- **AuditLog**
  - System-wide audit event logging
  - Change and operation tracking
  - Removable at compile time
- **MultiWarehouse**
  - Management of multiple warehouses
  - Grouping of locations into warehouses
  - Cross-warehouse data access and reporting
- **Interfaces**
  - **CLI** – command-line interface (as discribed under Core Features)


Each feature:
- Exists in its own package
- Registers its own CLI commands
- Can be fully removed at compile time

---

## Compile-Time Variability (Build Tags)

Examples:

```bash
go run ./cmd/app
go run -tags "reporting" ./cmd/app
go run -tags "reporting multiwarehouse" ./cmd/app
```

If a feature tag is missing:
- Its code is not compiled
- Its CLI commands do not exist
- It does not affect the binary
- Core Features are always part of CLI-Commands.

---

## Runtime Variability

Runtime variability is used **only for contextual and user-specific information** that cannot be decided at compile time.  
No features are enabled or disabled at runtime.

### Used Environment Variables

- **Authentication**
  - `WMS_USER` – username for CLI authentication
  - `WMS_PASS` – password for CLI authentication
  - Used to authenticate the current user without rebuilding the product

- **AuditLog**
  - `WMS_ACTOR` – identifies the actor responsible for an action
  - Used to annotate audit log entries with runtime context

### Scope of Runtime Variability

- Runtime variability is limited to:
  - user identity
  - authentication credentials
  - audit context
- All structural variability (feature selection) is handled at **compile time**


Some features use **runtime configuration** via environment variables.  
Runtime variability is applied **only where compile-time variability is not appropriate**, e.g. for user-specific or environment-specific behavior.

---

## Cross-Tree Constraints

Invalid feature combinations are prevented using **compile-time guards**.

Examples:

- `multiwarehouse ⇒ reporting`
- `notifications ⇒ reporting`
- `audit ⇒ auth`

Invalid combinations **fail at compile time**.

---

## Products (Conceptual)

Products are **purely conceptual SPL artifacts** used to describe and document **valid feature combinations** within the product line.  
They do **not** represent separate codebases, binaries, or Clone-and-Own variants.

Products exist **only as documentation and runtime output**:
- They are printed at application startup (e.g. in the CLI)
- They provide transparency about which feature set a build represents
- They help illustrate the product line structure

**Products do not control variability** and do not influence the build process at runtime.  
The **actual product derivation is performed exclusively via compile-time build tags**.

Examples:

| Product        | Enabled Features |
|---------------|------------------|
| MIN           | Core + CLI |
| PRO           | Core + CLI + Reporting + Logistics + Picking + Barcode|
| ENTERPRISE    | All features |

In summary:
- Products are **conceptual representations**, not concrete software artifacts
- Feature selection and product generation are handled entirely at **compile time**
- Products serve as a **documentation and presentation mechanism** only


---

## Running the System

### Basic Example

```bash
go run -tags "min" ./cmd/app --help
```

### Example with Multiple Features

```bash
go run -tags "min reporting multiwarehouse" ./cmd/app reporting.warehouse.summary W1
```

---

## Why This Is a Software Product Line

This project demonstrates SPL principles clearly:

- One code base, many products
- Explicit feature modularization
- Compile-time and runtime variability
- Enforced cross-tree constraints
- Reusable and extensible architecture

---

## Intended Use

This project is intended for:

- Academic SPL demonstration
- Teaching feature-oriented design in Go
- Demonstrating variability with build tags

It is **not** intended as a production-ready WMS.

---

## License

Academic / Educational Use

