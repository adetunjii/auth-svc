package db

import (
	"fmt"

	"github.com/adetunjii/auth-svc/internal/port"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	DatabaseUrl string `json:"url"`
}

type PostgresDB struct {
	instance *gorm.DB
	logger   port.AppLogger
}

var _ port.DB = (*PostgresDB)(nil)

type PreloadOption func(*gorm.DB)

func New(dbConfig DBConfig, logger port.AppLogger) *PostgresDB {
	db := &PostgresDB{
		instance: nil,
		logger:   logger,
	}

	if err := db.Connect(dbConfig); err != nil {
		logger.Fatal("connection to db failed: %v", err)
	}
	return db
}

func (p *PostgresDB) Connect(config DBConfig) error {

	var dsn string
	databaseUrl := config.DatabaseUrl
	fmt.Println(databaseUrl)
	if databaseUrl == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.Name)
	} else {
		dsn = databaseUrl
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	p.instance = db

	p.logger.Info(fmt.Sprintf("Database Connected Successfully %v...", dsn))
	return nil
}

func (p *PostgresDB) CloseConnection() error {
	if p.instance != nil {
		connection, err := p.instance.DB()
		if err != nil {
			return err
		}

		connection.Close()
	}
	return nil
}

func (p *PostgresDB) RestartConnection(config DBConfig) error {
	if p.instance != nil {
		p.CloseConnection()
	}

	return p.Connect(config)
}

func (p *PostgresDB) Save(arg interface{}) error {
	return p.instance.Create(arg).Error
}

func (p *PostgresDB) FindAll(dest interface{}, conditions map[string]interface{}) error {
	err := p.instance.Where(conditions).Find(dest).Error
	return err
}

func (p *PostgresDB) List(dest interface{}, conditions map[string]interface{}, limit int, offset int) error {
	err := p.instance.Limit(limit).Offset(offset).Where(conditions).Find(dest).Error
	return err
}

func (p *PostgresDB) FindWithPreload(dest interface{}, conditions map[string]interface{}, options ...func(*gorm.DB)) error {

	tx := p.instance.Where(conditions).Find(dest)

	for _, option := range options {
		option(tx)
	}

	return tx.Error
}

func (p *PostgresDB) Raw(dest interface{}, query string, values ...interface{}) error {
	return p.instance.Raw(query, values).Scan(dest).Error
}

func (p *PostgresDB) FindById(dest interface{}, id string) error {

	err := p.instance.Where("id = ?", id).First(dest).Error
	return err
}

func (p *PostgresDB) FindOne(dest interface{}, conditions map[string]interface{}) error {

	err := p.instance.Where(conditions).First(dest).Error
	return err
}

func (p *PostgresDB) Delete(model interface{}, id string) error {
	return p.instance.Where("id = ?", id).Delete(model).Error
}

func (p *PostgresDB) DeleteOne(model interface{}, condition map[string]interface{}) error {
	return p.instance.Where(condition).Delete(model).Error
}

func (p *PostgresDB) Update(model interface{}, condition map[string]interface{}, updates interface{}) error {
	return p.instance.Model(model).Where(condition).Updates(updates).Error
}

// a wrapper to add preload options to a query
func WithPreload(option string) PreloadOption {
	return func(tx *gorm.DB) {
		tx.Preload(option)
	}
}
