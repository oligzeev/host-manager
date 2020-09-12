package mapping

import (
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAggMappingService_GetAll_Success(t *testing.T) {
	assert := assert.New(t)

	var envMappings, osMappings []domain.Mapping
	envService := new(MockMappingService)
	envService.On("GetAll", testCtx, &envMappings).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*[]domain.Mapping)
		*arg = []domain.Mapping{
			{Id: "id1", Host: "host1"},
			{Id: "id2", Host: "host2"},
		}
	})
	osService := new(MockMappingService)
	osService.On("GetAll", testCtx, &osMappings).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*[]domain.Mapping)
		*arg = []domain.Mapping{
			{Id: "id3", Host: "host3"},
			{Id: "id4", Host: "host4"},
		}
	})

	service := NewAggMappingService(envService, osService)
	var result []domain.Mapping
	err := service.GetAll(testCtx, &result)

	assert.Nil(err)
	assert.Equal(4, len(result))
	assert.True(containsMapping(result, domain.Mapping{Id: "id1", Host: "host1"}))
	assert.True(containsMapping(result, domain.Mapping{Id: "id2", Host: "host2"}))
	assert.True(containsMapping(result, domain.Mapping{Id: "id3", Host: "host3"}))
	assert.True(containsMapping(result, domain.Mapping{Id: "id4", Host: "host4"}))
}

func TestAggMappingService_GetById_EnvSuccess(t *testing.T) {
	const (
		id   = "test-id"
		host = "test-host"
	)
	assert := assert.New(t)

	var mapping domain.Mapping
	envService := new(MockMappingService)
	envService.On("GetById", testCtx, id, &mapping).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*domain.Mapping)
		*arg = domain.Mapping{Id: id, Host: host}
	})

	service := NewAggMappingService(envService, nil)
	var result domain.Mapping
	err := service.GetById(testCtx, id, &result)

	assert.Nil(err)
	assert.Equal(id, result.Id)
	assert.Equal(host, result.Host)
}

func TestAggMappingService_GetById_OpenshiftSuccess(t *testing.T) {
	const (
		op   = "Test.GetById"
		id   = "test-id"
		host = "test-host"
	)
	assert := assert.New(t)

	var mapping domain.Mapping
	envService := new(MockMappingService)
	envService.On("GetById", testCtx, id, &mapping).Return(domain.E(op, domain.ErrNotFound))
	osService := new(MockMappingService)
	osService.On("GetById", testCtx, id, &mapping).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*domain.Mapping)
		*arg = domain.Mapping{Id: id, Host: host}
	})

	service := NewAggMappingService(envService, osService)
	var result domain.Mapping
	err := service.GetById(testCtx, id, &result)

	assert.Nil(err)
	assert.Equal(id, result.Id)
	assert.Equal(host, result.Host)
}

func TestAggMappingService_GetById_Error(t *testing.T) {
	const (
		envOp domain.ErrOp = "Test.EnvError"
		osOp  domain.ErrOp = "Test.OsError"
		notOp domain.ErrOp = "Test.NotFoundError"
		envId              = "env-id"
		osId               = "os-id"
	)
	assert := assert.New(t)

	var mapping domain.Mapping
	envService := new(MockMappingService)
	envService.On("GetById", testCtx, envId, &mapping).Return(domain.E(envOp))
	envService.On("GetById", testCtx, osId, &mapping).Return(domain.E(notOp, domain.ErrNotFound))
	osService := new(MockMappingService)
	osService.On("GetById", testCtx, osId, &mapping).Return(domain.E(osOp))

	service := NewAggMappingService(envService, osService)
	var result domain.Mapping
	var err error

	err = service.GetById(testCtx, envId, &result)
	assert.NotNil(err)
	assert.Equal(envOp, domain.EOp(err))

	err = service.GetById(testCtx, osId, &result)
	assert.NotNil(err)
	assert.Equal(osOp, domain.EOp(err))
}

func TestAggMappingService_GetAll_EnvError(t *testing.T) {
	const (
		op domain.ErrOp = "Test.Error"
	)
	assert := assert.New(t)

	var mappings []domain.Mapping
	envService := new(MockMappingService)
	envService.On("GetAll", testCtx, &mappings).Return(domain.E(op))

	service := NewAggMappingService(envService, nil)
	var result []domain.Mapping
	err := service.GetAll(testCtx, &result)

	assert.NotNil(err)
	assert.Equal(op, domain.EOp(err))
}

func TestAggMappingService_GetAll_OpenshiftError(t *testing.T) {
	const (
		op domain.ErrOp = "Test.Error"
	)
	assert := assert.New(t)

	var mappings []domain.Mapping
	envService := new(MockMappingService)
	envService.On("GetAll", testCtx, &mappings).Return(nil)
	osService := new(MockMappingService)
	osService.On("GetAll", testCtx, &mappings).Return(domain.E(op))

	service := NewAggMappingService(envService, osService)
	var result []domain.Mapping
	err := service.GetAll(testCtx, &result)

	assert.NotNil(err)
	assert.Equal(op, domain.EOp(err))
}
