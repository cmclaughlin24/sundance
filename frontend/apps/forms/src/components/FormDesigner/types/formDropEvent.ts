export interface PaletteDropEventData {
  source: "palette";
  parentId: string;
  position: number;
}

export type FormDropEventData = PaletteDropEventData;
