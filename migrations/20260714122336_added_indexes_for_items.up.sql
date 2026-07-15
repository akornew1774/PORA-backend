CREATE INDEX idx_items_next_reminder
ON items(next_reminder_at);

CREATE INDEX idx_items_checked_at
ON items(checked_at);