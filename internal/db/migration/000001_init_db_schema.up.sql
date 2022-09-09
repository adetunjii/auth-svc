CREATE TABLE "roles" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "title" varchar NOT NULL UNIQUE,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "permissions" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar NOT NULL,
  "description" varchar,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "role_permissions" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  "created_at" timestamptz DEFAULT(now()),
  "updated_at" timestamptz
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "username" varchar, 
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "phone_code" varchar NOT NULL, 
  "role_id" uuid NOT NULL,
  "password" varchar NOT NULL,
  "last_password_update" timestamptz,
  "is_email_verified" bool DEFAULT false,
  "is_phone_verified" bool DEFAULT false,
  "is_active" bool DEFAULT false,
  "address" varchar,
  "state" varchar,
  "country" varchar NOT NULL,
  "timezone" varchar, 
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "interests" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "title" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz 
);

CREATE TABLE "user_interests" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "interest_id" uuid NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "login_activities" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "token" varchar NOT NULL,
  "device" varchar NOT NULL,
  "time" timestamptz DEFAULT (now()),
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "countries" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "iso" varchar NOT NULL,
  "name" varchar NOT NULL UNIQUE,
  "phone_code" varchar NOT NULL UNIQUE,
  "nice_name" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "numcode" varchar NOT NULL,
  "imageUrl" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE INDEX idx_role_title ON "roles" ("title");

CREATE INDEX idx_role_permissions_role_id ON "role_permissions" ("role_id");

CREATE INDEX idx_role_permissions_permission_id ON "role_permissions" ("permission_id");

CREATE INDEX idx_users_email  ON "users" ("email");

CREATE INDEX idx_users_phone_number on "users" ("phone_number");

CREATE INDEX idx_users_role_id ON "users" ("role_id");

CREATE INDEX idx_user_interests_user_id ON "user_interests" ("user_id");

CREATE INDEX idx_login_activities_user_id ON "login_activities" ("user_id");

CREATE INDEX idx_countries_name ON "countries" ("name");

CREATE INDEX idx_countries_phone_code ON "countries" ("phone_code");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "user_interests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_interests" ADD FOREIGN KEY ("interest_id") REFERENCES "interests" ("id");

ALTER TABLE "login_activities" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
