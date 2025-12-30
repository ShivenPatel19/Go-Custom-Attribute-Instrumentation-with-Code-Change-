Below is a **clean, self-contained documentation-style write-up** covering **only the two methods** you asked for.
No extra concepts, no cross-references, no additional methods.

You can directly use this as a **design / Jira / internal product doc**.

---

# Custom Attribute Instrumentation in Go (OpenTelemetry)

## Scope

This document describes **two supported mechanisms** for attaching **business-specific custom attributes** to OpenTelemetry spans in **Go**, using only the OpenTelemetry APIs.

The focus is on:

* How attributes are attached
* Span and trace behavior
* When each method should be used

---

## Method 1: Active Span Enrichment (Same Span)

### Definition

Active Span Enrichment means **adding custom attributes to an already existing span** without creating a new span or modifying the trace structure.

The span must already be present in the `context.Context`, typically created by:

* Auto-instrumentation (HTTP, gRPC, DB, etc.)
* Upstream manual instrumentation

This method performs **pure enrichment** of the current span.

---

### How It Works

1. A span is already active in the execution context.
2. The span is retrieved from `context.Context`.
3. Custom business attributes are appended to that span.
4. Span lifecycle and timing remain unchanged.

---

### Behavior Characteristics

| Aspect                      | Behavior     |
| --------------------------- | ------------ |
| New span created            | No           |
| Trace ID                    | Unchanged    |
| Parent / Child relationship | Unchanged    |
| Span timing                 | Not affected |
| Trace structure             | Not modified |

---

### No-Span Scenario

If no span exists in the context:

* A **no-op span** is returned.
* Attribute calls are safely ignored.
* No errors or crashes occur.
* Attributes are **not recorded**.

This behavior is guaranteed by OpenTelemetry.

---

### Suitable Use Cases

* Adding business metadata (IDs, flags, classifications)
* Enriching auto-instrumented request spans
* Injecting KPIs or domain attributes
* Dynamic or rule-based attribute injection
* Zero-trace-shape modification scenarios

---

### Advantages

* Very low overhead
* No additional spans generated
* Safe to call anywhere
* Fully compatible with auto-instrumentation
* Ideal for dynamic instrumentation platforms

---

### Limitations

* Cannot measure execution time of a specific operation
* Depends on the presence of an active span

---

## Method 2: Explicit Span Creation with Attributes (New Child Span)

### Definition

Explicit Span Creation involves **creating a new span** to represent a logical business operation and attaching custom attributes to that span.

The newly created span becomes a **child span** of the currently active span, if one exists.

This method performs **explicit instrumentation**, not just enrichment.

---

### How It Works

1. A tracer is used to start a new span.
2. The new span is linked to the current context.
3. Custom attributes are attached at span creation or during execution.
4. The span is explicitly ended by the application.

---

### Behavior Characteristics

| Aspect           | Behavior       |
| ---------------- | -------------- |
| New span created | Yes            |
| Span type        | Child span     |
| Trace ID         | Same as parent |
| Execution timing | Captured       |
| Trace structure  | Modified       |

---

### No-Parent Scenario

If no active span exists:

* The created span becomes a **root span**
* A **new trace** is started automatically

---

### Suitable Use Cases

* Business operations requiring timing visibility
* Logical domain steps (e.g., payment, validation, rule execution)
* Performance analysis and debugging
* Operations not covered by auto-instrumentation

---

### Advantages

* Precise duration measurement
* Clear separation of business operations
* Improved trace readability for complex workflows
* Explicit control over span boundaries

---

### Limitations

* Increases total span volume
* Can clutter traces if used excessively
* Requires careful span naming and sampling strategy

---

## Comparison Summary

| Dimension              | Active Span Enrichment | Explicit Span Creation |
| ---------------------- | ---------------------- | ---------------------- |
| Creates new span       | No                     | Yes                    |
| Modifies trace shape   | No                     | Yes                    |
| Requires existing span | Yes                    | No                     |
| Captures timing        | No                     | Yes                    |
| Best for attributes    | Yes                    | No                     |
| Best for operations    | No                     | Yes                    |

---

## Conclusion

* **Active Span Enrichment** should be the **default approach** for adding business-specific attributes.
* **Explicit Span Creation** should be used **selectively** when a business operation requires its own visibility and timing.

Both methods are fully supported by OpenTelemetry Go and are safe for production use.