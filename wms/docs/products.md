# Products – Warehouse Management System (SPL)

This document describes the **product concept** of the Warehouse Management System Software Product Line.

Products in this project are **purely conceptual SPL artifacts**.  
They are used to **document and present valid feature combinations**, but they are **not concrete software products**.

Products:
- Do **not** represent separate code bases
- Are **not** Clone-and-Own variants
- Do **not** control variability at runtime
- Exist only for **documentation and CLI output**

All concrete products are derived from the **same code base** using **Go build tags** at compile time.  
Build tags determine which features are included in the binary; products only describe these combinations.

---

## Product Derivation

- Each feature corresponds to one or more Go build tags
- A concrete build enables a specific tag combination
- Disabled features are completely excluded from the compiled binary

Example:

```bash
go run -tags "pro reporting logistics picking barcode" ./cmd/app
