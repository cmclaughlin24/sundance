import { defineConfig } from "vite";
import { federation } from "@module-federation/vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";
import { dependencies } from "./package.json";

// https://vite.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  plugins: [
    devtools(),
    tanstackRouter({ target: "react", autoCodeSplitting: true }),
    react(),
    federation({
      dts: false,
      name: "authentication",
      filename: "remoteEntry.js",
      exposes: {
        "./bootstrap": "./src/bootstrap.tsx",
      },
      shared: {
        react: {
          singleton: true,
          requiredVersion: dependencies["react"],
        },
        "react-dom": {
          singleton: true,
          requiredVersion: dependencies["react-dom"],
        },
      },
    }),
  ],
  server: {
    proxy: {
      "/otlp": {
        target: process.env.OTLP_HOST || "http://localhost:4318",
        rewrite: (path) => path.replace(/^\/otlp/, ""),
      },
    },
    fs: {
      allow: [".", "../shared"],
    },
  },
});
