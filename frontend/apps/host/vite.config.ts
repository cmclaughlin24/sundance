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
      name: "container",
      remotes: {
        authentication: {
          type: "module",
          name: "authentication",
          entry: "/authentication-ui/remoteEntry.js",
          shareScope: "default",
        },
        forms: {
          type: "module",
          name: "forms",
          entry: "/forms-ui/remoteEntry.js",
          shareScope: "default",
        },
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
      "/authentication-ui": {
        target: process.env.AUTHENTICATION_HOST || "http://localhost:3001",
        rewrite: (path) => path.replace(/^\/authentication-ui/, ""),
      },
      "/forms-ui": {
        target: process.env.FORMS_HOST || "http://localhost:3002",
        rewrite: (path) => path.replace(/^\/forms-ui/, ""),
      },
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
