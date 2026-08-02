import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/designer/")({
  component: DesignerRouteComponent,
});

function DesignerRouteComponent() {
  return <div>Hello "/designer/"!</div>;
}
