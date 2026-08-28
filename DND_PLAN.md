# Drag and Drop Implementation Plan

## Context

### What's already built
- `DraggableCard` component with a `DragIndicator` icon — already used by `PaletteItem`, `ElementItem`, and `SectionItem` as a visual stub
- `Collapisble` component with `motion` (Framer Motion v12) already installed and in use
- `ItemTools` with `onReorder`, `onCopy`, `onDelete` props wired into `ElementItem` and `SectionItem` (currently no-ops)
- `MoveElement` and `MoveSection` event handlers are **fully implemented** in `eventHandlers/`
- Canvas renders a full `Page → Section → Element` hierarchy
- `useFormDesignerSelect` and `useFormDesignerDispatch` hooks available throughout

### What's stubbed / missing
- `onAddElement` and `onAddSection` event handlers return the aggregate unchanged
- `AddElementEvent` lacks the fields needed to create an element from the toolbox (`elementType`, `sectionId`)
- `AddSectionEvent` lacks `pageId`
- No drag-and-drop library is installed
- `ItemTools.onReorder` callbacks are empty arrow functions

---

## Dependencies to Install

```sh
pnpm add @dnd-kit/react
```

`@dnd-kit/react` is a single package that includes everything needed: `DragDropProvider`, `DragOverlay`,
`useDraggable`, `useDroppable`, and `useSortable` (via `@dnd-kit/react/sortable`). It automatically
pulls in `@dnd-kit/dom`, `@dnd-kit/abstract`, and other internals as peer dependencies — no separate
packages required. This is the new unified API that replaces the legacy `@dnd-kit/core` +
`@dnd-kit/sortable` + `@dnd-kit/utilities` combination.

---

## Implementation Steps

### ✅ Step 1 — Install @dnd-kit/react

```sh
pnpm --filter @sundance/forms add @dnd-kit/react
```

### Step 2 — Extend Event Types (`store/formDesigner/events.ts`)

`AddElementEvent` needs to carry enough information to construct a new element from a palette drag:

```ts
// Before
type AddElementEvent = {
  type: "AddElement";
  elementId: string;
  position: number;
};

// After
type AddElementEvent = {
  type: "AddElement";
  elementId: string;       // newly generated ID (caller's responsibility)
  elementType: ElementType;
  sectionId: string;       // target section to insert into
  position: number;        // index within that section
};
```

`AddSectionEvent` needs `pageId` so the handler knows which page to insert into:

```ts
// Before
type AddSectionEvent = { type: "AddSection"; sectionId: string; position: number };

// After
type AddSectionEvent = {
  type: "AddSection";
  sectionId: string;
  pageId: string;
  position: number;
};
```

### Step 3 — Create Element Factory (`store/formDesigner/elementFactory.ts`)

A new pure utility that maps an `ElementType` to a new `IElement` with a generated ID and
type-appropriate attribute defaults. Used by `onAddElement`.

```ts
export function createElementFromType(
  id: string,
  type: ElementType,
): IElement {
  // returns IElement with sensible defaults per type
}
```

### Step 4 — Implement `onAddElement` and `onAddSection`

**`eventHandlers/elementEventHandlers.ts`** — fill in the stub:

```ts
export function onAddElement(
  aggregate: IFormAggregate,
  event: AddElementEvent,
): IFormAggregate {
  const element = createElementFromType(event.elementId, event.elementType);
  const pages = aggregate.version.pages.map((page) => ({
    ...page,
    sections: page.sections.map((section) => {
      if (section.id !== event.sectionId) return section;
      return {
        ...section,
        elements: insertAtPosition(section.elements, element, event.position),
      };
    }),
  }));
  return { ...aggregate, version: { ...aggregate.version, pages } };
}
```

**`eventHandlers/sectionEventHandlers.ts`** — fill in the stub:

```ts
export function onAddSection(
  aggregate: IFormAggregate,
  event: AddSectionEvent,
): IFormAggregate {
  const section = createEmptySection(event.sectionId);
  const pages = aggregate.version.pages.map((page) => {
    if (page.id !== event.pageId) return page;
    return {
      ...page,
      sections: insertAtPosition(page.sections, section, event.position),
    };
  });
  return { ...aggregate, version: { ...aggregate.version, pages } };
}
```

### ✅ Step 5 — Set up `DragDropProvider` in `FormDesigner.tsx`

Extracted into `FormDesigner/providers/FormDesignerDragProvider.tsx` — a dedicated component that
sits inside `FormDesignerProvider` so it can access the Zustand store via `useFormDesignerDispatch`.
`FormDesigner.tsx` wraps the panels with `<FormDesignerDragProvider>` instead of inlining the
provider. `DragOverlay` uses the render-function form with a `switch` on `data.source`.
`onDragEnd` likewise switches on `dragData.source` and delegates to typed handler functions
(e.g. `handlePaletteDragEnd`) — body is stubbed pending Steps 2–4.

Two type files were added:
- `FormDesigner/types/formDragEvent.ts` — `FormDragEventData` discriminated union (`PaletteDragEventData`)
- `FormDesigner/types/formDropEvent.ts` — `FormDropEventData` discriminated union (`PaletteDropEventData`)

### ✅ Step 6 — Make `PaletteItem` Draggable

`PaletteItem` now accepts a `draggable?: boolean` prop (default `true`). `useDraggable` is always
called (hooks rules satisfied) with `disabled: !draggable`. The `type` field on `useDraggable` is
set to `item.dragType` — a new `PaletteItemDragType` enum (`Element` | `Section`) defined in
`palette.tsx` — so that droppables can use `accept` to restrict which palette types they receive.
When `draggable={false}` the `handleRef` is not forwarded to `DraggableCard`, making the overlay
clone purely visual. The `DragOverlay` renders `<PaletteItem item={...} draggable={false} />`.

> Note: `data.type` is used in the payload (not `data.elementType`). `handleRef` is passed as
> `null` when `draggable={false}` — consider changing to `undefined` to avoid a potential
> TypeScript error with `Ref<Element>`.

### ✅ Step 7 — Make `SectionItem` a Drop Zone

`SectionItem` now uses `useDroppable` with `accept: PaletteItemDragType.Element` — it only
activates for element-type palette drags, not section drags. The `ref` is applied to the root `li`.
Drop data is typed as `PaletteDropEventData` with `parentId: section.id`.

`SectionList` has `useDroppable` imported and is in progress — it will become the canvas-level
drop zone that accepts `PaletteItemDragType.Section` drags for creating new sections.

> `isDropTarget` visual highlight on `SectionItem` not yet wired into `getSectionItemStyles` — still pending.

### Step 8 — Handle Drop in `onDragEnd`

`handlePaletteDragEnd` stub is in place in `FormDesignerDragProvider`. Needs to dispatch the
correct event based on `dragData.type`:

```ts
const handlePaletteDragEnd = (
  dragData: PaletteDragEventData,
  dropData: PaletteDropEventData,
) => {
  if (dragData.type === "section") {
    dispatch({
      type: "AddSection",
      sectionId: crypto.randomUUID(),
      pageId: dropData.parentId,
      position: /* end of page */ 0,
    });
  } else {
    dispatch({
      type: "AddElement",
      elementId: crypto.randomUUID(),
      elementType: dragData.type,
      sectionId: dropData.parentId,
      position: /* end of section */ 0,
    });
  }
};
```

For future support of reordering existing elements within/between sections, add additional
`source` cases to the `onDragEnd` switch.

### Step 9 — Wire `ItemTools.onReorder` to Move Events

Now that `MoveElement` and `MoveSection` handlers are implemented, connect the stubs in
`ElementItem` and `SectionItem`:

```ts
// ElementItem
onReorder={(inc) =>
  dispatch({
    type: "MoveElement",
    elementId: element.id,
    targetSectionId: sectionId,        // needs sectionId passed as prop
    position: element.position + inc,
  })
}

// SectionItem
onReorder={(inc) =>
  dispatch({
    type: "MoveSection",
    sectionId: section.id,
    targetPageId: pageId,              // needs pageId passed as prop
    position: section.position + inc,
  })
}
```

This requires threading `sectionId` into `ElementItem` and `pageId` into `SectionItem` as props,
which means updating `ElementList` and `SectionList` to pass them down.

### Step 10 — Framer Motion Animations (low effort, already installed)

Wrap `ElementList` with `AnimatePresence` and convert `ElementItem`'s root to a `motion` component.
The `layout` prop handles reorder animations automatically — no manual calculation needed.

```tsx
// ElementList.tsx
import { AnimatePresence } from "motion/react";

<Box component="ul" sx={styles.list}>
  <AnimatePresence>
    {elements.map((element) => (
      <ElementItem element={element} key={element.id} />
    ))}
  </AnimatePresence>
</Box>

// ElementItem.tsx — change root Box to motion.li
<Box
  component={motion.li}
  layout
  initial={{ opacity: 0, y: -6 }}
  animate={{ opacity: 1, y: 0 }}
  exit={{ opacity: 0, y: -6 }}
  transition={{ type: "spring", bounce: 0.2, duration: 0.35 }}
  sx={styles.item}
  onClick={handleClk}
  role="button"
>
```

Apply the same pattern to `SectionItem` / `SectionList` for section-level animations.

---

## File Change Summary

| File | Type | Status | Change |
|---|---|---|---|
| `package.json` | Modify | ✅ | Add `@dnd-kit/react` |
| `FormDesigner/types/formDragEvent.ts` | **New** | ✅ | `FormDragEventData` / `PaletteDragEventData` discriminated union |
| `FormDesigner/types/formDropEvent.ts` | **New** | ✅ | `FormDropEventData` / `PaletteDropEventData` discriminated union |
| `FormDesigner/providers/FormDesignerDragProvider.tsx` | **New** | ✅ | `DragDropProvider` + `DragOverlay` + typed `onDragEnd` switch |
| `FormDesigner/FormDesigner.tsx` | Modify | ✅ | Wrap panels with `FormDesignerDragProvider` |
| `ToolboxPanel/palette.tsx` | Modify | ✅ | `PaletteItemDragType` enum + `dragType` field on every `IPaletteItem` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | ✅ | `useDraggable` with `type: item.dragType`; `draggable` prop pattern |
| `components/DraggableCard.tsx` | Modify | ✅ | Optional `handleRef` prop forwarded to root `Box` |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | 🔄 | `useDroppable` wired with `accept: PaletteItemDragType.Element`; `isDropTarget` style not yet applied |
| `CanvasPanel/lists/SectionList.tsx` | Modify | 🔄 | `useDroppable` imported; canvas-level section drop zone not yet wired |
| `store/formDesigner/events.ts` | Modify | | Extend `AddElementEvent` (`elementType`, `sectionId`); extend `AddSectionEvent` (`pageId`) |
| `store/formDesigner/elementFactory.ts` | **New** | | `createElementFromType()` factory function |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | | Implement `onAddElement` stub |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | | Implement `onAddSection` stub |
| `CanvasPanel/lists/SectionItem.style.ts` | Modify | | Add `isDropTarget` style variant |
| `CanvasPanel/lists/ElementList.tsx` | Modify | | Wrap with `AnimatePresence` |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | | Switch root to `motion.li` with `layout` + enter/exit animations; wire `onReorder` |

**Total: 11 modified files, 4 new files. 8 of 17 complete, 2 in progress.**

---

## Decisions / Open Questions

- **Drop position within a section**: Currently dispatches with `position: 0` as a placeholder. For precise insertion (before/after a specific element), `ElementItem` would also need `useDroppable` or the section would need to track hover position — this is a scope increase.
- **Cross-section drag of existing elements**: `MoveElement` supports cross-section moves. Wiring this up with `useSortable` (from `@dnd-kit/react/sortable`) across multiple droppable containers is the natural next step after the toolbox drop is working.
- **`handleRef` null vs undefined**: `DraggableCard` accepts `handleRef?: Ref<Element>` but `PaletteItem` passes `null` when `draggable={false}`. Change to `undefined` to avoid a potential TypeScript error.
- **`SectionList` as section drop zone**: Needs `pageId` threaded down as a prop so the droppable data can carry it for `AddSectionEvent` dispatch.
