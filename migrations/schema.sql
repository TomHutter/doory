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
, "is_active" bool NOT NULL DEFAULT 'true', "alarm" bool NOT NULL DEFAULT 'false');
CREATE TABLE IF NOT EXISTS "tokens" (
"id" TEXT PRIMARY KEY,
"token_id" TEXT NOT NULL,
"person_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "access_groups" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL,
"description" TEXT,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "tokens_access_groups" (
"id" TEXT PRIMARY KEY,
"token_id" char(36) NOT NULL,
"access_group_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "access_groups_doors" (
"id" TEXT PRIMARY KEY,
"access_group_id" char(36) NOT NULL,
"door_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "access_group_doors" (
"id" TEXT PRIMARY KEY,
"access_group_id" char(36) NOT NULL,
"door_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "token_access_groups" (
"id" TEXT PRIMARY KEY,
"token_id" char(36) NOT NULL,
"access_group_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "users" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL,
"email" TEXT,
"provider" TEXT NOT NULL,
"provider_id" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
