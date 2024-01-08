package database

import (
	"context"
	"fmt"

	"os"

	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Needed for GCP Cloud SQL postgres proxy connection
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/jo-fr/activityhub/backend/pkg/config"
	"github.com/jo-fr/activityhub/backend/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Module = fx.Options(
	fx.Provide(ProvideDatabase),
)

const (
	migrateFilesDir     = "pkg/database/migrations"
	schemaName          = "activityhub"
	migrationsTableName = "schema_migrations"
)

// Database holds all informations for database conn
type Database struct {
	*gorm.DB
}

func ProvideDatabase(lc fx.Lifecycle, config config.Config, logger *log.Logger) (*Database, error) {

	var _db *gorm.DB
	if config.Environment.IsLocal() {
		logger.Info("connecting to local database")
		dsn := getConnectionDSN(config.Database.Username, config.Database.Password, config.Database.Database, config.Database.Host, config.Database.Port)
		db, err := connectLocalDB(dsn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to local database")
		}
		_db = db

	} else {

		logger.Info("connecting to cloud sql database")

		uri := fmt.Sprintf("%s:%s:%s", config.GCP.ProjectID, config.GCP.Region, config.Database.Host)
		dsn := getConnectionDSN(config.Database.Username, config.Database.Password, config.Database.Database, uri, config.Database.Port)
		db, err := connectCloudSQL(dsn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to cloud sql database")
		}
		_db = db
	}

	if err := runMigrations(_db, config.Database.Database); err != nil {
		return nil, errors.Wrap(err, "failed to run migrations")
	}
	logger.Info("migrations ran successfully")

	db := &Database{_db}

	registerHooks(lc, db, logger)

	return db, nil

}

// registerHooks for uber fx
func registerHooks(lc fx.Lifecycle, db *Database, logger *log.Logger) {
	lc.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				logger.Info("closing database connection")

				dbConn, err := db.DB.DB()
				if err != nil {
					return err
				}

				return dbConn.Close()
			},
		},
	)
}

func runMigrations(db *gorm.DB, dbName string) error {

	// schema needs to be created before migrations can be run
	if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schemaName)).Error; err != nil {
		return errors.Wrap(err, "failed to create schema")
	}

	if err := db.Exec(fmt.Sprintf("set schema '%s'", schemaName)).Error; err != nil {
		return errors.Wrap(err, "failed to set schema")
	}

	if err := db.Exec(fmt.Sprintf("set search_path = %s, public", schemaName)).Error; err != nil {
		return errors.Wrap(err, "failed to  set search path")
	}

	dbConn, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get db connection")
	}

	driver, err := postgresMigrate.WithInstance(dbConn, &postgresMigrate.Config{
		DatabaseName:    dbName,
		SchemaName:      schemaName,
		MigrationsTable: migrationsTableName,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create db migration instance")
	}

	dir, err := getMigrationsFilePath()
	if err != nil {
		return errors.Wrap(err, "failed to get migrations file path")
	}

	m, err := migrate.NewWithDatabaseInstance(dir, dbName, driver)
	if err != nil {
		return errors.Wrap(err, "failed to create db migration")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to run db migration")
	}

	return nil
}

func connectCloudSQL(dsn string) (*gorm.DB, error) {

	postgresConfig := postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        dsn,
	})

	db, err := gorm.Open(postgresConfig, &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return db, nil
}

func connectLocalDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   schemaName,
			SingularTable: false,
		}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return db, nil
}

// getMigrationsFilePath returns the path to the migrations files
func getMigrationsFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get working directory")
	}

	return fmt.Sprintf("file://%s/%s", wd, migrateFilesDir), nil

}

// getConnectionDSN returns the connection dsn for postgres
func getConnectionDSN(username string, password string, database string, host string, port string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		username, password, database, host, port)
}
