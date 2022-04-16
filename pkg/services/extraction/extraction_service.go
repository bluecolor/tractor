package extraction

import (
	"encoding/json"
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/tasks"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/hibiken/asynq"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type Service struct {
	repo   *repo.Repository
	client *tasks.Client
}

func NewService(repo *repo.Repository, client *tasks.Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}

func (s *Service) OneExtraction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ext := models.Extraction{}
	if err := s.repo.First(&ext, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	utils.RespondwithJSON(w, http.StatusOK, ext)
}
func (s *Service) FindExtractions(w http.ResponseWriter, r *http.Request) {

	model := s.repo.
		Joins("SourceDataset").
		Joins("TargetDataset").
		Preload("SourceDataset.Connection").
		Preload("TargetDataset.Connection").
		Preload("Sessions", func(db *gorm.DB) *gorm.DB {
			// todo load only last sessions
			return db.Order("created_at desc")
		})

	if r.URL.Query().Get("sc") != "" {
		model = model.Where("SourceDataset.connection_id = ?", r.URL.Query().Get("sc"))
	}
	if r.URL.Query().Get("tc") != "" {
		model = model.Where("TargetDataset.connection_id = ?", r.URL.Query().Get("tc"))
	}
	if r.URL.Query().Get("q") != "" {
		q := r.URL.Query().Get("q")
		model = model.Where("extractions.name like ?", "%"+q+"%")
	}
	model = model.Find(&[]models.Extraction{})

	if model.Error != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, model.Error)
		return
	}

	result := paginate.New().Response(model, r, &[]models.Extraction{})
	utils.RespondwithJSON(w, http.StatusOK, result)
}
func (s *Service) DeleteExtraction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ext := models.Extraction{}
	if err := s.repo.First(&ext, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.repo.Delete(&ext).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, ext)
}
func (s *Service) CreateExtraction(w http.ResponseWriter, r *http.Request) {
	extraction := models.Extraction{}
	if err := json.NewDecoder(r.Body).Decode(&extraction); err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := s.repo.Create(&extraction).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, extraction)
}
func (s *Service) RunExtraction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	extraction := models.Extraction{}
	if err := s.repo.First(&extraction, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	session := extraction.NewSession()
	if err := s.repo.Create(session).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	if err := s.repo.
		Preload("Extraction").
		Preload("Extraction.SourceDataset").
		Preload("Extraction.SourceDataset.Connection").
		Preload("Extraction.SourceDataset.Connection.ConnectionType").
		Preload("Extraction.SourceDataset.Fields").
		Preload("Extraction.TargetDataset").
		Preload("Extraction.TargetDataset.Connection").
		Preload("Extraction.TargetDataset.Connection.ConnectionType").
		Preload("Extraction.TargetDataset.Fields").
		First(session, session.ID).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}

	ses, err := tasks.GetSession(session)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	task, err := tasks.NewSessionRunTask(ses)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if _, err := s.client.Enqueue(task, asynq.MaxRetry(0)); err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, session)
}
