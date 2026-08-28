# Drag and Drop Implementation Plan

## Context

### What's already built
- `DraggableCard` component with a `DragIndicator` icon — already used by `PaletteItem`, `ElementItem`, and `SectionItem` as a visual stub
- `Collapisble` component with `motion` (Framer Motion v12) already installed and in use
- `ItemTools` with `onReorder`, `onCopy`, `onDelete` props wired into `ElementItem` and `SectionItem` (currently no-ops)
- `MoveElement` and `MoveSection` event handlers are **fully implemented** in `eventHandlers/`
- Canvas renders a full `Page → Section → Element` hierarchy
- `useFormDesignerSelect` and `useFormDesignerDispatch` hooks available throughout

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

### ✅ Step 2 — Extend Event Types (`store/formDesigner/events.ts`)

`AddSectionEvent` and `AddElementEvent` carry caller-generated IDs (for auto-select after dispatch)
alongside the routing and position fields:

```ts
export type AddSectionEvent = {
  type: "AddSection";
  id: string;
  pageId: string;
  position: number;
};

export type AddElementEvent = {
  type: "AddElement";
  id: string;
  sectionId: string;
  elementType: ElementType;
  position: number;
};
```

### ✅ Step 3 — Element + Section Factories

**`src/factories/elementFactory.ts`** — Registry-based abstract factory. `registerElementFactory`
is private (not exported). All 9 element types registered inline with `satisfies` to verify
attribute shapes at compile time. `createElementFromType(elementType, id)` accepts a caller-provided
ID so the caller controls the ID for auto-select.

**`src/factories/sectionFactory.ts`** — `createEmptySection(id)` accepts a caller-provided ID,
returning a new `ISection` with empty `elements`, `rules`, `key`, `name: "Section"`, `position: 0`.

**`src/utils/id.ts`** — `generatedID()` + `isTemporaryID()`.

**`src/utils/position.ts`** — new file with two utilities:
- `sortPositioned<T extends HasPosition>(items)` — sorts by `position`
- `getNextPosition<T extends HasPosition>(items)` — returns last item's `position + 1`, or `0` if empty

`ElementAttributes` union extended to all 9 types.

`canDropItem` in both `PageItem` and `SectionItem` updated to use `PaletteItemDragType.Section`
enum instead of the `"section"` string literal.

### ✅ Step 4 — Drop Zone Position

`PaletteDropEventData` extended with `position: number`. Both `SectionItem` and `PageItem` compute
the append position at droppable registration time using `sortPositioned` + `getNextPosition`:

- **`SectionItem`** — `position: getNextPosition(sortPositioned(section.elements))`
- **`PageItem`** — `position: getNextPosition(sortPositioned(page.sections))`

`DropZoneIndicator` updated to own its `AnimatePresence` via `isVisible` prop — callers render it
unconditionally and pass `isVisible={canDrop}`.

### ✅ Step 5 — Implement `onAddElement` and `onAddSection`

**`elementEventHandlers.ts`** — imports `createElementFromType`; passes `event.id` as the element
ID; maps pages → sections to find target by `event.sectionId`; inserts via `insertAtPosition`.

**`sectionEventHandlers.ts`** — imports `createEmptySection`; passes `event.id` as the section ID;
maps pages to find target by `event.pageId`; inserts via `insertAtPosition`.

### ✅ Step 6 — Wire `handlePaletteDragEnd` Dispatch + Auto-select

`dispatch` and `select` are both destructured from their respective hooks in
`FormDesignerDragProvider`. `generatedID()` is called upfront to produce the ID, which is passed
into both the event (so the handler uses it as the entity ID) and `select` (so the new item is
immediately selected after dropping).

`handlePaletteDragEnd` uses a `switch` on `dragData.type` with `satisfies` type checking on each
constructed event. `"section"` routes to `AddSectionEvent`; all other types route to
`AddElementEvent`.

### ✅ Step 7 — Framer Motion Animations + DropZone Transition

**`ElementList`** — wrapped with `AnimatePresence` so enter/exit transitions are coordinated.

**`ElementItem`** — `motion.li` with:
- `layout="position"` — siblings slide smoothly when list changes
- `initial={{ opacity: 0, height: 0 }}` / `animate={{ opacity: 1, height: "auto" }}` / `exit={{ opacity: 0, height: 0 }}`
- `transition={{ type: "spring", bounce: 0, duration: 0.3 }}`

**`DropZoneIndicator` overlap fix** — exit variant updated with `opacity: 0`, `height: 0`, and a
fast `transition: { duration: 0.1, ease: "easeIn" }`. The indicator collapses in 100ms while the
new element expands over 300ms spring — the indicator is effectively gone before the element is
half-visible, eliminating the visual overlap without deferred dispatch or `onExitComplete`
complexity.

### Step 8 — Section-level Animations

Apply the same Framer Motion pattern as `ElementList` / `ElementItem` to sections:
- `SectionList` — `AnimatePresence` wrapper
- `SectionItem` — `motion.li` with `layout="position"`

### Step 9 — Wire `ItemTools.onReorder` to Move Events

Connect the stubs in `ElementItem` and `SectionItem`:

```ts
// ElementItem
onReorder={(inc) =>
  dispatch({
    type: "MoveElement",
    elementId: element.id,
    targetSectionId: sectionId,
    position: element.position + inc,
  })
}

// SectionItem
onReorder={(inc) =>
  dispatch({
    type: "MoveSection",
    sectionId: section.id,
    targetPageId: pageId,
    position: section.position + inc,
  })
}
```

Requires threading `sectionId` into `ElementItem` and `pageId` into `SectionItem` as props,
which means updating `ElementList` and `SectionList` to pass them down.

---

## File Change Summary

| File | Type | Status | Change |
|---|---|---|---|
| `package.json` | Modify | ✅ | Add `@dnd-kit/react` |
| `FormDesigner/types/formDragEvent.ts` | **New** | ✅ | `FormDragEventData` / `PaletteDragEventData` / `FormDragEventSource` enum |
| `FormDesigner/types/formDropEvent.ts` | **New** | ✅ | `FormDropEventData` / `PaletteDropEventData` with `position` |
| `FormDesigner/providers/FormDesignerDragProvider.tsx` | **New** | ✅ | `DragDropProvider` + `DragOverlay` + context + `useFormDragEvent` + typed switch dispatch + auto-select via caller-generated ID |
| `FormDesigner/FormDesigner.tsx` | Modify | ✅ | Wrap panels with `FormDesignerDragProvider` |
| `FormDesigner/FormDesigner.styles.ts` | Modify | ✅ | Right panel column widened to `28rem` |
| `ToolboxPanel/palette.tsx` | Modify | ✅ | `PaletteItemDragType` enum + `dragType` field on every `IPaletteItem` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | ✅ | `useDraggable` with `type: item.dragType`; `draggable` prop; `isDragging` → `onDrag` style |
| `components/DragDrop/DraggableCard.tsx` | Moved | ✅ | Moved from `components/` to `components/DragDrop/`; optional `handleRef` prop |
| `components/DragDrop/DropZoneIndicator.tsx` | **New** | ✅ | `isVisible` owns `AnimatePresence`; `isDropTarget` drives background animation; exit collapses in 100ms to avoid overlap with entering element |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | ✅ | `useDroppable` + `useFormDragEvent` + `canDropItem` + `sortPositioned`/`getNextPosition` + `DropZoneIndicator` |
| `CanvasPanel/lists/SectionItem.style.ts` | Modify | ✅ | `elements` style added for inner layout wrapper |
| `CanvasPanel/lists/PageItem.tsx` | Modify | ✅ | `useDroppable` (section drops) + `sortPositioned`/`getNextPosition` + `DropZoneIndicator` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | ✅ | Stale `useDroppable` import removed; `gap: 2.5` added |
| `src/types/elementAttributes.ts` | Modify | ✅ | `ElementAttributes` union extended to all 9 types |
| `src/factories/elementFactory.ts` | **New** | ✅ | Registry-based abstract factory; all 9 types registered inline; accepts caller-provided ID |
| `src/factories/sectionFactory.ts` | **New** | ✅ | `createEmptySection(id)` factory; accepts caller-provided ID |
| `src/utils/id.ts` | **New** | ✅ | `generatedID()` + `isTemporaryID()` |
| `src/utils/position.ts` | **New** | ✅ | `sortPositioned()` + `getNextPosition()` |
| `store/formDesigner/events.ts` | Modify | ✅ | `AddSectionEvent` + `AddElementEvent` have `id`, `pageId`/`sectionId`, `elementType` |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | ✅ | `onAddElement` implemented with caller-provided `event.id` |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | ✅ | `onAddSection` implemented with caller-provided `event.id` |
| `CanvasPanel/lists/ElementList.tsx` | Modify | ✅ | `AnimatePresence` wrapper added |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | ✅ | `motion.li` + `layout="position"` + `height` enter/exit + spring transition |
| `CanvasPanel/lists/SectionList.tsx` | Modify | | `AnimatePresence` wrapper |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | | `motion.li` + `layout="position"` |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | | Wire `onReorder` to `MoveElement` dispatch |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | | Wire `onReorder` to `MoveSection` dispatch |

**Total: 16 modified files, 7 new files, 1 moved. 24 of 28 complete.**

---

## Decisions / Open Questions

- **`handleRef` null vs undefined**: `PaletteItem` passes `null` when `draggable={false}`. Change to `undefined` to match `Ref<Element>` type.
- **Drop position**: `getNextPosition` appends to the end. Precise mid-list insertion would require per-element droppable hit-testing — out of scope for now.
- **Cross-section drag of existing elements**: Natural next step — `useSortable` from `@dnd-kit/react/sortable` with `group` for cross-section moves.
