package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
)

type SESService struct {
	domainRepository domain.DomainInfoRepository
	awsSesClient     *sesv2.Client
	snowflakeNode    *snowflake.Node
}

func NewSESService(domainRepository domain.DomainInfoRepository, awsConfig aws.Config, snowflakeNode *snowflake.Node) *SESService {
	return &SESService{
		snowflakeNode:    snowflakeNode,
		domainRepository: domainRepository,
		awsSesClient:     sesv2.NewFromConfig(awsConfig),
	}
}

func (s *SESService) CreateDomain(
	ctx context.Context,
	domainInfo *domain.DomainInfo,
) error {
	// create identity
	createIdentityResult, err := s.awsSesClient.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
		EmailIdentity:         &domainInfo.DomainName,
		DkimSigningAttributes: &types.DkimSigningAttributes{NextSigningKeyLength: types.DkimSigningKeyLengthRsa2048Bit},
		Tags: []types.Tag{
			{Key: aws.String("user"), Value: aws.String(domainInfo.UserId)},
		},
	}, func(o *sesv2.Options) { o.Region = domainInfo.Region })
	if err != nil {
		return err
	}

	// create mail from configuration
	_, err = s.awsSesClient.PutEmailIdentityMailFromAttributes(ctx, &sesv2.PutEmailIdentityMailFromAttributesInput{
		EmailIdentity:  &domainInfo.DomainName,
		MailFromDomain: aws.String(fmt.Sprintf("dispatch.%s", domainInfo.DomainName)),
	}, func(o *sesv2.Options) { o.Region = domainInfo.Region })
	if err != nil {
		s.DeleteIdentity(ctx, domainInfo)
		return err
	}
	setDNSRecords(createIdentityResult.DkimAttributes.Tokens, domainInfo)
	err = s.domainRepository.Save(ctx, domainInfo)
	if err != nil {
		s.DeleteIdentity(ctx, domainInfo)
		return err
	}
	return nil
}

func (s *SESService) DeleteIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
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
		Name:    "dispatch",
		Value:   fmt.Sprintf("10 feedback-smtp.%s.amazonses.com", domainInfo.Region),
		Status:  domain.DNS_RECORD_STATUS_PENDING,
	})
	records = append(records, domain.DnsRecord{
		Purpose: "SPF",
		Type:    "TXT",
		Name:    "dispatch",
		Value:   "v=spf1 include:amazonses.com ~all",
		Status:  domain.DNS_RECORD_STATUS_PENDING,
	})
	domainInfo.SetDnsRecords(records)
}