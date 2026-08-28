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

`AddSectionEvent` has `pageId`; `AddElementEvent` has `sectionId` and `elementType`:

```ts
export type AddSectionEvent = {
  type: "AddSection";
  pageId: string;
  position: number;
};

export type AddElementEvent = {
  type: "AddElement";
  sectionId: string;
  elementType: ElementType;
  position: number;
};
```

### ✅ Step 3 — Element + Section Factories

**`src/factories/elementFactory.ts`** — Registry-based abstract factory. `registerElementFactory`
is private (not exported). All 9 element types registered inline with `satisfies` to verify
attribute shapes at compile time. ID generation via `generatedID()`.

**`src/factories/sectionFactory.ts`** — `createEmptySection()` returns a new `ISection` with a
generated ID, empty `elements`, `rules`, `key`, `name: "Section"`, and `position: 0`.

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

**`elementEventHandlers.ts`** — imports `createElementFromType`; maps pages → sections to find
target by `event.sectionId`; inserts via `insertAtPosition`.

**`sectionEventHandlers.ts`** — imports `createEmptySection`; maps pages to find target by
`event.pageId`; inserts via `insertAtPosition`.

### ✅ Step 6 — Wire `handlePaletteDragEnd` Dispatch

`dispatch` is destructured from `useFormDesignerDispatch`. `handlePaletteDragEnd` uses a `switch`
on `dragData.type` with `satisfies` type checking on each constructed event before dispatching.
`"section"` routes to `AddSectionEvent`; all other types route to `AddElementEvent`. The
`elementFactory.ts` registry is triggered via the import of `createElementFromType` inside
`elementEventHandlers.ts` which is in the module graph.

### Step 7 — Wire `ItemTools.onReorder` to Move Events

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

### Step 8 — Framer Motion Animations

Remaining animation work:
- `ElementList` — `AnimatePresence` wrapper for element enter/exit
- `ElementItem` — `motion.li` with `layout="position"` + tween layout transition + enter/exit
- `SectionList` — `AnimatePresence` wrapper for section enter/exit
- `SectionItem` — `motion.li` with `layout="position"`

---

## File Change Summary

| File | Type | Status | Change |
|---|---|---|---|
| `package.json` | Modify | ✅ | Add `@dnd-kit/react` |
| `FormDesigner/types/formDragEvent.ts` | **New** | ✅ | `FormDragEventData` / `PaletteDragEventData` / `FormDragEventSource` enum |
| `FormDesigner/types/formDropEvent.ts` | **New** | ✅ | `FormDropEventData` / `PaletteDropEventData` with `position` |
| `FormDesigner/providers/FormDesignerDragProvider.tsx` | **New** | ✅ | `DragDropProvider` + `DragOverlay` + context + `useFormDragEvent` hook + `handlePaletteDragEnd` with typed switch dispatch |
| `FormDesigner/FormDesigner.tsx` | Modify | ✅ | Wrap panels with `FormDesignerDragProvider` |
| `FormDesigner/FormDesigner.styles.ts` | Modify | ✅ | Right panel column widened to `28rem` |
| `ToolboxPanel/palette.tsx` | Modify | ✅ | `PaletteItemDragType` enum + `dragType` field on every `IPaletteItem` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | ✅ | `useDraggable` with `type: item.dragType`; `draggable` prop; `isDragging` → `onDrag` style |
| `components/DragDrop/DraggableCard.tsx` | Moved | ✅ | Moved from `components/` to `components/DragDrop/`; optional `handleRef` prop |
| `components/DragDrop/DropZoneIndicator.tsx` | **New** | ✅ | `isVisible` owns `AnimatePresence`; `isDropTarget` drives background animation via Framer Motion variants + `useTheme` |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | ✅ | `useDroppable` + `useFormDragEvent` + `canDropItem` + `sortPositioned`/`getNextPosition` + `DropZoneIndicator` |
| `CanvasPanel/lists/SectionItem.style.ts` | Modify | ✅ | `elements` style added for inner layout wrapper |
| `CanvasPanel/lists/PageItem.tsx` | Modify | ✅ | `useDroppable` (section drops) + `sortPositioned`/`getNextPosition` + `DropZoneIndicator` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | ✅ | Stale `useDroppable` import removed; `gap: 2.5` added |
| `src/types/elementAttributes.ts` | Modify | ✅ | `ElementAttributes` union extended to all 9 types |
| `src/factories/elementFactory.ts` | **New** | ✅ | Registry-based abstract factory; all 9 types registered inline |
| `src/factories/sectionFactory.ts` | **New** | ✅ | `createEmptySection()` factory |
| `src/utils/id.ts` | **New** | ✅ | `generatedID()` + `isTemporaryID()` |
| `src/utils/position.ts` | **New** | ✅ | `sortPositioned()` + `getNextPosition()` |
| `store/formDesigner/events.ts` | Modify | ✅ | `AddSectionEvent` has `pageId`; `AddElementEvent` has `sectionId` + `elementType` |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | ✅ | `onAddElement` implemented |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | ✅ | `onAddSection` implemented |
| `CanvasPanel/lists/ElementList.tsx` | Modify | | `AnimatePresence` wrapper |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | | `motion.li` + `layout="position"` + tween transition + enter/exit; wire `onReorder` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | | `AnimatePresence` wrapper |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | | `motion.li` + `layout="position"` |

**Total: 15 modified files, 7 new files, 1 moved. 22 of 23 complete.**

---

## Decisions / Open Questions

- **`handleRef` null vs undefined**: `PaletteItem` passes `null` when `draggable={false}`. Change to `undefined` to match `Ref<Element>` type.
- **Drop position**: `getNextPosition` appends to the end. Precise mid-list insertion would require per-element droppable hit-testing — out of scope for now.
- **Cross-section drag of existing elements**: Natural next step — `useSortable` from `@dnd-kit/react/sortable` with `group` for cross-section moves.
