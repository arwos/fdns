package rules

import (
	"fmt"
	"testing"
)

func TestUnit_CheckRex(t *testing.T) {
	b1 := "!\n! This section contains the list of advertising networks domains, which are hosted on non advertising sites as subdomains.\n! Note, that we only put rules that block full subdomains here and not URL parts (there's `general_url.txt` for that).\n!\n! Good: ||ad.doubleclick.net^\n! Bad: /banner.jpg (should be in general_url.txt)\n! Bad: ||legitwebsite.com^$third-party (should be in adservers.txt)\n!\n!\n||dh8azcl753e1e.cloudfront.net^\n||d3apzcqz3ghyay.cloudfront.net^\n||ads-service.api.linkme.global^\n||d3m6crjuedf6o.cloudfront.net^\n||chuchle.all-usanomination.com^\n||bundle.ppas.monster^"
	rexResult1 := rex1.FindAllSubmatch([]byte(b1), -1)
	for _, b := range rexResult1 {
		fmt.Println(string(b[1]))
	}

	b2 := "\n#=====================================\n# Title: Hosts contributed by Steven Black\n# http://stevenblack.com\n\n0.0.0.0 ad-assets.futurecdn.net\n0.0.0.0 ck.getcookiestxt.com\n0.0.0.0 eu1.clevertap-prod.com\n0.0.0.0 wizhumpgyros.com\n0.0.0.0 coccyxwickimp.com\n0.0.0.0 webmail-who-int.000webhostapp.com\n0.0.0.0 010sec.com\n0.0.0.0 01mspmd5yalky8.com\n0.0.0.0 0byv9mgbn0.com\n0.0.0.0 ns6.0pendns.org\n0.0.0.0 dns.0pengl.com\n0.0.0.0 12724.xyz\n0.0.0.0 21736.xyz"
	rexResult2 := rex2.FindAllSubmatch([]byte(b2), -1)
	for _, b := range rexResult2 {
		fmt.Println(string(b[1]))
	}
}
