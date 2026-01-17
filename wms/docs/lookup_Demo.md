# CLI Commands Reference – Warehouse Management System (SPL)

This document is a **command reference / cheat sheet** for the Warehouse Management System Software Product Line.

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
go run ./cmd/app location.add L1 "Main Location"
```

---

### location.list
List locations.

```bash
location.list [limit]
```

Example:
```bash
go run ./cmd/app location.list
```

---

### item.add
Create a new item.

```bash
item.add <sku> <name> <description>
```

Example:
```bash
go run ./cmd/app item.add SKU1 "Tape" "Packing tape"
```

---

### item.list
List items.

```bash
item.list [limit]
```

Example:
```bash
go run ./cmd/app item.list
```

---

### stock.in
Increase stock at a location.

```bash
stock.in <sku> <locationCode> <quantity> [reference]
```

Example:
```bash
go run ./cmd/app stock.in SKU1 L1 10 ref-in
```

---

### stock.out
Decrease stock at a location.

```bash
stock.out <sku> <locationCode> <quantity> [reference]
```

Example:
```bash
go run ./cmd/app stock.out SKU1 L1 3 ref-out
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
go run -tags "reporting" ./cmd/app reporting.inventory
```

---

### reporting.movements
List stock movements.

```bash
reporting.movements [limit]
```

Example:
```bash
go run -tags "reporting" ./cmd/app reporting.movements
```

---

### reporting.warehouse.summary
Warehouse KPI summary (requires `multiwarehouse`).

```bash
reporting.warehouse.summary <warehouseCode>
```

Example:
```bash
go run -tags "reporting multiwarehouse" ./cmd/app reporting.warehouse.summary W1
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
List bins of a location.

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
Delete a bin.

```bash
logistics.bin.delete <locationCode> <binCode>
```

---

## Picking Feature (`picking`)

### picking.list.create
Create a pick list.

```bash
picking.list.create <orderRef>
```

---

### picking.list.show
Show pick list details.

```bash
picking.list.show <orderRef>
```

---

### picking.list.complete
Mark a pick list as completed.

```bash
picking.list.complete <orderRef>
```

---

## Barcode Feature (`barcode`)

### barcode.generate
Generate a barcode for an item.

```bash
barcode.generate <sku>
```

---

### barcode.scan
Scan a barcode.

```bash
barcode.scan <barcodeValue>
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
Show latest audit events.

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

## Notifications Feature (`notifications`)

### notifications.list
List notification rules.

```bash
notifications.list [limit]
```

---

### notifications.add
Add a notification rule.

```bash
notifications.add <event> <target>
```

---

### notifications.remove
Remove a notification rule.

```bash
notifications.remove <ruleId>
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

- Commands exist **only if the corresponding feature is compiled**
- Missing commands indicate disabled features
- Invalid feature combinations fail at **compile time**
- Core commands are always available
