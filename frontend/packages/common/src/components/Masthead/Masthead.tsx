import type { SundanceMastheadBrandingElement } from "./MastheadBranding.tsx";
import SundanceMastheadBranding from "./MastheadBranding.tsx";
import SundanceMastheadContent, {
  type SundanceMastheadContentElement,
} from "./MastheadContent.tsx";
import SundanceMastheadBrandingImg from "./MastheadBrandingImg.tsx";
import { Masthead } from "./Masthead.styles.ts";
import Container from "@mui/material/Container";
import Toolbar from "@mui/material/Toolbar";
import { PropsWithClassName } from "@/types/props-with-class-name.ts";
import SundanceMastheadLink from "./MastheadLink.tsx";

/**
 * A list of all the valid children of the `SundanceMasthead` component. Used to ensure
 * that only valid children are passed to the `SundanceMasthead` component.
 */
type SundanceMastheadChildren =
  | SundanceMastheadBrandingElement
  | SundanceMastheadContentElement;

interface SundanceMastheadComponent extends React.FC<
  PropsWithClassName<{
    children?: SundanceMastheadChildren | SundanceMastheadChildren[];
  }>
> {
  /**
   * A sub-component of `SundanceMasthead` that is used to render the branding of the app, such
   * as the logo and company name.
   *
   * @example
   * ```tsx
   * <SundanceMasthead.Branding>
   *    <img src="/path/to/image.jpg" alt="Super Awesome Brand Logo" />
   * </SundanceMasthead.Branding>
   * ```
   */
  Branding: typeof SundanceMastheadBranding;

  /**
   * A sub-component of `SundanceMasthead` that is used to render the image of the branding, such
   * as the logo of the company.
   *
   * @example
   * ```tsx
   * <SundanceMasthead.BrandingImg src="/path/to/image.jpg" alt="Super Awesome Brand Logo" />
   * ```
   */
  BrandingImg: typeof SundanceMastheadBrandingImg;

  /**
   * A sub-component of `SundanceMasthead` that is used to render the main content of the header
   * such as the main navigation, search bar, etc.
   *
   * @example
   * ```tsx
   * <SundanceMasthead.Content>
   *  <nav>
   *    <ul>
   *      <li>Home</li>
   *      <li>About</li>
   *    </ul>
   *  </nav>
   * </SundanceMasthead.Content>
   * ```
   */
  Content: typeof SundanceMastheadContent;

  Link: typeof SundanceMastheadLink;
}

const SundanceMasthead: SundanceMastheadComponent = function ({
  children,
  className,
}) {
  return (
    <Masthead className={className} position="relative">
      <Container maxWidth="xl">
        <Toolbar disableGutters>{children}</Toolbar>
      </Container>
    </Masthead>
  );
};

SundanceMasthead.Branding = SundanceMastheadBranding;
SundanceMasthead.BrandingImg = SundanceMastheadBrandingImg;
SundanceMasthead.Content = SundanceMastheadContent;
SundanceMasthead.Link = SundanceMastheadLink;

export default SundanceMasthead;
