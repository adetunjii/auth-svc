#!/bin/sh

# exit immediately if there is an error
set -e
echo "running db migration...."
/app/migrate -path /app/migration -database 'postgresql://sonar_user_checklos:8iu7*IU&@154.12.237.18:31373/dh_users?sslmode=disable' -verbose up

echo "migration ran successfully"
rm -rf /app/migrate
rm -rf /app/migration
