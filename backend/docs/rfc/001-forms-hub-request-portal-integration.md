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

This RFC proposes the backend integration between the Request Portal and Forms Hub to support form-based catalog items in the request cart. When an actor adds a Form Catalog Item to their cart, the Request Portal must coordinate with Forms Hub to synchronously normalize the form submission — validating field values and producing a canonical fact map — and link the result to the cart item before returning to the actor. This integration introduces a new execution path in the Request Portal alongside the existing Standard Catalog Item path, and establishes a unidirectional synchronous contract between the two systems via REST.

## Motivation

The Request Portal currently supports a single cart item type — the Standard Catalog Item — which requires no coordination with an external service at the time of cart creation. As the platform evolves to support form-based catalog items, a new execution path is needed: one where adding an item to the cart requires the Request Portal to validate and normalize form submission data through Forms Hub before the cart item is finalized.

Forms Hub exposes `POST /v1/api/submissions/normalize` — a synchronous endpoint that validates a submission against a form version, applies canonical tag mapping, and returns a canonical fact map inline without persisting any data. This enables the Request Portal to delegate all form validation and normalization concerns to Forms Hub within a single request cycle, providing immediate feedback to the actor at cart creation time.

Without this integration, Form Catalog Items cannot be supported in the request cart. With it, the Request Portal gains the ability to attach structured, validated canonical data to cart items — enabling downstream policy evaluation and fulfillment workflows to operate on a canonical, trusted representation of the request.

## Design

The integration between the Request Portal and Forms Hub follows a single-phase synchronous model. The Request Portal calls `POST /v1/api/submissions/normalize` on Forms Hub inline during cart item creation. Forms Hub validates the submission and returns a canonical fact map in the same request; the Request Portal links the result to the cart item and proceeds to policy evaluation. No data is persisted in Forms Hub as part of this flow.

### Sequence

![Forms Hub & Request Portal Backend Integration](../imgs/Forms%20Hub%20x%20Request%20Portal%20Backend%20Integration.png)

The diagram above illustrates the full interaction for Path B (Form Catalog Item). Path A (Standard Catalog Item) requires no coordination with Forms Hub and is included for contrast only.

1. **Create cart item**: The actor submits `POST /api/cart/items` to the Request Portal. The request payload identifies the catalog item type, branching execution into Path A or Path B.

2. **Normalize form submission** _(Path B only)_: The Request Portal calls `POST /v1/api/submissions/normalize` on Forms Hub, forwarding the `formId`, `versionId`, and `values` extracted from the item payload. Forms Hub validates the submission inline and returns `200 OK` with a canonical fact map. If validation fails, Forms Hub returns `400`; the Request Portal propagates this to the actor and the cart item is not created.

3. **Validate submission and produce canonical representation** _(Forms Hub internal)_: Forms Hub evaluates visibility and required rules, validates each resolved field, and maps field values to canonical tags. This step executes synchronously within the `POST /v1/api/submissions/normalize` request. No data is persisted.

4. **Link canonical submission to cart item**: The Request Portal attaches the canonical fact map returned by Forms Hub to the cart item and persists the record.

5. **Expand request and evaluate policy checks**: The Request Portal expands the request and runs policy evaluation against the now-complete cart item, consistent with Path A behavior.

### Changes Required

Explicit list of behavioral or structural changes required by each system involved. This is the consumer/provider obligation checklist.

#### Request Portal

- **Define `form` as a new value for the existing `items[].type` discriminator**: `POST /api/cart/items` already accepts a `type` field on each item. The Request Portal must define `form` as a new valid value alongside existing types, routing execution into Path B when present. No schema change is required; the routing logic in the handler must be extended.

- **Define the Path B `items[].payload` contract**: When `item[].type` is `form`, the freeform `payload` object must carry `formId`, `versionId`, and a `values` array (field ID → value pairs, with optional `collectionIndex` for repeating fields). These fields are extracted from `payload` and proxied directly to `POST /v1/api/submissions/normalize` on Forms Hub.

- **Call Forms Hub on cart item creation (Path B)**: The cart item creation handler must issue a synchronous `POST /v1/api/submissions/normalize` call to Forms Hub immediately after validating the request. If Forms Hub returns `400`, the handler must propagate the failure to the actor without creating the cart item. The `200 OK` response and its canonical fact map must be captured before proceeding.

- **Extend the cart item data model**: The cart item data model must be extended with one new field: `canonicalData` — the canonical fact map returned by Forms Hub and linked to the cart item on creation.

- **Define outbound auth strategy for Forms Hub calls**: The mechanism by which the Request Portal authenticates calls to `POST /v1/api/submissions/normalize` is not yet decided (see Unresolved Questions).

#### Forms Hub

No changes are required. `POST /v1/api/submissions/normalize` is already implemented and returns a canonical fact map synchronously without persisting data.

> **Note:** The canonical facts in the normalize response are built from dot-delimited canonical tag key paths (e.g. `applicant.address[].city`) defined in Forms Hub. The structure of the returned fact map is determined by the tag definitions configured for the form version.

### API Contracts

#### `POST /api/cart/items`

**Headers**

| Header              | Required | Description                      |
| ------------------- | -------- | -------------------------------- |
| `X-WF-REQUEST-DATE` | Yes      | Request date in ISO 8601 format  |
| `X-REQUEST-DATE`    | Yes      | Unique request identifier (UUID) |
| `X-CORRELATION-ID`  | Yes      | Correlation identifier (UUID)    |
| `X=WF-CLIENT-ID`    | Yes      | Client identifier (e.g. `iamx`   |

**Request**

```json
{
  "requesteeElid": "<string>",
  "parentCartDetails": {
    "<key>": "<string>"
  },
  "items": [
    {
      "type": "form",
      "payload": {
        "formId": "<string>",
        "versionId": "<string>",
        "values": [
          {
            "fieldId": "<string>",
            "value": "<any>",
            "collectionIndex": "<int | omitted>"
          }
        ]
      }
    }
  ]
}
```

> **Note**: The `payload` field is a freeform map in the existing schema (`"key": "string"`). If the Request Portal implementation stringifies the submission data on intake and unmarshals it before forwarding to Forms Hub, the wire format for the `payload` may remain a flat string map rather than the structured object shown above. The internal representation is an implementation detail of Request Portal; what matters for this contract is that the extracted `formId`, `versionId`, and `values` are forwarded correctly to `POST /v1/api/submissions/normalize`.

**Response — `200 OK`**

The response schema (`AddCartItemsResponse`) is defined by Request Portal.

**Error Responses**

| Status | Condition                                                                                   |
| ------ | ------------------------------------------------------------------------------------------- |
| `400`  | Missing or invalid request body or required headers; or Forms Hub normalization failure     |
| `500`  | Unexpected server error, including Forms Hub call failure                                   |

---

#### `POST /v1/api/submissions/normalize`

**Headers**

| Header          | Required | Description                                              |
| --------------- | -------- | -------------------------------------------------------- |
| `X-Tenant-ID`   | Yes      | Tenant identifier                                        |
| `Authorization` | Yes      | `Bearer <token>` — PingFederate client credentials token |

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

**Response — `200 OK`**

```json
{
  "message": "Ok!",
  "data": {
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

The `data` field is a nested fact map built from dot-delimited canonical tag key paths (e.g. `applicant.address[].city`). Array segments (`[]`) produce arrays of objects. Leaf values are the submitted field values after canonical mapping. No data is persisted by this endpoint.

**Error Responses**

| Status | Condition                                                                             |
| ------ | ------------------------------------------------------------------------------------- |
| `400`  | Missing/invalid body, missing `X-Tenant-ID`, validation failure, empty `values` array, or invalid form version status |
| `401`  | Missing or invalid bearer token                                                       |
| `500`  | Unexpected server error                                                               |

## Drawbacks

- **Forms Hub availability dependency:** Path B introduces a synchronous runtime dependency on Forms Hub at cart item creation time. If Forms Hub is unavailable or degraded, `POST /api/cart/items` will fail for Form Catalog Items. Standard Catalog Items (Path A) are unaffected. The Request Portal must handle this failure gracefully and surface a meaningful error to the actor.

- **I/O latency on cart creation:** `POST /api/cart/items` for Path B blocks on Forms Hub's inline validation and canonical mapping. Validation, data source lookups, and rule evaluation execute synchronously within the request. This latency is an accepted trade-off for real-time feedback and integration simplicity in this release.

## Rationale and Alternatives

The synchronous normalize integration — `POST /v1/api/submissions/normalize` inline during cart item creation — was chosen for the following reasons:

- **Real-time feedback:** The synchronous model surfaces validation results to the actor immediately at cart creation time. An actor submitting a Form Catalog Item with invalid or incomplete field values receives a `400` response in the same request, rather than discovering the failure asynchronously after the cart item has already been created.

- **First-release simplicity:** Eliminating async infrastructure — a broker consumer, outbox events, pending state management, and a result receiver — significantly reduces the implementation surface for both teams in the initial release of both products. Both the Request Portal and Forms Hub can deliver and validate the integration without a message broker dependency.

- **Reusability:** Forms Hub is purpose-built for form definition, validation, and canonical mapping. The legacy tooling it replaces undergoes approximately 100 changes per month, making a dedicated, independently deployable forms service a strategic necessity rather than an implementation convenience. Delegating form concerns to Forms Hub avoids duplicating this complexity in the Request Portal.

- **Separation of concerns:** Form validation, data source lookup, and canonical tag mapping are owned by Forms Hub. Embedding this logic in the Request Portal would couple two distinct domains and create a maintenance burden as form requirements evolve.

### Alternatives Considered

**Asynchronous submission via `POST /v1/api/submissions`**

An asynchronous integration was considered — the Request Portal calls `POST /v1/api/submissions`, which returns `202 Accepted` with a `referenceId`, and Forms Hub processes the submission in a background worker before publishing the result via a Kafka outbox event. This would decouple cart item creation from Forms Hub processing latency.

This was not chosen for the first release for two reasons: first, it delays validation feedback to the actor, who would not know whether their form submission was accepted or rejected until the event is processed and delivered; second, it requires both teams to implement async infrastructure (broker consumer, result receiver, pending state management, orphan recovery) that adds significant scope to an initial integration.

**Polling for submission status**

Rather than consuming result events, the Request Portal could poll `GET /v1/api/submissions/by-reference/{referenceId}/status` until a terminal status is reached. This would avoid implementing an event consumer.

This was not chosen because it still defers validation feedback to the actor, requires the Request Portal to manage polling intervals and timeouts, and provides no simplicity advantage over the synchronous normalize approach.

### Impact of Not Integrating

Without this integration, Form Catalog Items cannot be supported in the request cart. There is no current workaround — this is a feature gap that blocks the Request Portal from consuming form-based catalog items entirely.

## Unresolved Questions

1. **Outbound authentication from Request Portal to Forms Hub**: The mechanism by which the Request Portal authenticates its calls to `POST /v1/api/submissions/normalize` is not yet decided. Candidates include service-to-service JWT (client credentials flow via PingFederate), mTLS, or an API key. The chosen approach must align with the enterprise security posture and the authentication infrastructure available to the Request Portal.

## Future Possibilities

- **Asynchronous submission for high-latency forms:** If future form versions introduce long-running validation steps (e.g. external data source lookups with significant latency), a migration to the asynchronous `POST /v1/api/submissions` path — with broker-based result delivery — could be adopted without changing the `POST /api/cart/items` interface. The `normalize` endpoint would remain available for forms where synchronous latency is acceptable.

- **Submission amendments:** An actor may need to update form field values after a Form Catalog Item has been added to the cart. Under the synchronous model, an amendment is handled by re-calling `POST /v1/api/submissions/normalize` with the updated values and replacing the `canonicalData` on the cart item. The design of the cart item update endpoint must be defined.
