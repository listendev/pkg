package gen

import (
	"fmt"
	"sync"

	"github.com/garnet-org/pkg/analysisrequest"
	"github.com/garnet-org/pkg/rand"
)

var pool = sync.Pool{
	New: func() interface{} {
		return New()
	},
}

func New() *analysisrequest.NPMAnalysisRequest {
	snowflakeID := rand.String(19)
	name := rand.String(rand.Range(3, 20))
	vers := fmt.Sprintf("%d.%d.%d", rand.Range(0, 42), rand.Range(0, 42), rand.Range(0, 42))
	shasum := rand.String(40)
	ret := analysisrequest.NewNPMAnalysisRequest(snowflakeID, name, vers, shasum)

	return &ret
}

func Get(reuseProbability int) *analysisrequest.NPMAnalysisRequest {
	if reuseProbability > rand.Range(0, 100) {
		r := pool.Get().(*analysisrequest.NPMAnalysisRequest)
		// Put it back because Get() removes it from the pool
		defer pool.Put(r)

		return r
	}
	r := New()
	pool.Put(r)

	return r
}
