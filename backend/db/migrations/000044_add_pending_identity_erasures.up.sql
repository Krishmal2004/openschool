-- Backs a retry ledger for identity cleanup that fails after a student
-- profile has already been anonymised (S11): both the manual "erase person"
-- flow and the nightly retention purge scrub the local user record and
-- delete the identity-provider account as a second step, and either can fail
-- independently of the anonymisation that already committed. Without this
-- table a failure here was only logged, so the account's PII could survive
-- indefinitely with no way to know it needed a retry.
CREATE TABLE pending_identity_erasures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    local_done BOOLEAN NOT NULL DEFAULT FALSE,
    idp_done BOOLEAN NOT NULL DEFAULT FALSE,
    last_error TEXT,
    attempts INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
