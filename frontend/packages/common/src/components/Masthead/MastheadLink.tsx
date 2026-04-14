import Button from "@mui/material/Button";
import type { PropsWithChildren } from "react";

export type SundanceMastheadLinkProps = PropsWithChildren<{
  to: string;
  from: string;
  component: any;
}>;

const MastheadLink: React.FC<SundanceMastheadLinkProps> = function ({
  to,
  from,
  children,
  component: Component,
}) {
  return (
    <Button
      component={Component}
      to={to}
      from={from}
      sx={{ color: "white", textTransform: "none" }}
    >
      {children}
    </Button>
  );
};

export default MastheadLink;
