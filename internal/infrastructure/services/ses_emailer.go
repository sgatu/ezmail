package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/quotedprintable"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/sgatu/ezmail/internal/domain/services"
	"github.com/sgatu/ezmail/internal/infrastructure/services/thirdparty"
)

type SesEmailer struct {
	awsSesClient  thirdparty.SESClient
	snowflakeNode thirdparty.BaseSnowflakeNode
}

func NewSesEmailer(sesClient thirdparty.SESClient, snowflakeNode thirdparty.BaseSnowflakeNode) *SesEmailer {
	return &SesEmailer{
		awsSesClient:  sesClient,
		snowflakeNode: snowflakeNode,
	}
}

func (se *SesEmailer) SendEmail(ctx context.Context, email *services.PreparedEmail) error {
	var toAddresses []string
	err := json.Unmarshal([]byte(email.To), &toAddresses)
	if err != nil {
		return err
	}
	if len(toAddresses) == 0 {
		return fmt.Errorf("no to address")
	}
	var bccAddresses []string
	err = json.Unmarshal([]byte(email.BCC), &bccAddresses)
	if err != nil {
		return err
	}
	emailData := prepareEmail(email, toAddresses, se.snowflakeNode)
	replyTo := []string{}
	if email.ReplyTo != "" {
		replyTo = []string{email.ReplyTo}
	}
	_, err = se.awsSesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		ReplyToAddresses: replyTo,
		FromEmailAddress: &email.From,
		Destination: &types.Destination{
			ToAddresses:  toAddresses,
			BccAddresses: bccAddresses,
		},
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: emailData,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func prepareEmail(email *services.PreparedEmail, toAddresses []string, snowflakeNode thirdparty.BaseSnowflakeNode) []byte {
	var emailBuffer bytes.Buffer
	boundary := snowflakeNode.Generate().Base36()
	writeHeader(&emailBuffer, "From", email.From)
	writeHeader(&emailBuffer, "To", strings.Join(toAddresses, ", "))
	writeHeader(&emailBuffer, "Subject", email.Subject)
	writeHeader(&emailBuffer, "Content-Type", fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary))
	emailBuffer.WriteString(fmt.Sprintf("\n--%s\n", boundary))
	writeHeader(&emailBuffer, "Content-Type", fmt.Sprintf("multipart/alternative; boundary=\"sub_%s\"", boundary))
	emailBuffer.WriteString(fmt.Sprintf("\n--sub_%s\n", boundary))
	writeHeader(&emailBuffer, "Content-Type", "text/plain; charset=UTF-8")
	writeHeader(&emailBuffer, "Content-Transfer-Encoding", "quoted-printable")
	emailBuffer.WriteString("\n")
	writeQuotedPrintable(&emailBuffer, email.Text)
	emailBuffer.WriteString(fmt.Sprintf("\n--sub_%s\n", boundary))
	writeHeader(&emailBuffer, "Content-Type", "text/html; charset=UTF-8")
	writeHeader(&emailBuffer, "Content-Transfer-Encoding", "quoted-printable")
	emailBuffer.WriteString("\n")
	writeQuotedPrintable(&emailBuffer, email.Html)
	emailBuffer.WriteString(fmt.Sprintf("\n--sub_%s--\n\n--%s--\n", boundary, boundary))
	return emailBuffer.Bytes()
}

func writeHeader(buf *bytes.Buffer, headerName string, value string) {
	buf.WriteString(headerName)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteString("\n")
}

func writeQuotedPrintable(buf *bytes.Buffer, data string) {
	quotedWritter := quotedprintable.NewWriter(buf)
	// since bytes.Buffer wont return err neither will quotedWritter
	quotedWritter.Write([]byte(data))
	quotedWritter.Close()
}
