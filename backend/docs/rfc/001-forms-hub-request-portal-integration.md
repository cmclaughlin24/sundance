<div align="center">

# RFC-0001 Forms Hub & Request Portal Backend Integration

##### A proposal that describes a change, its motivation, and its implications.

</div>

---

- **Start Date:** 2026-07-08
- **Status:** Draft
- **Author(s):** Wm. Curtis McLaughlin
- **Related:** [ADR-003](../adr/003-asynchronous-submission-processing.md), [ADR-009](../adr/009-outbox-pattern-polling-vs-streaming.md)

---

## Summary

This RFC proposes the backend integration between the Request Portal and Forms Hub to support form-based catalog items in the request cart. When an actor adds a Form Catalog Item to their cart, the Request Portal must coordinate with Forms Hub to create a form submission, associate it with the cart item, and react to the asynchronous submission result (Accepted, Rejected, or Failed) published by Forms Hub. This integration introduces a new execution path in the Request Portal alongside the existing Standard Catalog Item path, and establishes a bidirectional contract between the two systems via synchronous REST and asynchronous events.

## Motivation

The Request Portal currently supports a single cart item type — the Standard Catalog Item — which requires no coordination with an external service at the time of cart creation. As the platform evolves to support form-based catalog items, a new execution path is needed: one where adding an item to the cart requires the Request Portal to initiate a form submission in Forms Hub and remain responsive to its asynchronous outcome.

Forms Hub processes submissions asynchronously by design (ADR-003). `POST /v1/api/submissions` returns immediately with a reference ID and a `pending` status; the actual validation, canonical tag mapping, and rule evaluation occur in a background worker. The final outcome — `accepted`, `rejected`, or `failed` — is published via a polling-based outbox relay (ADR-009). This means the Request Portal cannot treat form submission as a synchronous, fire-and-forget operation; it must store the submission reference, receive the result, and update the cart item accordingly.

Without this integration, Form Catalog Items cannot be supported in the request cart. With it, the Request Portal gains the ability to attach structured, validated form data to cart items — enabling downstream policy evaluation and fulfillment workflows to operate on a canonical, trusted representation of the request.

## Design

The integration between the Request Portal and Forms Hub follows a two-phase interaction model. In the first phase, the Request Portal synchronously initiates a form submission via the Forms Hub REST API and associates the returned reference with the cart item. In the second phase, Forms Hub asynchronously processes the submission and publishes the result, which the Request Portal receives and uses to finalize the cart item state and trigger downstream policy evaluation.

### Sequence

![Forms Hub & Request Portal Backend Integration](../imgs/Forms%20Hub%20x%20Request%20Portal%20Backend%20Integration.png)

The diagram above illustrates the full interaction for Path B (Form Catalog Item). Path A (Standard Catalog Item) requires no coordination with Forms Hub and is included for contrast only.

1. **Create cart item**: The actor submits `POST /api/cart/items` to the Request Portal. The request payload identifies the catalog item type, branching execution into Path A or Path B.

2. **Create form submission** _(Path B only)_: The Request Portal calls `POST /v1/api/submissions` on Forms Hub. Forms Hub persists the submission with a status of `pending` and returns `202 Accepted` with a submission reference ID. The Request Portal returns this response to the actor immediately.

3. **Associate cart item with submission reference**: The Request Portal stores the submission reference ID against the cart item, linking the two records internally.

4. **Validate submission and produce canonical representation** _(async, Forms Hub internal)_: A background worker validates the submission, resolves data source lookups, and maps fields to canonical tags. This step is internal to Forms Hub and runs concurrently with step 3.

5. **Publish submission result**: Forms Hub publishes the final outcome (`accepted`, `rejected`, or `failed`) to the Request Portal. The delivery mechanism for this step is unresolved — see Unresolved Questions.

6. **Finalize cart item**: On `accepted`, the Request Portal links the canonical submission data to the cart item. On `rejected` or `failed`, the cart item is marked as invalid.

7. **Expand request and evaluate policy checks** _(on `accepted` only)_: The Request Portal expands the request and runs policy evaluation against the now-complete cart item.

### Changes Required

Explicit list of behavioral or structural changes required by each
system involved. This is the consumer/provider obligation checklist.

#### Request Portal

- **Define `form` as a new value for the existing `items[].type` discriminator**: `POST /api/cart/items` already accepts a `type` field on each item. The Request Portal must define `form` as a new valid value alongside existing types, routing execution into Path B when present. No schema change is required; the routing logic in the handler must be extended.

- **Define the Path B `items[].payload` contract**: When `item[].type` is `form`, the freeform `payload` object must carry `formId`, `versionId`, and a `values` array (field ID -> value pairs, with optional `collectionIndex` for repeating fields). These fields are extracted from `payload` and proxied directly to `POST /v1/api/submissions` on Forms Hub.

- **Call Forms Hub on cart item creation (Path B)**: The cart item creation handler must issue a synchronous `POST /v1/api/submissions` call to Forms Hub immediately after validating the request. The `202 Accepted` response and its `referenceId` must be captured before the handler returns to the caller.

- **Extend the cart item data model**: The cart item data model must be extended with three new fields: `submissionReferenceId` (populated at creation from Forms Hub response), `submissionStatus` (`pending` | `accepted` | `rejected` | `failed`, updated when the result is received), and `canonicalData` (the canonical fact set linked to the cart item on acceptance.

- **Implement a result receiver for Forms Hub**: The Request Portal must implement the inbound side of step 5 to receive the submission result from Forms Hub. The delivery mechanism is unresolved (see Unresolved Questions); the receiver must handle all three terminal statuses: `accepted`, `rejected`, and `failed`.

- **Finalize cart item state on result receipt**: On `accepted`, the receiver must update `submissionStatus`, attach the canonical data to the cart item, and trigger policy evaluation. On `rejected` or `failed`, the receiver must mark the cart item as invalid and update `submissionStatus` accordingly.

- **Defer policy evaluation for Form Catalog Items**: Policy evaluation is currently triggered synchronously on cart item creation. For Path B, this trigger must be suppressed at creation time and rewired to fire only after an `accepted` result is received and the cart item is finalized.

- **Define outbound auth strategy for Forms Hub calls**: The mechanism by which the Request Portal authenticates calls to `POST /v1/api/submissions` is not yet decided (see Unresolved Questions).

#### Forms Hub

No changes are required. The delivery mechanism for submission results is unresolved — see [Unresolved Question #1](#unresolved-questions).

> **Note:** The canonical facts in the `submission.accepted` event payload are mapped into the Request Portal's data structure by the Forms Hub tagging mechanism.

### API Contracts

#### `POST /v1/api/submissions`

**Headers**

| Header            | Required | Description                                                      |
| ----------------- | -------- | ---------------------------------------------------------------- |
| `X-Tenant-ID`     | Yes      | Tenant identifier                                                |
| `Idempotency-Key` | Yes      | Client-generated UUIDv7; prevents duplicate submissions on retry |
| `Authorization`   | Yes      | `Bearer <token>` — PingFederate client credentials token         |

**Request**

```json
{
  "formId": "<UUIDv7>",
  "versionId": "<UUIDv7>",
  "values": [
    {
      "fieldId": "<string>",
      "value": "<any>",
      "collectionIndex": "<int | omitted>"
    }
  ]
}
```

**Response — `202 Accepted`**

```json
{
  "message": "Accepted!",
  "data": {
    "id": "<UUIDv7>",
    "tenantId": "<string>",
    "formId": "<UUIDv7>",
    "versionId": "<UUIDv7>",
    "referenceId": "<UUIDv7>",
    "status": "pending",
    "values": [
      {
        "fieldId": "<string>",
        "value": "<any>",
        "collectionIndex": "<int | omitted>"
      }
    ],
    "createdAt": "<RFC3339>",
    "updatedAt": "<RFC3339>"
  }
}
```

**Error Responses**

| Status | Condition                                                                                                        |
| ------ | ---------------------------------------------------------------------------------------------------------------- |
| `400`  | Missing/invalid body, missing `X-Tenant-ID`, missing `Idempotency-Key`, validation failure, empty `values` array |
| `401`  | Missing or invalid bearer token                                                                                  |
| `409`  | Duplicate submission — same `Idempotency-Key` already exists                                                     |
| `500`  | Unexpected server error                                                                                          |

### Event Contracts

Events are published to Kafka by the Forms Hub outbox relay. The partition key for all events is the `SubmissionID`.

#### `submission.accepted`

**Topic:** `submission.accepted`
**Partition key:** `SubmissionID`

```json
{
  "referenceId": "<UUIDv7>",
  "tenantId": "<string>",
  "formId": "<UUIDv7>",
  "versionId": "<UUIDv7>",
  "facts": {
    "applicant": {
      "firstName": "Jane",
      "lastName": "Doe",
      "address": [
        {
          "street": "123 Main St",
          "city": "Springfield"
        }
      ]
    }
  }
}
```

The `facts` field is a nested map built from dot-delimited canonical tag key paths (e.g. `applicant.address[].city`). Array segments (`[]`) produce arrays of objects. Leaf values are the submitted field values after canonical mapping via the Forms Hub tagging mechanism.

#### `submission.rejected`

**Topic:** `submission.rejected`
**Partition key:** `SubmissionID`

```json
{
  "referenceId": "<UUIDv7>",
  "tenantId": "<string>",
  "formId": "<UUIDv7>",
  "versionId": "<UUIDv7>",
  "reason": "<string>"
}
```

#### `submission.failed`

The `failed` status does not currently produce a Kafka event — see [Unresolved Question #3](#unresolved-questions).

## Drawbacks

- **Forms Hub availability dependency:** Path B introduces a runtime dependency on Forms Hub at cart item creation time. If Forms Hub is unavailable, `POST /api/cart/items` will fail for Form Catalog Items. Standard Catalog Items (Path A) are unaffected. The Request Portal must handle this failure gracefully and surface a meaningful error to the actor.

- **Orphaned submission risk:** If the Request Portal successfully calls `POST /v1/api/submissions` but fails to persist the `submissionReferenceId` against the cart item (step 3), the submission exists in Forms Hub but is unlinked from the cart. There is currently no recovery strategy for this state — see [Unresolved Question #4](#unresolved-questions).

- **Pending cart item timeout:** A Form Catalog Item cart item remains in a `pending` state until a submission result event is received. If the event is never delivered — due to a processing failure or delivery gap — the cart item remains pending indefinitely. A timeout strategy is not yet defined — see [Unresolved Question #5](#unresolved-questions).

## Rationale and Alternatives

The two-phase asynchronous integration — synchronous intake via `POST /v1/api/submissions` followed by event-driven result delivery — was chosen for the following reasons:

- **Reusability:** Forms Hub is purpose-built for form definition, validation, and canonical mapping. The legacy tooling it replaces undergoes approximately 100 changes per month, making a dedicated, independently deployable forms service a strategic necessity rather than an implementation convenience. Delegating form concerns to Forms Hub avoids duplicating this complexity in the Request Portal.

- **Separation of concerns:** Form validation, data source lookup, and canonical tag mapping are owned by Forms Hub. Embedding this logic in the Request Portal would couple two distinct domains and create a maintenance burden as form requirements evolve.

- **Decoupled throughput:** The asynchronous model ensures that Request Portal cart item creation is not blocked by the latency of submission processing. Validation, data source lookups, and canonical mapping can be slow and variable; holding the intake request open until processing completes would expose the Request Portal to that variability.

### Alternatives Considered

**Synchronous submission endpoint (`POST /v1/api/submissions/sync`)**

A synchronous variant of the submissions endpoint was considered — one that performs validation and canonical mapping inline before returning, rather than deferring to a background worker. This would simplify the Request Portal integration by eliminating the need to handle an asynchronous result.

This was not chosen for two reasons: first, if asynchronous validation steps (e.g. external data source lookups) are introduced in the future, a synchronous endpoint would couple intake throughput directly to that latency; second, it would tightly couple the Request Portal's cart creation flow to Forms Hub's processing time, introducing operational risk if Forms Hub degrades.

**Polling for submission status**

Rather than consuming result events, the Request Portal could poll `GET /v1/api/submissions/by-reference/{referenceId}/status` until a terminal status is reached. This would eliminate the need for the Request Portal to implement an event consumer or webhook receiver.

This was not chosen because it places a heavier burden on the Request Portal — it must manage polling intervals, handle timeouts, and deal with long-running submissions gracefully. Event-driven delivery inverts this burden: Forms Hub notifies the Request Portal when processing is complete, with no polling overhead.

### Impact of Not Integrating

Without this integration, Form Catalog Items cannot be supported in the request cart. There is no current workaround — this is a feature gap that blocks the Request Portal from consuming form-based catalog items entirely.

## Unresolved Questions

1. **Result delivery mechanism (step 5)**: The mechanism by which Forms Hub delivers the submission result to the Request Portal is not yet decided. The long-term goal is delivery via a message broker. A webhook-based approach may be required as a short-term solution pending broker infrastructure availability.

2. **Outbound authentication from Request Portal to Forms Hub**: The mechanism by which the Request Portal authenticates its calls to `POST /v1/api/submissions` is not yet decided. Candidates include service-to-service JWT (client credentials flow via PingFederate), mTLS, or an API key. The chosen approach must align with the enterprise security posture and the authentication infrastructure available to the Request Portal.

3. **`submission.failed` event**: A `failed` submission does not currently emit a Kafka event. Forms Hub records the failure internally and exposes a replay endpoint (`POST /v1/api/submissions/{submissionId}/replay`). For the Request Portal to react to a `failed` outcome, either a `submission.failed` Kafka event must be introduced, or an alternative notification mechanism must be defined.

4. **Orphaned submission recovery:** If the Request Portal fails to persist the `submissionReferenceId` after a successful `POST /v1/api/submissions` call, the submission is orphaned in Forms Hub with no corresponding cart item link. A recovery or reconciliation strategy for this state has not been defined.

5. **Pending cart item timeout:** A cart item in `pending` state will remain so indefinitely if a submission result event is never received. A timeout threshold and ownership (likely the Request Portal) have not been defined.

6. **Tag format requirements for `submission.accepted` facts:** The `facts` map in the `submission.accepted` payload is built from canonical tag key paths defined in Forms Hub. For the Request Portal to correctly consume this payload, the tag keys and structure must be agreed upon and standardized between the two teams. Specifically: which tag keys are required, what value types are expected at each leaf, and how collection segments (`[]`) map to the Request Portal's internal data model must be formally defined before implementation begins.

7. **Submission amendment:** An actor may need to update form field values after a Form Catalog Item has been added to the cart. For this release, an amendment is handled by creating a new submission via `POST /v1/api/submissions`, replacing the existing `submissionReferenceId` on the cart item, and returning the cart item to `pending` state while the new submission is processed. The design of the cart item update endpoint and the handling of the superseded submission in Forms Hub must be defined.

## Future Possibilities

- **Submission amendments:** For this release, amendments are handled by creating a new submission. A dedicated amendment endpoint in Forms Hub (e.g. `POST /v1/api/submissions/{submissionId}/amend`) could provide a more explicit contract in future releases, preserving the submission history and reducing ambiguity around which submission is authoritative for a given cart item.

- **Audit system consumers:** The `submission.accepted` and `submission.rejected` Kafka events are published to named topics and may be consumed by additional downstream systems in the future, such as audit or compliance services, without requiring changes to either Forms Hub or the Request Portal.
