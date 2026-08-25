import type { ElementType } from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";

export interface FormVersionRequest {
  metadata: Record<string, string>;
  pages: PageRequest[];
}

export interface PageRequest {
  id?: string;
  key: string;
  name: string;
  position: string;
  sections: SectionRequest[];
  rules: RuleRequest[];
}

export interface SectionRequest {
  id?: string;
  key: string;
  name: string;
  position: string;
  elements: ElementRequest[];
  rules: RuleRequest[];
}

export interface ElementRequest {
  id?: string;
  key: string;
  type: ElementType;
  name: string;
  description: string;
  position: string;
  attributes: ElementAttributes;
  tags: any[];
  rules: RuleRequest[];
}

export interface RuleRequest {}
