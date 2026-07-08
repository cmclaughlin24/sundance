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

1. **Create cart item** — The actor submits `POST /api/cart/items` to the Request Portal. The request payload identifies the catalog item type, branching execution into Path A or Path B.

2. **Create form submission** _(Path B only)_ — The Request Portal calls `POST /v1/api/submissions` on Forms Hub. Forms Hub persists the submission with a status of `pending` and returns `202 Accepted` with a submission reference ID. The Request Portal returns this response to the actor immediately.

3. **Associate cart item with submission reference** — The Request Portal stores the submission reference ID against the cart item, linking the two records internally.

4. **Validate submission and produce canonical representation** _(async, Forms Hub internal)_ — A background worker validates the submission, resolves data source lookups, and maps fields to canonical tags. This step is internal to Forms Hub and runs concurrently with step 3.

5. **Publish submission result** — Forms Hub publishes the final outcome (`accepted`, `rejected`, or `failed`) to the Request Portal. The delivery mechanism for this step is unresolved — see Unresolved Questions.

6. **Finalize cart item** — On `accepted`, the Request Portal links the canonical submission data to the cart item. On `rejected` or `failed`, the cart item is marked as invalid.

7. **Expand request and evaluate policy checks** _(on `accepted` only)_ — The Request Portal expands the request and runs policy evaluation against the now-complete cart item.

### Changes Required

Explicit list of behavioral or structural changes required by each
system involved. This is the consumer/provider obligation checklist.

#### Request Portal

- Change 1
- Change 2

#### Forms Hub

- Change 1
- Change 2

### API Contracts

Define the request/response schemas and error codes for each
synchronous API interaction.

### Event Contracts

Define the event shape, status values, and delivery guarantees for
each asynchronous interaction.

## Drawbacks

Why should we _not_ do this? What are the risks, costs, or
complexities introduced by this proposal?

## Rationale and Alternatives

- Why is this design the best in the space of possible designs?
- What other designs were considered and why were they not chosen?
- What is the impact of not doing this?

## Unresolved Questions

1. **Result delivery mechanism (step 5)** — The mechanism by which Forms Hub delivers the submission result to the Request Portal is not yet decided. The long-term goal is delivery via a message broker. A webhook-based approach may be required as a short-term solution pending broker infrastructure availability.

## Future Possibilities

What natural extensions or follow-on work does this proposal enable?
This is not a commitment — just a space to capture related ideas.
