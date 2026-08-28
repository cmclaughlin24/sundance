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
- `AddElementEvent` lacks `elementType` — the handler cannot know what kind of element to create
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

### ✅ Step 2 — Extend Event Types (`store/formDesigner/events.ts`)

`AddSectionEvent` and `AddElementEvent` updated:

```ts
export type AddSectionEvent = {
  type: "AddSection";
  pageId: string;
  position: number;
};

export type AddElementEvent = {
  type: "AddElement";
  sectionId: string;
  position: number;
  // ⚠️ elementType still missing — needed before onAddElement can be implemented
};
```

> `elementType: ElementType` still needs to be added to `AddElementEvent`.

### ✅ Step 3 — Element Factory (`src/factories/elementFactory.ts`)

Registry-based abstract factory. `registerElementFactory` is private (not exported). All 9 element
types registered inline with `satisfies` to verify attribute shapes at compile time. ID generation
handled internally via `generatedID()` from `src/utils/id.ts` (uses `crypto.randomUUID()` with a
`TEMP_` prefix; companion `isTemporaryID()` available for server-reconciliation later).

`ElementAttributes` union in `src/types/elementAttributes.ts` extended to include all 9 types:
`SegmentedElementAttributes`, `RadioElementAttributes`, `ToggleElementAttributes`,
`UserElementAttributes` were previously missing.

`canDropItem` in both `PageItem` and `SectionItem` updated to use `PaletteItemDragType.Section`
enum instead of the `"section"` string literal.

### Step 4 — Add `elementType` to `AddElementEvent`

```ts
export type AddElementEvent = {
  type: "AddElement";
  sectionId: string;
  elementType: ElementType;   // ← add this
  position: number;
};
```

### Step 5 — Implement `onAddElement` and `onAddSection`

**`eventHandlers/elementEventHandlers.ts`:**

```ts
export function onAddElement(
  aggregate: IFormAggregate,
  event: AddElementEvent,
): IFormAggregate {
  const element = createElementFromType(event.elementType);
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

**`eventHandlers/sectionEventHandlers.ts`:**

```ts
export function onAddSection(
  aggregate: IFormAggregate,
  event: AddSectionEvent,
): IFormAggregate {
  const section = createEmptySection(event.pageId);
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

A `createEmptySection` utility is also needed — returns a new `ISection` with a generated ID,
empty `elements`, `rules`, `key`, and `name`.

### Step 6 — Wire `handlePaletteDragEnd` Dispatch

Replace the `console.log` stub in `FormDesignerDragProvider`:

```ts
const handlePaletteDragEnd = (
  dragData: PaletteDragEventData,
  dropData: PaletteDropEventData,
) => {
  if (dragData.type === PaletteItemDragType.Section) {
    dispatch({
      type: "AddSection",
      pageId: dropData.parentId,
      position: 0,
    });
  } else {
    dispatch({
      type: "AddElement",
      elementType: dragData.type,
      sectionId: dropData.parentId,
      position: 0,
    });
  }
};
```

Also import `elementFactory.ts` here (or anywhere in the module graph above this call) so the
registry side effects run before `createElementFromType` is called.

### Step 7 — Wire `ItemTools.onReorder` to Move Events

Connect the stubs in `ElementItem` and `SectionItem`:

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

Requires threading `sectionId` into `ElementItem` and `pageId` into `SectionItem` as props,
which means updating `ElementList` and `SectionList` to pass them down.

### Step 8 — Framer Motion Animations

`DropZoneIndicator` now owns its own `AnimatePresence` via `isVisible` prop — callers render it
unconditionally and pass `isVisible={canDrop}`. The `motion` wrappers in `PageItem` and
`SectionItem` have been removed in favour of this self-contained approach.

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
| `FormDesigner/types/formDropEvent.ts` | **New** | ✅ | `FormDropEventData` / `PaletteDropEventData` |
| `FormDesigner/providers/FormDesignerDragProvider.tsx` | **New** | ✅ | `DragDropProvider` + `DragOverlay` + context + `useFormDragEvent` hook + typed `onDragEnd` switch |
| `FormDesigner/FormDesigner.tsx` | Modify | ✅ | Wrap panels with `FormDesignerDragProvider` |
| `FormDesigner/FormDesigner.styles.ts` | Modify | ✅ | Right panel column widened to `28rem` |
| `ToolboxPanel/palette.tsx` | Modify | ✅ | `PaletteItemDragType` enum + `dragType` field on every `IPaletteItem` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | ✅ | `useDraggable` with `type: item.dragType`; `draggable` prop; `isDragging` → `onDrag` style |
| `components/DragDrop/DraggableCard.tsx` | Moved | ✅ | Moved from `components/` to `components/DragDrop/`; optional `handleRef` prop |
| `components/DragDrop/DropZoneIndicator.tsx` | **New** | ✅ | `isVisible` owns `AnimatePresence`; `isDropTarget` drives background animation via Framer Motion variants + `useTheme` |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | ✅ | `useDroppable` + `useFormDragEvent` + `canDropItem` + `dragPaletteItem` memo + `DropZoneIndicator` |
| `CanvasPanel/lists/SectionItem.style.ts` | Modify | ✅ | `elements` style added for inner layout wrapper |
| `CanvasPanel/lists/PageItem.tsx` | Modify | ✅ | `useDroppable` (section drops) + same pattern as `SectionItem` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | ✅ | Stale `useDroppable` import removed; `gap: 2.5` added |
| `src/types/elementAttributes.ts` | Modify | ✅ | `ElementAttributes` union extended to all 9 types |
| `src/factories/elementFactory.ts` | **New** | ✅ | Registry-based abstract factory; all 9 types registered inline |
| `src/utils/id.ts` | **New** | ✅ | `generatedID()` + `isTemporaryID()` |
| `store/formDesigner/events.ts` | Modify | ✅ | `AddSectionEvent` has `pageId`; `AddElementEvent` has `sectionId` |
| `store/formDesigner/events.ts` | Modify | | `AddElementEvent` still needs `elementType: ElementType` |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | | Implement `onAddElement` (blocked on `elementType` field) |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | | Implement `onAddSection` + `createEmptySection` utility |
| `CanvasPanel/lists/ElementList.tsx` | Modify | | `AnimatePresence` wrapper |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | | `motion.li` + `layout="position"` + tween transition + enter/exit; wire `onReorder` |

**Total: 14 modified files, 6 new files, 1 moved. 17 of 23 complete.**

---

## Decisions / Open Questions

- **`handleRef` null vs undefined**: `PaletteItem` passes `null` when `draggable={false}`. Change to `undefined` to match `Ref<Element>` type.
- **Drop position**: Currently dispatches `position: 0`. Precise insertion would require additional droppable hit-testing per element/section.
- **Cross-section drag of existing elements**: Natural next step — `useSortable` from `@dnd-kit/react/sortable` with `group` for cross-section moves.
