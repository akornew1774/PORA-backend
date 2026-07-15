DROP INDEX IF EXISTS idx_devices_user_id;

DROP TABLE IF EXISTS devices;

ALTER TABLE items
DROP COLUMN IF EXISTS checked_at,
DROP COLUMN IF EXISTS next_reminder_at;