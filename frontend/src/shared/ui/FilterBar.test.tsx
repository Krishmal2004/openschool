import { describe, it } from "vitest";
import { render } from "@testing-library/react";
import FilterBar from "@/shared/ui/FilterBar";
import { expectNoA11yViolations } from "@/shared/testing/a11y";

describe("FilterBar", () => {
  it("has no accessibility violations with search and filter controls", async () => {
    const { container } = render(
      <FilterBar
        search={{ value: "", onChange: () => {}, placeholder: "Search…" }}
        controls={[
          { label: "Grade", node: <select aria-label="Grade"><option>All grades</option></select> },
          { label: "Status", node: <select aria-label="Status"><option>All statuses</option></select> },
        ]}
      />,
    );
    await expectNoA11yViolations(container);
  });

  it("has no accessibility violations with search only", async () => {
    const { container } = render(<FilterBar search={{ value: "", onChange: () => {} }} />);
    await expectNoA11yViolations(container);
  });
});
