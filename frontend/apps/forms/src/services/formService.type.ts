import type { ElementType } from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";
import type {
  RuleExpressionJoinOp,
  RuleExpressionOp,
  RuleExprSourceType,
  RuleType,
} from "@/types/rule";

export interface FormRequest {
  name: string;
  description: string;
}

export interface FormVersionRequest {
  metadata: Record<string, string>;
  pages: PageRequest[];
}

export interface PageRequest {
  id?: string;
  key: string;
  name: string;
  position: number;
  sections: SectionRequest[];
  rules: RuleRequest[];
}

export interface SectionRequest {
  id?: string;
  key: string;
  name: string;
  position: number;
  elements: ElementRequest[];
  rules: RuleRequest[];
}

export interface ElementRequest {
  id?: string;
  key: string;
  type: ElementType;
  name: string;
  description: string;
  position: number;
  attributes: ElementAttributes;
  rules: RuleRequest[];
  tags: ElementTagMappingRequest[];
}

export interface RuleRequest {
  id?: string;
  type: RuleType;
  expressions: RuleExpressionRequest[];
}

export interface RuleExpressionRequest {
  source: { type: RuleExprSourceType; key: string };
  operator: RuleExpressionOp;
  joinWithPrevious?: RuleExpressionJoinOp;
  value: any;
  position: number;
}

export interface ElementTagMappingRequest {
  tagVersionId: string;
  priority: number;
  hasStaticValue: boolean;
  staticValue: any;
}
