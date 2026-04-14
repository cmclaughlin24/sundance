import { createFileRoute } from "@tanstack/react-router";
import { bootstrap } from "forms/bootstrap";
import { MfeBootstrapComponent } from "@/components/MfeBootstrapComponent";

const Forms = MfeBootstrapComponent("/forms", bootstrap);

export const Route = createFileRoute("/forms/$")({
  component: Forms,
});
