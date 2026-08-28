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
as a static visual clone. `isDragging` drives an `onDrag` style (highlighted background + solid
border via `mergeSx`) applied to the source item while dragging.

> Note: `handleRef` is passed as `null` when `draggable={false}` — consider changing to
> `undefined` to avoid a potential TypeScript error with `Ref<Element>`.

### ✅ Step 7 — Drop Zones with Ambient + Active Indicators

Drop zones are wired with `useDroppable` and `accept`, and show a `<DropZoneIndicator>` as soon
as a drag starts — before the cursor is over them — using `useFormDragEvent()` context. When the
cursor is directly over the zone, `isDropTarget` (from `useDroppable`) is passed to
`DropZoneIndicator` which animates the background from `20%` to `60%` opacity using Framer Motion
`variants` and `useTheme` + `useMemo` to access the MUI palette inside motion variants.

- **`SectionItem`** — `accept: PaletteItemDragType.Element`; renders `<DropZoneIndicator text="Drop {label} here" isDropTarget={isDropTarget}>` when `canDropItem(dragData)` is true
- **`PageItem`** — `accept: PaletteItemDragType.Section`; same pattern for section drags

`dragPaletteItem` is memoised via `useMemo(() => findPaletteItem(dragData.type), [dragData])` in
both components so the indicator label reflects the specific element type being dragged.
`canDropItem()` is a local helper in each component that checks `dragData.source` and `dragData.type`.

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
Use `layout="position"` with a tween `layout` transition to prevent the stretching bounce that
occurs when `DropZoneIndicator` is inserted into the DOM alongside list items. Wrap
`DropZoneIndicator` in an `AnimatePresence` + `motion.div` for smooth height enter/exit.

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

// ElementItem.tsx
<Box
  component={motion.li}
  layout="position"
  initial={{ opacity: 0, y: -6 }}
  animate={{ opacity: 1, y: 0 }}
  exit={{ opacity: 0, y: -6 }}
  transition={{
    layout: { type: "tween", duration: 0.15, ease: "easeOut" },
    default: { type: "spring", bounce: 0.2, duration: 0.35 },
  }}
  sx={styles.item}
  onClick={handleClk}
  role="button"
>

// DropZoneIndicator wrapper in SectionItem / PageItem
<AnimatePresence>
  {canDrop && (
    <motion.div
      key="drop-indicator"
      initial={{ opacity: 0, height: 0 }}
      animate={{ opacity: 1, height: "4.25rem" }}
      exit={{ opacity: 0, height: 0 }}
      transition={{ type: "spring", bounce: 0, duration: 0.3 }}
    >
      <DropZoneIndicator text={...} isDropTarget={isDropTarget} />
    </motion.div>
  )}
</AnimatePresence>
```

Apply the same `layout="position"` pattern to `SectionItem` / `SectionList` for section-level animations.

---

## File Change Summary

| File | Type | Status | Change |
|---|---|---|---|
| `package.json` | Modify | ✅ | Add `@dnd-kit/react` |
| `FormDesigner/types/formDragEvent.ts` | **New** | ✅ | `FormDragEventData` / `PaletteDragEventData` / `FormDragEventSource` enum |
| `FormDesigner/types/formDropEvent.ts` | **New** | ✅ | `FormDropEventData` / `PaletteDropEventData` |
| `FormDesigner/providers/FormDesignerDragProvider.tsx` | **New** | ✅ | `DragDropProvider` + `DragOverlay` + context + `useFormDragEvent` hook + typed `onDragEnd` switch |
| `FormDesigner/FormDesigner.tsx` | Modify | ✅ | Wrap panels with `FormDesignerDragProvider` |
| `FormDesigner/FormDesigner.styles.ts` | Modify | ✅ | Right panel column widened to `28rem` |
| `ToolboxPanel/palette.tsx` | Modify | ✅ | `PaletteItemDragType` enum + `dragType` field on every `IPaletteItem` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | ✅ | `useDraggable` with `type: item.dragType`; `draggable` prop; `isDragging` → `onDrag` style via `mergeSx` |
| `components/DragDrop/DraggableCard.tsx` | Moved | ✅ | Moved from `components/` to `components/DragDrop/`; optional `handleRef` prop |
| `components/DragDrop/DropZoneIndicator.tsx` | **New** | ✅ | Ambient + active indicator; `isDropTarget` prop animates background via Framer Motion `variants` + `useTheme` |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | ✅ | `useDroppable` + `useFormDragEvent` + `canDropItem` + `dragPaletteItem` memo + `DropZoneIndicator` with `isDropTarget` |
| `CanvasPanel/lists/PageItem.tsx` | Modify | ✅ | `useDroppable` (section drops) + same pattern as `SectionItem` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | ✅ | Stale `useDroppable` import removed; `gap: 2.5` added to list styles |
| `store/formDesigner/events.ts` | Modify | | Extend `AddElementEvent` (`elementType`, `sectionId`); extend `AddSectionEvent` (`pageId`) |
| `store/formDesigner/elementFactory.ts` | **New** | | `createElementFromType()` factory function |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | | Implement `onAddElement` stub |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | | Implement `onAddSection` stub |
| `CanvasPanel/lists/ElementList.tsx` | Modify | | Wrap with `AnimatePresence`; `DropZoneIndicator` height animation wrapper |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | | `motion.li` with `layout="position"` + tween transition + enter/exit; wire `onReorder` |

**Total: 12 modified files, 5 new files, 1 moved. 13 of 19 complete.**

---

## Decisions / Open Questions

- **`canDropItem` brittleness**: Both `SectionItem` and `PageItem` use `data.type !== "section"` / `data.type === "section"`. Should store `source.type` (`PaletteItemDragType`) from `onDragStart` in context instead of deriving it from the `data` payload.
- **Drop position**: Currently dispatches `position: 0` as a placeholder. Precise insertion requires `ElementItem` / `SectionItem` to also act as droppables with position tracking.
- **`handleRef` null vs undefined**: `PaletteItem` passes `null` when `draggable={false}`. Change to `undefined` to match `Ref<Element>` type.
- **Cross-section drag of existing elements**: Natural next step after drop dispatch is working — use `useSortable` from `@dnd-kit/react/sortable` with `group` for cross-section moves.
