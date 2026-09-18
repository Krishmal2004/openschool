-- Destructive: permanently discards every stored kept_default_password value.
ALTER TABLE users DROP COLUMN kept_default_password;
