import Box from "@mui/material/Box";
import type { PropsWithChildren } from "react";
import { FieldElementLabel } from "./FieldElementLabel";
import type { IElement } from "@/types/element";

export const FieldElementContainer: React.FC<
  PropsWithChildren<{ element: IElement }>
> = function ({ element, children }) {
  return (
    <Box>
      <FieldElementLabel label={element.name} description="" htmlFor="" />
      <Box>{children}</Box>
    </Box>
  );
};
