# Autonomous IAM Service Desk Agent Workflow (MCP + RAG)

##### An architectural specification for conversational entitlement intake, dynamic form retrieval, pre-submission validation, and canonical cart fulfillment using Forms Hub and MCP.

---

## 1. Executive Summary

This document specifies the end-to-end design for an **Autonomous IAM / IT Service Desk Agent** that uses the **Forms Hub MCP Gateway** in conjunction with a **RAG-indexed Service Catalog** to guide users conversationally through form-backed access requests. 

By combining:
1. **RAG / Catalog Search:** Resolving natural-language user intent into specific catalog items and form versions.
2. **Forms Hub MCP Gateway:** Fetching live form schemas and executing dry-run validation (`POST /normalize`).
3. **Canonical Normalization (ADR-008 & RFC-0001):** Decoupling conversational presentation from downstream IAM facts.
4. **Request Service Integration:** Submitting validated cart items carrying certified `canonicalData` and `rawSubmission`.

The system allows employees to request complex, high-compliance entitlements (e.g., database write access, production deployment permissions) conversationally via Slack, Teams, or CLI without sacrificing schema validation, audit logging, or downstream data integrity.

---

## 2. System Architecture

```
 ┌────────────────────────────────────────────────────────────────────────┐
 │                              AI AGENT                                  │
 │                                                                        │
 │  - System Prompt / SOP: Guides 5-step lifecycle                        │
 │  - RAG Index: Entitlement embeddings & catalog metadata                │
 └──────┬──────────────────────┬────────────────────────┬─────────────┬───┘
        │ 1. Search            │ 2. Fetch Form          │ 3. Q&A      │ 4. Dry-Run
        ▼                      ▼                        ▼             ▼
 ┌──────────────┐     ┌────────────────┐     ┌──────────┐     ┌─────────────┐
 │   Service    │     │   Forms Hub    │     │   User   │     │  Forms Hub  │
 │   Catalog    │     │   MCP Gateway  │     │ (Chat)   │     │  Normalize  │
 └──────────────┘     └────────────────┘     └──────────┘     └──────┬──────┘
                                                                     │ 5. Cart
                                                                     ▼
                                                              ┌─────────────┐
                                                              │ IAM Request │
                                                              │   Portal    │
                                                              └─────────────┘
```

---

## 3. The 5-Step Execution Sequence

### **Step 1: Catalog Search & Intent Resolution (RAG)**
* **Actor:** User prompts: *"I need temporary access to the production customer database for incident triage."*
* **Agent Action:** Queries semantic catalog vector database.
* **Outcome:** Matches catalog item `AWS RDS Production Access` with:
  * `catalogItemId`: `"item_aws_rds_prod"`
  * `formId`: `"form_01a173c8-7fef-78c2-927f-58617134ca73"`
  * `versionId`: `"ver_01a173c9-0123-7a7c-ad31-37f91b195e47"`

---

### **Step 2: Fetch Active Form Definition (via MCP Gateway)**
* **Agent Action:** Invokes MCP tool `forms-api__get_api_v1_forms_formId_versions_versionId`:
  ```json
  {
    "path": {
      "formId": "form_01a173c8-7fef-78c2-927f-58617134ca73",
      "versionId": "ver_01a173c9-0123-7a7c-ad31-37f91b195e47"
    }
  }
  ```
* **Gateway Output:** Returns the complete page/section/element hierarchy, including validation rules and constraints:
  * Element `elem_jira`: Text, required, pattern `^INC-\d+$`
  * Element `elem_duration`: Number, required, min 1, max 8
  * Element `elem_reason`: Text, required, minLength 20

---

### **Step 3: Conversational User Interview**
* **Agent Action:** Inspects the form schema and identifies what required values have not yet been provided.
* **Conversation Turn:**
  * *Agent:* "I can prepare your request for AWS RDS Production Access. To proceed, please provide:
    1. The incident ticket number (e.g. INC-1234).
    2. The duration of access in hours (maximum 8).
    3. A business justification (at least 20 characters)."
  * *User:* "Incident is INC-4091, need it for 3 hours to inspect database locking issues."

---

### **Step 4: Dry-Run Submission Normalization (`POST /normalize`)**
* **Agent Action:** Before submitting to the cart, the agent executes synchronous normalization:
  ```json
  {
    "name": "forms-api__post_api_v1_submissions_normalize",
    "arguments": {
      "body": {
        "formId": "form_01a173c8-7fef-78c2-927f-58617134ca73",
        "versionId": "ver_01a173c9-0123-7a7c-ad31-37f91b195e47",
        "values": [
          { "elementId": "elem_jira", "value": "INC-4091" },
          { "elementId": "elem_duration", "value": 3 },
          { "elementId": "elem_reason", "value": "Inspect database locking issues" }
        ]
      }
    }
  }
  ```
* **Validation & Self-Correction:**
  * If a rule fails (e.g., duration was 12, exceeding max 8), the API returns a 400 error via MCP with `isError: true`. The agent reads the error and asks the user to adjust the value without failing the session.
* **Success Output:** Returns `canonicalData` (mapped semantic tags) and `rawSubmission`:
  ```json
  {
    "canonicalData": {
      "iam.entitlement.duration": 3,
      "compliance.ticket_id": "INC-4091",
      "compliance.justification": "Inspect database locking issues"
    }
  }
  ```

---

### **Step 5: Cart Creation & IAM Submission**
* **Agent Action:** Forwards the validated package to the Request Service (as specified in RFC-0001):
  ```json
  {
    "catalogItemId": "item_aws_rds_prod",
    "type": "form",
    "canonicalData": {
      "iam.entitlement.duration": 3,
      "compliance.ticket_id": "INC-4091",
      "compliance.justification": "Inspect database locking issues"
    },
    "rawSubmission": [
      { "elementId": "elem_jira", "value": "INC-4091" },
      { "elementId": "elem_duration", "value": 3 },
      { "elementId": "elem_reason", "value": "Inspect database locking issues" }
    ]
  }
  ```
* **Outcome:** The Request Service links the item to the user's cart and initiates approval workflows.

---

## 4. Security & Least-Privilege Policy Configuration

To ensure safety in production, the MCP Gateway policy for an IAM Agent should be scoped to **read-only and normalization tools only**:

```json
{
  "name": "forms-api",
  "policy": {
    "allow": [
      "GET /api/v1/forms/**",
      "POST /api/v1/submissions/normalize"
    ],
    "deny": [
      "POST /api/v1/forms/**",
      "PUT /**",
      "DELETE /**"
    ]
  }
}
```
* The agent cannot mutate forms, delete resources, or publish schemas.
* All requests automatically carry the tenant context (`X-Tenant-ID`) and Bearer token injected by the gateway.
