# CLI Commands Reference – Warehouse Management System (SPL)

This document is a **command reference / cheat sheet**.

Each CLI command is listed with:
- Purpose
- Syntax
- Minimal example

Use this file as a **quick lookup** during demos or grading.

---

## Core Commands (Always Available)

### location.add
Create a new location.

```bash
location.add <code> <name>
```

Example:
```bash
go run -tags "min" ./cmd/app location.add L1 "Main Location"
```

---

### location.list
List locations.

```bash
location.list [limit]
```

Example:
```bash
go run -tags "min" ./cmd/app location.list
```

---

### item.add
Create a new item.

```bash
item.add <sku> <name> <description>
```

Example:
```bash
go run -tags "min" ./cmd/app item.add SKU1 "Tape" "Packing tape"
```

---

### item.list
List items.

```bash
item.list [limit]
```

---

### stock.in
Increase stock at a location.

```bash
stock.in <sku> <locationCode> <quantity> [reference]
```

Example:
```bash
go run -tags "min" ./cmd/app stock.in SKU1 L1 10 ref-in
```

---

### stock.out
Decrease stock at a location.

```bash
stock.out <sku> <locationCode> <quantity> [reference]
```

Example:
```bash
go run -tags "min" ./cmd/app stock.out SKU1 L1 3 ref-out
```

---

## Reporting Feature (`reporting`)

### reporting.inventory
Show inventory overview.

```bash
reporting.inventory
```

Example:
```bash
go run -tags "min reporting" ./cmd/app reporting.inventory
```

---

### reporting.movements
List stock movements.

```bash
reporting.movements [limit]
```

---

### reporting.warehouse.summary
Warehouse KPI summary (requires `multiwarehouse`).

```bash
reporting.warehouse.summary <warehouseCode>
```

Example:
```bash
go run -tags "min reporting multiwarehouse" ./cmd/app reporting.warehouse.summary W1
```

---

## Logistics Feature (`logistics`)

### logistics.zone.add
Create a zone in a location.

```bash
logistics.zone.add <locationCode> <zoneCode> [name]
```

---

### logistics.zone.list
List zones of a location.

```bash
logistics.zone.list <locationCode> [limit]
```

---

### logistics.zone.delete
Delete a zone.

```bash
logistics.zone.delete <locationCode> <zoneCode>
```

---

### logistics.bin.add
Create a bin inside a zone.

```bash
logistics.bin.add <locationCode> <zoneCode> <binCode> [name]
```

---

### logistics.bin.list
List bins of a location (optionally filtered by zone).

```bash
logistics.bin.list <locationCode> [zoneCode] [limit]
```

---

### logistics.bin.assign
Assign an item to a bin.

```bash
logistics.bin.assign <locationCode> <binCode> <sku>
```

---

### logistics.bin.unassign
Unassign an item from a bin.

```bash
logistics.bin.unassign <locationCode> <binCode> <sku>
```

---

### logistics.bin.delete
Delete a bin (fails if items are assigned).

```bash
logistics.bin.delete <locationCode> <binCode>
```

---

## Tracking Feature (`tracking`)

### tracking.set
Attach tracking information to an outbound order.

```bash
tracking.set <orderRef> <trackingId> [url] [carrier]
```

---

### tracking.get
Query tracking information.

```bash
tracking.get <orderRef>
```

---

### tracking.clear
Remove tracking information.

```bash
tracking.clear <orderRef>
```

---

## Audit Feature (`audit`)

### audit.list
List audit events.

```bash
audit.list [limit]
```

---

### audit.filter
Filter audit events.

```bash
audit.filter action=<a> entity=<e> actor=<u> limit=<n>
```

---

### audit.tail
Show last audit events.

```bash
audit.tail
```

---

## Authentication Feature (`auth`)

### auth.user.add
Create a user (Admin only).

```bash
auth.user.add <username> <role> <password>
```

---

### auth.user.list
List users (Admin only).

```bash
auth.user.list [limit]
```

---

### auth.user.disable
Disable a user (Admin only).

```bash
auth.user.disable <username>
```

---

### auth.whoami
Show authenticated user.

```bash
auth.whoami
```

---

## MultiWarehouse Feature (`multiwarehouse`)

### warehouse.add
Create a warehouse.

```bash
warehouse.add <warehouseCode> [name]
```

---

### warehouse.list
List warehouses.

```bash
warehouse.list [limit]
```

---

### warehouse.location.assign
Assign a location to a warehouse.

```bash
warehouse.location.assign <warehouseCode> <locationCode>
```

---

### warehouse.location.list
List locations of a warehouse.

```bash
warehouse.location.list <warehouseCode> [limit]
```

---

## Notes

- Commands only exist if the corresponding feature is compiled.
- Missing commands indicate disabled features.
- Invalid feature combinations fail at compile time.

