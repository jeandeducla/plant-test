package plants

import (
	"github.com/jeandeducla/api-plant/internal/models"
	"gorm.io/gorm"
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

func (db *PlantsDB) GetEnergyManagerById(id uint) (*models.EnergyManager, error) {
    var em models.EnergyManager
    result := db.gorm.Find(&em, id)
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return nil, nil
    }
    return &em, nil
}
    
