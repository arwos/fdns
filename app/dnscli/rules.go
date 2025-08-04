package dnscli

import (
	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/random"
	"go.osspkg.com/validate"
)

const UniversalZone = "*"

type Rules struct {
	data cache.Cache[string, []string]
}

func NewRules() *Rules {
	return &Rules{
		data: cache.New[string, []string](),
	}
}

func (v *Rules) Replace(vals map[string][]string) {
	v.data.Replace(vals)
}

func (v *Rules) Resolve(zone string) (result []string) {
	defer func() {
		random.Shuffle(result)
	}()

	for i := 2; i >= 0; i-- {
		vv := validate.GetDomainLevel(zone, i)

		if ips, ok := v.data.Get(vv); ok {
			result = append(result, ips...)
			return
		}
	}

	if ips, ok := v.data.Get(UniversalZone); ok {
		result = append(result, ips...)
		return
	}

	return
}
