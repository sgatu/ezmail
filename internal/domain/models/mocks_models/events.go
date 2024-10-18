package mocks_models

import (
	"context"
	"time"

	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type EventBusMock struct {
	pushReturn error
	PushCalls  int
}

type ScheduledEventRepositoryMock struct {
	pushReturn       error
	removeNextReturn error
	getNextReturn    struct {
		evt events.Event
		err error
	}
	PushCalls       int
	GetNextCalls    int
	RemoveNextCalls int
}

func MockEventBus() *EventBusMock {
	return &EventBusMock{}
}

func MockScheduledEventRepository() *ScheduledEventRepositoryMock {
	return &ScheduledEventRepositoryMock{
		getNextReturn: struct {
			evt events.Event
			err error
		}{},
	}
}

func (ebm *EventBusMock) Push(ctx context.Context, e events.Event) error {
	ebm.PushCalls++
	return ebm.pushReturn
}

func (ebm *EventBusMock) SetPushReturn(e error) {
	ebm.pushReturn = e
}

func (ser *ScheduledEventRepositoryMock) Push(ctx context.Context, when time.Time, evt events.Event) error {
	ser.PushCalls++
	return ser.pushReturn
}

func (ser *ScheduledEventRepositoryMock) GetNext(ctx context.Context) (events.Event, error) {
	ser.GetNextCalls++
	return ser.getNextReturn.evt, ser.getNextReturn.err
}

func (ser *ScheduledEventRepositoryMock) RemoveNext(ctx context.Context) error {
	ser.RemoveNextCalls++
	return ser.removeNextReturn
}

func (ser *ScheduledEventRepositoryMock) SetPushReturn(e error) {
	ser.pushReturn = e
}

func (ser *ScheduledEventRepositoryMock) SetRemoveNextReturn(e error) {
	ser.removeNextReturn = e
}

func (ser *ScheduledEventRepositoryMock) SetGetNextReturn(evt events.Event, e error) {
	ser.getNextReturn = struct {
		evt events.Event
		err error
	}{
		evt: evt,
		err: e,
	}
}
