CREATE TABLE "groups" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"owner_id" uuid,
"desc" text NOT NULL,
"language" varchar(25) DEFAULT 'en-US',
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT "group_name" UNIQUE ("name")
)
WITHOUT OIDS;
COMMENT ON COLUMN "groups"."owner_id" IS 'default_value=>created_by';

CREATE TABLE "categories" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(64) NOT NULL,
"owner_id" uuid,
"path" varchar(255),
"priority" int4 NOT NULL DEFAULT 100,
"group_id" uuid NOT NULL,
"group_name" varchar(64),
"parent_id" uuid,
"icon_style" varchar(255),
"language" varchar(25) DEFAULT 'en-US',
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "category_name_owner_parent_id" UNIQUE ("name", "owner_id", "parent_id"),
CONSTRAINT "category_name_created_by_parent_id" UNIQUE ("name", "parent_id", "created_by")
)
WITHOUT OIDS;
COMMENT ON COLUMN "categories"."owner_id" IS 'default_value=>created_by';
COMMENT ON COLUMN "categories"."path" IS 'Category item path or url';
COMMENT ON COLUMN "categories"."icon_style" IS 'Icon style for item-logo indicator';
COMMENT ON COLUMN "categories"."language" IS 'Language, default => English US';

CREATE TABLE "audits" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"table_name" varchar(255) NOT NULL,
"log_records" jsonb NOT NULL,
"new_log_records" jsonb,
"log_type" varchar(255) NOT NULL,
"log_by" uuid,
"log_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
COMMENT ON TABLE "audits" IS 'Transactions audit log table to capture CRUD, access and systems related logs';


ALTER TABLE "categories" ADD CONSTRAINT "group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "categories" ADD CONSTRAINT "group_name" FOREIGN KEY ("group_name") REFERENCES "groups" ("name") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "categories" ADD CONSTRAINT "parent_id" FOREIGN KEY ("parent_id") REFERENCES "categories" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;

