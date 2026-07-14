<div align="center">

# RFC-0001 Forms Hub & Request Portal Integration

##### A proposal that describes a change, its motivation, and its implications.

</div>

---

- **Start Date:** 2026-07-08
- **Status:** Draft
- **Author(s):** Wm. Curtis McLaughlin
- **Related:** [ADR-003](../adr/003-asynchronous-submission-processing.md), [ADR-009](../adr/009-outbox-pattern-polling-vs-streaming.md)

---

## Summary

This RFC proposes the end-to-end integration between the Request Portal and Forms Hub to support form-based catalog items. The integration spans two phases: a frontend phase in which the Request Portal loads the Forms Hub Micro Frontend (MFE) and mounts a `<CyberFormElement>` web component to render the form and collect user input, and a backend phase in which the Request Portal creates a cart item carrying the canonical fact map and raw submission returned by the Form Component.

In the frontend phase, the Form Component fetches the published form definition from the Forms Service, renders the form, calls `POST /api/v1/submissions/normalize` directly to validate and normalize the submission, and notifies the host application with the resulting `canonicalData` and `rawSubmission`. In the backend phase, the Request Portal forwards both to `POST /api/cart/items` on the Request Service, which links them to the cart item and proceeds to policy evaluation.

This integration introduces a host/guest MFE contract between the two frontend applications and a new execution path in the Request Portal's cart item creation handler, while requiring no changes to the Forms Hub backend.

## Motivation

The Request Portal currently supports a single cart item type — the Standard Catalog Item — which requires no coordination with an external service at the time of cart creation. As the platform evolves to support form-based catalog items, two new capabilities are needed: a form rendering experience embedded within the Request Portal, and a mechanism to validate and normalize form submission data before a cart item is finalized.

Forms Hub addresses both needs. On the frontend, Forms Hub exposes an MFE bundle that the Request Portal can load on demand, mounting a `<CyberFormElement>` web component to render any active form version without the Request Portal implementing any form rendering logic. On the backend, Forms Hub exposes `POST /api/v1/submissions/normalize` — a synchronous endpoint that validates a submission against a form version, applies canonical tag mapping, and returns a canonical fact map inline without persisting any data. The Form Component calls this endpoint directly, providing immediate validation feedback to the actor before notifying the host.

Without this integration, Form Catalog Items cannot be supported in the request cart. With it, the Request Portal gains embedded form rendering through a decoupled MFE contract, and the ability to attach both structured canonical data and the raw submission to cart items — enabling downstream policy evaluation and fulfillment workflows to operate on a trusted, normalized representation of the request.

## Design

The integration between the Request Portal and Forms Hub follows a two-phase model. In the frontend phase, the Request Portal loads the Forms Hub MFE bundle on demand and mounts a `<CyberFormElement>` web component, which independently fetches the form definition, renders the form, validates and normalizes the submission via `POST /api/v1/submissions/normalize`, and notifies the host application with the result. In the backend phase, the Request Portal forwards the notification payload to `POST /api/cart/items` on the Request Service, which links the `canonicalData` and `rawSubmission` to the cart item and proceeds to policy evaluation. No data is persisted in Forms Hub at any point in this flow.

### Sequence

![Forms Hub & Request Portal Integration](../imgs/Forms%20Hub%20x%20Request%20Portal%20Integration.png)

1. **Select form-backed catalog item**: The actor selects a form-backed catalog item in the Request Portal.

2. **Load Forms Hub MFE bundle**: The Request Portal loads the Forms Hub MFE bundle. The bundle is returned and made available to the host application.

3. **Initialize Form Component**: The Request Portal mounts `<CyberFormElement formId versionId />`, passing the form and version identifiers for the selected catalog item.

4. **Load published form definition**: The Form Component calls `GET /api/v1/form/{formId}/version/{versionId}` on the Forms Service to fetch the published form definition.

5. **Render form and collect user input**: The Form Component renders the form from its metadata definition and collects field values from the actor.

6. **Validate and normalize submission**: The Form Component calls `POST /api/v1/submissions/normalize` on the Forms Service, forwarding the `formId`, `versionId`, and collected `values`. The Forms Service evaluates visibility and required rules, validates each resolved field, and maps field values to canonical tags. If validation fails, the Form Component surfaces the error to the actor inline; the flow does not proceed. On success, the Forms Service returns a canonical fact map.

7. **Notify host application**: The Form Component fires a host notification event carrying the `canonicalData` (canonical fact map) and `rawSubmission` (the collected field values) to the Request Portal.

8. **Create cart item**: The Request Portal calls `POST /api/cart/items` on the Request Service, forwarding the `canonicalData`, `rawSubmission`, and catalog item metadata. The Request Service links both to the cart item, persists the record, expands the request, and evaluates policy checks.

### Changes Required

Explicit list of behavioral or structural changes required by each system involved. This is the consumer/provider obligation checklist.

#### Request Portal

- **Load the Forms Hub MFE bundle on demand**: When an actor selects a form-backed catalog item, the Request Portal must dynamically load the Forms Hub MFE bundle (see Unresolved Questions for bundle delivery mechanism).

- **Mount `<CyberFormElement>` with required props**: The Request Portal must mount the `<CyberFormElement>` web component, passing `formId`, `versionId`, and `tenantId` for the selected catalog item.

- **Handle the host notification event**: The Request Portal must listen for the host notification event fired by the Form Component on successful submission. The event payload carries `canonicalData` and `rawSubmission`; both must be captured before proceeding to cart item creation.

- **Define `form` as a new value for the existing `items[].type` discriminator**: `POST /api/cart/items` already accepts a `type` field on each item. The Request Portal must define `form` as a new valid value alongside existing types, routing execution into the form-backed path when present.

- **Define the form-backed `items[].payload` contract**: When `items[].type` is `form`, the payload must carry `canonicalData` and `rawSubmission` forwarded from the Form Component notification event, along with any catalog item metadata required by the Request Service.

- **Extend the cart item data model**: The cart item data model must be extended with two new fields: `canonicalData` — the canonical fact map produced by the Forms Service — and `rawSubmission` — the raw field values as submitted by the actor (`[{ fieldId, value, collectionIndex }]`).

#### Forms Hub

- **Expose the MFE bundle**: Forms Hub must expose a loadable MFE bundle that the Request Portal can dynamically load without authentication (see Unresolved Questions for bundle delivery mechanism).

- **Implement `<CyberFormElement>` web component**: The MFE must expose a `<CyberFormElement>` web component that accepts `formId`, `versionId`, and `tenantId` as props.

- **Fetch the published form definition**: On mount, the Form Component must call `GET /api/v1/form/{formId}/version/{versionId}` on the Forms Service to load the active form definition, passing `tenantId` as the `X-Tenant-ID` header.

- **Render the form and collect user input**: The Form Component must render the form from its metadata definition and collect field values from the actor, applying visibility and required rules dynamically.

- **Call `POST /api/v1/submissions/normalize` client-side**: On actor submission, the Form Component must call `POST /api/v1/submissions/normalize` on the Forms Service, forwarding `formId`, `versionId`, `tenantId` as `X-Tenant-ID`, and the collected `values`. Validation errors must be surfaced inline to the actor; the host is not notified on failure.

- **Fire host notification event on success**: On a successful normalize response, the Form Component must fire a host notification event to the Request Portal carrying `canonicalData` (the canonical fact map) and `rawSubmission` (`[{ fieldId, value, collectionIndex }]`). The event name is unresolved — see Unresolved Questions.

> **Note:** No changes are required to the Forms Hub backend. `POST /api/v1/submissions/normalize` and `GET /api/v1/form/{formId}/version/{versionId}` are already implemented.

### API Contracts

#### `GET /api/v1/form/{formId}/version/{versionId}`

Called by the Form Component on mount to fetch the published form definition.

**Headers**

| Header          | Required | Description                                                              |
| --------------- | -------- | ------------------------------------------------------------------------ |
| `X-Tenant-ID`   | Yes      | Tenant identifier, passed from the `tenantId` prop                       |
| `Authorization` | Yes      | `Bearer <token>` — sourced from the authenticated user's session context |

**Response — `200 OK`**

The response schema is defined by Forms Hub. Returns the full form version definition including pages, sections, fields, and validation rules.

**Error Responses**

| Status | Condition                       |
| ------ | ------------------------------- |
| `400`  | Missing `X-Tenant-ID`           |
| `401`  | Missing or invalid bearer token |
| `404`  | Form or version not found       |
| `500`  | Unexpected server error         |

---

#### `POST /api/v1/submissions/normalize`

Called by the Form Component after the actor submits the form. Validates and normalizes the submission inline; no data is persisted.

**Headers**

| Header          | Required | Description                                                              |
| --------------- | -------- | ------------------------------------------------------------------------ |
| `X-Tenant-ID`   | Yes      | Tenant identifier, passed from the `tenantId` prop                       |
| `Authorization` | Yes      | `Bearer <token>` — sourced from the authenticated user's session context |

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

| Status | Condition                                                                                                             |
| ------ | --------------------------------------------------------------------------------------------------------------------- |
| `400`  | Missing/invalid body, missing `X-Tenant-ID`, validation failure, empty `values` array, or invalid form version status |
| `401`  | Missing or invalid bearer token                                                                                       |
| `500`  | Unexpected server error                                                                                               |

---

#### Host Notification Event

Fired by the Form Component to the Request Portal on successful normalization. The event name is unresolved — see Unresolved Questions.

**Payload**

```json
{
  "canonicalData": {
    "applicant": {
      "firstName": "Jane",
      "lastName": "Doe"
    }
  },
  "rawSubmission": [
    {
      "fieldId": "<string>",
      "value": "<any>",
      "collectionIndex": "<int | omitted>"
    }
  ]
}
```

---

#### `POST /api/cart/items`

**Headers**

| Header              | Required | Description                      |
| ------------------- | -------- | -------------------------------- |
| `X-WF-REQUEST-DATE` | Yes      | Request date in ISO 8601 format  |
| `X-REQUEST-DATE`    | Yes      | Unique request identifier (UUID) |
| `X-CORRELATION-ID`  | Yes      | Correlation identifier (UUID)    |
| `X-WF-CLIENT-ID`    | Yes      | Client identifier (e.g. `iamx`)  |

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
        "canonicalData": { "<key>": "<any>" },
        "rawSubmission": [
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

**Response — `200 OK`**

The response schema (`AddCartItemsResponse`) is defined by Request Portal.

**Error Responses**

| Status | Condition                                           |
| ------ | --------------------------------------------------- |
| `400`  | Missing or invalid request body or required headers |
| `500`  | Unexpected server error                             |

## Drawbacks

- **MFE bundle load latency:** Loading the Forms Hub MFE bundle on demand adds latency to the form-backed catalog item selection flow. The bundle must be fetched and parsed before the Form Component can be mounted and the form definition loaded. This is an accepted trade-off for keeping form rendering logic out of the Request Portal.

- **Forms Hub availability dependency:** The form-backed catalog item flow introduces a synchronous runtime dependency on both the Forms Hub MFE and the Forms Service. If the MFE bundle is unavailable, the form cannot be rendered. If the Forms Service is unavailable or degraded, normalization will fail and the cart item cannot be created. The Request Portal must handle both failure modes gracefully and surface meaningful errors to the actor.

- **Client-side normalization latency:** The Form Component calls `POST /api/v1/submissions/normalize` synchronously within the actor's browser before notifying the host. Validation, data source lookups, and rule evaluation execute within this request. This latency is an accepted trade-off for real-time inline feedback to the actor before cart item creation.

- **Host/guest MFE coupling:** The Request Portal and the Forms Hub MFE are coupled through the host notification event contract — its name, payload shape, and firing conditions. Changes to the event contract require coordination between both teams. The event name is currently unresolved (see Unresolved Questions), and its shape must be treated as a versioned interface once agreed.

- **Client-side bearer token dependency:** Both `GET /api/v1/form/{formId}/version/{versionId}` and `POST /api/v1/submissions/normalize` require a bearer token sourced from the authenticated user's session context. The Form Component depends on the host application having an active, valid session at the time of form rendering and submission. If the session has expired, both calls will fail with `401`.

## Rationale and Alternatives

This design — a two-phase MFE + synchronous normalize integration — was chosen for the following reasons:

- **Real-time feedback:** The synchronous normalize call within the Form Component surfaces validation results to the actor inline, before the host is notified and before the cart item is created. An actor submitting a form with invalid or incomplete field values receives immediate feedback without leaving the form, rather than discovering the failure after cart creation.

- **Decoupled form rendering:** The MFE approach keeps all form rendering, field collection, and validation logic inside Forms Hub. The Request Portal mounts a single web component and handles a single notification event — it has no knowledge of form structure, field types, or validation rules. This boundary allows both teams to evolve their respective surfaces independently.

- **First-release simplicity:** Eliminating async infrastructure — a broker consumer, outbox events, pending state management, and a result receiver — significantly reduces the implementation surface for both teams in the initial release. Both the Request Portal and Forms Hub can deliver and validate the integration without a message broker dependency.

- **Reusability:** Forms Hub is purpose-built for form definition, validation, and canonical mapping. The legacy tooling it replaces undergoes approximately 100 changes per month, making a dedicated, independently deployable forms service a strategic necessity rather than an implementation convenience. Delegating form concerns to Forms Hub avoids duplicating this complexity in the Request Portal.

- **Separation of concerns:** Form validation, data source lookup, and canonical tag mapping are owned by Forms Hub. Embedding this logic in the Request Portal would couple two distinct domains and create a maintenance burden as form requirements evolve.

### Alternatives Considered

**Server-side normalize — Request Portal backend calls `POST /api/v1/submissions/normalize`**

The original design proposed in this RFC had the Request Portal backend call `POST /api/v1/submissions/normalize` server-to-server during `POST /api/cart/items`, with the actor submitting form values directly to the Request Portal. This avoided an MFE dependency and kept the frontend integration surface minimal.

This was not chosen because it requires the Request Portal to own form rendering — collecting, validating, and forwarding field values — without any knowledge of the form structure. It also defers validation feedback: the actor submits the entire cart item payload before discovering whether the form values are valid, rather than receiving inline feedback within the form itself.

**Asynchronous submission via `POST /api/v1/submissions`**

An asynchronous integration was considered — the Form Component or Request Portal calls `POST /api/v1/submissions`, which returns `202 Accepted` with a `referenceId`, and Forms Hub processes the submission in a background worker before publishing the result via a Kafka outbox event.

This was not chosen for the first release for two reasons: first, it delays validation feedback to the actor, who would not know whether their form submission was accepted or rejected until the event is processed and delivered; second, it requires both teams to implement async infrastructure (broker consumer, result receiver, pending state management, orphan recovery) that adds significant scope to an initial integration.

**Polling for submission status**

Rather than consuming result events, the Request Portal could poll `GET /api/v1/submissions/by-reference/{referenceId}/status` until a terminal status is reached. This would avoid implementing an event consumer.

This was not chosen because it still defers validation feedback to the actor, requires the Request Portal to manage polling intervals and timeouts, and provides no simplicity advantage over the synchronous normalize approach.

### Impact of Not Integrating

Without this integration, Form Catalog Items cannot be supported in the request cart. There is no current workaround — this is a feature gap that blocks the Request Portal from consuming form-based catalog items entirely.

## Unresolved Questions

1. **Host notification event name**: The name of the event fired by the Form Component to the Request Portal has not been agreed. The event name must be stable once set, as it forms part of the versioned host/guest contract between the two applications.

2. **MFE bundle delivery mechanism**: The mechanism by which the Request Portal discovers and loads the Forms Hub MFE bundle is not yet decided. Candidates include a well-known static CDN URL, a module federation manifest, or a configuration value supplied to the Request Portal at deploy time. The chosen approach affects how the Request Portal resolves the bundle URL and how Forms Hub manages bundle versioning and cache invalidation.

3. **Session expiry handling in the Form Component**: The Form Component calls both `GET /api/v1/form/{formId}/version/{versionId}` and `POST /api/v1/submissions/normalize` using the authenticated user's session token. If the session expires mid-flow, both calls will fail with `401`. The Form Component's behavior on a `401` response — whether to surface an error, attempt a token refresh, or delegate session recovery to the host — has not been defined.

## Future Possibilities

- **Asynchronous submission for high-latency forms:** If future form versions introduce long-running validation steps (e.g. external data source lookups with significant latency), a migration to the asynchronous `POST /api/v1/submissions` path — with broker-based result delivery — could be adopted without changing the `POST /api/cart/items` interface. The `normalize` endpoint would remain available for forms where synchronous latency is acceptable.

- **Submission amendments:** An actor may need to update form field values after a Form Catalog Item has been added to the cart. Under the current design, an amendment would re-trigger the MFE flow — remounting `<CyberFormElement>` with the existing `rawSubmission` pre-populated — re-calling `POST /api/v1/submissions/normalize` with the updated values, and replacing both `canonicalData` and `rawSubmission` on the cart item. The design of the cart item update endpoint must be defined.
