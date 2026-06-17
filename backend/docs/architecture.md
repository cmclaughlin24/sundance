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
