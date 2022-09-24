package sqlstore

import (
	"os"
	"testing"

	"gitlab.com/dh-backend/auth-service/internal/db"
	"gitlab.com/dh-backend/auth-service/internal/port"
	"gitlab.com/dh-backend/auth-service/pkg/logging"
)

var testDB *db.PostgresDB
var sqlStore *SqlStore
var logger port.AppLogger

const databaseUrl = "postgresql://teej4y:password@localhost:5432/auth-service?sslmode=disable"

func TestMain(m *testing.M) {

	dbConfig := db.DBConfig{
		DatabaseUrl: databaseUrl,
	}

	sugarLogger := logging.NewZapSugarLogger()
	logger = logging.NewLogger(sugarLogger)

	testDB = db.New(dbConfig, logger)
	sqlStore = NewSqlStore(testDB, logger)

	os.Exit(m.Run())

}
