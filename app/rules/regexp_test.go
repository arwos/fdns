package rules

import (
	"testing"

	"github.com/miekg/dns"
	"go.osspkg.com/casecheck"
)

func TestUnit_NewRegexpRule1(t *testing.T) {
	rr, err := newRegexpRule(`(?U)(?P<key>\d+)\.aaa\.ru\.`, dns.TypeA, []string{"1.1.1.$key", "1.1.2.$key"})
	casecheck.NoError(t, err)
	casecheck.True(t, rr.Match(`12.aaa.ru.`))
	casecheck.Equal(t,
		[]string{"12.aaa.ru.\t3600\tIN\tA\t1.1.1.12", "12.aaa.ru.\t3600\tIN\tA\t1.1.2.12"},
		rr.Compile(dns.TypeA, `12.aaa.ru.`),
	)
}

func TestUnit_NewRegexpRule2(t *testing.T) {
	rr, err := newRegexpRule(`(?U)(.*)\.aaa\.ru\.`, dns.TypeA, []string{"1.1.1.1", "1.1.2.1"})
	casecheck.NoError(t, err)
	casecheck.True(t, rr.Match(`12.aaa.ru.`))
	casecheck.Equal(t,
		[]string{"12.aaa.ru.\t3600\tIN\tA\t1.1.1.1", "12.aaa.ru.\t3600\tIN\tA\t1.1.2.1"},
		rr.Compile(dns.TypeA, `12.aaa.ru.`),
	)
}
