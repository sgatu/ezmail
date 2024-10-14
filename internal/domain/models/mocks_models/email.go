package mocks_models

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/email"
)

type emailRepositoryMock struct {
	returnGetById struct {
		email *email.Email
		err   error
	}
	returnSave   error
	GetByIdCalls int
	SaveCalls    int
}

type templateRepositoryMock struct {
	returnSave    error
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

func MockEmailRepository() *emailRepositoryMock {
	return &emailRepositoryMock{
		returnGetById: struct {
			email *email.Email
			err   error
		}{},
	}
}

func (erm *emailRepositoryMock) GetById(ctx context.Context, id int64) (*email.Email, error) {
	erm.GetByIdCalls++
	return erm.returnGetById.email, erm.returnGetById.err
}

func (erm *emailRepositoryMock) Save(ctx context.Context, email *email.Email) error {
	erm.SaveCalls++
	return erm.returnSave
}

func (erm *emailRepositoryMock) SetGetByIdReturn(em *email.Email, err error) {
	erm.returnGetById = struct {
		email *email.Email
		err   error
	}{
		email: em,
		err:   err,
	}
}

func (erm *emailRepositoryMock) SetSaveReturn(err error) {
	erm.returnSave = err
}

func (trm *templateRepositoryMock) GetAll(ctx context.Context) ([]email.CompactTemplate, error) {
	trm.GetAllCalls++
	return trm.returnGetAll.tpls, trm.returnGetAll.err
}

func (trm *templateRepositoryMock) GetById(ctx context.Context, id int64) (*email.Template, error) {
	trm.GetByIdCalls++
	return trm.returnGetById.tpl, trm.returnGetById.err
}

func (trm *templateRepositoryMock) Save(ctx context.Context, tpl *email.Template) error {
	trm.SaveCalls++
	return trm.returnSave
}

func (trm *templateRepositoryMock) SetGetByIdReturn(t *email.Template, err error) {
	trm.returnGetById = struct {
		err error
		tpl *email.Template
	}{
		tpl: t,
		err: err,
	}
}

func (trm *templateRepositoryMock) SetGetAllReturn(tpls []email.CompactTemplate, err error) {
	trm.returnGetAll = struct {
		err  error
		tpls []email.CompactTemplate
	}{
		tpls: tpls,
		err:  err,
	}
}

func (trm *templateRepositoryMock) SetSaveReturn(err error) {
	trm.returnSave = err
}
