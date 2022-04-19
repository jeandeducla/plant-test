package plants

import (
	"github.com/jeandeducla/api-plant/internal/models"
)

type Service struct {
    DB *PlantsDB
}

func NewPlantsService(plantsDB *PlantsDB) *Service {
    return &Service{DB: plantsDB}
}

func (s *Service) GetAllEnergyManagers() ([]models.EnergyManager, error) {
    return s.DB.GetAllEnergyManagers()
}

func (s *Service) GetEnergyManager(id uint) (*models.EnergyManager, error) {
    return s.DB.GetEnergyManagerById(id)
}

func (s *Service) DeleteEnergyManager(id uint) error {
    return s.DB.DeleteEnergyManagerById(id)
}
