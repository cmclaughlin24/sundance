<div align="center">

# RFC-0002 Forms Hub & Service Catalog Integration

##### A proposal that describes a change, its motivation, and its implications.

</div>

---

- **Start Date:** 2026-07-16
- **Status:** Draft
- **Author(s):** Wm. Curtis McLaughlin
- **Related:** [RFC-0001](001-forms-hub-request-portal-integration.md)

---

## Summary

This RFC proposes the integration between Forms Hub and Service Catalog to keep the service catalog synchronized with the published state of forms. When a form designer publishes or retires a form version in Forms Hub, the Forms Service emits a Kafka event — `PublishFormVersion` or `RetireFormVersion` respectively — on a topic owned by Forms Hub. Service Catalog consumes both events: on `PublishFormVersion` it creates or updates a catalog entry that links to the form and version, making it selectable as a form-backed catalog item (see RFC-0001); on `RetireFormVersion` it marks the corresponding catalog entry inactive or removes it, preventing it from being selected for new requests.

This integration requires Forms Hub to emit two new domain events via its existing outbox infrastructure and requires Service Catalog to implement a consumer for both, while introducing no synchronous coupling between the two systems at runtime.

## Motivation

Forms Hub is purpose-built for form definition, versioning, and submission processing. Service Catalog is the system through which actors discover and initiate form-backed service requests. For these two systems to work together, Service Catalog must know which form versions are available and which are no longer valid — but as independently developed services with no shared data store, there is no mechanism today by which a form version's lifecycle state is reflected in the catalog.

Without this integration, form-backed catalog items cannot exist in Service Catalog. A form designer has no way to surface a published form version as a selectable catalog item, and retiring a form version has no effect on what actors can see or request. This is the prerequisite integration that enables the request flow described in RFC-0001 — a form-backed catalog item cannot be selected, rendered, or submitted against a form version that Service Catalog does not know about.

With this integration, the form version lifecycle drives catalog state automatically and asynchronously. Publishing a form version creates or updates the corresponding catalog entry, making it immediately available for actors to discover and initiate requests. Retiring a form version removes or deactivates the catalog entry, ensuring actors cannot initiate new requests against a stale form — protecting both data quality and the integrity of downstream fulfillment workflows.

## Design

The integration between Forms Hub and Service Catalog is entirely asynchronous. Forms Hub owns the `form` Kafka topic and publishes domain events via its existing outbox infrastructure when a form version is published or retired. Service Catalog consumes both event types and updates its catalog state accordingly. No synchronous HTTP calls are introduced between the two systems.

### Sequence

![Forms Hub & Service Catalog Integration](../imgs/Forms%20Hub%20x%20Service%20Catalog%20Integration.png)

1. **Publish or retire a form version**: A form designer calls `POST /api/v1/forms/{id}/versions/{versionId}/publish` or `POST /api/v1/forms/{id}/versions/{versionId}/retire` via the Forms Hub MFE.

2. **Transition version and write outbox event**: The Forms Service transitions the form version to `active` (or `retired`), marshals a `PublishFormVersionPayload` (or `RetireFormVersionPayload`), and writes it as a domain event to the `outbox` collection transactionally with the version status update.

3. **Relay event to Kafka**: The Outbox Relay Worker claims the event and publishes it to the `form` Kafka topic, keyed by `FormID`, with an `eventType` message header of `published` (or `retired`).

4. **Create, update, or deactivate catalog entry**: Service Catalog consumes the event from the `form` topic. On `published`, it creates or updates a catalog entry linking to the form and version, making it selectable for new requests. On `retired`, it removes or deactivates the corresponding catalog entry, preventing actors from initiating new requests against it.

### Changes Required

#### Forms Hub

- **No backend changes required.** Both `PublishFormVersion` and `RetireFormVersion` events are already emitted via the existing outbox infrastructure when a form version is published or retired. This RFC formalizes and commits to the event contracts below as a stable, versioned interface. Changes to the payload shape of either event require a new RFC.

#### Service Catalog

- **Implement a consumer for the `form` Kafka topic**: Service Catalog must subscribe to the `form` topic and route messages by the `eventType` header.

- **Handle `eventType: published`**: On receipt, create or update a catalog entry carrying the `formId`, `versionId`, `version`, `name`, `description`, `metadata`, and `tenantId` from the event payload. The catalog entry must link to the form and version such that it can be selected as a form-backed catalog item (see RFC-0001).

- **Handle `eventType: retired`**: On receipt, remove or deactivate the catalog entry for the given `formId` and `versionId`, preventing actors from initiating new requests against it. See Unresolved Questions for whether this is a soft deactivation or a hard delete.

- **Implement idempotent event handling**: The outbox relay may deliver the same event more than once. Both handlers must be safe to apply repeatedly without producing duplicate or inconsistent catalog state.

### API Contracts

This integration is entirely asynchronous. No synchronous HTTP calls are introduced between Forms Hub and Service Catalog.

### Event Contracts

Both events are published to the `form` Kafka topic by the Forms Hub Outbox Relay Worker. The event type is communicated via a Kafka message header.

**Topic:** `form`
**Partition key:** `FormID`
**Event type header:** `eventType: published | retired`

---

#### `PublishFormVersion`

Emitted by the Forms Service when a form version is transitioned from `draft` to `active`.

**Headers**

| Header      | Value       |
| ----------- | ----------- |
| `eventType` | `published` |

**Payload**

```json
{
  "tenantId": "<string>",
  "formId": "<UUIDv7>",
  "versionId": "<UUIDv7>",
  "version": "<int>",
  "name": "<string>",
  "description": "<string>",
  "metadata": { "<key>": "<string>" },
  "publishedBy": "<string>"
}
```

| Field         | Type                | Description                                                  |
| ------------- | ------------------- | ------------------------------------------------------------ |
| `tenantId`    | `string`            | The tenant that owns the form                                |
| `formId`      | `string` (UUIDv7)   | The form identifier                                          |
| `versionId`   | `string` (UUIDv7)   | The specific version that was published                      |
| `version`     | `int`               | The auto-incremented version number                          |
| `name`        | `string`            | The form name at the time of publish                         |
| `description` | `string`            | The form description at the time of publish                  |
| `metadata`    | `map[string]string` | Arbitrary key-value metadata attached to the form            |
| `publishedBy` | `string`            | The subject identifier of the user who triggered the publish |

---

#### `RetireFormVersion`

Emitted by the Forms Service when a form version is transitioned from `active` to `retired`.

**Headers**

| Header      | Value     |
| ----------- | --------- |
| `eventType` | `retired` |

**Payload**

```json
{
  "tenantId": "<string>",
  "formId": "<UUIDv7>",
  "versionId": "<UUIDv7>",
  "version": "<int>",
  "retiredBy": "<string>"
}
```

| Field       | Type              | Description                                                 |
| ----------- | ----------------- | ----------------------------------------------------------- |
| `tenantId`  | `string`          | The tenant that owns the form                               |
| `formId`    | `string` (UUIDv7) | The form identifier                                         |
| `versionId` | `string` (UUIDv7) | The specific version that was retired                       |
| `version`   | `int`             | The auto-incremented version number                         |
| `retiredBy` | `string`          | The subject identifier of the user who triggered the retire |

## Drawbacks

## Rationale and Alternatives

### Alternatives Considered

### Impact of Not Integrating

## Unresolved Questions

- When the `RetireFormVersion` event is received, should Service Catalog soft-deactivate the catalog entry (preserving it for historical reference) or hard-delete it? The answer determines the data model and idempotency requirements on the Service Catalog side.

## Future Possibilities
