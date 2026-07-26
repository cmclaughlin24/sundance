import Box from "@mui/material/Box";
import type { PropsWithChildren } from "react";
import { FieldElementLabel } from "./FieldElementLabel";
import type { IElement } from "@/types/element";
import { fieldElementContainerStyles } from "./FieldElementContainer.style";

export const FieldElementContainer: React.FC<
  PropsWithChildren<{ element: IElement }>
> = function ({ element, children }) {
  return (
    <Box sx={fieldElementContainerStyles["container"]}>
      <Box sx={fieldElementContainerStyles["label"]}>
        <FieldElementLabel
          label={element.name}
          description=""
          htmlFor={element.id}
        />
      </Box>
      <Box sx={fieldElementContainerStyles["input"]}>{children}</Box>
    </Box>
  );
};
