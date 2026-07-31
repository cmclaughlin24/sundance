import Box from "@mui/material/Box";
import type { PropsWithChildren } from "react";
import { FieldElementLabel } from "./FieldElementLabel";
import type { IElement } from "@/types/element";
import { fieldElementContainerStyles } from "./FieldElementContainer.style";
import { useElementRuleState } from "@/store/submission/useSubmissionContext";

export const FieldElementContainer: React.FC<
  PropsWithChildren<{ element: IElement }>
> = function ({ element, children }) {
  const ruleState = useElementRuleState(element);

  return (
    <Box sx={fieldElementContainerStyles["container"]}>
      <Box sx={fieldElementContainerStyles["label"]}>
        <FieldElementLabel
          label={element.name}
          description="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
          htmlFor={element.id}
          required={ruleState.required}
        />
      </Box>
      <Box sx={fieldElementContainerStyles["input"]}>{children}</Box>
    </Box>
  );
};
