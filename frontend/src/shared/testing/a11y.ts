import axe from "axe-core";
import { expect } from "vitest";

// jsdom has no CSS layout engine, so color-contrast (and a few other
// layout-dependent rules) can't run meaningfully here - this catches
// missing/invalid labels, bad ARIA and landmark/structure issues instead.
const DISABLED_RULES = ["color-contrast"];

export async function expectNoA11yViolations(container: Element) {
  const results = await axe.run(container, {
    rules: Object.fromEntries(DISABLED_RULES.map((id) => [id, { enabled: false }])),
  });
  const summary = results.violations.map((v) => `${v.id}: ${v.help} (${v.nodes.length} node(s))`).join("\n");
  expect(results.violations, summary).toHaveLength(0);
}
