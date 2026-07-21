import type { HasDataSourceRef } from "./data";

export interface BaseFieldElementAttributes {
  isRequired: boolean;
  isReadOnly: boolean;
}

export interface TextElementAttributes extends BaseFieldElementAttributes {
  minLength?: number;
  maxLength?: number;
  pattern?: string;
  placeholder?: string;
}

export interface NumberElementAttributes extends BaseFieldElementAttributes {
  min?: number;
  max?: number;
  step?: number;
}

export interface SelectElementAttributes
  extends BaseFieldElementAttributes, HasDataSourceRef {
  data: any[];
  multiple: boolean;
  minSelected?: number;
  maxSelected?: number;
}

export interface CheckboxElementAttributes
  extends BaseFieldElementAttributes, HasDataSourceRef {
  isCheckedByDefault: boolean;
  data: any[];
}

export interface DateElementAttributes extends BaseFieldElementAttributes {
  minDate?: string;
  maxDate?: string;
}

export type ElementAttributes =
  | TextElementAttributes
  | NumberElementAttributes
  | SelectElementAttributes
  | CheckboxElementAttributes
  | DateElementAttributes;
