ALTER TABLE items
ADD COLUMN checked_at TIMESTAMPTZ,
ADD COLUMN next_reminder_at TIMESTAMPTZ;

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,

    device_token TEXT NOT NULL UNIQUE,
    device_type TEXT NOT NULL,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_devices_user_id ON devices(user_id);