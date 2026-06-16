<div align="center">

# Forms Hub Architecture

##### A high-level overview of the architecture of Forms Hub, a platform for building and rendering forms.

</div>

## 1. Introduction and Goals

### 1.1 Purpose

Forms Hub is a multi-tenant Software as a Service (SaaS) system designed to standardize how data is captured, validated, and consumed across enterprise workflows.

It provides a metadata-driven approach to form definition, enabling forms to be created and modified independently of application code. The system is responsible for redering forms dynamically, validating submissions, and transforming submissions into canonical formats that can be reliably consumed by downstream systems.

The platform is intended to:

1. Provide a centralized system for form design and managing dyanmic forms.
2. Enable form rendering without requiring custom frontend or backend development per form.
3. Support configurable validation rules to ensure data quality and consistency.
4. Normalize submission data into canonical strucctures independent of form layout or field naming.
5. Decouple form design from downstream integration contracts.
6. Reduce duplication of form logic across enterprise applications.
7. Improve consistency, maintainability, and reusability of form-based workflows.

### 1.2 Business Problem

Enterprise organizations often rely on fors to initiate and drive operational workflows, such as access requests, service requests, and approvals. Over time, multiple form-building solutions havbe evolved independently to support these needs, leading to a fragmented landscape and one-off implementations to meet specific requirements.

At Wells Fargo, the current landscape includes:

- **Access Request Tool (A.R.T.)**: A legacy suite of applications for managing the request, approval, and fulfillment of requests via automated and manual process. A.R.T. includes a form builder that allows users to create forms for various requests. However, the form builder is tightly coupled to requests, which leads to several issues:
  1. **Embedded Business Logic**: The form builder embeds business logic (service request items (SRIs)) directly into the form definition, making it difficult to reuse logic across different forms and leading to duplication of logic.
  2. **Tight Coupling**: Changes to the form definitions impact the downstream processing logic, making it difficult to evolve forms independently of the underlying workflow or integration contracts.

- **WorkX**: A newer platform designed and developed to support Consumer Lending workflows. WorkX includes a form builder that allows users to create forms that are associated with specific tasks in a workflow. However, the form builder is tightly coupled to the WorkX task (embedded in the task definition), which reduces the reusability of forms and leads to duplication of form definitions across different tasks.

The Forms Hub initiative aims to address these limitations by providing a centralized, metadata-driven form platform that decouples form definition from business logic execution and downstream integration contracts.

### 1.3 goals

### 1.4 Non-Goals

### 1.5 Stakeholders

### 1.6 Success Criteria

## 2. Constraints

## 3. Context and Scope

## 4. Solution and Strategy

## 5. Building Blocks and Components

## 6. Runtime View

## 7. Deployment View

## 8. Cross-Cutting Concerns

## 9. Architectural Decisions

## 10. Quality Requirements

## 11. Risks and Technical Debt

## 12. Glossary
