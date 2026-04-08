CREATE TABLE task_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    is_private BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);



/*
TaskList :
	id (int) (обязательный уникальный)
	user_id (int) (обязательный)
	title string (обязательный)
	description (string) (необязательный)
	is_private (bool) (обязательный, база - false)
	created_at (DateTime) (обязательный, при создании)
	updated_at (DateTime) (необязательный)
*/