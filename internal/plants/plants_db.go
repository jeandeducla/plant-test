package plants

import (
	"errors"

	"github.com/jeandeducla/api-plant/internal/models"
	"gorm.io/gorm"
)

var (
    ErrEmptyResult = errors.New("Empty result")
)

type PlantsDB struct  {
    gorm *gorm.DB
}

func NewPlantsDB(db *gorm.DB) *PlantsDB {
    return &PlantsDB{gorm: db}
}

func (db *PlantsDB) GetAllEnergyManagers() ([]models.EnergyManager, error) {
    var ems []models.EnergyManager
    if err := db.gorm.Find(&ems).Error; err != nil {
        return nil, err
    }
    return ems, nil
}

func (db *PlantsDB) CreateEnergyManager(em *models.EnergyManager) error {
    result := db.gorm.Create(em)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) GetEnergyManagerById(id uint) (*models.EnergyManager, error) {
    var em models.EnergyManager
    result := db.gorm.Find(&em, id)
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return nil, ErrEmptyResult
    }
    return &em, nil
}
    
func (db *PlantsDB) DeleteEnergyManagerById(id uint) error {
    result := db.gorm.Delete(&models.EnergyManager{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}
