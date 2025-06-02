package service

import "gorm.io/gorm"

type CalculationRepository interface {
	CreateCalc(calc Calculation) error
	GetAllCalcs() ([]Calculation, error)
	GetCalcById(id string) (Calculation, error)
	UpdateCalcById(calc Calculation) error
	DeleteCalcById(id string) error
}

//TODO интерфейсы должны быть объявлены в том пакете, где он б использ, а в пакете реализации

type calcRepository struct {
	db *gorm.DB
}

func NewCalcRepository(db *gorm.DB) CalculationRepository {
	return &calcRepository{db: db}
}

func (r *calcRepository) CreateCalc(calc Calculation) error {
	return r.db.Create(&calc).Error
}

func (r *calcRepository) GetAllCalcs() ([]Calculation, error) {
	var calcs []Calculation
	err := r.db.Find(&calcs).Error
	return calcs, err
}

func (r *calcRepository) GetCalcById(id string) (Calculation, error) {
	var calc Calculation
	err := r.db.First(&calc, "id=?", id).Error
	return calc, err
}

func (r *calcRepository) UpdateCalcById(calc Calculation) error {
	return r.db.Save(&calc).Error
}

func (r *calcRepository) DeleteCalcById(id string) error {
	return r.db.Delete(&Calculation{}, "id=?", id).Error
}
