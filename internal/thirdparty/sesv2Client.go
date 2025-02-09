package thirdparty

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type SESClient interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
	CreateEmailIdentity(ctx context.Context, params *sesv2.CreateEmailIdentityInput, optFns ...func(*sesv2.Options)) (*sesv2.CreateEmailIdentityOutput, error)
	PutEmailIdentityMailFromAttributes(
		ctx context.Context,
		params *sesv2.PutEmailIdentityMailFromAttributesInput,
		optFns ...func(*sesv2.Options),
	) (*sesv2.PutEmailIdentityMailFromAttributesOutput, error)
	DeleteEmailIdentity(ctx context.Context, params *sesv2.DeleteEmailIdentityInput, optFns ...func(*sesv2.Options)) (*sesv2.DeleteEmailIdentityOutput, error)
	GetEmailIdentity(ctx context.Context, params *sesv2.GetEmailIdentityInput, optFns ...func(*sesv2.Options)) (*sesv2.GetEmailIdentityOutput, error)
}

type AWSSesV2Client struct {
	*sesv2.Client
}
