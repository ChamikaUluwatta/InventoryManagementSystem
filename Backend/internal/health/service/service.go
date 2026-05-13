package service

import "context"

type service struct {
	checkers []HealthChecker
}

func NewService(checkers ...HealthChecker) *service {
	return &service{checkers: checkers}
}

type Service interface {
	Check(ctx context.Context) ([]CheckResult, string)
}

func (s *service) Check(ctx context.Context) ([]CheckResult, string) {
	results := make([]CheckResult, 0, len(s.checkers))
	status := "healthy"
	for _, checker := range s.checkers {
		err := checker.Check(ctx)
		result := CheckResult{
			Name:   checker.Name(),
			Status: "healthy",
		}
		if err != nil {
			result.Status = "unhealthy"
			result.Error = err.Error()
			status = "unhealthy"
		}
		results = append(results, result)
	}
	return results, status
}
