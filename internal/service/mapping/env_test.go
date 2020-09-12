package mapping

import (
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewEnvMappingService(t *testing.T) {
	assert := assert.New(t)

	cfg := domain.MappingConfig{Prefix: "test_"}
	service := NewEnvMappingService(cfg, []string{
		// Yes
		"TEST_KEY1=val1",
		"test_key2=VAL2",
		"Test_Key3=Val3",
		// No
		"TESTKEY4=val4",
		"APP_key5=VAL5",
		"Key6Test_=Val6",
	})
	var result []domain.Mapping
	err := service.GetAll(testCtx, &result)
	assert.Nil(err)
	assert.Equal(3, len(result))
	assert.True(containsMapping(result, domain.Mapping{Id: "key1", Host: "val1"}))
	assert.True(containsMapping(result, domain.Mapping{Id: "key2", Host: "VAL2"}))
	assert.True(containsMapping(result, domain.Mapping{Id: "key3", Host: "Val3"}))

	var mapping domain.Mapping

	err = service.GetById(testCtx, "key1", &mapping)
	assert.Nil(err)
	assert.True(reflect.DeepEqual(domain.Mapping{Id: "key1", Host: "val1"}, mapping))

	err = service.GetById(testCtx, "KEY1", &mapping)
	assert.Nil(err)
	assert.True(reflect.DeepEqual(domain.Mapping{Id: "key1", Host: "val1"}, mapping))

	err = service.GetById(testCtx, "Key1", &mapping)
	assert.Nil(err)
	assert.True(reflect.DeepEqual(domain.Mapping{Id: "key1", Host: "val1"}, mapping))

	err = service.GetById(testCtx, "key2", &mapping)
	assert.Nil(err)
	assert.True(reflect.DeepEqual(domain.Mapping{Id: "key2", Host: "VAL2"}, mapping))

	err = service.GetById(testCtx, "key3", &mapping)
	assert.Nil(err)
	assert.True(reflect.DeepEqual(domain.Mapping{Id: "key3", Host: "Val3"}, mapping))

	err = service.GetById(testCtx, "key4", &mapping)
	assert.NotNil(err)
	assert.True(domain.ECode(err) == domain.ErrNotFound)
}
