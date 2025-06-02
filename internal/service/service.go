package service

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

//TODO add valid

type CalcilationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationsById(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}

func NewCalculationService(r CalculationRepository) CalcilationService {
	return &calcService{repo: r}
}

func (s *calcService) calcExpr(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	res, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", res), nil
}

func (s *calcService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAllCalcs()
}

func (s *calcService) CreateCalculation(expression string) (Calculation, error) {
	res, err := s.calcExpr(expression)
	if err != nil {
		return Calculation{}, err
	}
	calc := Calculation{
		Id:         uuid.NewString(),
		Expression: expression,
		Result:     res,
	}
	if err := s.repo.CreateCalc(calc); err != nil {
		return Calculation{}, err
	}
	return calc, nil
}

func (s *calcService) GetCalculationsById(id string) (Calculation, error) {
	return s.repo.GetCalcById(id)
}

func (s *calcService) UpdateCalculation(id, expression string) (Calculation, error) {
	calc, err := s.repo.GetCalcById(id)
	if err != nil {
		return Calculation{}, err
	}
	res, err := s.calcExpr(expression)
	if err != nil {
		return Calculation{}, err
	}
	calc.Expression = expression
	calc.Result = res

	if err := s.repo.UpdateCalcById(calc); err != nil {
		return Calculation{}, err
	}
	return calc, nil
}

func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalcById(id)
}
