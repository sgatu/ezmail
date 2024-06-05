package ses

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/user"
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
	user *user.User,
	domainName string,
	region string,
) (*domain.DomainInfo, error) {
	_, err := s.awsSesClient.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
		EmailIdentity:         &domainName,
		DkimSigningAttributes: &types.DkimSigningAttributes{NextSigningKeyLength: types.DkimSigningKeyLengthRsa2048Bit},
		Tags: []types.Tag{
			{Key: aws.String("user"), Value: aws.String(user.Id)},
		},
	}, func(o *sesv2.Options) { o.Region = region })
	if err != nil {
		return nil, err
	}
	_, err = s.awsSesClient.PutEmailIdentityMailFromAttributes(ctx, &sesv2.PutEmailIdentityMailFromAttributesInput{
		EmailIdentity: &domainName,
	}, func(o *sesv2.Options) { o.Region = region })
	if err != nil {
		s.awsSesClient.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{EmailIdentity: &domainName}, func(o *sesv2.Options) { o.Region = region })
		return nil, err
	}
	identity, err := s.awsSesClient.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{EmailIdentity: &domainName}, func(o *sesv2.Options) { o.Region = region })
	if err != nil {
		s.awsSesClient.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{EmailIdentity: &domainName}, func(o *sesv2.Options) { o.Region = region })
		return nil, err
	}
	domainInfo := &domain.DomainInfo{
		Id:         s.snowflakeNode.Generate().String(),
		DomainName: domainName,
		UserId:     user.Id,
		Region:     region,
		Created:    time.Now().UTC(),
	}
	records := make([]domain.DnsRecord, 0, 3)
	for _, domainToken := range identity.DkimAttributes.Tokens {
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
		Value:   fmt.Sprintf("10 feedback-smtp.%s.amazonses.com", region),
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
	err = s.domainRepository.Save(ctx, domainInfo)
	if err != nil {
		s.awsSesClient.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{EmailIdentity: &domainName}, func(o *sesv2.Options) { o.Region = region })
		return nil, err
	}
	return domainInfo, nil
}
