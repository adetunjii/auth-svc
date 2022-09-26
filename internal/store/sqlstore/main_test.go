package sqlstore

import (
	"os"
	"testing"

	"github.com/adetunjii/auth-svc/internal/db"
	"github.com/adetunjii/auth-svc/internal/port"
	"github.com/adetunjii/auth-svc/pkg/logging"
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
