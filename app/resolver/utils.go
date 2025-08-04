package resolver

import (
	"fmt"

	"github.com/miekg/dns"
)

func ParseRR(list ...string) ([]dns.RR, error) {
	out := make([]dns.RR, 0, len(list))
	for _, s := range list {
		val, err := dns.NewRR(s)
		if err != nil {
			return nil, fmt.Errorf("parse `%s`: %w", s, err)
		}
		out = append(out, val)
	}
	return out, nil
}
