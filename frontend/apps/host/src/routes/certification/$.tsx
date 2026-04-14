import { createFileRoute } from "@tanstack/react-router";
import { bootstrap } from "certification/bootstrap";
import { MfeBootstrapComponent } from "@/components/MfeBootstrapComponent";

const Certification = MfeBootstrapComponent("/certification", bootstrap);

export const Route = createFileRoute("/certification/$")({
  component: Certification,
});
