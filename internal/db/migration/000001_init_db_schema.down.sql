-- drop indexes
DROP INDEX IF EXISTS idx_role_title;
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone_number;
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_user_interests_user_id;
DROP INDEX IF EXISTS idx_login_activities_user_id;
DROP INDEX IF EXISTS idx_countries_name;
DROP INDEX IF EXISTS idx_countries_phone_code;

-- drop foreign keys
ALTER TABLE role_permissions DROP CONSTRAINT role_permissions_role_id_fkey;
ALTER TABLE role_permissions DROP CONSTRAINT role_permissions_permission_id_fkey;
ALTER TABLE users DROP CONSTRAINT users_role_id_fkey;
ALTER TABLE user_interests DROP CONSTRAINT user_interests_user_id_fkey;
ALTER TABLE user_interests DROP CONSTRAINT user_interests_interest_id_fkey;
ALTER TABLE login_activities DROP CONSTRAINT login_activities_user_id_fkey;

-- drop tables
DROP TABLE IF EXISTS "countries";
DROP TABLE IF EXISTS "role_permissions";
DROP TABLE IF EXISTS "user_interests";
DROP TABLE IF EXISTS "interests";
DROP TABLE IF EXISTS "login_activities";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "roles";