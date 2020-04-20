package detecting_MX_records

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestMXRecordsLookup(t *testing.T) {
	DNS := "msn.com"
	mxRecords, err := net.LookupMX(DNS)
	assert.Nil(t, err)
	for _, mxRecord := range mxRecords {
		log.Println(mxRecord)
	}
}
