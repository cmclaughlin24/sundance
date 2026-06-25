<div align="center">

# #005 Strategy Pattern for Data Source Lookup and Field Validation

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Forms Hub supports multiple data source types (`static`, `scheduled`, `webhook`, `data-lake`) for resolving lookup values, and multiple field types (`text`, `number`, `select`, `checkbox`, `date`) for validating submitted values. Each type has distinct behaviour. A switch-statement approach would centralise this logic, making it harder to extend, test in isolation, and reason about as the number of types grows.

## Decision

Both concerns are implemented using the Strategy pattern via a generic `StrategyRegistry[K, S]` in `pkg/common/stratreg`. Each data source type has a dedicated `LookupStrategy` implementation; each field type has a dedicated `FieldValidatorStrategy` implementation. Strategies are registered at startup and resolved at runtime by key. New types are added by implementing the relevant interface and registering the strategy — no existing resolution logic is modified.

## Consequences

- Each strategy is independently testable in isolation with no dependency on other strategies or the registry.
- An unregistered key returns `ErrStrategyNotFound` at runtime. For field validation this is treated as non-retryable — the submission transitions to `failed` immediately rather than retrying indefinitely.
- New data source types and field types can be added without modifying existing strategies or the resolution path.
- The registry is populated at startup; misconfiguration (missing strategy registration) surfaces as a runtime error on first use rather than at compile time.
