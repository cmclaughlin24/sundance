import { createBrowserHistory } from "@tanstack/react-router";
import { initTelemetry } from "./telemetry.ts";
import { bootstrap } from "./bootstrap.tsx";
import "./index.css";

initTelemetry();

bootstrap(document.getElementById("app"), {
  defaultHistory: createBrowserHistory(),
});
