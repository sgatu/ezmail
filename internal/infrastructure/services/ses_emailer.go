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
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/services"
)

type SesEmailer struct {
	awsSesClient  *sesv2.Client
	snowflakeNode *snowflake.Node
}

func (se *SesEmailer) SendEmail(ctx context.Context, email *services.PreparedEmail) error {
	var toAddresses []string
	err := json.Unmarshal([]byte(email.To), &toAddresses)
	if err != nil {
		return err
	}
	var bccAddresses []string
	err = json.Unmarshal([]byte(email.BCC), &bccAddresses)
	if err != nil {
		return err
	}
	emailData, err := prepareEmail(email, toAddresses, se.snowflakeNode)
	if err != nil {
		return err
	}
	se.awsSesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		ReplyToAddresses: []string{email.ReplyTo},
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
	return nil
}

func prepareEmail(email *services.PreparedEmail, toAddresses []string, snowflakeNode *snowflake.Node) ([]byte, error) {
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
	err := writeQuotedPrintable(&emailBuffer, email.Text)
	if err != nil {
		return nil, err
	}
	emailBuffer.WriteString(fmt.Sprintf("\n--sub_%s\n", boundary))
	writeHeader(&emailBuffer, "Content-Type", "text/html; charset=UTF-8")
	writeHeader(&emailBuffer, "Content-Transfer-Encoding", "quoted-printable")
	emailBuffer.WriteString("\n")
	err = writeQuotedPrintable(&emailBuffer, email.Html)
	if err != nil {
		return nil, err
	}
	emailBuffer.WriteString(fmt.Sprintf("\n--sub_%s--\n\n--%s--\n", boundary, boundary))
	return emailBuffer.Bytes(), nil
}

func writeHeader(buf *bytes.Buffer, headerName string, value string) {
	buf.WriteString(headerName)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteString("\n")
}

func writeQuotedPrintable(buf *bytes.Buffer, data string) error {
	quotedWritter := quotedprintable.NewWriter(buf)
	_, err := quotedWritter.Write([]byte(data))
	if err != nil {
		return err
	}
	return quotedWritter.Close()
}
