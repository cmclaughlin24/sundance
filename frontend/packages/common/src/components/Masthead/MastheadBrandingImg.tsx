import type { ComponentPropsWithoutRef } from "react";
import { BrandingImg } from "./Masthead.styles";

export type SundanceMastheadBrandingImgProps = ComponentPropsWithoutRef<"img">;

const SundanceMastheadBrandingImg: React.FC<SundanceMastheadBrandingImgProps> = function ({
  className,
  ...props
}) {
  return <BrandingImg className={className} {...props} />;
};

export default SundanceMastheadBrandingImg;
