import { Link } from "react-router";
import { Button, InlineNotification, InlineLoading } from "@carbon/react";
import { CheckmarkFilled, ArrowRight, Layers, Book, UserMultiple, EventSchedule, ChevronRight } from "@carbon/icons-react";

const NEXT_STEPS = [
  {
    title: "Curriculum",
    body: "Define levels (a grade, a stream, an exam class) and selection groups to control which subjects students can pick.",
    path: "/curriculum",
    icon: Layers,
  },
  {
    title: "Subjects",
    body: "Build your subject catalogue, then offer subjects to students through curriculum selection groups.",
    path: "/subjects",
    icon: Book,
  },
  {
    title: "Students & Teachers",
    body: "Enrol students and add teachers - each gets their own sign-in and profile automatically.",
    path: "/students",
    icon: UserMultiple,
  },
  {
    title: "Attendance",
    body: "Once classes have students, teachers can create sessions and mark attendance from the class page.",
    path: "/attendance",
    icon: EventSchedule,
  },
];

interface Props {
  submitting: boolean;
  submitError: string | null;
  submitted: boolean;
  onRetry: () => void;
  onGoToDashboard: () => void;
}

export default function DoneStep({ submitting, submitError, submitted, onRetry, onGoToDashboard }: Props) {
  return (
    <div>
      {submitting && (
        <div style={{ display: "flex", justifyContent: "center", padding: "3.5rem 0" }}>
          <InlineLoading description="Saving your school setup…" />
        </div>
      )}

      {!submitting && submitError && (
        <div style={{ textAlign: "center" }}>
          <InlineNotification
            kind="error"
            title="Could not finish setup"
            subtitle={submitError}
            hideCloseButton
            lowContrast
            style={{ marginBottom: "1.25rem", maxWidth: "100%", textAlign: "left" }}
          />
          <Button onClick={onRetry} style={{ width: "100%", maxWidth: "100%" }}>
            Retry
          </Button>
        </div>
      )}

      {!submitting && !submitError && submitted && (
        <>
          <div style={{ textAlign: "center", marginBottom: "1.75rem" }}>
            <div
              style={{
                width: "3.5rem",
                height: "3.5rem",
                margin: "0 auto 1rem",
                borderRadius: "50%",
                background: "var(--os-status-present-bg)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <CheckmarkFilled size={28} style={{ fill: "var(--os-success)" }} />
            </div>
            <h2 style={{ margin: "0 0 0.35rem", fontSize: "1.25rem", fontWeight: 600, color: "var(--os-text-primary)" }}>
              Your school is ready
            </h2>
            <p style={{ margin: 0, fontSize: "0.875rem", color: "var(--os-text-secondary)" }}>
              Here's what OpenSchool helps you run day to day, and where to go next.
            </p>
          </div>

          <div style={{ display: "grid", gap: "0.75rem", marginBottom: "1.5rem" }}>
            {NEXT_STEPS.map((item) => (
              <Link
                key={item.title}
                to={item.path}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "0.875rem",
                  padding: "0.875rem 1rem",
                  border: "1px solid var(--os-border-subtle)",
                  background: "var(--os-layer-hover)",
                  textDecoration: "none",
                  transition: "border-color 0.15s ease, background 0.15s ease",
                }}
              >
                <div
                  style={{
                    width: "2.25rem",
                    height: "2.25rem",
                    borderRadius: "50%",
                    background: "var(--os-accent-light)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    flexShrink: 0,
                  }}
                >
                  <item.icon size={18} style={{ fill: "var(--os-accent)" }} />
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <p style={{ margin: "0 0 0.25rem", fontSize: "0.875rem", fontWeight: 600, color: "var(--os-text-primary)" }}>
                    {item.title}
                  </p>
                  <p style={{ margin: 0, fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>{item.body}</p>
                </div>
                <ChevronRight size={16} style={{ fill: "var(--os-text-tertiary)", flexShrink: 0 }} />
              </Link>
            ))}
          </div>

          <Button renderIcon={ArrowRight} kind="primary" onClick={onGoToDashboard} style={{ width: "100%", maxWidth: "100%" }}>
            Go to Dashboard
          </Button>
        </>
      )}
    </div>
  );
}
