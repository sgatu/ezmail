package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type DomainInfo struct {
	bun.BaseModel `bun:"table:domain,alias:di"`
	Id            string      `bun:",pk"`
	DomainName    string      `bun:",notnull"`
	UserId        string      `bun:",notnull"`
	Created       time.Time   `bun:",notnull"`
	DnsRecords    []DnsRecord `bun:"records,notnull,msgpack"`
	Validated     bool        `bun:",notnull"`
}

type DnsRecordStatus int

const (
	DNS_RECORD_STATUS_PENDING = iota
	DNS_RECORD_STATUS_VERIFIED
	DNS_RECORD_STATUS_FAILED
)

type DnsRecord struct {
	Type   string
	Value  string
	Status DnsRecordStatus
}

type DomainInfoRepository interface {
	Save(ctx context.Context, di *DomainInfo) error
}
