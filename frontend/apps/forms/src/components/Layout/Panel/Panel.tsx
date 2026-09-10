import Box from "@mui/material/Box";
import type { SxProps, Theme } from "@mui/material/styles";
import { PanelTitle } from "./PanelTitle";
import { mergeSx } from "merge-sx";
import type { RefObject } from "react";

const styles: Readonly<SxProps<Theme>> = {
  padding: 2.5,
};

export interface PanelProps
  extends React.PropsWithChildren, React.HTMLAttributes<HTMLDivElement> {
  sx?: SxProps<Theme>;
  ref?: RefObject<HTMLDivElement | null>;
}

interface PanelComponent extends React.FC<PanelProps> {
  Title: typeof PanelTitle;
}

const Panel: PanelComponent = function ({ sx, children, ref, ...props }) {
  return (
    <Box sx={mergeSx(styles, sx)} ref={ref} {...props}>
      {children}
    </Box>
  );
};

Panel.Title = PanelTitle;

export default Panel;
