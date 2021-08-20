CREATE TABLE "policies" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"app_id" varchar(255) NOT NULL,
"app_name" varchar(255),
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "policy_name_app_id" UNIQUE ("name", "app_id")
)
WITHOUT OIDS;
COMMENT ON CONSTRAINT "policy_name_app_id" ON "policies" IS 'policy app_id and name must be unique';
COMMENT ON COLUMN "policies"."app_id" IS 'Centrally assigned app_id (for app-accessKey-controls)';
COMMENT ON COLUMN "policies"."app_name" IS 'Centrally optional/assigned app_name (for app-accessKey-controls)';

CREATE TABLE "items" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"app_id" uuid NOT NULL,
"value_type" varchar(64) NOT NULL,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "item_name_app_id" UNIQUE ("name", "app_id")
)
WITHOUT OIDS;
COMMENT ON CONSTRAINT "item_name_app_id" ON "items" IS 'policy item name must be unique';

CREATE TABLE "groups" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"name" varchar(255) NOT NULL,
"app_id" uuid NOT NULL,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "group_name_app_id" UNIQUE ("name", "app_id")
)
WITHOUT OIDS;
CREATE TABLE "group_items" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"item_id" uuid NOT NULL,
"group_id" uuid NOT NULL,
"order" int8 NOT NULL,
"value" jsonb NOT NULL,
"operator" varchar(25) NOT NULL,
"relation" varchar(25) DEFAULT AND,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "order_itemId_groupId" UNIQUE ("order", "item_id", "group_id"),
CONSTRAINT "item_id_group_id" UNIQUE ("item_id", "group_id")
)
WITHOUT OIDS;
COMMENT ON CONSTRAINT "order_itemId_groupId" ON "group_items" IS 'order, itemId and groupId uniqueness';
COMMENT ON CONSTRAINT "item_id_group_id" ON "group_items" IS 'item_id and group_id uniqueness';
COMMENT ON COLUMN "group_items"."order" IS 'item order within the group';
COMMENT ON COLUMN "group_items"."value" IS 'item-json object: valueType and value key-value pair';
COMMENT ON COLUMN "group_items"."operator" IS 'the requested-item operation w.r.t to this itemValue, i.e requested-item-value > itemValue';
COMMENT ON COLUMN "group_items"."relation" IS 'item relation between this and next item in the order-specified';

CREATE TABLE "app_policies" (
"id" uuid NOT NULL DEFAULT gen_random_uuid(),
"app_id" uuid NOT NULL,
"policy_id" uuid NOT NULL,
"language" varchar(25) NOT NULL DEFAULT en-US,
"desc" text,
"is_active" bool NOT NULL DEFAULT true,
"created_by" uuid,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id") ,
CONSTRAINT "app_name" UNIQUE ("app_id"),
CONSTRAINT "app_name_access" UNIQUE ("app_id", "policy_id")
)
WITHOUT OIDS;
COMMENT ON CONSTRAINT "app_name" ON "app_policies" IS 'organization app or solution name must be unique';

CREATE TABLE "policy_groups" (
"policy_id" uuid NOT NULL,
"group_id" uuid NOT NULL,
"relation" varchar(25) DEFAULT group relation,
"order" int8 NOT NULL DEFAULT group order,
"created_by" uuid,
"created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
"updated_by" uuid,
"updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
"is_active" bool DEFAULT true
)
WITHOUT OIDS;

ALTER TABLE "group_items" ADD CONSTRAINT "group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "group_items" ADD CONSTRAINT "item_id" FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "policy_groups" ADD CONSTRAINT "policy_id" FOREIGN KEY ("policy_id") REFERENCES "policies" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "policy_groups" ADD CONSTRAINT "group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE RESTRICT ON UPDATE CASCADE;

