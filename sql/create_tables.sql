CREATE TABLE if not exists task_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE if not exists task_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE if not exists task_template_task_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_tag_id INTEGER NOT NULL,
    task_template_id INTEGER NOT NULL,

    FOREIGN KEY (task_tag_id)
    REFERENCES task_tags(id)
    ON DELETE CASCADE,

    FOREIGN KEY (task_template_id)
    REFERENCES task_templates(id)
    ON DELETE CASCADE,

    UNIQUE (task_tag_id, task_template_id)
);

CREATE TABLE if not exists list_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE if not exists list_template_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    list_template_id INTEGER NOT NULL,
    task_template_id INTEGER NOT NULL,
    position INTEGER NOT NULL,

    FOREIGN KEY (list_template_id)
    REFERENCES list_templates(id)
    ON DELETE CASCADE,

    FOREIGN KEY (task_template_id)
    REFERENCES task_templates(id)
    ON DELETE CASCADE,

    UNIQUE (list_template_id, task_template_id)
);

CREATE TABLE if not exists active_lists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    list_template_id INTEGER NULL,
    name TEXT NOT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL,

    FOREIGN KEY (list_template_id)
    REFERENCES list_templates(id)
    ON DELETE CASCADE
);

CREATE TABLE if not exists active_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    active_list_id INTEGER NOT NULL,
    task_template_id INTEGER NULL,

    name TEXT NOT NULL,
    description TEXT,

    completed INTEGER NOT NULL DEFAULT 0,
    position INTEGER NOT NULL,

    created_at TEXT DEFAULT current_timestamp NOT NULL,
    completed_at TEXT DEFAULT current_timestamp,
    deleted_at TEXT DEFAULT NULL,

    FOREIGN KEY (active_list_id)
    REFERENCES active_lists(id)
    ON DELETE CASCADE,

    FOREIGN KEY (task_template_id)
    REFERENCES task_templates(id)
    ON DELETE SET NULL
);

