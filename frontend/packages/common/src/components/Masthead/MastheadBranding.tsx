import type { PropsWithChildren, ReactElement } from "react";
import { Branding } from "./Masthead.styles";

export type SundanceMastheadBrandingProps = PropsWithChildren<{
  className?: string;
}>;

export type SundanceMastheadBrandingElement =
  ReactElement<SundanceMastheadBrandingProps>;

const SundanceMastheadBranding: React.FC<SundanceMastheadBrandingProps> = ({
  children,
  className,
}) => {
  return <Branding className={className}>{children}</Branding>;
};

export default SundanceMastheadBranding;
