package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/thirdparty"
)

const DOMAIN_PREFIX = "dispatch"

type SESIdentityManager struct {
	awsSesClient thirdparty.SESClient
}

func NewSESIdentityManager(domainRepository domain.DomainInfoRepository, sesClient thirdparty.SESClient, snowflakeNode *snowflake.Node) *SESIdentityManager {
	return &SESIdentityManager{
		awsSesClient: sesClient,
	}
}

func (s *SESIdentityManager) CreateIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
	// create identity
	createIdentityResult, err := s.awsSesClient.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
		EmailIdentity:         &domainInfo.DomainName,
		DkimSigningAttributes: &types.DkimSigningAttributes{NextSigningKeyLength: types.DkimSigningKeyLengthRsa2048Bit},
	}, func(o *sesv2.Options) { o.Region = domainInfo.Region })
	if err != nil {
		return err
	}

	// create mail from configuration
	_, err = s.awsSesClient.PutEmailIdentityMailFromAttributes(ctx, &sesv2.PutEmailIdentityMailFromAttributesInput{
		EmailIdentity:  &domainInfo.DomainName,
		MailFromDomain: aws.String(fmt.Sprintf("%s.%s", DOMAIN_PREFIX, domainInfo.DomainName)),
	}, func(o *sesv2.Options) { o.Region = domainInfo.Region })
	if err != nil {
		s.DeleteIdentity(ctx, domainInfo)
		return err
	}
	setDNSRecords(createIdentityResult.DkimAttributes.Tokens, domainInfo)
	return nil
}

func (s *SESIdentityManager) DeleteIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
	_, err := s.awsSesClient.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{EmailIdentity: &domainInfo.DomainName}, func(o *sesv2.Options) { o.Region = domainInfo.Region })
	return err
}

func setDNSRecords(dkimTokens []string, domainInfo *domain.DomainInfo) {
	records := make([]domain.DnsRecord, 0, 3)
	for _, domainToken := range dkimTokens {
		records = append(records, domain.DnsRecord{
			Purpose: "DKIM",
			Type:    "CNAME",
			Name:    fmt.Sprintf("%s._domainkey", domainToken),
			Value:   fmt.Sprintf("%s.dkim.amazonses.com", domainToken),
			Status:  domain.DNS_RECORD_STATUS_PENDING,
		})
	}

	records = append(records, domain.DnsRecord{
		Purpose: "SPF",
		Type:    "MX",
		Name:    DOMAIN_PREFIX,
		Value:   fmt.Sprintf("10 feedback-smtp.%s.amazonses.com", domainInfo.Region),
		Status:  domain.DNS_RECORD_STATUS_PENDING,
	})
	records = append(records, domain.DnsRecord{
		Purpose: "SPF",
		Type:    "TXT",
		Name:    DOMAIN_PREFIX,
		Value:   "v=spf1 include:amazonses.com ~all",
		Status:  domain.DNS_RECORD_STATUS_PENDING,
	})
	domainInfo.SetDnsRecords(records)
}
