package test

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/models/mocks_models"
	"github.com/sgatu/ezmail/internal/domain/services"
)

func getDefaultEmailStorageService() (
	services.EmailStoreService,
	*mocks_models.EmailRepositoryMock,
	*mocks_models.TemplateRepositoryMock,
	*mocks_models.DomainRepositoryMock,
	*mocks_models.EventBusMock,
	*mocks_models.ScheduledEventRepositoryMock,
	*snowflake.Node,
) {
	sn, _ := snowflake.NewNode(0)
	em := mocks_models.MockEmailRepository()
	tp := mocks_models.MockTemplateRepository()
	dm := mocks_models.MockDomainRepository()
	ev := mocks_models.MockEventBus()
	sc := mocks_models.MockScheduledEventRepository()

	srv := services.NewDefaultEmailStoreService(em, tp, dm, ev, sn, sc)
	return srv, em, tp, dm, ev, sc, sn
}

func getDefaultEmailStorageService_NoScheduler() (
	services.EmailStoreService,
	*mocks_models.EmailRepositoryMock,
	*mocks_models.TemplateRepositoryMock,
	*mocks_models.DomainRepositoryMock,
	*mocks_models.EventBusMock,
	*snowflake.Node,
) {
	sn, _ := snowflake.NewNode(0)
	em := mocks_models.MockEmailRepository()
	tp := mocks_models.MockTemplateRepository()
	dm := mocks_models.MockDomainRepository()
	ev := mocks_models.MockEventBus()

	srv := services.NewDefaultEmailStoreService(em, tp, dm, ev, sn, nil)
	return srv, em, tp, dm, ev, sn
}

func TestDefaultGetById(t *testing.T) {
	svc, em, _, _, _, _, _ := getDefaultEmailStorageService()
	svc.GetById(context.TODO(), 0)
	if em.GetByIdCalls != 1 {
		t.Errorf("Expected 1 but found %d", em.GetByIdCalls)
	}
}

func TestDefaultMarkEmailAsSent(t *testing.T) {
	svc, em, _, _, _, _, sn := getDefaultEmailStorageService()
	dt := models.DateTime(time.Now())
	eml, _ := email.NewEmail(sn, 1, 1, "test", "test", "", "", make(map[string]string), &dt)
	em.SetGetByIdReturn(eml, nil)
	svc.MarkEmailAsSent(context.TODO(), eml.Id)
	if em.LastSave == nil {
		t.Errorf("Email save not called")
	} else {
		if em.LastSave.ProcessedDate.IsZero() || time.Now().UTC().Before(em.LastSave.ProcessedDate.Time) {
			t.Errorf("Expected processed date to be set")
		}
	}
}

func TestDefaultMarkEmailAsSentFail(t *testing.T) {
	svc, em, _, _, _, _, _ := getDefaultEmailStorageService()
	err := fmt.Errorf("")
	em.SetGetByIdReturn(nil, err)
	errReturn := svc.MarkEmailAsSent(context.TODO(), 1)
	if errReturn != err {
		t.Error("Expected error return")
	}
}

func TestDefaultPrepareEmail_RetrieveFailEmail(t *testing.T) {
	svc, em, _, _, _, _, _ := getDefaultEmailStorageService()
	err := fmt.Errorf("")
	em.SetGetByIdReturn(nil, err)
	_, errReturn := svc.PrepareEmail(context.TODO(), 1)
	if errReturn != err {
		t.Error("Expected error return")
	}
}

func TestDefaultPrepareEmail_RetrieveFailTemplate(t *testing.T) {
	svc, em, tpr, _, _, _, sn := getDefaultEmailStorageService()
	dt := models.DateTime(time.Now())
	eml, _ := email.NewEmail(sn, 1, 1, "test", "test", "", "", make(map[string]string), &dt)
	em.SetGetByIdReturn(eml, nil)
	err := fmt.Errorf("")
	tpr.SetGetByIdReturn(nil, err)
	_, errReturn := svc.PrepareEmail(context.TODO(), 1)
	if errReturn != err {
		t.Error("Expected error return")
	}
}

func TestDefaultPrepareEmail_RetrieveFailDomain(t *testing.T) {
	svc, em, tpr, dir, _, _, sn := getDefaultEmailStorageService()
	dt := models.DateTime(time.Now())
	eml, _ := email.NewEmail(sn, 1, 1, "test", "test", "", "", make(map[string]string), &dt)
	em.SetGetByIdReturn(eml, nil)
	tp := email.NewTemplate(sn, "", "", "")
	tpr.SetGetByIdReturn(tp, nil)
	err := fmt.Errorf("")
	dir.SetGetByIdReturn(nil, err)
	_, errReturn := svc.PrepareEmail(context.TODO(), 1)
	if errReturn != err {
		t.Error("Expected error return")
	}
}

func TestDefaultPrepareEmail_Ok(t *testing.T) {
	svc, em, tpr, dir, _, _, sn := getDefaultEmailStorageService()
	dt := models.DateTime(time.Now())
	ctx := map[string]string{"test_key": "test_value"}
	eml, _ := email.NewEmail(sn, 1, 1, "test_from", "test_to", "test_reply_to", "test_bcc", ctx, &dt)
	em.SetGetByIdReturn(eml, nil)
	tp := email.NewTemplate(sn, "test-txt-[[test_key]]", "test-html-[[test_key]]", "subject-[[test_key]]")
	tpr.SetGetByIdReturn(tp, nil)
	di := &domain.DomainInfo{
		Id:         1,
		DomainName: "test.tld",
		Validated:  true,
	}
	dir.SetGetByIdReturn(di, nil)
	pe, errReturn := svc.PrepareEmail(context.TODO(), 1)
	if errReturn != nil {
		t.Fatal("Unexpected error")
	}
	if pe.Html != "test-html-test_value" || pe.Text != "test-txt-test_value" || pe.Subject != "subject-test_value" {
		t.Fatal("Email variable replacement is failing while preparing")
	}
	if pe.BCC != "test_bcc" {
		t.Fatal("Email preparation failed to set bcc")
	}
	if pe.From != "test_from" {
		t.Fatal("Email preparation failed to set from")
	}
	if pe.To != "test_to" {
		t.Fatal("Email preparation failed to set to")
	}
	if pe.ReplyTo != "test_reply_to" {
		t.Fatal("Email preparation failed to set reply to")
	}
}

func TestDefaultCreateNewEmail_NoTemplateId(t *testing.T) {
	svc, _, _, _, _, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 0,
	}
	err := svc.CreateEmail(context.TODO(), &crReq)
	if err == nil || err.Error() != email.ErrTemplateNotFound("none").Error() {
		t.Fatal("Expected error on email creation with no template id")
	}
}

func TestDefaultCreateNewEmail_TemplateNotFound(t *testing.T) {
	svc, _, tpr, _, _, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
	}
	err := fmt.Errorf("error")
	tpr.SetGetByIdReturn(nil, err)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate == nil || !errors.Is(err, errCreate) {
		t.Fatal("Expected error on email creation when template not found")
	}
}

func TestDefaultCreateNewEmail_InvalidFromEmail(t *testing.T) {
	svc, _, tpr, _, _, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
	}
	tpr.SetGetByIdReturn(&email.Template{}, nil)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate == nil || errCreate.Error() != "invalid from value" {
		t.Fatal("Expected error on email creation when invalid from email is set")
	}
}

func TestDefaultCreateNewEmail_InvalidToEmails(t *testing.T) {
	svc, _, tpr, _, _, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
		To:         make([]models.EmailAddress, 0),
		From: models.EmailAddress{ParsedEmail: mail.Address{
			Name:    "test",
			Address: "test@test.tld",
		}},
	}
	tpr.SetGetByIdReturn(&email.Template{}, nil)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate == nil || errCreate.Error() != "no email destinataries" {
		t.Fatal("Expected error on email creation when no to emails are set")
	}
}

func TestDefaultCreateNewEmail_DomainNotFound(t *testing.T) {
	svc, _, tpr, dir, _, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
		To: []models.EmailAddress{
			{StringEmail: "to@other.tld"},
		},
		From: models.EmailAddress{ParsedEmail: mail.Address{
			Name:    "test",
			Address: "test@test.tld",
		}},
	}
	tpr.SetGetByIdReturn(&email.Template{}, nil)
	err := fmt.Errorf("error")
	dir.SetGetByNameReturn(nil, err)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate == nil || !errors.Is(err, errCreate) {
		t.Fatal("Expected error on email creation when domain not found")
	}
}

func TestDefaultCreateNewEmail_DomainNotValidated(t *testing.T) {
	svc, _, tpr, dir, _, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
		To: []models.EmailAddress{
			{StringEmail: "to@other.tld"},
		},
		From: models.EmailAddress{ParsedEmail: mail.Address{
			Name:    "test",
			Address: "test@test.tld",
		}},
	}
	tpr.SetGetByIdReturn(&email.Template{}, nil)
	dir.SetGetByNameReturn(&domain.DomainInfo{Validated: false}, nil)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate == nil || errCreate.Error() != "domain not validated" {
		t.Fatal("Expected error on email creation when domain not validated")
	}
}

func TestDefaultCreateNewEmail_NoScheduler(t *testing.T) {
	svc, _, tpr, dir, _, _ := getDefaultEmailStorageService_NoScheduler()
	dtWhen := models.DateTime(time.Now())
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
		To: []models.EmailAddress{
			{StringEmail: "to@other.tld"},
		},
		From: models.EmailAddress{ParsedEmail: mail.Address{
			Name:    "test",
			Address: "test@test.tld",
		}},
		When: &dtWhen,
	}
	tpr.SetGetByIdReturn(&email.Template{}, nil)
	dir.SetGetByNameReturn(&domain.DomainInfo{Validated: true}, nil)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate == nil || errCreate.Error() != "scheduled email without scheduling configuration" {
		t.Fatal("Expected error on email creation when scheduled time set but no scheduler defined")
	}
}

func TestDefaultCreateNewEmail_OkNow(t *testing.T) {
	svc, emr, tpr, dir, evBus, _, _ := getDefaultEmailStorageService()
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
		To: []models.EmailAddress{
			{StringEmail: "to@other.tld"},
		},
		From: models.EmailAddress{ParsedEmail: mail.Address{
			Name:    "test",
			Address: "test@test.tld",
		}},
	}
	tpr.SetGetByIdReturn(&email.Template{Id: 2}, nil)
	dir.SetGetByNameReturn(&domain.DomainInfo{Validated: true, Id: 3}, nil)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate != nil {
		t.Fatal("Unexpected error while creating email")
	}
	if emr.SaveCalls != 1 {
		t.Fatal("Expected call to save email on email repository")
	}
	if evBus.PushCalls != 1 {
		t.Fatal("Expected event of new email to be pushed to event bus")
	}
}

func TestDefaultCreateNewEmail_OkScheduled(t *testing.T) {
	svc, emr, tpr, dir, _, schRepo, _ := getDefaultEmailStorageService()
	dtWhen := models.DateTime(time.Now().Add(5000 * time.Second))
	crReq := email.CreateNewEmailRequest{
		TemplateId: 1,
		To: []models.EmailAddress{
			{StringEmail: "to@other.tld"},
		},
		From: models.EmailAddress{ParsedEmail: mail.Address{
			Name:    "test",
			Address: "test@test.tld",
		}},
		BCC: []models.EmailAddress{
			{StringEmail: "to@other.tld"},
		},
		When: &dtWhen,
	}
	tpr.SetGetByIdReturn(&email.Template{Id: 2}, nil)
	dir.SetGetByNameReturn(&domain.DomainInfo{Validated: true, Id: 3}, nil)
	errCreate := svc.CreateEmail(context.TODO(), &crReq)
	if errCreate != nil {
		t.Fatal("Unexpected error while creating email")
	}
	if emr.SaveCalls != 1 {
		t.Fatal("Expected call to save email on email repository")
	}
	if schRepo.PushCalls != 1 {
		t.Fatal("Expected event of new email to be pushed to scheduler queue")
	}
}
