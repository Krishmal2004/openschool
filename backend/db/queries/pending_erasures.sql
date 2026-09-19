-- name: UpsertPendingErasure :exec
-- Records that a user still needs identity cleanup after its profile was
-- already anonymised. local_done/idp_done only ever move true → stays true
-- (the OR keeps a step already confirmed done from being un-done by a
-- retry that only re-attempts the other step).
INSERT INTO pending_identity_erasures (user_id, local_done, idp_done, last_error, attempts)
VALUES ($1, $2, $3, $4, 1)
ON CONFLICT (user_id) DO UPDATE SET
    local_done = pending_identity_erasures.local_done OR EXCLUDED.local_done,
    idp_done   = pending_identity_erasures.idp_done OR EXCLUDED.idp_done,
    last_error = EXCLUDED.last_error,
    attempts   = pending_identity_erasures.attempts + 1,
    updated_at = now();

-- name: ListPendingErasures :many
SELECT * FROM pending_identity_erasures
WHERE NOT (local_done AND idp_done)
ORDER BY created_at ASC
LIMIT $1;

-- name: DeletePendingErasure :exec
DELETE FROM pending_identity_erasures WHERE id = $1;
