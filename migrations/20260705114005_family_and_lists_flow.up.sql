CREATE TABLE families (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,

    invite_code VARCHAR(8) UNIQUE NOT NULL,

    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (owner_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_families_owner_id ON families(owner_id);

CREATE TABLE family_members (
    user_id UUID NOT NULL,
    family_id UUID NOT NULL,

    role TEXT NOT NULL,
    color TEXT NOT NULL,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, family_id),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    FOREIGN KEY (family_id)
        REFERENCES families(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_family_members_family_id ON family_members(family_id);

CREATE TABLE lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID,
    family_id UUID,

    name TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_lists_owner
    CHECK (
        (user_id IS NOT NULL AND family_id IS NULL)
        OR
        (user_id IS NULL AND family_id IS NOT NULL)
    ),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    FOREIGN KEY (family_id)
        REFERENCES families(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_lists_user_id ON lists(user_id);

CREATE INDEX idx_lists_family_id ON lists(family_id);

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    list_id UUID NOT NULL,

    name TEXT NOT NULL,
    section TEXT,

    quantity DOUBLE PRECISION NOT NULL,
    unit TEXT,

    priority INTEGER NOT NULL DEFAULT 1,
    urgent BOOLEAN NOT NULL DEFAULT FALSE,

    checked BOOLEAN NOT NULL DEFAULT FALSE,
    remind_every_days INTEGER,

    added_by_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (list_id)
        REFERENCES lists(id)
        ON DELETE CASCADE,

    FOREIGN KEY (added_by_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_items_list_id ON items(list_id);

CREATE INDEX idx_items_added_by_id ON items(added_by_id);