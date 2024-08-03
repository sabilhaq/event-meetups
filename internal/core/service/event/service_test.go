package event_test

/*
	The purpose of testing the Service component is to ensure it has correct
	implementation of business logic.

	The common pitfall when creating test for Service component is we tend to use
	concrete implementation for the dependency components (e.g actual EventStorage
	for MySQL). Not only this will increase the test complexity but also it will
	increase the possibility of getting false test result. The reason is simply
	because service such as MySQL has its own constraints & has much higher chance
	of failing rather than its mock counterpart (e.g network failure).

	So to avoid this pitfall, our first go to choice is to use mock implementation
	for the dependency when testing the Service component. This way we can control
	more the behavior of the dependency components to fit our test scenarios.
*/

import (
	"context"
	"errors"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/event"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	eventStorage := newMockEventStorage(nil)

	// define test cases
	testCases := []struct {
		Name    string
		Config  event.ServiceConfig
		IsError bool
	}{
		{
			Name: "Test Missing Event Storage",
			Config: event.ServiceConfig{
				EventStorage: nil,
			},
			IsError: true,
		},
		{
			Name: "Test Valid Config",
			Config: event.ServiceConfig{
				EventStorage: eventStorage,
			},
			IsError: false,
		},
	}
	// execute test cases
	for _, testcase := range testCases {
		t.Run(testcase.Name, func(t *testing.T) {
			_, err := event.NewService(testcase.Config)
			require.Equal(t, testcase.IsError, (err != nil), "unexpected error")
		})
	}
}

func TestServiceGetEvents(t *testing.T) {
	// initialize new service
	output := newService()

	// get available events
	retEvents, err := output.Service.GetEvents(context.Background())
	require.NoError(t, err, "unexpected error")

	// check returned events
	require.ElementsMatch(t, output.Events, retEvents, "mismatch events")

	// set error on get available events
	output.EventStorage.SetRetErrOnGetEvents(true)

	// get available events, should return error
	_, err = output.Service.GetEvents(context.Background())
	require.Error(t, err, "expected error")
}

func newService() *newServiceOutput {
	// generate events
	events := []entity.Event{
		*(testutil.NewTestEvent()),
		*(testutil.NewTestEvent()),
		*(testutil.NewTestEvent()),
		*(testutil.NewTestEvent()),
		*(testutil.NewTestEvent()),
	}

	// initialize dependencies
	eventStorage := newMockEventStorage(events)

	// initialize service
	cfg := event.ServiceConfig{
		EventStorage: eventStorage,
	}
	svc, _ := event.NewService(cfg)

	return &newServiceOutput{
		Service:      svc,
		EventStorage: eventStorage,
		Events:       events,
	}
}

type newServiceOutput struct {
	Service      event.Service
	EventStorage *mockEventStorage
	Events       []entity.Event
}

type mockEventStorage struct {
	data              map[int]entity.Event
	retErrOnGetEvent  bool
	retErrOnGetEvents bool
}

func (gs *mockEventStorage) SetRetErrOnGetEvent(retErr bool) {
	gs.retErrOnGetEvent = retErr
}

func (gs *mockEventStorage) SetRetErrOnGetEvents(retErr bool) {
	gs.retErrOnGetEvents = retErr
}

func (gs *mockEventStorage) GetEvents(ctx context.Context) ([]entity.Event, error) {
	if gs.retErrOnGetEvents {
		return nil, ErrIntentionalError
	}
	var events []entity.Event
	for _, v := range gs.data {
		events = append(events, v)
	}
	return events, nil
}

func newMockEventStorage(events []entity.Event) *mockEventStorage {
	data := map[int]entity.Event{}
	for _, event := range events {
		data[event.ID] = event
	}

	return &mockEventStorage{data: data}
}

var ErrIntentionalError = errors.New("intentional error")
