import { Link, Outlet, createRootRoute } from "@tanstack/react-router";
import logo from "../assets/train-bandits.png";
import { SundanceMasthead } from "@sundance/common";
import Box from "@mui/material/Box";
import Sidebar from "@/components/ui/Sidebar";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <SundanceMasthead>
        <SundanceMasthead.Branding>
          <Link to="/">
            <SundanceMasthead.BrandingImg src={logo} alt="Sundance" />
          </Link>
        </SundanceMasthead.Branding>
        <SundanceMasthead.Content>
          <Box sx={{ flex: 1 }}>
            <SundanceMasthead.Link
              to="/authentication"
              from="/"
              component={Link}
            >
              Authentication
            </SundanceMasthead.Link>
            <SundanceMasthead.Link
              to="/certification"
              from="/"
              component={Link}
            >
              Certification
            </SundanceMasthead.Link>
          </Box>
          <Box sx={{ display: "flex", justifyContent: "end" }}>
            <SundanceMasthead.Link
              to="/authentication/login"
              from="/"
              component={Link}
            >
              Login
            </SundanceMasthead.Link>
          </Box>
        </SundanceMasthead.Content>
      </SundanceMasthead>
      <div className="flex-1 flex overflow-hidden h-full">
        <Sidebar className="shrink-0" />
        <main className="flex-1">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
