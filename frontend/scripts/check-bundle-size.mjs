// Size budget check (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 5):
// fails if the app's own JS entry chunk exceeds 100 KB or total CSS exceeds
// 400 KB. Deliberately dependency-free (no vite-bundle-visualizer) so it
// runs with nothing beyond `pnpm build` having already produced dist/.
//
// Usage: node scripts/check-bundle-size.mjs [--warn-only]
import { readdirSync, statSync } from "node:fs";
import { join } from "node:path";

const DIST_ASSETS = join(import.meta.dirname, "..", "dist", "assets");
const APP_CHUNK_BUDGET_BYTES = 100 * 1024;
const TOTAL_CSS_BUDGET_BYTES = 400 * 1024;
const warnOnly = process.argv.includes("--warn-only");

function sizeOf(file) {
  return statSync(join(DIST_ASSETS, file)).size;
}

let entries;
try {
  entries = readdirSync(DIST_ASSETS);
} catch {
  console.error(`Could not read ${DIST_ASSETS} — run "pnpm build" first.`);
  process.exit(1);
}

// Vite names the entry chunk after the app's own main.tsx; vendor/feature
// chunks get their own hashed names, so "index-*.js" is specifically ours.
const appChunks = entries.filter((f) => /^index-.*\.js$/.test(f));
const cssFiles = entries.filter((f) => f.endsWith(".css"));

const appChunkBytes = appChunks.reduce((sum, f) => sum + sizeOf(f), 0);
const totalCssBytes = cssFiles.reduce((sum, f) => sum + sizeOf(f), 0);

const kb = (bytes) => (bytes / 1024).toFixed(1);
console.log(`App entry chunk: ${kb(appChunkBytes)} KB (budget ${kb(APP_CHUNK_BUDGET_BYTES)} KB)`);
console.log(`Total CSS: ${kb(totalCssBytes)} KB (budget ${kb(TOTAL_CSS_BUDGET_BYTES)} KB)`);

const overBudget = appChunkBytes > APP_CHUNK_BUDGET_BYTES || totalCssBytes > TOTAL_CSS_BUDGET_BYTES;
if (overBudget) {
  const message = "One or more bundle size budgets were exceeded.";
  if (warnOnly) {
    console.warn(`::warning::${message}`);
  } else {
    console.error(message);
    process.exit(1);
  }
}
