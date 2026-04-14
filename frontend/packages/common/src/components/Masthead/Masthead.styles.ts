import styled from "@emotion/styled";
import AppBar from "@mui/material/AppBar";
import { SundanceTheme } from "../ThemeProvider/ThemeProvider";

export const Masthead = styled(AppBar)`
  border-bottom: 0.25rem solid ${SundanceTheme.palette.secondary.main};
`;

export const Branding = styled.div``;

export const BrandingImg = styled.img`
  max-height: 3rem;
`;

export const Content = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
`;
