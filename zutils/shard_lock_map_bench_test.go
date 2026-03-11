package zutils

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkItems(b *testing.B) {
	slm := NewShardLockMaps()

	for i := 0; i < 10000; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		slm.Items()
	}
}

func BenchmarkKeys(b *testing.B) {
	slm := NewShardLockMaps()

	for i := 0; i < 10000; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		slm.Keys()
	}
}

func BenchmarkMarshalJson(b *testing.B) {
	slm := NewShardLockMaps()

	for i := 0; i < 10000; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		_, err := slm.MarshalJSON()
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}

func BenchmarkSingleInsertAbsent(b *testing.B) {
	slm := NewShardLockMaps()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slm.Set(strconv.Itoa(i), "value")
	}
}

func BenchmarkSingleInsertAbsentSyncMap(b *testing.B) {
	var m sync.Map
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), "value")
	}
}

func BenchmarkSingleInsertPresent(b *testing.B) {
	slm := NewShardLockMaps()
	slm.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slm.Set("key", "value")
	}
}

func BenchmarkSingleInsertPresentSyncMap(b *testing.B) {
	var m sync.Map
	m.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store("key", "value")
	}
}

func BenchmarkMultiInsertDifferentSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	_, set := GetSetSyncMap(&m, finished)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_1_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(1)
	benchmarkMultiInsertDifferentWithMap(slm, b)
}
func BenchmarkMultiInsertDifferent_16_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(16)
	benchmarkMultiInsertDifferentWithMap(slm, b)
}
func BenchmarkMultiInsertDifferent_32_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(32)
	benchmarkMultiInsertDifferentWithMap(slm, b)
}
func BenchmarkMultiInsertDifferent_256_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(256)
	benchmarkMultiInsertDifferentWithMap(slm, b)
}

func BenchmarkMultiInsertSame(b *testing.B) {
	slm := NewShardLockMaps()
	finished := make(chan struct{}, b.N)
	_, set := GetSet(slm, finished)
	slm.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertSameSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	_, set := GetSetSyncMap(&m, finished)
	m.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSame(b *testing.B) {
	slm := NewShardLockMaps()
	finished := make(chan struct{}, b.N)
	get, _ := GetSet(slm, finished)
	slm.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSameSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	get, _ := GetSetSyncMap(&m, finished)
	m.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferentSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSetSyncMap(&m, finished)
	m.Store("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_1_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(1)
	benchmarkMultiGetSetDifferentWithMap(slm, b)
}
func BenchmarkMultiGetSetDifferent_16_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(16)
	benchmarkMultiGetSetDifferentWithMap(slm, b)
}
func BenchmarkMultiGetSetDifferent_32_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(32)
	benchmarkMultiGetSetDifferentWithMap(slm, b)
}
func BenchmarkMultiGetSetDifferent_256_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(256)
	benchmarkMultiGetSetDifferentWithMap(slm, b)
}

func BenchmarkMultiGetSetBlockSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSetSyncMap(&m, finished)
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_1_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(1)
	benchmarkMultiGetSetBlockWithMap(slm, b)
}
func BenchmarkMultiGetSetBlock_16_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(16)
	benchmarkMultiGetSetBlockWithMap(slm, b)
}
func BenchmarkMultiGetSetBlock_32_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(32)
	benchmarkMultiGetSetBlockWithMap(slm, b)
}
func BenchmarkMultiGetSetBlock_256_Shard(b *testing.B) {
	slm := NewShardLockMapsWithCount(256)
	benchmarkMultiGetSetBlockWithMap(slm, b)
}

func benchmarkMultiInsertDifferent(b *testing.B) {
	slm := NewShardLockMaps()
	finished := make(chan struct{}, b.N)
	_, set := GetSet(slm, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func benchmarkMultiGetSetDifferent(b *testing.B) {
	slm := NewShardLockMaps()
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(slm, finished)
	slm.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func benchmarkMultiGetSetBlock(b *testing.B) {
	slm := NewShardLockMaps()
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(slm, finished)
	for i := 0; i < b.N; i++ {
		slm.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func GetSet(slm ShardLockMaps, finished chan struct{}) (set func(key, value string), get func(key, value string)) {
	return func(key, value string) {
			for i := 0; i < 10; i++ {
				slm.Get(key)
			}
			finished <- struct{}{}
		}, func(key, value string) {
			for i := 0; i < 10; i++ {
				slm.Set(key, value)
			}
			finished <- struct{}{}
		}
}

func GetSetSyncMap(m *sync.Map, finished chan struct{}) (get func(key, value string), set func(key, value string)) {
	get = func(key, value string) {
		for i := 0; i < 10; i++ {
			m.Load(key)
		}
		finished <- struct{}{}
	}
	set = func(key, value string) {
		for i := 0; i < 10; i++ {
			m.Store(key, value)
		}
		finished <- struct{}{}
	}
	return
}

func benchmarkMultiInsertDifferentWithMap(slm ShardLockMaps, b *testing.B) {
	finished := make(chan struct{}, b.N)
	_, set := GetSet(slm, finished)
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i), strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func benchmarkMultiGetSetDifferentWithMap(slm ShardLockMaps, b *testing.B) {
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(slm, finished)
	for i := 0; i < b.N; i++ {
		go get(strconv.Itoa(i), "value")
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func benchmarkMultiGetSetBlockWithMap(slm ShardLockMaps, b *testing.B) {
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(slm, finished)
	for i := 0; i < 100; i++ {
		set(strconv.Itoa(i), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}
