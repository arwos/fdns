package vars

import (
	"fmt"

	"go.osspkg.com/goppy/v2/xdns"
)

const DefaultTTl uint32 = 3600

const dnsRecordTmpl = "%s\t%d\tIN\t%s\t%s"

func BuildDNSRecord(domain string, qtype uint16, value string) string {
	return fmt.Sprintf(dnsRecordTmpl, domain, DefaultTTl, xdns.QTypeString(qtype), value)
}
