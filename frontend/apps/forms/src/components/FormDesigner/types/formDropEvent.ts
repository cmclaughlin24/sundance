export interface PaletteDropEventData {
  source: "palette";
  parentId: string;
  position: number;
}

export interface CanvasDropEventData {
  source: "canvas";
  parentId: string;
  position: number;
}

export type FormDropEventData = PaletteDropEventData | CanvasDropEventData;
