package mocks_models

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/email"
)

type EmailRepositoryMock struct {
	returnGetById struct {
		email *email.Email
		err   error
	}
	returnSave   error
	LastSave     *email.Email
	GetByIdCalls int
	SaveCalls    int
}

type TemplateRepositoryMock struct {
	returnSave    error
	LastSave      *email.Template
	returnGetById struct {
		err error
		tpl *email.Template
	}
	returnGetAll struct {
		err  error
		tpls []email.CompactTemplate
	}
	GetByIdCalls int
	SaveCalls    int
	GetAllCalls  int
}

func MockEmailRepository() *EmailRepositoryMock {
	return &EmailRepositoryMock{
		returnGetById: struct {
			email *email.Email
			err   error
		}{},
	}
}

func MockTemplateRepository() *TemplateRepositoryMock {
	return &TemplateRepositoryMock{
		returnGetById: struct {
			err error
			tpl *email.Template
		}{},
		returnGetAll: struct {
			err  error
			tpls []email.CompactTemplate
		}{
			err:  nil,
			tpls: make([]email.CompactTemplate, 0),
		},
	}
}

func (erm *EmailRepositoryMock) GetById(ctx context.Context, id int64) (*email.Email, error) {
	erm.GetByIdCalls++
	return erm.returnGetById.email, erm.returnGetById.err
}

func (erm *EmailRepositoryMock) Save(ctx context.Context, email *email.Email) error {
	erm.SaveCalls++
	erm.LastSave = email
	return erm.returnSave
}

func (erm *EmailRepositoryMock) SetGetByIdReturn(em *email.Email, err error) {
	erm.returnGetById = struct {
		email *email.Email
		err   error
	}{
		email: em,
		err:   err,
	}
}

func (erm *EmailRepositoryMock) SetSaveReturn(err error) {
	erm.returnSave = err
}

func (trm *TemplateRepositoryMock) GetAll(ctx context.Context) ([]email.CompactTemplate, error) {
	trm.GetAllCalls++
	return trm.returnGetAll.tpls, trm.returnGetAll.err
}

func (trm *TemplateRepositoryMock) GetById(ctx context.Context, id int64) (*email.Template, error) {
	trm.GetByIdCalls++
	return trm.returnGetById.tpl, trm.returnGetById.err
}

func (trm *TemplateRepositoryMock) Save(ctx context.Context, tpl *email.Template) error {
	trm.SaveCalls++
	trm.LastSave = tpl
	return trm.returnSave
}

func (trm *TemplateRepositoryMock) SetGetByIdReturn(t *email.Template, err error) {
	trm.returnGetById = struct {
		err error
		tpl *email.Template
	}{
		tpl: t,
		err: err,
	}
}

func (trm *TemplateRepositoryMock) SetGetAllReturn(tpls []email.CompactTemplate, err error) {
	trm.returnGetAll = struct {
		err  error
		tpls []email.CompactTemplate
	}{
		tpls: tpls,
		err:  err,
	}
}

func (trm *TemplateRepositoryMock) SetSaveReturn(err error) {
	trm.returnSave = err
}
