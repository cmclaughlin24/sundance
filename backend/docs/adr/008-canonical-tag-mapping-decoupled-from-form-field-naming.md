<div align="center">

# #008 Canonical Tag Mapping Decoupled from Form Field Naming

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Downstream systems that consume submission data need a stable, consistent data structure regardless of which form version collected the data or how its fields were named. Without a normalisation layer, downstream consumers would need to understand the structure of every form version they consume, coupling them tightly to form design decisions and making form evolution expensive.

## Decision

Each form field carries one or more `FieldTagMapping` entries that associate it with a versioned canonical tag. During submission processing, submitted field values are mapped through these tag mappings to produce a `CanonicalFact` keyed by tag key rather than form field key. Tags follow an independent versioning lifecycle (`draft → active → deprecated → retired`), and the tag resolution policy selects the winning tag version at processing time. This means downstream systems consume a consistent structure defined by the tag model, not the form model.

## Consequences

- Downstream systems are decoupled from form field naming and layout; forms can evolve independently of integration contracts.
- Form designers must explicitly map fields to tag versions; unmapped fields produce no canonical output.
- The tag lifecycle enforces a managed deprecation path before a tag version is retired, giving downstream consumers advance notice of breaking changes.
- At most one `active` tag version per tag can exist at any given time, ensuring a deterministic resolution at submission processing time.
