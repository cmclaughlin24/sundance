<div align="center">

# #006 MongoDB as Primary Datastore

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Forms Hub's core data structure — the form definition — is a deeply nested, polymorphic hierarchy of pages, sections, fields, and per-field-type attributes. Representing this in a relational schema would require multiple joined tables and complex mapping logic to reconstruct the hierarchy at read time. A document-oriented model maps more naturally to this structure without requiring a rigid schema that must be migrated with each form definition change.

## Decision

MongoDB is used as the primary datastore for all domain data across both services. Each aggregate is persisted as a document in its own collection. A MongoDB replica set is required in production to support multi-document transactions, which are used for operations that must be atomic across multiple aggregates (e.g. publishing a tag version while atomically deprecating the existing active version). All access goes through the `Repository` port interfaces defined in `core/ports/secondary.go`; no service code imports the MongoDB driver directly.

The Ports and Adapters architecture (ADR-1) mitigates the risk of tight coupling to MongoDB specifically. The `Repository` port interfaces are the only contracts the domain and application layers depend on; swapping MongoDB for a different datastore requires only a new adapter implementation with no changes to domain or application logic.

## Consequences

- The deeply nested, polymorphic form definition hierarchy is persisted and retrieved as a single document read, avoiding complex joins or recursive queries.
- A MongoDB replica set is a hard production dependency. Single-node deployments do not support multi-document transactions.
- The `(form_id, version)` and `(tag_id, version)` unique indexes enforce version uniqueness at the database level, providing a safety net beyond the application-layer checks.
- All MongoDB concerns are isolated behind port interfaces; the in-memory repository implementations allow the domain and application layers to be tested without a running MongoDB instance.
