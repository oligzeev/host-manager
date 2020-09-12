package mapping

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkGetAll10(b *testing.B) {
	benchmarkGetAll(10, b)
}

func BenchmarkGetAll100(b *testing.B) {
	benchmarkGetAll(100, b)
}

func BenchmarkGetAll1000(b *testing.B) {
	benchmarkGetAll(1000, b)
}

func BenchmarkGetAll10000(b *testing.B) {
	benchmarkGetAll(10000, b)
}

func benchmarkGetAll(count int, b *testing.B) {
	ctx := context.Background()
	hostMap := make(map[string]domain.Mapping)
	var hosts []domain.Mapping
	for i := 0; i < 100; i++ {
		idx := strconv.Itoa(i)
		mapping := domain.Mapping{Id: "id" + idx, Host: "host" + idx}
		hostMap["id"+idx] = mapping
		hosts = append(hosts, mapping)
	}
	envService := &EnvMappingService{hostMap: hostMap, hosts: hosts}
	osService := &OpenshiftMappingService{hosts: hostMap, mutex: sync.RWMutex{}}
	service := NewAggMappingService(envService, osService)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []domain.Mapping
		if err := service.GetAll(ctx, &result); err != nil {
			b.Fatalf("can't get all mappings: %v", err)
		}
	}
}
