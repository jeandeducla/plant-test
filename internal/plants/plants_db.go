package plants

import (
	"errors"

	"github.com/jeandeducla/api-plant/internal/models"
	"gorm.io/gorm"
)

var (
    ErrEmptyResult = errors.New("Empty result")
)

type DB interface {
    GetAllEnergyManagers() ([]models.EnergyManager, error)
    CreateEnergyManager(em *models.EnergyManager) error
    GetEnergyManagerById(id uint) (*models.EnergyManager, error)
    DeleteEnergyManagerById(id uint) error
    UpdateEnergyManager(em *models.EnergyManager) error

    GetAllPlants() ([]models.Plant, error)
    CreatePlant(plant *models.Plant) error
    GetPlantById(id uint) (*models.Plant, error)
    DeletePlantById(id uint) error
    UpdatePlant(plant *models.Plant) error

    GetPlantsByEnergyManagerId(id uint) ([]models.Plant, error)

    GetAssetsByPlantId(id uint) ([]models.Asset, error)
    CreateAsset(asset *models.Asset) error
    GetAssetByPlantId(plant_id uint, asset_id uint) (*models.Asset, error)
    DeleteAssetById(asset_id uint) error
    UpdateAsset(asset *models.Asset) error
}

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

func (db *PlantsDB) UpdateEnergyManager(em *models.EnergyManager) error {
    result := db.gorm.Model(em).Updates(em)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) GetAllPlants() ([]models.Plant, error) {
    var ems []models.Plant
    if err := db.gorm.Find(&ems).Error; err != nil {
        return nil, err
    }
    return ems, nil
}

func (db *PlantsDB) CreatePlant(plant *models.Plant) error {
    result := db.gorm.Create(plant)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) GetPlantById(id uint) (*models.Plant, error) {
    var plant models.Plant
    result := db.gorm.Find(&plant, id)
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return nil, ErrEmptyResult
    }
    return &plant, nil
}

func (db *PlantsDB) DeletePlantById(id uint) error {
    result := db.gorm.Delete(&models.Plant{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) UpdatePlant(plant *models.Plant) error {
    result := db.gorm.Model(plant).Updates(plant)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) GetPlantsByEnergyManagerId(id uint) ([]models.Plant, error) {
    var plants []models.Plant
    result := db.gorm.Where("energy_manager_id = ?", id).Find(&plants)
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return plants, ErrEmptyResult
    }
    return plants, nil
}

func (db *PlantsDB) GetAssetsByPlantId(id uint) ([]models.Asset, error) {
    var assets []models.Asset
    result := db.gorm.Where("plant_id = ?", id).Find(&assets)
    if result.Error != nil {
        return assets, result.Error
    }
    if result.RowsAffected == 0 {
        return assets, ErrEmptyResult
    }
    return assets, nil
}

func (db *PlantsDB) CreateAsset(asset *models.Asset) error {
    result := db.gorm.Create(asset)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) GetAssetByPlantId(plant_id uint, asset_id uint) (*models.Asset, error) {
    var asset models.Asset
    result := db.gorm.Where("plant_id = ?", plant_id).Find(&asset, asset_id)
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return nil, ErrEmptyResult
    }
    return &asset, nil
}

func (db *PlantsDB) DeleteAssetById(asset_id uint) error {
    result := db.gorm.Delete(&models.Asset{}, asset_id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}

func (db *PlantsDB) UpdateAsset(asset *models.Asset) error {
    result := db.gorm.Model(asset).Updates(asset)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrEmptyResult
    }
    return nil
}
