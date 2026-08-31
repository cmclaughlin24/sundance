import type { HasPosition } from "./hasPosition";

export interface HasID {
  id: string;
}

export interface HasKey {
  key: string;
}

export interface HasName {
  name: string;
}

export interface IFormObject extends HasID, HasKey, HasName, HasPosition {}
