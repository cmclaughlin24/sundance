<div align="center">

# Forms Hub Architecture

##### A high-level overview of the architecture of Forms Hub, a platform for building and rendering forms.

</div>

## 1. Introduction and Goals

This document describes Forms Hub, a multi-tenant Software as a Service (SaaS) system designed to standardize how forms are built, rendered, validated, and consumed across enterprise workflows. The system uses a metadata-driven approach to form definition, enabling forms to be created and modified independently of application code.

### 1.1 Requirements Overview

#### 1.1.1 Business Problem

Enterprise organizations often rely on forms to initiate and drive operational workflows, such as access requests, service requests, and approvals. Over time, multiple form-building solutions have evolved independently to support these needs, leading to a fragmented landscape and one-off implementations to meet specific requirements.

At Wells Fargo, current solutions include:

- **Access Request Tool (ART)**: A legacy suite of applications for managing the request, approval, and fulfillment of requests via automated and manual processes. ART includes a form builder that allows users to create forms for various requests. However, the form builder is tightly coupled to requests, which leads to several issues:
  1. **Embedded Business Logic**: The form builder embeds business logic (service request items or SRIs) directly into the form definition, making it difficult to reuse logic across different forms and leading to duplication of logic.
  2. **Tight Coupling**: Changes to the form definitions impact the downstream processing logic, making it difficult to evolve forms independently of the underlying workflow or integration contracts.

- **WorkX**: A newer platform designed and developed to support Consumer Lending workflows. WorkX includes a form builder that allows users to create forms that are associated with specific tasks in a workflow. In WorkX, a form definition is embedded in the task definition, which reduces the reusability of forms and leading to duplication of form definitions across different tasks.

#### 1.1.2 Goals

| Goal                          | Description                                                                                                                                                              |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Centralized Form Management   | Provide a centralized system for form design and management, enabling users to create, modify, and manage forms in a single location.                                    |
| Dynamic Form Rendering        | Provide a low-code form rendering without requiring custom frontend or backend development per form.                                                                     |
| Configurable Validation       | Support configurable validation rules to ensure data quality and consistency; allow users to define validation rules that can be applied to form fields and submissions. |
| Canonical Data Transformation | Normalize submission data into canonical structures independent of form layout or field naming, enabling downstream systems to consume data in a consistent format.      |
| Decoupled Form Design         | Decouple form design from downstream integration contracts, allowing forms to evolve independently of the underlying workflow or integration requirements.               |

#### 1.1.3 Key Use Cases

| Use Case               | Actor(s)          | Description                                                                                                                                                                                          |
| ---------------------- | ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Design a Form          | Form Designer     | A form designer creates a new form version with pages, sections, and fields. She configures validation rules for form elements, and publishes the form before rendering.                             |
| Render a Form          | End User          | An end user loads an active form version; the system evaluates visibility, required, and read-only rules dynamically and renders the form accordingly.                                               |
| Submit a Form          | End User          | An end user submits a form; the system validates the submission, persists it idempotently, and returns a reference identifier for tracking.                                                          |
| Consume Canonical Data | Downstream System | A downstream system consumes canonical submission data via a message broker or REST API; the system transforms the submission data into a canonical format and delivers it to the downstream system. |
| Manage Lookup Data     | Form Designer     | A form designer configures a data source to provide dynamic lookup data for form fields across one or more forms.                                                                                    |

#### 1.1.4 Functional Requirements

| Requirement                   | Description                                                                                                                                                                         |
| ----------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Tenant Management             | The system shall support multiple isolated tenants. Each tenant's forms, submissions, tags, and data sources shall be scoped and inaccessible to other tenants.                     |
| Form Definition               | The system shall allow users to define forms as versioned, metadata-driven schemas composed of pages, sections, and typed fields without requiring application code changes.        |
| Form Versioning               | The system shall support a managed lifecycle for form versions: `draft → active → retired`. Only one version may be active at a time.                                               |
| Dynamic Form Rendering        | The system shall render forms dynamically from their metadata definition without requiring custom frontend or backend development per form.                                         |
| Configurable Validation       | The system shall support configurable visibility, required, and read-only rules on pages, sections, and fields, evaluated at runtime against submitted values.                      |
| Submission Intake             | The system shall accept form submissions asynchronously and idempotently, associating each with a tenant, form version, and a unique reference ID for external tracking.            |
| Canonical Data Transformation | The system shall map submitted field values to canonical tag versions, normalizing data into consistent structures independent of form layout or field naming.                      |
| Data Source Management        | The system shall support configurable lookup data sources (static, scheduled HTTP, webhook, and data lake) to populate dynamic select fields across forms.                          |
| Tag Management                | The system shall support definition and versioning of canonical semantic tags (`draft → active → deprecated → retired`) that field values are mapped to for downstream consumption. |

### 1.2 Quality Goals

### 1.3 Stakeholders

## 3. Context & Scope

### 3.1 System Context

Forms Hub is a multi-tenant SaaS platform that sits at the center of form design, rendering, submission intake, and canonical data delivery. The diagram below shows the system boundary and all external actors and systems that interact with it.

```mermaid
C4Context
  title System Context — Forms Hub

  Person(formDesigner, "Form Designer", "Creates and manages tenants, data sources, forms, tags, and validation rules.")
  Person(endUser, "End User", "Renders active forms and submits responses.")
  Person_Ext(downstreamSystem, "Downstream System", "Subscribes to and consumes canonical submission events.")

  System(frontend, "Frontend Application", "Web UI through which Form Designers and End Users interact with Forms Hub.")

  System_Boundary(formsHub, "Forms Hub") {
    System(tenantsService, "Tenants Service", "Manages tenant identities and data sources that supply dynamic lookup data for form fields.")
    System(formsService, "Forms Service", "Manages form definitions, versioning, submissions, and canonical tag mappings.")
  }

  System_Ext(messageBroker, "Message Broker", "Receives canonical submission events published by Forms Hub for downstream consumption (e.g. Kafka).")
  System_Ext(mongodb, "MongoDB", "Primary datastore for all domain data across both services.")
  System_Ext(redis, "Redis", "Distributed cache for lookup data and leader election locking for background workers.")
  System_Ext(pingfederate, "PingFederate", "OAuth2 identity provider. Forms Hub validates inbound JWT bearer tokens against PingFederate's JWKS endpoint.")
  System_Ext(externalHTTPAPIs, "External HTTP APIs", "Third-party endpoints polled by the Tenants Service to refresh scheduled and webhook data source lookups.")
  System_Ext(bigquery, "Google BigQuery", "Data lake queried by the Tenants Service to resolve data lake lookup sources.")

  Rel(formDesigner, frontend, "Uses")
  Rel(endUser, frontend, "Uses")
  Rel(frontend, tenantsService, "Manages tenants and data sources", "REST/JSON HTTPS")
  Rel(frontend, formsService, "Manages forms, tags; renders and submits forms", "REST/JSON HTTPS")
  Rel(formsService, messageBroker, "Publishes canonical submission events", "async")
  Rel(downstreamSystem, messageBroker, "Subscribes to submission events", "async")
  Rel(tenantsService, mongodb, "Reads/writes tenant and data source data", "MongoDB wire protocol")
  Rel(formsService, mongodb, "Reads/writes form, submission, and tag data", "MongoDB wire protocol")
  Rel(tenantsService, redis, "Caches lookup data; distributed leader election", "Redis protocol")
  Rel(formsService, redis, "Distributed leader election for submission worker", "Redis protocol")
  Rel(frontend, pingfederate, "Authenticates users, obtains JWT", "OAuth2/HTTPS")
  Rel(tenantsService, pingfederate, "Validates JWT bearer tokens", "JWKS/HTTPS")
  Rel(formsService, pingfederate, "Validates JWT bearer tokens", "JWKS/HTTPS")
  Rel(tenantsService, externalHTTPAPIs, "Fetches lookup data for scheduled and webhook sources", "HTTP/HTTPS")
  Rel(tenantsService, bigquery, "Queries lookup data for data lake sources", "BigQuery API")
```

### 3.2 External Interfaces

| External System          | Direction | Protocol                     | Initiator       | Purpose                                                                                                                                                                                                                   | Status              |
| ------------------------ | --------- | ---------------------------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| **Frontend Application** | Inbound   | REST/JSON over HTTPS         | Frontend        | Form Designers and End Users interact with both services through a generic web UI.                                                                                                                                        | Active              |
| **MongoDB**              | Outbound  | MongoDB wire protocol        | Both services   | Primary datastore. Tenants Service stores tenant and data source records; Forms Service stores forms, form versions, submissions, and tags.                                                                               | Active              |
| **Redis**                | Outbound  | Redis protocol               | Both services   | Tenants Service caches resolved lookup data. Both services use Redis-backed distributed locking (`SetNX` + Lua scripts) for background worker leader election.                                                            | Active              |
| **PingFederate**         | Outbound  | JWKS over HTTPS              | Both services   | Both services validate inbound JWT bearer tokens by fetching signing keys from PingFederate's JWKS URI. Audience, issuer, expiry, and issued-at claims are verified.                                                      | Active              |
| **External HTTP APIs**   | Outbound  | HTTP/HTTPS                   | Tenants Service | The Tenants Service calls arbitrary third-party HTTP endpoints to fetch lookup key-value pairs for `scheduled` and `webhook` data sources. Supports `GET`, `POST`, `PUT`, and `PATCH` with configurable headers and body. | Active              |
| **Google BigQuery**      | Outbound  | BigQuery API                 | Tenants Service | The Tenants Service queries a configured data lake to resolve lookup data for `data-lake` data sources, using configurable catalog, schema, query, and field mappings.                                                    | Planned (stub only) |
| **Message Broker**       | Outbound  | Async messaging (e.g. Kafka) | Forms Service   | After a submission is accepted and canonical tag mapping is applied, the Forms Service publishes a canonical submission event for downstream consumption.                                                                 | Planned             |
