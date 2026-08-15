import type { HasDataSourceRef, ILookup, LookupValue } from "./data";

export interface BaseFieldElementAttributes<T> {
  isRequired: boolean;
  isReadOnly: boolean;
  defaultValue: T | undefined;
}

export interface TextElementAttributes extends BaseFieldElementAttributes<string> {
  minLength?: number;
  maxLength?: number;
  pattern?: string;
  placeholder?: string;
}

export interface NumberElementAttributes extends BaseFieldElementAttributes<number> {
  min?: number;
  max?: number;
  step?: number;
}

export interface SelectElementAttributes
  extends BaseFieldElementAttributes<LookupValue>, HasDataSourceRef {
  data: ILookup[];
  multiple: boolean;
  minSelected?: number;
  maxSelected?: number;
}

export interface CheckboxElementAttributes
  extends BaseFieldElementAttributes<LookupValue>, HasDataSourceRef {
  isCheckedByDefault: boolean;
  data: ILookup[];
}

export interface DateElementAttributes extends BaseFieldElementAttributes<any> {
  minDate?: string;
  maxDate?: string;
}

export interface SegmentedElementAttributes extends BaseFieldElementAttributes<LookupValue> {
  data: ILookup[];
}

export interface RadioElementAttributes extends BaseFieldElementAttributes<LookupValue> {
  data: ILookup[];
  orientation: "horizontal" | "vertical";
}

export interface ToggleElementAttributes extends BaseFieldElementAttributes<boolean> {}

export interface UserElementAttributes extends BaseFieldElementAttributes<void> {
  canIncludeSelf: boolean;
  multiple: boolean;
}

export type ElementAttributes =
  | TextElementAttributes
  | NumberElementAttributes
  | SelectElementAttributes
  | CheckboxElementAttributes
  | DateElementAttributes;
