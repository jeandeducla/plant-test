package plants

import (
	"github.com/jeandeducla/api-plant/internal/models"
)

type ServiceInterface interface {
    GetAllEnergyManagers() ([]models.EnergyManager, error)
    CreateEnergyManager(input CreateEnergyManagerInput) error
    GetEnergyManager(id uint) (*models.EnergyManager, error)
    DeleteEnergyManager(id uint) error
    UpdateEnergyManager(id uint, input UpdateEnergyManagerInput) error
}

type Service struct {
    DB *PlantsDB
}

func NewPlantsService(plantsDB *PlantsDB) *Service {
    return &Service{DB: plantsDB}
}

func (s *Service) GetAllEnergyManagers() ([]models.EnergyManager, error) {
    return s.DB.GetAllEnergyManagers()
}

type CreateEnergyManagerInput struct {
    Name    string `json:"name" binding:"required"`
    Surname string `json:"surname" binding:"required"`
}

func (s *Service) CreateEnergyManager(input CreateEnergyManagerInput) error {
    em := models.EnergyManager{
        Name: input.Name,
        Surname: input.Surname,
    }
    return s.DB.CreateEnergyManager(&em)
}

func (s *Service) GetEnergyManager(id uint) (*models.EnergyManager, error) {
    return s.DB.GetEnergyManagerById(id)
}

func (s *Service) DeleteEnergyManager(id uint) error {
    return s.DB.DeleteEnergyManagerById(id)
}

type UpdateEnergyManagerInput struct {
    Name    string `json:"name" binding:"required"`
    Surname string `json:"surname" binding:"required"`
}

func (s *Service) UpdateEnergyManager(id uint, input UpdateEnergyManagerInput) error {
    em, err := s.DB.GetEnergyManagerById(id)
    if err != nil {
        return err
    }
    em.Name = input.Name
    em.Surname = input.Surname
    return s.DB.UpdateEnergyManager(em)
}
