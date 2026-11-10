import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/data-sources/$dataSourceId/")({
  beforeLoad: async ({ params }) => {
    throw redirect({
      to: "/data-sources/$dataSourceId/{-$tab}",
      params: { dataSourceId: params.dataSourceId },
      replace: true,
    });
  },
});
