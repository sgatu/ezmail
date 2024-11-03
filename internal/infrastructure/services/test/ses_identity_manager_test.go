package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/mocks_models"
	"github.com/sgatu/ezmail/internal/infrastructure/services"
	"github.com/sgatu/ezmail/internal/thirdparty/mock"
)

func getDomainInfo() *domain.DomainInfo {
	return &domain.DomainInfo{
		DomainName: "test.com",
		Region:     "region",
	}
}

func TestCreateIdentity_ErrorCreate(t *testing.T) {
	client := mock.MockSesV2Client{}
	domRepo := mocks_models.MockDomainRepository()
	sn, _ := snowflake.NewNode(1)
	im := services.NewSESIdentityManager(domRepo, &client, sn)
	errCreate := fmt.Errorf("err")
	client.SetCreateEmailIdentityResponse(nil, errCreate)
	err := im.CreateIdentity(context.TODO(), getDomainInfo())
	if err != errCreate {
		t.Fatal("Expected create identity error return")
	}
}

func TestCreateIdentity_ErrorPut(t *testing.T) {
	client := mock.MockSesV2Client{}
	domRepo := mocks_models.MockDomainRepository()
	sn, _ := snowflake.NewNode(1)
	im := services.NewSESIdentityManager(domRepo, &client, sn)
	errPut := fmt.Errorf("err")
	client.SetCreateEmailIdentityResponse(&sesv2.CreateEmailIdentityOutput{}, nil)
	client.SetPutEmailIdentityMailFromAttributesResponse(nil, errPut)
	err := im.CreateIdentity(context.TODO(), getDomainInfo())
	if err != errPut {
		t.Fatal("Expected create identity error return on put step")
	}
	if client.DeleteIdentityCalls != 1 {
		t.Fatal("Expected a delete call after put fail")
	}
}

func TestCreateIdentity_Ok(t *testing.T) {
	client := mock.MockSesV2Client{}
	domRepo := mocks_models.MockDomainRepository()
	sn, _ := snowflake.NewNode(1)
	im := services.NewSESIdentityManager(domRepo, &client, sn)
	client.SetCreateEmailIdentityResponse(&sesv2.CreateEmailIdentityOutput{
		DkimAttributes: &types.DkimAttributes{
			Tokens: []string{"tok1", "tok2", "tok3"},
		},
	}, nil)
	client.SetPutEmailIdentityMailFromAttributesResponse(nil, nil)
	di := getDomainInfo()
	_ = im.CreateIdentity(context.TODO(), di)
	if client.DeleteIdentityCalls > 0 {
		t.Fatal("Unexpected a delete call after put")
	}
	records, _ := di.GetDnsRecords()
	if len(records) != 5 {
		t.Fatalf("Expected 5 DNS records, found %d", len(records))
	}
	dkim1 := records[0]
	if dkim1.Purpose != "DKIM" || dkim1.Name != "tok1._domainkey" || dkim1.Value != "tok1.dkim.amazonses.com" || dkim1.Type != "CNAME" {
		t.Fatal("Invalid 1st dkim record")
	}
	dkim2 := records[1]
	if dkim2.Purpose != "DKIM" || dkim2.Name != "tok2._domainkey" || dkim2.Value != "tok2.dkim.amazonses.com" || dkim2.Type != "CNAME" {
		t.Fatal("Invalid 2nd dkim record")
	}
	dkim3 := records[2]
	if dkim3.Purpose != "DKIM" || dkim3.Name != "tok3._domainkey" || dkim3.Value != "tok3.dkim.amazonses.com" || dkim3.Type != "CNAME" {
		t.Fatal("Invalid 3rd dkim record")
	}
	spfmx := records[3]
	if spfmx.Purpose != "SPF" || spfmx.Name != services.DOMAIN_PREFIX || spfmx.Value != "10 feedback-smtp.region.amazonses.com" || spfmx.Type != "MX" {
		t.Fatal("Invalid spfmx record")
	}
	spftxt := records[4]
	if spftxt.Purpose != "SPF" || spftxt.Name != services.DOMAIN_PREFIX || spftxt.Value != "v=spf1 include:amazonses.com ~all" || spftxt.Type != "TXT" {
		t.Fatal("Invalid spftxt record")
	}
}

func TestDeleteIdentity_Error(t *testing.T) {
	client := mock.MockSesV2Client{}
	domRepo := mocks_models.MockDomainRepository()
	sn, _ := snowflake.NewNode(1)
	im := services.NewSESIdentityManager(domRepo, &client, sn)
	errDelete := fmt.Errorf("err")
	client.SetDeleteEmailIdentityOutput(nil, errDelete)
	err := im.DeleteIdentity(context.TODO(), getDomainInfo())
	if err != errDelete {
		t.Fatal("Expected delete identity error return")
	}
}

func TestDeleteIdentity_Ok(t *testing.T) {
	client := mock.MockSesV2Client{}
	domRepo := mocks_models.MockDomainRepository()
	sn, _ := snowflake.NewNode(1)
	im := services.NewSESIdentityManager(domRepo, &client, sn)
	client.SetDeleteEmailIdentityOutput(nil, nil)
	err := im.DeleteIdentity(context.TODO(), getDomainInfo())
	if err != nil {
		t.Fatal("Unexpected delete identity error rturn")
	}
}
