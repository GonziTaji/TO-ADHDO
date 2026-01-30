CREATE TABLE if not exists articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    attributes JSONB DEFAULT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE if not exists articles_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    alt TEXT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE if not exists articles_prices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    description TEXT DEFAULT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL
);

CREATE TABLE if not exists tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE if not exists articles_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    created_at TEXT DEFAULT current_timestamp NOT NULL,
    updated_at TEXT DEFAULT current_timestamp NOT NULL,
    deleted_at TEXT DEFAULT NULL,

    FOREIGN KEY (article_id)
    REFERENCES articles(id)
    ON DELETE CASCADE,

    FOREIGN KEY (tag_id)
    REFERENCES tags(id)
    ON DELETE CASCADE,

    UNIQUE (tag_id, article_id, deleted_at)
);

