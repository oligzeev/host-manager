package mapping

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/stretchr/testify/mock"
	"reflect"
)

type MockCtx string

var testCtx = context.WithValue(context.Background(), MockCtx("mock"), "test")

func containsMapping(mappings []domain.Mapping, mapping domain.Mapping) bool {
	for _, element := range mappings {
		if reflect.DeepEqual(element, mapping) {
			return true
		}
	}
	return false
}

type MockMappingService struct {
	mock.Mock
}

func (m *MockMappingService) GetAll(ctx context.Context, result *[]domain.Mapping) error {
	arguments := m.Called(ctx, result)
	return arguments.Error(0)
}

func (m *MockMappingService) GetById(ctx context.Context, id string, result *domain.Mapping) error {
	arguments := m.Called(ctx, id, result)
	return arguments.Error(0)
}
