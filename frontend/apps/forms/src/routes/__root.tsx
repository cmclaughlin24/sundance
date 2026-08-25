import { Outlet, createRootRoute } from "@tanstack/react-router";
import { MainContainer } from "@/components/layout/MainContainer";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <MainContainer>
      <Outlet />
    </MainContainer>
  );
}
