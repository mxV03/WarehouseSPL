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

### Modular Feature Structure

- Each feature is implemented as its **own Go package**
- The **Core** is always present
- Optional features are only compiled when their build tag is enabled

```
internal/
  features/
    reporting/
    logistics/
    tracking/
    audit/
    auth/
    multiwarehouse/
  core/
```

---

## Core Features (Always Active)

- Item management (SKU, name, description)
- Location management
- Stock movements (IN / OUT)
- Order handling (inbound / outbound)
- (Modular CLI)

---

## Optional Features

- **Reporting** – inventory and movement reports
- **Logistics Management** – zones, bins, item placement
- **Tracking** – shipment tracking as a separate entity
- **AuditLog** – system-wide audit events (optional, compile-time)
- **Authentication** – users, roles, CLI guards
- **MultiWarehouse** – grouping of locations into warehouses

Each feature:
- Exists in its own package
- Registers its own CLI commands
- Can be fully removed at compile time

---

## Compile-Time Variability (Build Tags)

Examples:

```bash
go run -tags "min" ./cmd/app
go run -tags "min reporting" ./cmd/app
go run -tags "min reporting multiwarehouse" ./cmd/app
```

If a feature tag is missing:
- Its code is not compiled
- Its CLI commands do not exist
- It does not affect the binary

---

## Runtime Variability

Some features use **runtime configuration** via environment variables:

- `WMS_ACTOR` – audit actor
- `WMS_USER`, `WMS_PASS` – authentication

Runtime variability is used **only where compile-time variability is not appropriate**.

---

## Cross-Tree Constraints

Invalid feature combinations are prevented using **compile-time guards**.

Examples:

- `multiwarehouse ⇒ reporting`
- `audit ⇒ cli`
- `auth ⇒ audit` (conceptual dependency)

Invalid combinations **fail at compile time**.

---

## Products (Conceptual)

Products are **conceptual SPL artifacts** that document valid feature combinations.
They are not Clone-and-Own variants.

Examples:

| Product        | Enabled Features |
|---------------|------------------|
| MIN           | Core + CLI |
| PRO           | Core + CLI + Reporting |
| ENTERPRISE    | All features |

Products are displayed at startup and documented separately.

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

