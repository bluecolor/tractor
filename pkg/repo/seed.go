package repo

import (
	"encoding/csv"
	"os"
	"path"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/rs/zerolog/log"
)

func (r *Repository) SeedExtractionModes(basePath string) (err error) {
	filename := path.Join(basePath, "extraction_modes.csv")
	log.Info().Msg("Seeding extraction modes from " + filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to open " + filename)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to read " + filename)
	}

	records := []models.ExtractionMode{}
	for i, row := range data {
		if i == 0 {
			continue
		}
		record := models.ExtractionMode{
			Name: row[0],
		}
		records = append(records, record)
	}
	err = r.Save(&records).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to seed extraction modes")
		return
	}
	return
}
func (r *Repository) SeedFileTypes(basePath string) (err error) {
	filename := path.Join(basePath, "file_types.csv")
	log.Info().Msg("seeding file types from " + filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to open " + filename)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to read " + filename)
	}

	records := []models.FileType{}
	for i, row := range data {
		if i == 0 {
			continue
		}
		record := models.FileType{
			Name: row[0],
		}
		records = append(records, record)
	}
	err = r.Save(&records).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to seed file types")
		return
	}
	return
}
func (r *Repository) SeedProviderTypes(basePath string) (err error) {
	filename := path.Join(basePath, "provider_types.csv")
	log.Info().Msg("seeding provider types from " + filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to open " + filename)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to read " + filename)
	}

	records := []models.ProviderType{}
	for i, row := range data {
		if i == 0 {
			continue
		}
		record := models.ProviderType{
			Name: row[0],
		}
		records = append(records, record)
	}
	err = r.Save(&records).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to seed provider types")
		return
	}
	return
}
func (r *Repository) SeedConnectionTypes(basePath string) (err error) {
	filename := path.Join(basePath, "connection_types.csv")
	log.Info().Msg("seeding connection types from " + filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to open " + filename)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to read " + filename)
	}

	records := []models.ConnectionType{}
	for i, row := range data {
		if i == 0 {
			continue
		}
		record := models.ConnectionType{
			Name: row[0],
		}
		records = append(records, record)
	}
	err = r.Save(&records).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to seed connection types")
		return
	}
	return
}
