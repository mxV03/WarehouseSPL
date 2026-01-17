# CLI Demo – Warehouse Management System (SPL)

This document shows a **complete, end-to-end CLI demo** of the Warehouse Management System Software Product Line.

The demo is designed to:
- Be executable step by step
- Demonstrate **feature variability**
- Demonstrate **cross-tree constraints**
- Show realistic WMS usage

All examples use `go run` to make **compile-time variability visible**.

---

## Preconditions

- Go installed
- Project root directory
- No existing database (or start with a clean DB)

---

## 1. Minimal Product (MIN)

### Build & Inspect

```bash
go run -tags " ./cmd/app --help
```

Expected:
- Core commands only
- No optional feature commands

---

## 2. Core Workflow Demo

### 2.1 Create Locations

```bash
go run ./cmd/app location.add L1 "Main Location"
go run ./cmd/app location.add L2 "Secondary Location"
```

### 2.2 Create Item

```bash
go run ./cmd/app item.add SKU1 "Packing Tape" "Standard tape"
```

### 2.3 Stock In / Out

```bash
go run ./cmd/app stock.in  SKU1 L1 10 ref-in-1
go run ./cmd/app stock.out SKU1 L1 3  ref-out-1
```

---

## 3. Feature Variability Demo (Reporting)

### Enable Reporting

```bash
go run -tags "reporting" ./cmd/app --help
```

Expected:
- Reporting commands are visible

### Inventory Report

```bash
go run -tags "reporting" ./cmd/app reporting.inventory
```

---

## 4. MultiWarehouse + Reporting (Constraint Demo)

### Invalid Combination (Compile-Time Error)

```bash
go run -tags "multiwarehouse" ./cmd/app
```

Expected:
- Compile-time error
- Explanation: `multiwarehouse` requires `reporting`

---

### Valid Combination

```bash
go run -tags "reporting multiwarehouse" ./cmd/app --help
```

### Create Warehouse

```bash
go run -tags "reporting multiwarehouse" ./cmd/app warehouse.add W1 "Warehouse One"
```

### Assign Locations

```bash
go run -tags "reporting multiwarehouse" ./cmd/app warehouse.location.assign W1 L1
go run -tags "reporting multiwarehouse" ./cmd/app warehouse.location.assign W1 L2
```

### Warehouse Summary Report

```bash
go run -tags "reporting multiwarehouse" ./cmd/app reporting.warehouse.summary W1
```

Expected:
- Number of locations
- Stock movement totals

---

## 5. Logistics Feature Demo

### Enable Logistics

```bash
go run -tags "logistics" ./cmd/app --help
```

### Create Zones and Bins

```bash
go run -tags "logistics" ./cmd/app logistics.zone.add L1 Z1 "Zone A"
go run -tags "logistics" ./cmd/app logistics.bin.add  L1 Z1 B1 "Bin 1"
```

### Assign Item to Bin

```bash
go run -tags "logistics" ./cmd/app logistics.bin.assign L1 B1 SKU1
```

---

## 6. Tracking Feature Demo

### Enable Tracking

```bash
go run -tags "tracking" ./cmd/app --help
```

### Set and Query Tracking

```bash
go run -tags "tracking" ./cmd/app tracking.set OUT-1 TRK123 https://track.example/TRK123 DHL
go run -tags "tracking" ./cmd/app tracking.get OUT-1
```

---

## 7. Audit Feature Demo

### Enable Audit

```bash
$env:WMS_ACTOR="demo-user"
go run -tags "audit auth" ./cmd/app audit.tail
```

### Trigger Audit Events

```bash
go run -tags "audit auth" ./cmd/app location.add L3 "Audit Test"
go run -tags "audit auth" ./cmd/app audit.list
```

---

## 8. Authentication Demo

### Enable Auth

```bash
go run -tags "auth" ./cmd/app --help
```


### Login via Environment

```bash
$env:WMS_USER="worker"
$env:WMS_PASS="worker1234"

go run -tags "auth" ./cmd/app auth.whoami
```

Hint:
- Before this step, at least **one Admin user must already exist**.
- User creation is **restricted to Admin users only**.
- If no Admin exists yet, you must **temporarily comment out the `RequireRole(Admin)` check**
  in the `auth.user.add` CLI command to bootstrap the first Admin account.

---

## 9. SPL Summary

This demo shows:

- Compile-time feature inclusion
- Runtime configuration
- Feature-dependent CLI commands
- Enforced cross-tree constraints
- Multiple valid product variants

The system behavior changes **at compile time**, not via configuration flags.


