package uid_test

import (
	"sync"
	"testing"

	"github.com/hoodie-ninja/uid"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentGenerationUnique(t *testing.T) {
	const workers, perWorker = 64, 1000
	results := make([][]uid.UUID, workers)
	var wg sync.WaitGroup
	for w := range workers {
		wg.Go(func() {
			ids := make([]uid.UUID, 0, perWorker)
			for i := range perWorker {
				if i%2 == 0 {
					ids = append(ids, uid.NewV4())
				} else {
					ids = append(ids, uid.NewV7())
				}
			}
			results[w] = ids // distinct index per goroutine: no shared write
		})
	}
	wg.Wait()
	// collect after the barrier and assert global uniqueness
	seen := make(map[uid.UUID]bool, workers*perWorker)
	for _, ids := range results {
		for _, id := range ids {
			assert.False(t, seen[id], "duplicate UUID generated concurrently: %s", id)
			seen[id] = true
		}
	}
	assert.Len(t, seen, workers*perWorker)
}
