# Form Rendering Engine — Implementation Plan

## Overview

The `FormElement` component in `frontend/apps/forms/` is currently a scaffold that fetches form data but renders an empty `<form>`. This document outlines the full implementation plan for a functional rendering engine supporting multi-page wizard navigation, client-side rule evaluation, dynamic data source loading, and submission handling.

---

## Phase 1 — Type System Refinements

### 1a. Strongly-typed element attributes ✓

Refactor `frontend/apps/forms/src/types/element.ts` from a loose `Record<string, any>` to a discriminated union of concrete attribute types matching the backend REST DTOs.

```ts
// text
interface TextElementAttributes {
  isRequired: boolean;
  isReadOnly: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: string;
  placeholder?: string;
}

// number
interface NumberElementAttributes {
  isRequired: boolean;
  isReadOnly: boolean;
  min?: number;
  max?: number;
  step?: number;
}

// select
interface SelectElementAttributes {
  isRequired: boolean;
  isReadOnly: boolean;
  data: any[];
  dataSourceRef?: IDataSourceRef;
  multiple: boolean;
  minSelected?: number;
  maxSelected?: number;
}

// checkbox
interface CheckboxElementAttributes {
  isRequired: boolean;
  isReadOnly: boolean;
  isCheckedByDefault: boolean;
}

// date
interface DateElementAttributes {
  isRequired: boolean;
  isReadOnly: boolean;
  minDate?: string;
  maxDate?: string;
}

type ElementAttributes =
  | TextElementAttributes
  | NumberElementAttributes
  | SelectElementAttributes
  | CheckboxElementAttributes
  | DateElementAttributes;
```

Also remove `boolean` from the `ElementType` union — it exists in the frontend but not in the backend domain. Map it to `checkbox` or remove entirely. ✓

### 1b. New file: `types/dataSource.ts` ✓

```ts
interface ILookup {
  key: string;
  value: string;
}

interface IBindingSource {
  type: "field" | "static";
  key: string;
  value: any;
}

interface IDataSourceRef {
  dataSourceId: string;
  bindings: Record<string, IBindingSource>;
}
```

### 1c. Tighten submission value types

Update `ISubmissionValue.value` in `types/submission.ts` to be typed per element where feasible, rather than `any`.

**Files to modify:**

- `frontend/apps/forms/src/types/element.ts`
- `frontend/apps/forms/src/types/submission.ts`

**Files to create:**

- `frontend/apps/forms/src/types/dataSource.ts`

---

## Phase 2 — Service Layer

### 2a. Implement `DataSourcesService.getLookups()` ✓

Add the missing method to `frontend/apps/forms/src/services/dataSourcesService.ts`:

```ts
async getLookups(
  dataSourceId: string,
  params: Record<string, any>,
  options: DefaultRequestOptions,
): Promise<ILookup[]> {
  return this._get(`/data-sources/${dataSourceId}/look-ups`, options);
}
```

Calls `GET /data-sources/{dataSourceId}/look-ups` on the tenants backend (`http://localhost:8080`).

### 2b. Environment-based API URLs

Replace hardcoded `localhost` URLs in `frontend/apps/forms/src/hooks/useHttpService.ts` with Vite environment variables.

```ts
// Before
resolveHttpService(FormsService, "http://localhost:8081");

// After
resolveHttpService(FormsService, import.meta.env.VITE_FORMS_API_URL);
```

Create `.env` and `.env.example` at `frontend/apps/forms/`:

```
VITE_FORMS_API_URL=http://localhost:8081
VITE_TENANTS_API_URL=http://localhost:8080
```

**Files to modify:**

- `frontend/apps/forms/src/services/dataSourcesService.ts`
- `frontend/apps/forms/src/hooks/useHttpService.ts`

**Files to create:**

- `frontend/apps/forms/.env`
- `frontend/apps/forms/.env.example`

---

## Phase 3 — Form State Architecture

### 3a. New file: `context/FormContext.ts` ✓

Define the context shape, reducer actions, and initial state. _Created as `store/formContext.ts`, `store/formReducer.ts`._

```ts
interface FormState {
  values: Record<string, any>; // elementId → value
  errors: Record<string, string[]>; // elementId → error messages
  ruleStates: Record<
    string,
    {
      // elementId → computed rule output
      visible: boolean;
      required: boolean;
      readonly: boolean;
    }
  >;
}

type FormAction =
  | { type: "SET_VALUE"; elementId: string; value: any }
  | { type: "SET_ERROR"; elementId: string; errors: string[] }
  | { type: "SET_RULE_STATES"; ruleStates: FormState["ruleStates"] }
  | { type: "INITIALIZE"; values: Record<string, any> };
```

### 3b. New file: `context/FormProvider.tsx` ✓

Provider component that:

- Accepts `rawSubmission` prop and initializes state from it via `useReducer` initializer function — no extra render cycle ✓
- Wraps children with state and dispatch contexts ✓
- _Created as `store/FormProvider.tsx`; consumer hooks in `store/useFormContext.ts`_ ✓

**Files to create:**

- `frontend/apps/forms/src/store/formContext.ts` ✓
- `frontend/apps/forms/src/store/formReducer.ts` ✓
- `frontend/apps/forms/src/store/FormProvider.tsx` ✓
- `frontend/apps/forms/src/store/useFormContext.ts` ✓

---

## Phase 4 — Client-Side Rule Evaluator

### 4a. New file: `hooks/useRuleEvaluator.ts`

A pure utility (not reliant on React state/effects) that mirrors the backend's `ExprRuleEvaluator` logic in TypeScript.

Given the current `values` map and a list of elements, for each element:

1. Iterate `element.rules` grouped by `rule.type` (`visible`, `required`, `readonly`)
2. For each rule, iterate `rule.expressions` in ascending `position` order
3. Evaluate `values[expression.fieldKey] <operator> expression.value`
4. Chain expression results using `expression.joinWithPrevious` (`and` / `or`)
5. Return the boolean result for each rule type

**Operators to implement** (matching `RuleExpressionOp`):

| Enum value | Operation |
| ---------- | --------- |
| `equal`    | `===`     |
| `nequal`   | `!==`     |
| `lt`       | `<`       |
| `gt`       | `>`       |
| `lte`      | `<=`      |
| `gte`      | `>=`      |

**Default rule states** (when no rule of that type exists on an element):

- `visible`: `true`
- `required`: use `attributes.isRequired`
- `readonly`: use `attributes.isReadOnly`

The server remains authoritative — client-side evaluation is for UX responsiveness only. The backend re-evaluates all rules during submission processing.

**Files to create:**

- `frontend/apps/forms/src/hooks/useRuleEvaluator.ts`

---

## Phase 5 — Field Components

One component per `ElementType`. All field components:

- Consume `FormContext` via `useContext` to read and write values
- Read `ruleStates[element.id]` to determine visibility, required state, and read-only state
- Return `null` when `ruleStates[element.id].visible === false`
- Dispatch `SET_VALUE` on change

| Component       | MUI Input                     | Key attributes respected                                          |
| --------------- | ----------------------------- | ----------------------------------------------------------------- |
| `TextField`     | `MUI TextField`               | `minLength`, `maxLength`, `pattern`, `placeholder`                |
| `NumberField`   | `MUI TextField type="number"` | `min`, `max`, `step`                                              |
| `SelectField`   | `MUI Select` / `Autocomplete` | `data`, `dataSourceRef`, `multiple`, `minSelected`, `maxSelected` |
| `CheckboxField` | `MUI Checkbox`                | `isCheckedByDefault` (initializes value on mount)                 |
| `DateField`     | `MUI TextField type="date"`   | `minDate`, `maxDate`                                              |

**`SelectField` specifics:**

- If `attributes.data` is non-empty and `attributes.dataSourceRef` is absent, render inline options directly
- If `attributes.dataSourceRef` is present, fetch options via `useDataSourceLookups` hook (Phase 8)
- Show a loading spinner while lookups are fetching
- Show an inline error if the lookup call fails
- Respect `multiple` for multi-select behavior

**Files to create:**

- `frontend/apps/forms/src/components/fields/TextField.tsx`
- `frontend/apps/forms/src/components/fields/NumberField.tsx`
- `frontend/apps/forms/src/components/fields/SelectField.tsx`
- `frontend/apps/forms/src/components/fields/CheckboxField.tsx`
- `frontend/apps/forms/src/components/fields/DateField.tsx`

---

## Phase 6 — Layout Components

### 6a. `ElementRenderer`

Dispatches to the correct field component by `element.type`. Acts as a single switch point — consumers never import field components directly.

```tsx
switch (element.type) {
  case "text":
    return <TextField element={element} />;
  case "number":
    return <NumberField element={element} />;
  case "select":
    return <SelectField element={element} />;
  case "checkbox":
    return <CheckboxField element={element} />;
  case "date":
    return <DateField element={element} />;
}
```

### 6b. `SectionRenderer`

- Sorts `section.elements` by `position` ascending
- Maps each element through `ElementRenderer`
- Evaluates section-level rules from `FormContext` for section visibility
- Renders the section label/title

### 6c. `PageRenderer`

- Sorts `page.sections` by `position` ascending
- Maps each section through `SectionRenderer`
- Evaluates page-level rules for page visibility (skips invisible pages in wizard navigation)

### 6d. Multi-page wizard in `FormElement`

Refactor `FormElement.tsx` to:

- Sort `formVersion.pages` by `position` ascending
- Track `currentPageIndex` in local state
- Render only the current page via `PageRenderer`
- Provide Next and Back buttons with boundary guards
- Skip pages where all elements are invisible (due to rules)
- Show a progress indicator (e.g. "Page 2 of 4")

**Files to create:**

- `frontend/apps/forms/src/components/layout/ElementRenderer.tsx`
- `frontend/apps/forms/src/components/layout/SectionRenderer.tsx`
- `frontend/apps/forms/src/components/layout/PageRenderer.tsx`

**Files to modify:**

- `frontend/apps/forms/src/components/FormElement/FormElement.tsx`

---

## Phase 7 — Submission Flow

### 7a. `useActionState` for submit

Use React 19's `useActionState` in `FormElement` for the async submit action:

```ts
const [state, submitAction, isPending] = useActionState(
  async (prevState, _formData) => {
    // 1. Collect visible element IDs from context.ruleStates
    // 2. Build ISubmissionValue[] from context.values for visible elements only
    // 3. Call submissionsService.submit(formId, versionId, values, options)
    // 4. Return { status: "success" } or { status: "error", message }
  },
  { status: "idle" },
);
```

### 7b. Idempotency key

The backend enforces an `Idempotency-Key` header on `POST /submissions`. Generate a UUID on `FormElement` mount and thread it through to `SubmissionsService.submit()`. Add an optional `idempotencyKey` parameter to the `_post` call or handle it as a dedicated header in `SubmissionsService`.

### 7c. Surface submit errors

Display server-side validation errors returned from the backend (e.g. required field failures, type mismatches from the submission pipeline) back in `FormContext` via `SET_ERROR` dispatches.

### 7d. `onSubmit` callback

After a successful submission, call `props.onSubmit({ raw: values, normalized: normalizedResult })`. The `normalized` value should come from the `POST /submissions` response or a preceding `normalize` call.

**Files to modify:**

- `frontend/apps/forms/src/components/FormElement/FormElement.tsx`
- `frontend/apps/forms/src/services/submissionService.ts`

---

## Phase 8 — DataSources for Select Fields

### 8a. New hook: `useDataSourceLookups`

```ts
function useDataSourceLookups(
  dataSourceRef: IDataSourceRef | undefined,
  formValues: Record<string, any>,
  options: DefaultRequestOptions,
): { lookups: ILookup[]; isLoading: boolean; error: unknown };
```

- Uses `useAsyncData` internally
- If `dataSourceRef` is undefined, returns `{ lookups: [], isLoading: false, error: null }`
- Resolves `field`-type bindings from `dataSourceRef.bindings` against `formValues` to build the query params
- Re-fetches when any referenced field value changes (dependency array includes the resolved binding values)
- Calls `DataSourcesService.getLookups(dataSourceRef.dataSourceId, resolvedParams, options)`

### 8b. Static binding support

For `BindingSource.type === "static"`, pass the literal `value` directly as the query param — no form value lookup needed.

### 8c. Dynamic binding support (field-referenced)

For `BindingSource.type === "field"`, resolve `formValues[bindingSource.key]` as the param value. Include the resolved field values in `useAsyncData`'s dependency array so lookups re-fetch when upstream fields change.

**Files to create:**

- `frontend/apps/forms/src/hooks/useDataSourceLookups.ts`

---

## Phase 9 — Route Integration

### 9a. New route: form viewer

Add a file-based route at `frontend/apps/forms/src/routes/forms/$formId/versions/$versionId.tsx`:

```tsx
export const Route = createFileRoute("/forms/$formId/versions/$versionId")({
  component: FormViewerPage,
});

function FormViewerPage() {
  const { formId, versionId } = Route.useParams();
  return (
    <FormElement
      tenantId={/* from context or route search param */}
      formId={formId}
      versionId={versionId}
      onSubmit={(event) => {
        /* handle */
      }}
    />
  );
}
```

### 9b. Root layout

Update `frontend/apps/forms/src/routes/__root.tsx` to remove the placeholder `<div>Hello "__root"!</div>` and replace with a proper `<Outlet />` only (chrome is handled by the host shell).

**Files to modify:**

- `frontend/apps/forms/src/routes/__root.tsx`

**Files to create:**

- `frontend/apps/forms/src/routes/forms/$formId/versions/$versionId.tsx`

---

## Full File Inventory

### Files to create

| File                                                                   | Purpose                                                                                          |
| ---------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `frontend/apps/forms/src/types/elementAttributes.ts`                   | Discriminated union of all element attribute types ✓                                             |
| `frontend/apps/forms/src/types/dataSource.ts`                          | `ILookup`, `IDataSourceRef`, `IBindingSource`, `HasDataSourceRef` types — created as `data.ts` ✓ |
| `frontend/apps/forms/src/store/formContext.ts`                         | State and dispatch context objects ✓                                                             |
| `frontend/apps/forms/src/store/formReducer.ts`                         | Reducer, actions, `FormState`, `initialFormState`, `initializeForm` ✓                            |
| `frontend/apps/forms/src/store/FormProvider.tsx`                       | Provider component, wires `useReducer` with `rawSubmission` initializer ✓                        |
| `frontend/apps/forms/src/store/useFormContext.ts`                      | `useFormState` and `useFormDispatch` consumer hooks ✓                                            |
| `frontend/apps/forms/src/hooks/useRuleEvaluator.ts`                    | Client-side rule evaluator utility                                                               |
| `frontend/apps/forms/src/hooks/useDataSourceLookups.ts`                | Async lookup fetcher for select fields                                                           |
| `frontend/apps/forms/src/components/fields/TextField.tsx`              | Text field component                                                                             |
| `frontend/apps/forms/src/components/fields/NumberField.tsx`            | Number field component                                                                           |
| `frontend/apps/forms/src/components/fields/SelectField.tsx`            | Select field component                                                                           |
| `frontend/apps/forms/src/components/fields/CheckboxField.tsx`          | Checkbox field component                                                                         |
| `frontend/apps/forms/src/components/fields/DateField.tsx`              | Date field component                                                                             |
| `frontend/apps/forms/src/components/layout/ElementRenderer.tsx`        | Dispatches to field component by element type                                                    |
| `frontend/apps/forms/src/components/layout/SectionRenderer.tsx`        | Renders a section and its elements                                                               |
| `frontend/apps/forms/src/components/layout/PageRenderer.tsx`           | Renders a page and its sections                                                                  |
| `frontend/apps/forms/src/routes/forms/$formId/versions/$versionId.tsx` | Form viewer route                                                                                |
| `frontend/apps/forms/.env`                                             | Local environment variable defaults                                                              |
| `frontend/apps/forms/.env.example`                                     | Documented environment variable template                                                         |

### Files to modify

| File                                                                 | Change                                                                                        |
| -------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `frontend/apps/forms/src/types/element.ts`                           | Replace loose attributes type with discriminated union; remove `boolean` from `ElementType` ✓ |
| `frontend/apps/forms/src/types/submission.ts`                        | Tighten `ISubmissionValue.value` typing                                                       |
| `frontend/apps/forms/src/services/dataSourcesService.ts`             | Implement `getLookups()` method ✓                                                             |
| `frontend/apps/forms/src/services/submissionService.ts`              | Add idempotency key header support                                                            |
| `frontend/apps/forms/src/hooks/useHttpService.ts`                    | Use `import.meta.env` for base URLs                                                           |
| `frontend/apps/forms/src/components/FormElement/FormElement.tsx`     | Full rewrite — multi-page wizard, `FormProvider`, `useActionState` — _partial: `FormProvider` and `rawSubmission` wired up_ |
| `frontend/apps/forms/src/components/FormElement/FormElement.type.ts` | Add `token` prop if auth is wired up later                                                    |
| `frontend/apps/forms/src/routes/__root.tsx`                          | Remove placeholder content                                                                    |

---

## Open Questions

1. **MUI DatePicker vs. native input** — `@mui/x-date-pickers` is not in the current dependencies. Decision needed: add the package for a richer date picker experience, or use a native `<input type="date">` wrapped in MUI styling.

2. **Dynamic binding re-fetch scope** — `DataSourceRef.bindings` can reference other field values (`type: "field"`), meaning a select's options depend on what the user entered elsewhere. Confirm whether re-fetching lookups on upstream field changes is in scope, or whether only static bindings need to be handled initially.

3. **Auth token source** — `accessToken` is currently `"placeholder"`. The `authentication` MFE is a stub. This will need to be resolved before the renderer can be used against a real backend. Options: a prop on `FormElementProps`, a shared React context from the host shell, or a dedicated auth hook.
