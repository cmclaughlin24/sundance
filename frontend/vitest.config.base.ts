import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "happy-dom",
    coverage: {
      provider: "v8",
      reporter: ["html", "json", "text", "json-summary"],
      thresholds: {
        branches: 75,
        lines: 75,
        functions: 75,
        statements: 75,
      },
    },
  },
});
