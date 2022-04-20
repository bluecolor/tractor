package repo

import (
	"fmt"
	"os"
	"time"

	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	*gorm.DB
}

func New(config conf.DB) (*Repository, error) {
	zlog := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}).With().Caller().Logger()
	newLogger := logger.New(
		&zlog, // IO.writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	gormConfig := &gorm.Config{Logger: newLogger}
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}
func (r *Repository) Drop() (err error) {
	migrator := r.Migrator()
	for _, model := range models.Models {
		log.Info().Msgf("dropping %T", model)
		if migrator.HasTable(model) {
			err = migrator.DropTable(model)
			if err != nil {
				log.Error().Err(err).Msg("failed to drop model: " + fmt.Sprintf("%T", model))
				return
			}
		}
	}
	return
}
func (r *Repository) Migrate() (err error) {
	for _, model := range models.Models {
		log.Info().Msgf("migrating %T", model)
		err = r.DB.AutoMigrate(model)
		if err != nil {
			log.Error().Err(err).Msg("failed to migrate model: " + fmt.Sprintf("%T", model))
			return
		}
	}
	return
}
func (r *Repository) Seed(basePath string, reset bool) (err error) {
	param := models.Param{}
	err = r.Where("name = ?", "db.seed.state").First(&param).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			param.Value = "1"
			param.Name = "db.seed.state"
			r.Create(&param)
		} else {
			log.Error().Err(err).Msg("failed to get db.seed.state")
			return
		}
	} else if param.Value == "1" && !reset {
		log.Info().Msg("db.seed.state is already set to 1")
		return
	}

	if reset {
		if err = r.Drop(); err != nil {
			return
		}
		if err = r.Migrate(); err != nil {
			return
		}
	}
	functions := []func(string) error{
		r.SeedFileTypes,
		r.SeedProviders,
		r.SeedConnectionTypes,
	}
	for _, f := range functions {
		err = f(basePath)
		if err != nil {
			return
		}
	}

	return
}
