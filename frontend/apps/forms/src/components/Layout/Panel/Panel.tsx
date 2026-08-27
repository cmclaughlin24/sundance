import Box from "@mui/material/Box";
import type { SxProps, Theme } from "@mui/material/styles";
import { PanelTitle } from "./PanelTitle";
import { mergeSx } from "merge-sx";

const styles: Readonly<SxProps<Theme>> = {
  padding: 2.5,
};

export type PanelProps = React.PropsWithChildren<{ sx?: SxProps<Theme> }>;

interface PanelComponent extends React.FC<PanelProps> {
    Title: typeof PanelTitle;
}

const Panel: PanelComponent = function ({ sx, children }) {
  return <Box sx={mergeSx(styles, sx)}>{children}</Box>;
};

Panel.Title = PanelTitle;

export default Panel;
