<div align="center">

# #001 Hexagonal (Ports and Adapters) Architecture

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Forms Hub services need to be testable without external dependencies and adaptable to infrastructure changes over time. Embedding infrastructure concerns (MongoDB, Redis, HTTP clients) directly in application or domain logic makes both testing and evolution expensive — a change to the persistence layer would require changes to domain code, and running tests would require live infrastructure.

## Decision

Each service applies a Ports and Adapters (Hexagonal) structure. Domain and application logic lives in `core/` and has no imports from `adapters/` or any infrastructure package. Inbound port interfaces in `core/ports/primary.go` define the contracts that REST handlers drive into the application. Outbound port interfaces in `core/ports/secondary.go` define the contracts that persistence, client, and worker adapters fulfil. All adapter implementations live in `adapters/` and are wired at startup in `core/core.go`. Every infrastructure concern has an in-memory substitute selectable via configuration, requiring no external dependencies for local development or testing.

## Consequences

- The domain and application layers are independently testable using in-memory adapters with no external dependencies.
- Adding a new infrastructure adapter requires implementing the relevant port interface only; no domain or application code changes.
- The boundary is enforced by Go package structure — `core/` packages do not import `adapters/` packages.
- All infrastructure wiring is centralised in `core/core.go`, making the dependency graph explicit and auditable at startup.
