package gen

import (
	"fmt"
	"strings"
	"sync"

	"github.com/garnet-org/pkg/analysisrequest"
	"github.com/garnet-org/pkg/rand"
)

var (
	RealPackages = false

	pool = sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
	realPackageTriples = []string{
		"react 18.0.0 b468736d1f4a5891f38585ba8e8fb29f91c3cb96",
		"@vue/devtools 6.5.0 32b09b0e3ca7b757802f9a0b9ded8c2035ce7874",
	}
)

func New() analysisrequest.AnalysisRequest {
	snowflakeID := rand.String(19)
	name := rand.String(rand.Range(3, 20))
	vers := fmt.Sprintf("%d.%d.%d", rand.Range(0, 42), rand.Range(0, 42), rand.Range(0, 42))
	shasum := rand.String(40)
	if RealPackages {
		elems := strings.Split(rand.Elem(realPackageTriples), " ")
		name = elems[0]
		vers = elems[1]
		shasum = elems[2]
	}

	// Not perturbating these
	priority := uint8(0)
	force := false

	randomType := rand.Range(1, int(analysisrequest.MaxType()))
	switch analysisrequest.Type(randomType) {
	case analysisrequest.Nop:
		break

	case analysisrequest.NPMInstallWhileFalco:
		ret, _ := analysisrequest.NewNPM(analysisrequest.NPMInstallWhileFalco, snowflakeID, priority, force, name, vers, shasum)

		return ret

	case analysisrequest.NPMTestWhileFalco:
		ret, _ := analysisrequest.NewNPM(analysisrequest.NPMTestWhileFalco, snowflakeID, priority, force, name, vers, shasum)

		return ret

	case analysisrequest.NPMDepsDev:
		ret, _ := analysisrequest.NewNPM(analysisrequest.NPMDepsDev, snowflakeID, priority, force, name, vers, shasum)

		return ret

	}

	return analysisrequest.NewNOP(snowflakeID, priority, force)
}

func Get(reuseProbability int) analysisrequest.AnalysisRequest {
	if reuseProbability > rand.Range(0, 100) {
		r := pool.Get().(analysisrequest.AnalysisRequest)
		// Put it back because Get() removes it from the pool
		defer pool.Put(r)

		return r
	}
	r := New()
	pool.Put(r)

	return r
}
