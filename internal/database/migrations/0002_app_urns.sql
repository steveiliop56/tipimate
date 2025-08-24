ALTER TABLE "schemas" RENAME TO "apps_old";

CREATE TABLE IF NOT EXISTS "apps" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "created_at" TEXT DEFAULT (DATETIME('now')) NOT NULL,
    "updated_at" TEXT DEFAULT (DATETIME('now')) NOT NULL,
    "deleted_at" TEXT NULL,
    "urn" TEXT UNIQUE NOT NULL,
    "version" INTEGER DEFAULT 0 NOT NULL,
    "latest_version" INTEGER DEFAULT 0 NOT NULL
);

INSERT INTO "apps" ("created_at", "updated_at", "deleted_at", "urn", "version", "latest_version")
SELECT "created_at", "updated_at", "deleted_at", "id" || ":migrated", "version", "latest_version"
FROM "apps_old";

DROP TABLE "apps_old";