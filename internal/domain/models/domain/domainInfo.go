package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type DomainInfo struct {
	bun.BaseModel `bun:"table:domain,alias:di"`
	DomainName    string    `bun:",notnull"`
	RawDnsRecords string    `bun:"records,notnull"`
	Region        string    `bun:",notnull"`
	Created       time.Time `bun:",notnull"`
	dnsRecords    []DnsRecord
	Id            int64 `bun:",pk"`
	Validated     bool  `bun:",notnull"`
}

func (di *DomainInfo) GetDnsRecords() ([]DnsRecord, error) {
	if len(di.RawDnsRecords) > 0 && len(di.dnsRecords) == 0 {
		var dnsRecords []DnsRecord
		err := json.Unmarshal([]byte(di.RawDnsRecords), &dnsRecords)
		if err != nil {
			return []DnsRecord{}, err
		}
		di.dnsRecords = dnsRecords
	}
	return di.dnsRecords, nil
}

func (di *DomainInfo) SetDnsRecords(records []DnsRecord) error {
	strDnsRecords, err := json.Marshal(records)
	if err != nil {
		return err
	}
	di.dnsRecords = records
	di.RawDnsRecords = string(strDnsRecords)
	return nil
}

var ErrDomainInfoNotFound error = fmt.Errorf("domain info not found")

type DnsRecordStatus int

const (
	DNS_RECORD_STATUS_PENDING = iota
	DNS_RECORD_STATUS_VERIFIED
	DNS_RECORD_STATUS_FAILED
)

type DnsRecord struct {
	Purpose string          `json:"purpose"` // SPF, DKIM, VALIDATION
	Type    string          `json:"type"`    // MX, CNAME, TXT
	Name    string          `json:"name"`
	Value   string          `json:"value"`
	Status  DnsRecordStatus `json:"status"`
}

type DomainInfoRepository interface {
	GetAll(ctx context.Context) ([]DomainInfo, error)
	GetDomainInfoById(ctx context.Context, id int64) (*DomainInfo, error)
	GetDomainInfoByName(ctx context.Context, name string) (*DomainInfo, error)
	Save(ctx context.Context, di *DomainInfo) error
}
