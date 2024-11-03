package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type MockSesV2Client struct {
	sendResponse struct {
		out *sesv2.SendEmailOutput
		err error
	}
	createIdentityResponse struct {
		out *sesv2.CreateEmailIdentityOutput
		err error
	}
	putEmailIdentityMailFromAttributesResponse struct {
		out *sesv2.PutEmailIdentityMailFromAttributesOutput
		err error
	}
	deleteIdentityResponse struct {
		out *sesv2.DeleteEmailIdentityOutput
		err error
	}
	sendLastInput struct {
		Ctx            context.Context
		SendEmailInput *sesv2.SendEmailInput
		Opts           []func(*sesv2.Options)
	}
	createIdentityLastInput struct {
		Ctx                 context.Context
		CreateIdentityInput *sesv2.CreateEmailIdentityInput
		Opts                []func(*sesv2.Options)
	}
	putEmailIdentityMailFromAttributesInput struct {
		Ctx                   context.Context
		PutEmailIdentityInput *sesv2.PutEmailIdentityMailFromAttributesInput
		Opts                  []func(*sesv2.Options)
	}
	deleteIdentityInput struct {
		Ctx                 context.Context
		DeleteIdentityInput *sesv2.DeleteEmailIdentityInput
		Opts                []func(*sesv2.Options)
	}

	SendCalls                               int
	CreateIdentityCalls                     int
	PutIdentityEmailMailFromAttributesCalls int
	DeleteIdentityCalls                     int
}

func (ms *MockSesV2Client) SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	ms.sendLastInput = struct {
		Ctx            context.Context
		SendEmailInput *sesv2.SendEmailInput
		Opts           []func(*sesv2.Options)
	}{
		Ctx:            ctx,
		SendEmailInput: params,
		Opts:           optFns,
	}
	ms.SendCalls++
	return ms.sendResponse.out, ms.sendResponse.err
}

func (ms *MockSesV2Client) SetSendEmailResponse(Out *sesv2.SendEmailOutput, err error) {
	ms.sendResponse = struct {
		out *sesv2.SendEmailOutput
		err error
	}{
		out: Out,
		err: err,
	}
}

func (ms *MockSesV2Client) GetLastSendParams() struct {
	Ctx            context.Context
	SendEmailInput *sesv2.SendEmailInput
	Opts           []func(*sesv2.Options)
} {
	return ms.sendLastInput
}

func (ms *MockSesV2Client) CreateEmailIdentity(ctx context.Context, params *sesv2.CreateEmailIdentityInput, optFns ...func(*sesv2.Options)) (*sesv2.CreateEmailIdentityOutput, error) {
	ms.createIdentityLastInput = struct {
		Ctx                 context.Context
		CreateIdentityInput *sesv2.CreateEmailIdentityInput
		Opts                []func(*sesv2.Options)
	}{
		Ctx:                 ctx,
		CreateIdentityInput: params,
		Opts:                optFns,
	}
	ms.CreateIdentityCalls++
	return ms.createIdentityResponse.out, ms.createIdentityResponse.err
}

func (ms *MockSesV2Client) SetCreateEmailIdentityResponse(out *sesv2.CreateEmailIdentityOutput, err error) {
	ms.createIdentityResponse = struct {
		out *sesv2.CreateEmailIdentityOutput
		err error
	}{
		out: out,
		err: err,
	}
}

func (ms *MockSesV2Client) PutEmailIdentityMailFromAttributes(
	ctx context.Context,
	params *sesv2.PutEmailIdentityMailFromAttributesInput,
	optFns ...func(*sesv2.Options),
) (*sesv2.PutEmailIdentityMailFromAttributesOutput, error) {
	ms.PutIdentityEmailMailFromAttributesCalls++
	ms.putEmailIdentityMailFromAttributesInput = struct {
		Ctx                   context.Context
		PutEmailIdentityInput *sesv2.PutEmailIdentityMailFromAttributesInput
		Opts                  []func(*sesv2.Options)
	}{
		Opts:                  optFns,
		Ctx:                   ctx,
		PutEmailIdentityInput: params,
	}
	return ms.putEmailIdentityMailFromAttributesResponse.out, ms.putEmailIdentityMailFromAttributesResponse.err
}

func (ms *MockSesV2Client) SetPutEmailIdentityMailFromAttributesResponse(
	out *sesv2.PutEmailIdentityMailFromAttributesOutput,
	err error,
) {
	ms.putEmailIdentityMailFromAttributesResponse = struct {
		out *sesv2.PutEmailIdentityMailFromAttributesOutput
		err error
	}{
		out: out,
		err: err,
	}
}

func (ms *MockSesV2Client) GetLastPutEmailIdentityMailFromAttributesParams() struct {
	Ctx                   context.Context
	PutEmailIdentityInput *sesv2.PutEmailIdentityMailFromAttributesInput
	Opts                  []func(*sesv2.Options)
} {
	return ms.putEmailIdentityMailFromAttributesInput
}

func (ms *MockSesV2Client) DeleteEmailIdentity(ctx context.Context, params *sesv2.DeleteEmailIdentityInput, optFns ...func(*sesv2.Options)) (*sesv2.DeleteEmailIdentityOutput, error) {
	ms.DeleteIdentityCalls++
	ms.deleteIdentityInput = struct {
		Ctx                 context.Context
		DeleteIdentityInput *sesv2.DeleteEmailIdentityInput
		Opts                []func(*sesv2.Options)
	}{
		Ctx:                 ctx,
		DeleteIdentityInput: params,
		Opts:                optFns,
	}
	return ms.deleteIdentityResponse.out, ms.deleteIdentityResponse.err
}

func (ms *MockSesV2Client) GetLastDeleteEmailIdentityInput() struct {
	Ctx                 context.Context
	DeleteIdentityInput *sesv2.DeleteEmailIdentityInput
	Opts                []func(*sesv2.Options)
} {
	return ms.deleteIdentityInput
}

func (ms *MockSesV2Client) SetDeleteEmailIdentityOutput(out *sesv2.DeleteEmailIdentityOutput, err error) {
	ms.deleteIdentityResponse = struct {
		out *sesv2.DeleteEmailIdentityOutput
		err error
	}{
		out: out,
		err: err,
	}
}
