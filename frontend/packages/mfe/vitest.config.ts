import { mergeConfig, defineConfig } from "vitest/config";
import viteConfig from "./vite.config";
import baseConfig from "../../vitest.config.base";

export default mergeConfig(
  mergeConfig(viteConfig, baseConfig),
  defineConfig({
    test: {
      coverage: {
        thresholds: {
          branches: 75,
          lines: 75,
          functions: 75,
          statements: 75,
        },
      },
    },
  }),
);
