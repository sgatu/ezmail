package test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/services"
	iservices "github.com/sgatu/ezmail/internal/infrastructure/services"
	"github.com/sgatu/ezmail/internal/thirdparty/mock"
)

func getPreparedEmail() *services.PreparedEmail {
	return &services.PreparedEmail{
		Id:      1,
		Html:    "html",
		Text:    "text",
		From:    "from",
		To:      "[\"test@test.com\"]",
		BCC:     "[\"testbcc@test.com\"]",
		ReplyTo: "replyto",
		Subject: "subject",
		Domain: struct {
			Name string
			Id   int64 "json:\",string\""
		}{
			Name: "domain",
			Id:   1,
		},
	}
}

func getPreparedEmailInvalidTo() *services.PreparedEmail {
	pe := getPreparedEmail()
	pe.To = "invalid"
	return pe
}

func getPreparedEmailEmptyTo() *services.PreparedEmail {
	pe := getPreparedEmail()
	pe.To = "[]"
	return pe
}

func getPreparedEmailInvalidBcc() *services.PreparedEmail {
	pe := getPreparedEmail()
	pe.BCC = "invalid"
	return pe
}

func TestSendEmailInvalidTo(t *testing.T) {
	client := mock.MockSesV2Client{}
	sn, _ := snowflake.NewNode(8)
	s := iservices.NewSesEmailer(&client, sn)
	err := s.SendEmail(context.TODO(), getPreparedEmailInvalidTo())
	if err == nil {
		t.Fatal("Expected error thrown, invalid to email")
	}
}

func TestSendEmailEmptyTo(t *testing.T) {
	client := mock.MockSesV2Client{}
	sn, _ := snowflake.NewNode(8)
	s := iservices.NewSesEmailer(&client, sn)
	err := s.SendEmail(context.TODO(), getPreparedEmailEmptyTo())
	if err == nil {
		t.Fatal("Expected error thrown, empty to email")
	}
}

func TestSendEmailInvalidBcc(t *testing.T) {
	client := mock.MockSesV2Client{}
	sn, _ := snowflake.NewNode(8)
	s := iservices.NewSesEmailer(&client, sn)
	err := s.SendEmail(context.TODO(), getPreparedEmailInvalidBcc())
	if err == nil {
		t.Fatal("Expected error thrown, invalid to email")
	}
}

func TestSendEmailSesFail(t *testing.T) {
	client := mock.MockSesV2Client{}
	errSend := fmt.Errorf("err")
	client.SetSendEmailResponse(nil, errSend)
	sn, _ := snowflake.NewNode(8)
	s := iservices.NewSesEmailer(&client, sn)
	err := s.SendEmail(context.TODO(), getPreparedEmail())
	if err != errSend {
		t.Fatal("Expected error thrown, ses fail")
	}
}

func TestSendEmailSesOk(t *testing.T) {
	client := mock.MockSesV2Client{}
	client.SetSendEmailResponse(nil, nil)
	sn := &mock.MockSnowflakeNode{}
	idBoundary := snowflake.ID(99999999999)
	boundaryB36 := idBoundary.Base36()
	sn.SetNextId(idBoundary)
	s := iservices.NewSesEmailer(&client, sn)
	err := s.SendEmail(context.TODO(), getPreparedEmail())
	if err != nil {
		t.Fatal("Unexpected error thrown, send email ok")
	}
	callParams := client.GetLastSendParams()
	if !slices.Equal(callParams.SendEmailInput.Destination.ToAddresses, []string{"test@test.com"}) {
		t.Fatal("Expected to emails to be sent to aws")
	}
	if !slices.Equal(callParams.SendEmailInput.Destination.BccAddresses, []string{"testbcc@test.com"}) {
		t.Fatal("Expected bcc emails to be sent to aws")
	}
	emailSent := string(callParams.SendEmailInput.Content.Raw.Data)
	lines := strings.Split(emailSent, "\n")
	if lines[0] != "From: from" {
		t.Fatal("not valid from in email preparation")
	}
	if lines[1] != "To: test@test.com" {
		t.Fatal("not valid to in email preparation")
	}
	if lines[2] != "Subject: subject" {
		t.Fatal("not valid subject in email preparation")
	}
	if lines[3] != fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"", boundaryB36) {
		t.Fatal("not valid content type and boundary in email preparation")
	}
	if lines[4] != "" {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[5] != fmt.Sprintf("--%s", boundaryB36) {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[6] != fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"sub_%s\"", boundaryB36) {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[7] != "" {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[8] != fmt.Sprintf("--sub_%s", boundaryB36) {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[9] != "Content-Type: text/plain; charset=UTF-8" {
		t.Fatal("not valid content type for plain/text in email preparation")
	}
	if lines[10] != "Content-Transfer-Encoding: quoted-printable" {
		t.Fatal("not valid content transfer encoding for plain/text in email preparation")
	}
	if lines[11] != "" {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[12] != "text" {
		t.Fatal("not valid text body in email preparation")
	}
	if lines[13] != fmt.Sprintf("--sub_%s", boundaryB36) {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[14] != "Content-Type: text/html; charset=UTF-8" {
		t.Fatal("not valid content type for text/html in email preparation")
	}
	if lines[15] != "Content-Transfer-Encoding: quoted-printable" {
		t.Fatal("not valid content transfer encoding for plain/html in email preparation")
	}
	if lines[16] != "" {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[17] != "html" {
		t.Fatal("not valid html body in email preparation")
	}
	if lines[18] != fmt.Sprintf("--sub_%s--", boundaryB36) {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[19] != "" {
		t.Fatal("not valid separation in email preparation")
	}
	if lines[20] != fmt.Sprintf("--%s--", boundaryB36) {
		t.Fatal("not valid separation in email preparation")
	}
}
