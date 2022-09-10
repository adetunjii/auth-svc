#!/bin/sh

# exit immediately if there is an error
set -e
echo "running db migration...."
/app/migrate -path /app/migration -database "postgresql://sonar_user_checklos:S0N4RQUB3!@#\$@154.12.237.18:31373/dh-user?sslmode=enable" -verbose up

echo "migration ran successfully"
