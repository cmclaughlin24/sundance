import type { PropsWithChildren, ReactElement } from "react";
import { Content } from "./Masthead.styles";

export type SundanceMastheadContentProps = PropsWithChildren<{
  className?: string;
}>;

export type SundanceMastheadContentElement =
  ReactElement<SundanceMastheadContentProps>;

const SundanceMastheadContent: React.FC<SundanceMastheadContentProps> = ({
  children,
  className,
}) => {
  return <Content className={className}>{children}</Content>;
};

export default SundanceMastheadContent;
