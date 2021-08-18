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
"company_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "companies" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL,
"description" TEXT,
"contact_person_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "people" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL,
"surname" TEXT NOT NULL,
"company_id" char(36) NOT NULL,
"email" TEXT NOT NULL,
"phone" TEXT NOT NULL,
"id_number" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
