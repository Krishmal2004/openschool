-- Distinguishes "chose to keep the default password" from "changed it" so
-- must_change_password's meaning stays a single bit while the backend can
-- still expire that choice after a week (S1): must_change_password already
-- goes FALSE both when a user sets a real password and when they click
-- "keep this password", so nothing previously told those two cases apart.
-- kept_default_password is set TRUE only by the latter, and cleared again
-- if the user later does set a real password.
ALTER TABLE users ADD COLUMN kept_default_password BOOLEAN NOT NULL DEFAULT FALSE;
