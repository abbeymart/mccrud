CREATE TABLE "users" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"username" varchar(255) NOT NULL,
"email" varchar(255) NOT NULL,
"password" varchar(1000) NOT NULL,
"api_token" varchar(1000),
"is_admin" bool NOT NULL DEFAULT false,
"accept_term" bool NOT NULL DEFAULT false,
"verified" bool NOT NULL DEFAULT false,
"groups" jsonb,
"can_upload" bool DEFAULT false,
"profile" jsonb NOT NULL,
"desc" text,
"is_active" bool NOT NULL DEFAULT false,
"created_by" uuid,
"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "username" UNIQUE ("username"),
CONSTRAINT "email" UNIQUE ("email"),
CONSTRAINT "user_api_token" UNIQUE ("username", "api_token"),
CONSTRAINT "email_api_token" UNIQUE ("email", "api_token"),
CONSTRAINT "api_token" UNIQUE ("api_token")
)
WITHOUT OIDS;
CREATE UNIQUE INDEX "user_name" ON "users" ("username" ASC NULLS LAST);
COMMENT ON TABLE "users" IS 'User records for access control';
COMMENT ON COLUMN "users"."api_token" IS 'Assigned token for API access';
COMMENT ON COLUMN "users"."groups" IS 'User groups';
COMMENT ON COLUMN "users"."can_upload" IS 'If user is allowed to upload documents / files';
COMMENT ON COLUMN "users"."profile" IS 'User profile, updatable by user';
COMMENT ON COLUMN "users"."is_active" IS 'false on registration, true after verification';

CREATE TABLE "access_keys" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"user_id" uuid NOT NULL,
"login_name" varchar(255) NOT NULL,
"token" varchar(255) NOT NULL,
"expire" int8 NOT NULL,
"created_by" uuid,
"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "user_token" UNIQUE ("user_id", "token")
)
WITHOUT OIDS;
COMMENT ON COLUMN "access_keys"."user_id" IS 'Login user user_id';
COMMENT ON COLUMN "access_keys"."login_name" IS 'Login name may be email or username';
COMMENT ON COLUMN "access_keys"."token" IS 'Login token for client-UI authorization';
COMMENT ON COLUMN "access_keys"."expire" IS 'Login expiration in seconds since 1970';

CREATE TABLE "services" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"category" varchar(64) NOT NULL,
"path" varchar(255),
"cost" decimal,
"priority" int8 DEFAULT 100,
"icon_style" varchar(2000),
"parent_id" uuid,
"language" varchar(12) NOT NULL DEFAULT en-US,
"desc" text NOT NULL,
"status" varchar(64),
"status_update" varchar(255),
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "name_category" UNIQUE ("name", "category"),
CONSTRAINT "name_parent" UNIQUE ("name", "parent_id"),
CONSTRAINT "name_path_language" UNIQUE ("name", "path", "language")
)
WITHOUT OIDS;
COMMENT ON COLUMN "services"."category" IS 'service item category';

CREATE TABLE "roles" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"role_id" uuid NOT NULL,
"service_category" varchar(64) NOT NULL,
"service_id" uuid NOT NULL,
"can_create" bool NOT NULL DEFAULT false,
"can_update" bool NOT NULL DEFAULT false,
"can_delete" bool NOT NULL DEFAULT false,
"can_read" bool NOT NULL DEFAULT false,
"can_crud" bool NOT NULL DEFAULT false,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid NOT NULL,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
COMMENT ON COLUMN "roles"."role_id" IS 'user group id';
COMMENT ON COLUMN "roles"."service_category" IS 'service item category';
COMMENT ON COLUMN "roles"."service_id" IS 'service item id';

CREATE TABLE "groups" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"owner_id" uuid,
"desc" text NOT NULL,
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
"language" varchar(25) DEFAULT en-US,
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

CREATE TABLE "settings" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"category" varchar(255) NOT NULL,
"name" varchar(255) NOT NULL,
"value_type" varchar(25) NOT NULL,
"relation" varchar(64) NOT NULL,
"value" jsonb NOT NULL,
"unit" varchar(25),
"parentId" uuid,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool DEFAULT false,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
COMMENT ON COLUMN "settings"."value_type" IS 'Type of setting value: string, number, bool, date etc.';
COMMENT ON COLUMN "settings"."relation" IS 'Setting relation';
COMMENT ON COLUMN "settings"."unit" IS 'Unit of measure: unit, km, seconds etc.';

CREATE TABLE "locations" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"code" varchar(25) NOT NULL,
"name" varchar(255) NOT NULL,
"category" varchar(255),
"phone_code" varchar(255),
"currency" varchar(12) NOT NULL DEFAULT USD,
"latitude" float8,
"longitude" float8,
"time_zone" date,
"language" varchar(255) NOT NULL DEFAULT en-US,
"languages" jsonb NOT NULL,
"parent_id" int4,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
UNIQUE ("code", "category", "parent_id"),
UNIQUE ("code", "name", "parent_id"),
UNIQUE ("code", "category", "phone_code", "parent_id")
)
WITHOUT OIDS;
CREATE TABLE "contacts" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"first_name" varchar(255) NOT NULL,
"last_name" varchar(255) NOT NULL,
"middle_name" varchar(255),
"email" varchar(255),
"phone" varchar(30),
"address" varchar(255),
"ownerId" uuid NOT NULL,
"category" varchar(100),
"job" varchar(100) DEFAULT Job Title,
"org" varchar(255),
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"create_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "owner_email_category" UNIQUE ("email", "ownerId", "category")
)
WITHOUT OIDS;
COMMENT ON COLUMN "contacts"."org" IS 'Organization';

CREATE TABLE "contents" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"title" varchar(255) NOT NULL,
"sub_title" varchar(255),
"access" varchar(25),
"body" text NOT NULL,
"type" varchar(64) NOT NULL,
"category" varchar(255) NOT NULL,
"tag" varchar(25),
"priority" int4,
"parent_id" uuid,
"language" varchar(25) NOT NULL DEFAULT en-US,
"owner_id" uuid,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "addresses" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"category" varchar(64) NOT NULL,
"contact_id" uuid NOT NULL,
"owner_id" uuid NOT NULL,
"street_number" int4 NOT NULL,
"street_name" varchar(512) NOT NULL,
"street_type" varchar(25),
"street_direction" varchar(25),
"city" varchar(255) NOT NULL,
"state" varchar(255) NOT NULL,
"country" varchar(255) NOT NULL,
"phone" varchar(30),
"postal_code" varchar(25),
"latitude" float8,
"longitude" float8,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
UNIQUE ("category", "owner_id", "phone"),
UNIQUE ("category", "street_number", "city", "state", "country")
)
WITHOUT OIDS;
CREATE TABLE "phones" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"category" varchar(100) NOT NULL DEFAULT Business,
"phone" varchar(30) NOT NULL,
"contact_id" uuid NOT NULL,
"owner_id" uuid NOT NULL,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestampz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
UNIQUE ("phone", "category", "contact_id")
)
WITHOUT OIDS;
CREATE TABLE "taxes" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"category" varchar(255) NOT NULL,
"juri" varchar(255) NOT NULL,
"city" varchar(255),
"state" varchar(255),
"country" varchar(255),
"region" varchar(255),
"location" varchar(255) NOT NULL,
"amount" decimal,
"unit" varchar(25),
"range_from" decimal,
"range_to" decimal DEFAULT en-US,
"language" varchar(25) NOT NULL,
"desc" text,
"is_active" bool DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "payments" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "organizations" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"number" int4,
"category" varchar(100),
"industry" varchar(100),
"phone" varchar(30),
"address" varchar(255),
"city" varchar(100),
"state" varchar(100),
"country" varchar(100),
"owner_id" uuid NOT NULL,
"email" varchar(255),
"website" varchar(255),
"reg_date" date,
"latitude" float8,
"longitude" float8,
"language" varchar(25) NOT NULL DEFAULT en-US,
"parent_id" uuid,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "locales" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"code" varchar(25),
"name" varchar(255),
"value" jsonb,
"category" varchar(100),
"language" varchar(25) DEFAULT en-US,
"desc" text,
"is_active`" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
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
COMMENT ON TABLE "audits" IS 'Transactions audit log table to capature CRUD, access and systems related logs';

CREATE TABLE "verify_users" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"email" varchar(255) NOT NULL,
"username" varchar(255) NOT NULL,
"first_name" varchar(255) NOT NULL,
"last_name" varchar(255) NOT NULL,
"full_name" varchar(1000),
"msg_from" varchar(500) NOT NULL,
"expire" int8 NOT NULL,
"req_url" varchar(1000) NOT NULL,
"token" varchar(1000) NOT NULL,
"task_type" varchar(64) NOT NULL,
"login_name" varchar(255),
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "verify_email" UNIQUE ("email"),
CONSTRAINT "verify_username" UNIQUE ("username")
)
WITHOUT OIDS;
COMMENT ON COLUMN "verify_users"."msg_from" IS 'Message from email address';

CREATE TABLE "documents" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"title" varchar(255) NOT NULL,
"owner_id" uuid NOT NULL,
"activity_id" uuid,
"file_name" varchar(255) NOT NULL,
"file_type" varchar(255) NOT NULL,
"file_id" varchar(255),
"file_url" varchar(3000),
"file_size" float8,
"file_size_unit" varchar(25),
"file_content" bytea,
"new_file_name" varchar(255),
"ref" varchar(255),
"folder_id" uuid,
"category" varchar(255),
"parent_id" uuid,
"language" varchar(25) NOT NULL DEFAULT en-US,
"file_response" jsonb,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"create_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "files" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"file_type" varchar(255) NOT NULL,
"file_name" varchar(255) NOT NULL,
"new_file_name" varchar(255),
"file_content" bytea NOT NULL,
"folder" varchar(255),
"file_id" varchar(512) NOT NULL,
"language" varchar(25) NOT NULL DEFAULT en-US,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "file_id" UNIQUE ("file_id"),
UNIQUE ("file_type", "file_name", "created_by")
)
WITHOUT OIDS;
CREATE TABLE "folders" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255),
"owner_id" uuid,
"group" varchar(255),
"category" varchar(255),
"parent_id" uuid,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
COMMENT ON COLUMN "folders"."group" IS 'Bookmark, Cart, Order etc.';

CREATE TABLE "messages" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"subject" varchar(255) NOT NULL,
"owner_id" uuid NOT NULL,
"from" varchar(255) NOT NULL,
"to" jsonb NOT NULL,
"cc" jsonb,
"bcc" jsonb,
"delegate" varchar(255),
"content" bytea NOT NULL,
"attachment" jsonb,
"category" varchar(255),
"parent_id" uuid,
"language" varchar(255),
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz(255) NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
COMMENT ON COLUMN "messages"."cc" IS 'copied email addresses';
COMMENT ON COLUMN "messages"."bcc" IS 'blind-copied email addresses';
COMMENT ON COLUMN "messages"."delegate" IS 'sender email on behalf of from email address';

CREATE TABLE "user_profiles" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"user_id" uuid NOT NULL,
"first_name" varchar(120) NOT NULL,
"last_name" varchar(120) NOT NULL,
"middle_name" varchar(120),
"language" varchar(12) NOT NULL DEFAULT en-US,
"phone" varchar(30),
"rec_email" varchar(500),
"group" uuid,
"date_of_birth" date,
"two_factor_auth" bool,
"postal_code" varchar(12),
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "recovery_email" UNIQUE ("user_id", "rec_email"),
CONSTRAINT "phone_number" UNIQUE ("phone")
)
WITHOUT OIDS;
CREATE UNIQUE INDEX ON "user_profiles" ("user_id" ASC NULLS LAST);
COMMENT ON TABLE "user_profiles" IS 'User records for access control';
COMMENT ON COLUMN "user_profiles"."language" IS 'English US';

CREATE TABLE "support" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"subject" varchar(255) NOT NULL,
"owner_id" uuid NOT NULL,
"sender" varchar(255) NOT NULL,
"responder" varchar(255),
"category" varchar(255),
"status" varchar(255) NOT NULL,
"message" text NOT NULL,
"prior_message" text,
"language" varchar(25) NOT NULL DEFAULT en-US,
"parent_id" uuid,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "bookmarks" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"folder" varchar(255),
"category" varchar(255),
"book_url" varchar(3000),
"book_id" uuid,
"owner_id" uuid,
"language" varchar(25) DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "todos" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"owner_id" uuid NOT NULL,
"category" varchar(255) NOT NULL,
"start_date" date,
"end_date" date,
"effort" float8,
"effort_unit" varchar(25),
"duration" float4,
"duration_unit" varchar(25),
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"parent_id" uuid,
"lead" varchar(255),
"resources" jsonb,
"is_active" bool NOT NULL DEFAULT false,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
UNIQUE ("name", "category", "owner_id")
)
WITHOUT OIDS;
CREATE TABLE "resources" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"first_name" varchar(255) NOT NULL,
"last_name" varchar(255) NOT NULL,
"job_title" varchar(255) NOT NULL,
"user_id" uuid,
"org" varchar(255),
"email" varchar(255),
"phone" varchar(255),
"address" varchar(255),
"rate" decimal,
"rate_unit" varchar(25),
"rate_currency" varchar(12) NOT NULL DEFAULT USD,
"language" varchar(25) NOT NULL DEFAULT en-US,
"parent_id" uuid,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
)
WITHOUT OIDS;
CREATE TABLE "todos_resources" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"todo_id" uuid NOT NULL,
"resource_id" uuid NOT NULL,
"role" varchar(255),
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") 
)
WITHOUT OIDS;
CREATE TABLE "issues" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"category" varchar(64) NOT NULL,
"start_date" date,
"end_date" date,
"likelihood" varchar(12),
"impact" varchar(12),
"severity" varchar(12),
"owner_id" uuid,
"lead" uuid,
"status" varchar(25) NOT NULL,
"language" varchar(25) DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
)
WITHOUT OIDS;
CREATE TABLE "apps" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"app_name" varchar(100) NOT NULL,
"access_key" uuid NOT NULL DEFAULT gen_random_uuid(),
"desc" text NOT NULL,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT "access_key" UNIQUE ("access_key"),
CONSTRAINT "app_name" UNIQUE ("app_name")
)
WITHOUT OIDS;
COMMENT ON TABLE "apps" IS 'Application table for mConnect solutions access control management';
COMMENT ON COLUMN "apps"."id" IS 'app_id';
COMMENT ON COLUMN "apps"."access_key" IS 'access_key';


ALTER TABLE "services" ADD CONSTRAINT "parent_id" FOREIGN KEY ("parent_id") REFERENCES "services" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "user_profiles" ADD CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "access_keys" ADD CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "verify_users" ADD CONSTRAINT "user_email" FOREIGN KEY ("email") REFERENCES "users" ("email") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "roles" ADD CONSTRAINT "service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "contacts" ADD CONSTRAINT "owner_id" FOREIGN KEY ("ownerId") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "addresses" ADD CONSTRAINT "contact_id" FOREIGN KEY ("contact_id") REFERENCES "contacts" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "phones" ADD FOREIGN KEY ("contact_id") REFERENCES "contacts" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "documents" ADD FOREIGN KEY ("parent_id") REFERENCES "documents" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "documents" ADD FOREIGN KEY ("folder_id") REFERENCES "folders" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "locations" ADD FOREIGN KEY ("parent_id") REFERENCES "locations" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "todos_resources" ADD CONSTRAINT "todo_id" FOREIGN KEY ("todo_id") REFERENCES "todos" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "todos_resources" ADD FOREIGN KEY ("resource_id") REFERENCES "resources" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "categories" ADD CONSTRAINT "group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "categories" ADD CONSTRAINT "group_name" FOREIGN KEY ("group_name") REFERENCES "groups" ("name") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "categories" ADD CONSTRAINT "parent_id" FOREIGN KEY ("parent_id") REFERENCES "categories" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "roles" ADD CONSTRAINT "role_id" FOREIGN KEY ("role_id") REFERENCES "categories" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;

