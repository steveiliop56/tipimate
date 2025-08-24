CREATE TABLE IF NOT EXISTS "schemas" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "created_at" TEXT DEFAULT (DATETIME('now')) NOT NULL,
    "updated_at" TEXT DEFAULT (DATETIME('now')) NOT NULL,
    "deleted_at" TEXT NULL,
    "version" INTEGER DEFAULT 0 NOT NULL,
    "latest_version" INTEGER DEFAULT 0 NOT NULL
);
