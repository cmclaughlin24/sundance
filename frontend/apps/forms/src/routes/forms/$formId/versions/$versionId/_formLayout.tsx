import Box from "@mui/material/Box";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/forms/$formId/versions/$versionId/_formLayout")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <Box
      sx={{
        marginX: "auto",
        maxWidth: 1440,
      }}
    >
      <Outlet />
    </Box>
  );
}
