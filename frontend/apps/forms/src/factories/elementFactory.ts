import type { ElementType, IElement } from "@/types/element";
import type {
  CheckboxElementAttributes,
  DateElementAttributes,
  NumberElementAttributes,
  RadioElementAttributes,
  SegmentedElementAttributes,
  SelectElementAttributes,
  TextElementAttributes,
  ToggleElementAttributes,
  UserElementAttributes,
} from "@/types/elementAttributes";

export type ElementFactory = (id: string) => IElement;

const registry = new Map<ElementType, ElementFactory>();

export function createElementFromType(
  elementType: ElementType,
  id: string,
): IElement {
  const factory = registry.get(elementType);

  if (!factory) {
    throw new Error(`No factory registered for element type: ${elementType}`);
  }

  return factory(id);
}

function registerElementFactory(
  elementType: ElementType,
  factory: ElementFactory,
) {
  registry.set(elementType, factory);
}

const base = { isRequired: false, isReadOnly: false, defaultValue: undefined };

registerElementFactory("text", (id) => ({
  id,
  type: "text",
  name: "Text",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: { ...base } satisfies TextElementAttributes,
}));

registerElementFactory("number", (id) => ({
  id,
  type: "number",
  name: "Number",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: { ...base } satisfies NumberElementAttributes,
}));

registerElementFactory("date", (id) => ({
  id,
  type: "date",
  name: "Date",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: { ...base } satisfies DateElementAttributes,
}));

registerElementFactory("toggle", (id) => ({
  id,
  type: "toggle",
  name: "Toggle",
  description: "",
  key: "",
  position: 0,
  rules: [],
  tags: [],
  attributes: { ...base } satisfies ToggleElementAttributes,
}));

registerElementFactory("checkbox", (id) => ({
  id,
  type: "checkbox",
  name: "Checkbox",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: {
    ...base,
    isCheckedByDefault: false,
    data: [],
  } satisfies CheckboxElementAttributes,
}));

registerElementFactory("radio", (id) => ({
  id,
  type: "radio",
  name: "Radio",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: {
    ...base,
    data: [],
    orientation: "vertical",
  } satisfies RadioElementAttributes,
}));

registerElementFactory("select", (id) => ({
  id,
  type: "select",
  name: "Select",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: {
    ...base,
    data: [],
    multiple: false,
  } satisfies SelectElementAttributes,
}));

registerElementFactory("segmented", (id) => ({
  id,
  type: "segmented",
  name: "Segmented",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: { ...base, data: [] } satisfies SegmentedElementAttributes,
}));

registerElementFactory("user", (id) => ({
  id,
  type: "user",
  name: "User",
  description: "",
  key: id,
  position: 0,
  rules: [],
  tags: [],
  attributes: {
    ...base,
    canIncludeSelf: false,
    multiple: false,
  } satisfies UserElementAttributes,
}));
