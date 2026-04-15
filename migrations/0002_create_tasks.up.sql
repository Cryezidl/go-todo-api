CREATE TABLE tasks (
    id UUID PRIMARY KEY gen_random_uuid(),
    list_id UUID NOT NULL,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'new',
    priority INT CHECK (priority >= 0),
    deadline TIMESTAMPTZ,
    remind_at TIMESTAMPTZ,
    tags TEXT[],
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    
    CONSTRAINT fk_list FOREIGN KEY (list_id) REFERENCES task_lists(id) ON DELETE CASCADE,
    CONSTRAINT fk_creator FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);