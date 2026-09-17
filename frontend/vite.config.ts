/// <reference types="vitest/config" />
import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Vendor groups change rarely; keeping them in their own chunks keeps the app cache warm across deploys.
const vendorGroups = [
  { name: "react", test: /node_modules[\\/](react|react-dom|react-router|scheduler)[\\/]/ },
  { name: "carbon-icons", test: /node_modules[\\/]@carbon[\\/]icons-react[\\/]/ },
  { name: "carbon", test: /node_modules[\\/]@carbon[\\/]/ },
  { name: "query", test: /node_modules[\\/](@tanstack|axios)[\\/]/ },
  { name: "thunderid", test: /node_modules[\\/]@thunderid[\\/]/ },
];

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
  },
  test: {
    environment: "jsdom",
    include: ["src/**/*.test.{ts,tsx}"],
  },
  build: {
    target: "es2022",
    rollupOptions: {
      output: {
        codeSplitting: { groups: vendorGroups },
      },
    },
  },
});
