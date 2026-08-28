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
(`handlePaletteDragEnd`) — body is stubbed pending Steps 2–4. `onDragStart` sets `activeDragData`
in context; `onDragEnd` clears it to `null`.

A `FormDesignerDragContext` is co-located in the provider file, exposed via `useFormDragEvent()`.
The hook returns `FormDragEventData | null` — `null` means no drag is active (no error thrown).
`FormDragEventSource` enum replaces the `"palette"` string literal throughout.

Type files added:
- `FormDesigner/types/formDragEvent.ts` — `FormDragEventData` / `PaletteDragEventData` / `FormDragEventSource` enum
- `FormDesigner/types/formDropEvent.ts` — `FormDropEventData` / `PaletteDropEventData`

### ✅ Step 6 — Make `PaletteItem` Draggable

`PaletteItem` accepts a `draggable?: boolean` prop (default `true`). `useDraggable` is always
called (hooks rules satisfied) with `disabled: !draggable`. The `type` field is set to
`item.dragType` (`PaletteItemDragType` enum) so droppables can use `accept` to restrict which
palette types they receive. The `DragOverlay` renders `<PaletteItem item={...} draggable={false} />`
as a static visual clone.

> Note: `handleRef` is passed as `null` when `draggable={false}` — consider changing to
> `undefined` to avoid a potential TypeScript error with `Ref<Element>`.

### ✅ Step 7 — Drop Zones with Ambient Indicators

Drop zones are wired with `useDroppable` and `accept`, and show a `<DropZoneIndicator>` as soon
as a drag starts — before the cursor is over them — using `useFormDragEvent()` context:

- **`SectionItem`** — `accept: PaletteItemDragType.Element`; renders `<DropZoneIndicator text="Drop element here">` when `canDropItem(dragData)` is true (drag active + type is not `"section"`)
- **`PageItem`** — `accept: PaletteItemDragType.Section`; renders `<DropZoneIndicator text="Drop your section here">` when drag active + type is `"section"`. This is the drop zone for new sections, not `SectionList`.

`canDropItem()` is a local helper in each component that checks `dragData.source` and `dragData.type`.

> `isDropTarget` visual highlight (stronger style when cursor is directly over the zone) not yet
> wired into `getSectionItemStyles` — still pending.
>
> `SectionList.tsx` has a stale unused `useDroppable` import — can be cleaned up.
>
> `canDropItem` TODO noted in code: logic should use `PaletteItemDragType` rather than checking
> `data.type !== "section"` to avoid future brittleness.

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
| `FormDesigner/types/formDragEvent.ts` | **New** | ✅ | `FormDragEventData` / `PaletteDragEventData` / `FormDragEventSource` enum |
| `FormDesigner/types/formDropEvent.ts` | **New** | ✅ | `FormDropEventData` / `PaletteDropEventData` |
| `FormDesigner/providers/FormDesignerDragProvider.tsx` | **New** | ✅ | `DragDropProvider` + `DragOverlay` + context + `useFormDragEvent` hook + typed `onDragEnd` switch |
| `FormDesigner/FormDesigner.tsx` | Modify | ✅ | Wrap panels with `FormDesignerDragProvider` |
| `ToolboxPanel/palette.tsx` | Modify | ✅ | `PaletteItemDragType` enum + `dragType` field on every `IPaletteItem` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | ✅ | `useDraggable` with `type: item.dragType`; `draggable` prop pattern; `FormDragEventSource` enum |
| `components/DragDrop/DraggableCard.tsx` | Moved | ✅ | Moved from `components/` to `components/DragDrop/`; optional `handleRef` prop |
| `components/DragDrop/DropZoneIndicator.tsx` | **New** | ✅ | Ambient drop zone indicator component shown on drag start |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | ✅ | `useDroppable` + `useFormDragEvent` + `canDropItem` + `DropZoneIndicator` |
| `CanvasPanel/lists/PageItem.tsx` | Modify | ✅ | `useDroppable` (section drops) + `useFormDragEvent` + `canDropItem` + `DropZoneIndicator` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | 🔄 | Stale `useDroppable` import — needs cleanup |
| `CanvasPanel/lists/SectionItem.style.ts` | Modify | | Add `isDropTarget` stronger hover style variant |
| `store/formDesigner/events.ts` | Modify | | Extend `AddElementEvent` (`elementType`, `sectionId`); extend `AddSectionEvent` (`pageId`) |
| `store/formDesigner/elementFactory.ts` | **New** | | `createElementFromType()` factory function |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | | Implement `onAddElement` stub |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | | Implement `onAddSection` stub |
| `CanvasPanel/lists/ElementList.tsx` | Modify | | Wrap with `AnimatePresence` |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | | Switch root to `motion.li` with `layout` + enter/exit animations; wire `onReorder` |

**Total: 11 modified files, 5 new files, 1 moved. 11 of 19 complete, 1 in progress.**

---

## Decisions / Open Questions

- **`isDropTarget` active hover style**: `isDropTarget` from `useDroppable` is available in both `SectionItem` and `PageItem` but not yet passed into their style functions. Should show a stronger highlight (e.g. solid border) when the cursor is directly over the zone, on top of the ambient `DropZoneIndicator`.
- **`canDropItem` brittleness**: Both `SectionItem` and `PageItem` use `data.type !== "section"` / `data.type === "section"` as the condition. Should use `PaletteItemDragType` from `event.operation.source.type` (available in `onDragStart`) instead of the `data` payload for robustness.
- **Drop position**: Currently dispatches `position: 0` as a placeholder. Precise insertion requires `ElementItem` / `SectionItem` to also act as droppables with position tracking.
- **`SectionList` cleanup**: Stale `useDroppable` import needs removing.
- **`handleRef` null vs undefined**: `PaletteItem` passes `null` when `draggable={false}`. Change to `undefined` to match `Ref<Element>` type.
- **Cross-section drag of existing elements**: Natural next step after drop dispatch is working — use `useSortable` from `@dnd-kit/react/sortable` with `group` for cross-section moves.
