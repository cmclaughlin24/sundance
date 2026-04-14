import { MfeBootstrapComponent } from "@/components/MfeBootstrapComponent";
import { createFileRoute } from "@tanstack/react-router";
import { bootstrap } from "authentication/bootstrap";

const Authentication = MfeBootstrapComponent("/authentication", bootstrap);

export const Route = createFileRoute("/authentication/$")({
  component: Authentication,
});
