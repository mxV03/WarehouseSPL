# Products – Warehouse Management System (SPL)

This document explains the **product concept** of the Warehouse Management System Software Product Line.

It describes:
- Why products exist
- How products are derived from features
- Which products are currently defined
- How products relate to build tags and constraints

---

## 1. Why Products Exist in This SPL

In this project, **products are conceptual SPL artifacts**.

They serve to:
- Document **valid feature combinations**
- Make **product derivation explicit**
- Provide **named variants** for demonstration and evaluation
- Explain the impact of variability on system behavior

Products are **not** implemented via Clone-and-Own.
Instead, all products are derived from the **same code base** using build tags.

---

## 2. Product Derivation Mechanism

Products are derived by enabling a **set of features** at compile time.

Technically:
- Each feature corresponds to a Go build tag
- A product is defined by a specific tag combination
- Disabled features do not exist in the compiled binary

Example:

```bash
go run -tags "min reporting multiwarehouse" ./cmd/app
```

This derives a product that includes:
- Core
- CLI
- Reporting
- MultiWarehouse

---

## 3. Feature Overview (for Product Mapping)

### Mandatory Features
- **Core** – inventory, locations, stock movements
- **CLI** – modular command-line interface

These features are always present.

### Optional Features
- **AuditLog** – system-wide audit events
- **Authentication** – users, roles, access control
- **Barcode** - print, scan (barcode)
- **Logistics** – zones, bins, item placement
- **MultiWarehouse** – grouping locations into warehouses
- **Notification** - send notifications
- **Picking** - picklist operations
- **Reporting** – inventory & movement reports
- **Tracking** – shipment tracking entity


---

## 4. Defined Products

### 4.1 MIN Product

**Purpose:**
Minimal, core-only product.

**Included Features:**
- Core
- CLI

**Build Tags:**
```text
min
```

**Use Case:**
- Basic inventory management
- Teaching the core domain
- Baseline SPL product

---

### 4.2 PRO Product

**Purpose:**
Adds new buisness-opertion-features and analytical capabilities.

**Included Features:**
- Core
- CLI
- Reporting
- Logistics
- Picking
- Barcode

**Build Tags:**
```text
pro reporting logistics picking barcode
```

**Use Case:**
- Inventory analysis
- New operation features
- Management overview
- Demonstrates optional feature inclusion

---

### 4.3 ENTERPRISE Product

**Purpose:**
Full-featured enterprise-grade variant.

**Included Features:**
- Core
- CLI
- AuditLog
- Authentification
- Barcode
- Logistics
- Multiwarehouse
- Notifications
- Picking
- Reporting
- Tracking

**Build Tags:**
```text
ent audit auth barcode logistics multiwarehouse notifications picking reporting tracking 
```

**Use Case:**
- Large warehouse networks
- Auditable operations
- Authentification 
- Multi-site reporting

---

## 5. Cross-Tree Constraints and Products

Some products are only valid if **constraints are satisfied**.

Examples:

- **MultiWarehouse ⇒ Reporting**
  - Warehouse-wide KPIs require reporting functionality

- **Authentication ⇒ AuditLog** (conceptual)
  - Security-relevant actions must be auditable

Invalid products:

```bash
go run -tags "min multiwarehouse" ./cmd/app
```

→ fails at compile time

---

## 6. Notes for Extension

- Additional products can be defined without code duplication
- New features automatically become available for product composition
- Products remain documentation artifacts, not forks


