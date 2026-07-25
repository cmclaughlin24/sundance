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

### 2b. Environment-based API URLs ✓

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

## Phase 3 — Form State Architecture ✓ _partial_

### 3a. Zustand store ✓

State management migrated from `useReducer` + dual contexts to **Zustand**. Created as `store/formStore.ts`:

- `createFormStore(form, version, raw)` — initializes store with `form`, `version`, `values` (pre-filled from `raw`) ✓
- `setValue(elementId, value)` — updates `values` record ✓
- `setError(elementId, errors)` — stub, not yet implemented
- `errors` not yet in store state

### 3b. `FormStoreProvider` ✓

`store/FormStoreProvider.tsx` — uses `useRef` to create the store once on mount, publishes via `FormStoreContext`. Accepts `form`, `version`, `rawSubmission` props ✓

### 3c. Granular selectors ✓

`store/useFormStoreContext.ts` — selector hooks using `useStore`:

- `useForm()` — subscribes to `form`
- `useFormVersion()` — subscribes to `version`
- `useFormValues()` — subscribes to full `values` map
- `useElementValue(elementId)` — subscribes to a single field value
- `useTenantId()` — subscribes to `form.tenantId`
- `useFormDispatch()` — returns `{ setValue, setError }` via `useShallow`
- `useElementRuleState(element)` — reads `evalCtx`, seeds defaults from `element.attributes`

**Files created:**

- `frontend/apps/forms/src/store/formStore.ts` ✓
- `frontend/apps/forms/src/store/FormStoreProvider.tsx` ✓
- `frontend/apps/forms/src/store/useFormStoreContext.ts` ✓
- `frontend/apps/forms/src/store/evalContext.ts` ✓

---

## Phase 4 — Client-Side Rule Evaluator ✓

### 4a. `utils/evaluate.ts` ✓

Pure evaluation utility — no React dependency, no `Function` constructor. Implemented as:

- `evaluateRules(rules, evalCtx, defaultState?)` — takes `IRule[]` and `EvalContext` (`Record<string, any>`), returns `Readonly<IRuleState>`
- `evaluateRule(rule, evalCtx)` — evaluates a single rule's expressions with join chaining
- `buildEvalContext(pages, values)` — builds `element.key → value` map from the full page hierarchy ✓
- Registry-based operator dispatch via `Map<RuleExpressionOp, EvaluatorFn>`
- `IRuleState` type added to `types/rule.ts` ✓
- `EvalContext` published via `store/evalContext.ts`; computed once in `FormRenderer` via `useMemo`, consumed by all field components via `useEvalContext()` ✓
- `useElementRuleState(element)` in `useFormContext.ts` — reads `evalCtx`, seeds defaults from `element.attributes` ✓
- `utils/filter.ts` — `filterVisible<T extends HasRules>(items, evalCtx)` utility ✓

**Files created:**

- `frontend/apps/forms/src/utils/evaluate.ts` ✓
- `frontend/apps/forms/src/utils/filter.ts` ✓
- `frontend/apps/forms/src/types/rule.ts` — `IRuleState` added ✓

---

## Phase 5 — Field Components

One component per `ElementType`. All field components:

- Consume `FormContext` via `useContext` to read and write values
- Read `ruleStates[element.id]` to determine visibility, required state, and read-only state
- Return `null` when `ruleStates[element.id].visible === false`
- Dispatch `SET_VALUE` on change

| Component       | MUI Input                     | Key attributes respected                                          | Status |
| --------------- | ----------------------------- | ----------------------------------------------------------------- | ------ |
| `TextField`     | `MUI TextField`               | `minLength`, `maxLength`, `pattern`, `placeholder`                | ✓ _partial — renders, `onChange` dispatches `SET_VALUE`, consumes `ruleState` (`required`, `readonly`)_ |
| `NumberField`   | `MUI TextField type="number"` | `min`, `max`, `step`                                              | ✓ _partial — renders, type guard fixed, `onChange` dispatches `SET_VALUE`, consumes `ruleState` (`required`, `readonly`)_ |
| `SelectField`   | `MUI Select` / `Autocomplete` | `data`, `dataSourceRef`, `multiple`, `minSelected`, `maxSelected` | ✓ _partial — static and dynamic (`dataSourceRef`) rendering, `BaseSelectFieldElement` extracted, disables while loading; `multiple` pending_ |
| `CheckboxField` | `MUI Checkbox`                | `isCheckedByDefault` (initializes value on mount)                 | ✓ _static and dynamic (`dataSourceRef`) rendering, `BaseCheckboxFieldElement` extracted, per-lookup onChange, disables while loading_ |
| `DateField`     | `MUI TextField type="date"`   | `minDate`, `maxDate`                                              | ✓ |

**`SelectField` specifics:**

- If `attributes.data` is non-empty and `attributes.dataSourceRef` is absent, render inline options directly
- If `attributes.dataSourceRef` is present, fetch options via `useDataSourceLookups` hook (Phase 8)
- Show a loading spinner while lookups are fetching
- Show an inline error if the lookup call fails
- Respect `multiple` for multi-select behavior

**Files to create:**

- `frontend/apps/forms/src/components/FormElement/Elements/TextFieldElement.tsx` ✓ _partial_
- `frontend/apps/forms/src/components/FormElement/Elements/NumberFieldElement.tsx` ✓ _partial_
- `frontend/apps/forms/src/components/FormElement/Elements/SelectFieldElement.tsx` ✓ _partial — `multiple` pending_
- `frontend/apps/forms/src/components/FormElement/Elements/CheckboxFieldElement.tsx` ✓
- `frontend/apps/forms/src/components/FormElement/Elements/DateFieldElement.tsx` ✓
- `frontend/apps/forms/src/components/FormElement/Elements/BaseFieldElement.tsx` ✓
- `frontend/apps/forms/src/components/FormElement/Elements/FieldElementLabel.tsx` ✓

---

## Phase 6 — Layout Components

### 6a. `ElementRenderer` ✓

Dispatches to the correct field component by `element.type` via a `Map` registry. `useFormDispatch` wired — dispatches `SET_VALUE` on change. All five types registered: `text`, `number`, `select`, `checkbox`, `date` ✓

### 6b. `SectionRenderer` ✓

- Sorts `section.elements` by `position` ascending via `sortPositioned`
- Filters elements by `filterVisible(elements, evalCtx)` via `useEvalContext()` ✓
- Maps visible elements through `ElementRenderer`

### 6c. `PageRenderer` ✓

- Sorts `page.sections` by `position` ascending
- Filters sections by `filterVisible(sections, evalCtx)` via `useEvalContext()` ✓
- Maps visible sections through `SectionRenderer`

### 6d. `FormRenderer` ✓ _partial_

- Reads `form`, `version`, `values` via Zustand selectors ✓
- Publishes `EvalContextContext` via `useMemo` keyed on `[version, values]` ✓
- Filters pages by `filterVisible(pages, evalCtx)` ✓
- `handleSubmit` correctly collects `Object.entries(values)` into `ISubmissionValue[]` ✓
- Multi-page wizard (Next/Back, single page at a time, progress indicator) not yet implemented

### 6e. `FormElement` updated ✓ _partial_

- `FormStoreProvider` wired with `form`, `version`, `rawSubmission` ✓
- `submitType` prop controls submission mode ✓
- `async` → calls `submissionService.submit()`, returns `{ referenceId }` ✓
- `sync` / default → calls `submissionService.normalize()`, returns `{ raw, normalized }` ✓
- Multi-page wizard not yet implemented

**Files created:**

- `frontend/apps/forms/src/components/FormElement/Renderer/ElementRenderer.tsx` ✓
- `frontend/apps/forms/src/components/FormElement/Renderer/SectionRenderer.tsx` ✓
- `frontend/apps/forms/src/components/FormElement/Renderer/PageRenderer.tsx` ✓
- `frontend/apps/forms/src/components/FormElement/Renderer/FormRenderer.tsx` ✓ _partial_
- `frontend/apps/forms/src/utils/sort.ts` ✓

---

## Phase 7 — Submission Flow ✓ _partial_

### 7a. Submit handler ✓

`FormElement` supports two submission modes via `submitType` prop:
- `"async"` → `submissionService.submit()` → returns `{ referenceId }` ✓
- `"sync"` / default → `submissionService.normalize()` → returns `{ raw, normalized }` ✓

`onSubmit` callback called with typed result on success ✓. Errors silently caught — no error surfacing yet.

### 7b. Idempotency key

The backend enforces an `Idempotency-Key` header on `POST /submissions`. Not yet wired — `submissionService.submit()` sends no idempotency key. Generate a UUID on `FormElement` mount and pass as a header in `SubmissionsService.submit()`.

### 7c. Surface submit errors

Server-side validation errors not yet surfaced. `setError` in the Zustand store is a stub. Needs implementation in `formStore.ts` and dispatch on submission failure.

**Files to modify:**

- `frontend/apps/forms/src/services/submissionService.ts` — add idempotency key header
- `frontend/apps/forms/src/store/formStore.ts` — implement `setError`

---

## Phase 8 — DataSources for Select Fields ✓

### 8a. New hook: `useDataSource` ✓

Created as `hooks/useDataSource.ts`:

- `useTenantId()` for tenant header — no direct `useFormState` call needed ✓
- `useEvalContext()` for binding resolution — `evalCtx` keyed by `element.key` ✓
- `resolveBindings(bindings, evalCtx)` resolves both `static` and `field` binding types ✓
- `serializeFilters(filters)` — sorts keys before serializing for stable dependency comparison ✓
- `serialized` string used as `useAsyncData` dependency — prevents unnecessary re-fetches ✓
- Guards on `tenantId` before making API call ✓

### 8b. Static binding support ✓

`BindingSource.type === "static"` — literal `binding.value` passed directly as filter param.

### 8c. Dynamic binding support ✓

`BindingSource.type === "field"` — resolved via `evalCtx[binding.key]`; re-fetches when bound field values change via serialized dependency.

**Files created:**

- `frontend/apps/forms/src/hooks/useDataSource.ts` ✓

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

Update `frontend/apps/forms/src/routes/__root.tsx` to remove the placeholder `<div>Hello "__root"!</div>` and replace with a proper `<Outlet />` only (chrome is handled by the host shell). _Still contains placeholder content._

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
| `frontend/apps/forms/src/store/formStore.ts`                           | Zustand store — `createFormStore`, `setValue`, `setError` (stub) ✓                               |
| `frontend/apps/forms/src/store/FormStoreProvider.tsx`                  | Provider — `useRef` store creation, publishes `FormStoreContext` ✓                               |
| `frontend/apps/forms/src/store/useFormStoreContext.ts`                 | Granular selectors — `useForm`, `useFormVersion`, `useFormValues`, `useElementValue`, `useTenantId`, `useFormDispatch`, `useElementRuleState` ✓ |
| `frontend/apps/forms/src/utils/evaluate.ts`                            | Pure rule evaluator — `evaluateRules`, `evaluateRule`, `buildEvalContext`, operator registry ✓    |
| `frontend/apps/forms/src/utils/filter.ts`                              | `filterVisible<T extends HasRules>(items, evalCtx)` utility ✓                                    |
| `frontend/apps/forms/src/store/evalContext.ts`                         | `EvalContextContext` and `useEvalContext` hook ✓                                                  |
| `frontend/apps/forms/src/hooks/useDataSource.ts`                       | Async lookup fetcher — `resolveBindings`, `serializeFilters`, `useTenantId` integration ✓        |
| `frontend/apps/forms/src/components/FormElement/Elements/TextFieldElement.tsx`     | Text field — placeholder, onChange, ruleState ✓ _partial_                           |
| `frontend/apps/forms/src/components/FormElement/Elements/NumberFieldElement.tsx`   | Number field — min/max/step, onChange, ruleState ✓ _partial_                        |
| `frontend/apps/forms/src/components/FormElement/Elements/SelectFieldElement.tsx`   | Select field — static + dynamic dataSourceRef, BaseSelectFieldElement ✓ _partial: multiple pending_ |
| `frontend/apps/forms/src/components/FormElement/Elements/CheckboxFieldElement.tsx` | Checkbox group — static + dynamic dataSourceRef, BaseCheckboxFieldElement ✓         |
| `frontend/apps/forms/src/components/FormElement/Elements/DateFieldElement.tsx`     | Date field component ✓                                                               |
| `frontend/apps/forms/src/components/FormElement/Renderer/ElementRenderer.tsx`      | Dispatches to field component by element type ✓                                      |
| `frontend/apps/forms/src/components/FormElement/Renderer/SectionRenderer.tsx`      | Renders a section and its elements ✓                                                 |
| `frontend/apps/forms/src/components/FormElement/Renderer/PageRenderer.tsx`         | Renders a page and its sections ✓                                                    |
| `frontend/apps/forms/src/routes/forms/$formId/versions/$versionId.tsx` | Form viewer route                                                                                |
| `frontend/apps/forms/.env`                                             | Local environment variable defaults                                                              |
| `frontend/apps/forms/.env.example`                                     | Documented environment variable template                                                         |

### Files to modify

| File                                                                 | Change                                                                                        |
| -------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `frontend/apps/forms/src/types/element.ts`                           | Replace loose attributes type with discriminated union; remove `boolean` from `ElementType` ✓ |
| `frontend/apps/forms/src/types/rule.ts`                              | Add `IRuleState` interface ✓                                                                  |
| `frontend/apps/forms/src/types/submission.ts`                        | Tighten `ISubmissionValue.value` typing                                                       |
| `frontend/apps/forms/src/services/dataSourcesService.ts`             | Implement `getLookups()` method ✓                                                             |
| `frontend/apps/forms/src/services/submissionService.ts`              | Add idempotency key header support                                                            |
| `frontend/apps/forms/src/hooks/useHttpService.ts`                    | Use `import.meta.env` for base URLs ✓                                                         |
| `frontend/apps/forms/src/components/FormElement/FormElement.tsx`     | `FormStoreProvider`, `submitType`, `asyncSubmit`/`syncSubmit` wired ✓ — multi-page wizard pending |
| `frontend/apps/forms/src/components/FormElement/FormElement.type.ts` | `ISyncSubmitEvent`, `IAsyncSubmitEvent`, `submitType` discriminated union ✓                    |
| `frontend/apps/forms/src/routes/index.tsx`                           | Updated with real `FormElement` usage ✓                                                       |
| `frontend/apps/forms/src/routes/__root.tsx`                          | Remove placeholder content                                                                    |

---

## Open Questions

1. **MUI DatePicker vs. native input** — `@mui/x-date-pickers` is not in the current dependencies. Decision needed: add the package for a richer date picker experience, or use a native `<input type="date">` wrapped in MUI styling.

2. **Auth token source** — `accessToken` is currently `"placeholder"`. The `authentication` MFE is a stub. This will need to be resolved before the renderer can be used against a real backend. Options: a prop on `FormElementProps`, a shared React context from the host shell, or a dedicated auth hook.
