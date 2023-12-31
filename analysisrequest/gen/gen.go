package gen

import (
	"fmt"
	"strings"
	"sync"

	"github.com/listendev/pkg/analysisrequest"
	"github.com/listendev/pkg/rand"
)

var (
	RealPackages = false

	pool = sync.Pool{
		New: func() interface{} {
			return New()
		},
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

	randomType := rand.Range(1, int(analysisrequest.LastType()))
	switch analysisrequest.Type(randomType) {
	case analysisrequest.Nop:
		break

	case analysisrequest.NPMInstallWhileDynamicInstrumentation:
		ret, _ := analysisrequest.NewNPM(analysisrequest.NPMInstallWhileDynamicInstrumentation, snowflakeID, priority, force, name, vers, shasum)

		return ret

	// case analysisrequest.NPMTestWhileDynamicInstrumentation:
	// 	ret, _ := analysisrequest.NewNPM(analysisrequest.NPMTestWhileDynamicInstrumentation, snowflakeID, priority, force, name, vers, shasum)

	// 	return ret

	case analysisrequest.NPMAdvisory:
		ret, _ := analysisrequest.NewNPM(analysisrequest.NPMAdvisory, snowflakeID, priority, force, name, vers, shasum)

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

func GetWithAttrs(force bool, priority uint8, probability int) analysisrequest.AnalysisRequest {
	msg := Get(probability)
	// Since the generator always generate messages with priority 0 and we take the priority from options/flags...
	if priority > 0 {
		msg.SetPrio(priority)
	}
	if force {
		msg.SetForce(force)
	}

	return msg
}
