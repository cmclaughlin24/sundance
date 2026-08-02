import { type ElementType, type IElement } from "@/types/element";
import { TextFieldElement } from "../Elements/TextFieldElement";
import {
  useElementRuleState,
  useSubmissionDispatch,
} from "@/store/submission/useSubmissionContext";
import { NumberFieldElement } from "../Elements/NumberFieldElement";
import type { IRuleState } from "@/types/rule";
import { SelectFieldElement } from "../Elements/SelectFieldElement";
import { CheckboxFieldElement } from "../Elements/CheckboxFieldElement";
import { DateFieldElement } from "../Elements/DateFieldElement";
import { motion } from "motion/react";
import { elementVariants } from "./renderer.animations";

export type ElementComponent = React.FC<{
  element: IElement;
  ruleState: IRuleState;
  onChange: (value: any) => void;
}>;

const registry = new Map<ElementType, ElementComponent>([
  ["checkbox", CheckboxFieldElement],
  ["date", DateFieldElement],
  ["number", NumberFieldElement],
  ["select", SelectFieldElement],
  ["text", TextFieldElement],
]);

export const ElementRenderer: React.FC<{ element: IElement }> = function ({
  element,
}) {
  const { setValue } = useSubmissionDispatch();
  const ruleState = useElementRuleState(element);

  const handleChange = (value: any) => {
    setValue(element.id, value);
  };

  const Component = registry.get(element.type);

  if (!Component) {
    return <>Element not defined!</>;
  }

  return (
    <motion.div variants={elementVariants} initial="initial" animate="animate" exit="exit">
      <Component
        element={element}
        ruleState={ruleState}
        onChange={handleChange}
      />
    </motion.div>
  );
};
