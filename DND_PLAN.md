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
pnpm add @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
```

| Package | Purpose |
|---|---|
| `@dnd-kit/core` | `DndContext`, `useDraggable`, `useDroppable`, sensors, `DragOverlay` |
| `@dnd-kit/sortable` | `SortableContext`, `useSortable` for reordering within lists |
| `@dnd-kit/utilities` | CSS transform helpers (`CSS.Translate.toString`) |

---

## Implementation Steps

### Step 1 — Install @dnd-kit

```sh
pnpm --filter @sundance/forms add @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
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

### Step 5 — Set up `DndContext` in `FormDesigner.tsx`

Wrap the entire designer in `DndContext` at the top level so drag state is shared across both panels.
Configure a `PointerSensor` with a small activation distance (8px) to preserve click behavior on palette items.
Render a `DragOverlay` to show a floating clone of the dragged card.

```tsx
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";

export const FormDesigner: React.FC<FormDesignerProps> = ({ form, version }) => {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } })
  );

  const handleDragEnd = ({ active, over }) => {
    if (!over) return;
    const source = active.data.current;
    const target = over.data.current;

    if (source?.source === "palette" && target?.type === "section") {
      dispatch({
        type: "AddElement",
        elementId: crypto.randomUUID(),
        elementType: source.elementType,
        sectionId: target.sectionId,
        position: target.elementCount,
      });
    }
  };

  return (
    <FormDesignerProvider form={form} version={version}>
      <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
        <Box sx={formDesignerStyles.container}>
          {/* panels */}
        </Box>
        <DragOverlay>
          {/* floating clone rendered here while dragging */}
        </DragOverlay>
      </DndContext>
    </FormDesignerProvider>
  );
};
```

> Note: `dispatch` must be obtained inside the provider via a child component or context bridge since `FormDesignerProvider` wraps `DndContext`.

### Step 6 — Make `PaletteItem` Draggable

Use `useDraggable` in `PaletteItem`. Apply `listeners` to the `DraggableCard`'s drag handle icon
and `setNodeRef` to the root element. The drag data payload identifies this as a palette drag:

```ts
const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
  id: `palette-${item.type}`,
  data: { source: "palette", elementType: item.type },
});
```

- Apply `isDragging` to reduce opacity on the source item while dragging
- Pass `listeners` down to `DraggableCard` so only the drag handle activates the drag (not the whole card)

### Step 7 — Make `SectionItem` a Drop Zone

Use `useDroppable` in `SectionItem`:

```ts
const { setNodeRef, isOver } = useDroppable({
  id: `section-${section.id}`,
  data: { type: "section", sectionId: section.id, elementCount: section.elements.length },
});
```

Pass `isOver` into `getSectionItemStyles` alongside `isSelected` to render a distinct "ready to drop"
highlight (e.g. blue dashed border) when a palette item is dragged over.

### Step 8 — Handle Drop in `onDragEnd`

The `onDragEnd` callback (in `DndContext` — Step 5) checks:

1. Was the drag source a palette item? (`active.data.current.source === "palette"`)
2. Was the drop target a section? (`over.data.current.type === "section"`)
3. If yes to both → dispatch `AddElementEvent`

For future support of reordering existing elements within/between sections, also handle:
- Source is an existing element (`source === "element"`) + target is a section → dispatch `MoveElementEvent`

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

| File | Type | Change |
|---|---|---|
| `package.json` | Modify | Add `@dnd-kit/core`, `@dnd-kit/sortable`, `@dnd-kit/utilities` |
| `store/formDesigner/events.ts` | Modify | Extend `AddElementEvent` (`elementType`, `sectionId`); extend `AddSectionEvent` (`pageId`) |
| `store/formDesigner/elementFactory.ts` | **New** | `createElementFromType()` factory function |
| `store/formDesigner/eventHandlers/elementEventHandlers.ts` | Modify | Implement `onAddElement` stub |
| `store/formDesigner/eventHandlers/sectionEventHandlers.ts` | Modify | Implement `onAddSection` stub |
| `FormDesigner/FormDesigner.tsx` | Modify | Add `DndContext`, `DragOverlay`, `onDragEnd` handler, `PointerSensor` |
| `ToolboxPanel/PaletteItem.tsx` | Modify | Add `useDraggable`, pass `listeners` to `DraggableCard` |
| `components/DraggableCard.tsx` | Modify | Accept and forward drag `listeners` + `attributes` props |
| `CanvasPanel/lists/SectionItem.tsx` | Modify | Add `useDroppable`, pass `isOver` to styles |
| `CanvasPanel/lists/SectionItem.style.ts` | Modify | Add `isOver` style variant |
| `CanvasPanel/lists/ElementList.tsx` | Modify | Wrap with `AnimatePresence` |
| `CanvasPanel/lists/ElementItem.tsx` | Modify | Switch root to `motion.li` with `layout` + enter/exit animations; wire `onReorder` |
| `CanvasPanel/lists/SectionList.tsx` | Modify | Thread `pageId` prop down to `SectionItem` |

**Total: 11 modified files, 1 new file.**

---

## Decisions / Open Questions

- **Drop position within a section**: The plan above drops to the end of the section (`elementCount`). For precise insertion (before/after a specific element), `ElementItem` would also need `useDroppable` or the section would need to track hover position — this is a scope increase.
- **Cross-section drag of existing elements**: `MoveElement` supports cross-section moves. Wiring this up with @dnd-kit's `SortableContext` across multiple droppable containers is the natural next step after the toolbox drop is working.
- **Section creation from toolbox**: The palette has a "Section" type. Dropping it onto the canvas (outside any section) could dispatch `AddSectionEvent`. This requires the canvas area itself to also be a droppable zone.
- **`DraggableCard` drag handle forwarding**: The `listeners` from `useDraggable` need to be applied to the `DragIndicator` icon specifically (not the whole card) to avoid conflicts with click selection. This requires `DraggableCard` to accept and forward handle props.
