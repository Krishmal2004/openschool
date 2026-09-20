import { readdirSync } from "node:fs";
import js from "@eslint/js";
import globals from "globals";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import tseslint from "typescript-eslint";
import { defineConfig, globalIgnores } from "eslint/config";

const features = readdirSync(new URL("./src/features", import.meta.url), { withFileTypes: true })
  .filter((d) => d.isDirectory())
  .map((d) => d.name);

const restrict = (patterns) => ["error", { patterns }];
const RULE = "@typescript-eslint/no-restricted-imports";

const BRITISH_SPELLING_WORDS = [
  ["Enrollment", "Enrolment"], ["enrollment", "enrolment"],
  ["Enrollments", "Enrolments"], ["enrollments", "enrolments"],
  ["Enroll", "Enrol"], ["enroll", "enrol"],
  ["Enrolls", "Enrols"], ["enrolls", "enrols"],
  ["Color", "Colour"], ["color", "colour"],
  ["Colors", "Colours"], ["colors", "colours"],
  ["Colored", "Coloured"], ["colored", "coloured"],
  ["Coloring", "Colouring"], ["coloring", "colouring"],
  ["Organization", "Organisation"], ["organization", "organisation"],
  ["Organizations", "Organisations"], ["organizations", "organisations"],
  ["Organize", "Organise"], ["organize", "organise"],
  ["Organized", "Organised"], ["organized", "organised"],
  ["Organizing", "Organising"], ["organizing", "organising"],
  ["Center", "Centre"], ["center", "centre"],
  ["Centers", "Centres"], ["centers", "centres"],
  ["Favorite", "Favourite"], ["favorite", "favourite"],
  ["Favorites", "Favourites"], ["favorites", "favourites"],
  ["Behavior", "Behaviour"], ["behavior", "behaviour"],
  ["Behaviors", "Behaviours"], ["behaviors", "behaviours"],
  ["Canceled", "Cancelled"], ["canceled", "cancelled"],
  ["Canceling", "Cancelling"], ["canceling", "cancelling"],
  ["Analyze", "Analyse"], ["analyze", "analyse"],
  ["Analyzed", "Analysed"], ["analyzed", "analysed"],
  ["Analyzing", "Analysing"], ["analyzing", "analysing"],
];

const BRITISH_SPELLING_LABEL_KEYS = new Set([
  "label", "labelText", "title", "subtitle", "placeholder", "description",
  "text", "aria-label", "ariaLabel", "fallback", "heading",
]);

function reportBritishSpelling(context, node, value) {
  for (const [american, british] of BRITISH_SPELLING_WORDS) {
    if (new RegExp(`\\b${american}\\b`).test(value)) {
      context.report({ node, message: `Use British spelling "${british}" instead of "${american}" in user-facing text.` });
    }
  }
}

const britishSpellingRule = {
  meta: { type: "suggestion" },
  create(context) {
    return {
      JSXText(node) {
        reportBritishSpelling(context, node, node.value);
      },
      JSXAttribute(node) {
        if (!BRITISH_SPELLING_LABEL_KEYS.has(node.name.name)) return;
        const value = node.value;
        if (value?.type === "Literal" && typeof value.value === "string") {
          reportBritishSpelling(context, node, value.value);
        } else if (value?.type === "JSXExpressionContainer" && value.expression.type === "Literal" && typeof value.expression.value === "string") {
          reportBritishSpelling(context, node, value.expression.value);
        }
      },
      Property(node) {
        const keyName = node.key.type === "Identifier" ? node.key.name : node.key.type === "Literal" ? node.key.value : null;
        if (!keyName || !BRITISH_SPELLING_LABEL_KEYS.has(keyName)) return;
        if (node.value.type === "Literal" && typeof node.value.value === "string") {
          reportBritishSpelling(context, node, node.value.value);
        }
      },
    };
  },
};

// Layer contract from docs/FRONTEND_REFACTOR_PLAYBOOK.md section 2.2, enforced per folder.
const layerRules = [
  {
    files: ["src/shared/**/*.{ts,tsx}"],
    rules: {
      [RULE]: restrict([
        { group: ["../*"], message: "Use the @/ alias instead of relative parent paths." },
        { group: ["@/features/*", "@/app/*", "@/layouts/*"], message: "shared/ must not depend on features, app or layouts." },
      ]),
    },
  },
  {
    files: ["src/features/*/api/**/*.ts"],
    rules: {
      [RULE]: restrict([
        { group: ["../*"], message: "Use the @/ alias instead of relative parent paths." },
        { group: ["react", "react-dom", "react-router", "@carbon/*", "@tanstack/*"], message: "api/ files are plain HTTP wrappers with no UI or hook imports." },
        { group: ["@/features/*/queries/*", "@/features/*/components/*", "@/features/*/pages/*", "@/shared/ui/*", "@/shared/hooks/*"], message: "api/ may import only other api types and shared/api." },
        { group: ["@/features/*/api/*"], allowTypeImports: true, message: "api/ files share types only, never call each other." },
      ]),
    },
  },
  {
    files: ["src/features/**/*.{ts,tsx}"],
    ignores: ["src/features/*/api/**"],
    rules: {
      [RULE]: restrict([
        { group: ["../*"], message: "Use the @/ alias instead of relative parent paths." },
        { group: ["axios", "@/shared/api/client"], message: "Only api/ files call the HTTP client." },
      ]),
    },
  },
  // A feature may use another feature's queries and components, never its pages or api.
  ...features.map((name) => ({
    files: [`src/features/${name}/**/*.{ts,tsx}`],
    ignores: [`src/features/${name}/api/**`],
    rules: {
      [RULE]: restrict([
        { group: ["../*"], message: "Use the @/ alias instead of relative parent paths." },
        { group: ["axios", "@/shared/api/client"], message: "Only api/ files call the HTTP client." },
        {
          group: features.filter((f) => f !== name).flatMap((f) => [`@/features/${f}/pages/*`, `@/features/${f}/api/*`]),
          allowTypeImports: true,
          message: "Import another feature's queries or components, not its pages or api (types are fine).",
        },
      ]),
    },
  })),
  {
    files: ["src/app/**/*.{ts,tsx}", "src/layouts/**/*.{ts,tsx}"],
    rules: {
      [RULE]: restrict([
        { group: ["../*"], message: "Use the @/ alias instead of relative parent paths." },
        { group: ["@/features/*/api/*"], allowTypeImports: true, message: "app/ and layouts/ compose pages and queries, never api." },
        { group: ["axios"], message: "Only shared/api and feature api/ files use axios." },
      ]),
    },
  },
];

export default defineConfig([
  globalIgnores(["dist", "coverage"]),
  {
    files: ["**/*.{ts,tsx}"],
    plugins: { local: { rules: { "british-spelling": britishSpellingRule } } },
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      globals: globals.browser,
    },
    rules: {
      [RULE]: restrict([{ group: ["../*"], message: "Use the @/ alias instead of relative parent paths." }]),
      // Pages and components over this limit are split into the feature's components/ or hooks/.
      "max-lines": ["error", { max: 250, skipBlankLines: true, skipComments: true }],
      "no-restricted-properties": [
        "error",
        { property: "toLocaleDateString", message: "Use formatDate/formatMonth/formatLongDate/... from @/shared/lib/date instead." },
        { property: "toLocaleString", message: "Use formatDateTime from @/shared/lib/date instead." },
        { property: "toLocaleTimeString", message: "Add a helper to @/shared/lib/date instead of calling this directly." },
      ],
      "local/british-spelling": "error",
    },
  },
  {
    files: ["src/shared/lib/date.ts"],
    rules: {
      "no-restricted-properties": "off",
    },
  },
  ...layerRules,
]);
