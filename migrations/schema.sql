CREATE TABLE IF NOT EXISTS "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "schema_migration_version_idx" ON "schema_migration" (version);
CREATE TABLE IF NOT EXISTS "doors" (
"id" TEXT PRIMARY KEY,
"room" TEXT NOT NULL,
"floor" TEXT NOT NULL,
"building" TEXT NOT NULL,
"description" TEXT,
"company" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
