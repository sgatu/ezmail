package test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/models/mocks_models"
	"github.com/sgatu/ezmail/internal/domain/services/mocks_services"
	"github.com/sgatu/ezmail/internal/worker/processors"
)

func getNewEmailProcessor() ([]string, *processors.NewEmailProcessor, *processors.RunningContext) {
	rx := processors.RunningContext{
		EventBus:            mocks_models.MockEventBus(),
		EmailStoreService:   mocks_services.MockEmailStoreService(),
		EmailerService:      mocks_services.MockEmailer(),
		ScheduledEventsRepo: mocks_models.MockScheduledEventRepository(),
		Rc:                  nil,
	}
	initMethod := processors.InitNewEmailProcessor()
	msgTypes, proc := initMethod(&rx)
	return msgTypes, proc.(*processors.NewEmailProcessor), &rx
}

func TestInvalidEvent(t *testing.T) {
	_, prc, rx := getNewEmailProcessor()
	res := prc.Process(context.TODO(), &EventInvalid{})
	if res != nil {
		t.Fatal("Unexpected error while processing invalid event for NewEmailProcessor")
	}
	if rx.EmailStoreService.(*mocks_services.EmailerStoreServiceMock).PrepareEmailCalls > 0 {
		t.Fatal("Unexpected call on EmailStoreService for invalid event")
	}
}

func TestEmailPrepareTest(t *testing.T) {
	_, prc, rx := getNewEmailProcessor()
	evt := events.CreateNewEmailEvent(44)
	err := fmt.Errorf("error prepare")
	rx.EmailStoreService.(*mocks_services.EmailerStoreServiceMock).SetPrepareEmailResult(nil, err)
	res := prc.Process(context.TODO(), evt)
	if rx.EmailStoreService.(*mocks_services.EmailerStoreServiceMock).PrepareEmailCalls != 1 {
		t.Fatal("Expected call on EmailStoreService for valid event")
	}
	if !errors.Is(res, err) {
		t.Fatal("Invalid error or no error returned on prepare email for valid event")
	}
}
