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
		"alpinejs-component 1.1.2 31339edcd3eb3f9c9bb5cedb775cd33a141efaed",
		"@alpinejs/ui 3.10.5-beta.8 a4b0edc67b3c8bcd6deafdb417579e5bbdfbedc4",
		"alpinejs-notify 1.0.2 6e7475154cf3b171602d161aa45bd77cf33070bf",
		"alpinejs-ray 1.1.1 8c65b5cf52f157f280bb0b84b1ddc0c2d4936693",
		"@alpinejs/trap 3.2.3 a4900ecf3ebc797345380bc506fb8f4203c9e267",
		"chalk 4.0.0 6e98081ed2d17faab615eb52ac66ec1fe6209e72",
		"postcss-clean 1.2.1 0b1636c2a961fc1862856616484689be7eaff417",
		"debug 4.2.0 7f150f93920e94c58f5574c2fd01a3110effe7f1",
		"kube 1.2.4 301f0f72e3bf18d7ace131e78acc16fc69c09da2",
		"vue 0.8.2 c1d30517b5160982a48ea22022b6974bd1bbde6a",
		"@babel/highlight 7.9.0 4e9b45ccb82b79607271b2979ad82c7b68163079",
		"support-color 7.1.0 6e6e25a258e16a0148cdc92f0950e60c9c24617c",
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
