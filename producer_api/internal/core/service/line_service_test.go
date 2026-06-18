package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"producer/internal/core/domain"
	"producer/internal/core/service"
)

type mockPublisher struct{ mock.Mock }

func (m *mockPublisher) Publish(ctx context.Context, queue string, body []byte) error {
	args := m.Called(ctx, queue, body)
	return args.Error(0)
}

type mockDedup struct{ mock.Mock }

func (m *mockDedup) IsDuplicate(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func validMessage() domain.LineMessage {
	return domain.LineMessage{
		To:       "U123",
		Messages: []domain.Message{{Type: domain.MessageTypeText, Text: "Hello, world!"}},
	}
}

func TestSendMessage_HappyPath(t *testing.T) {
	pub := new(mockPublisher)
	dedup := new(mockDedup)
	raw := []byte(`{"to":"U123"}`)

	dedup.On("IsDuplicate", mock.Anything, mock.Anything).Return(false, nil)

	var captured []byte
	pub.On("Publish", mock.Anything, service.QueueName, mock.Anything).
		Run(func(args mock.Arguments) { captured = args.Get(2).([]byte) }).
		Return(nil)

	svc := service.NewLineService(pub, dedup)
	res, err := svc.SendMessage(context.Background(), raw, validMessage())

	require.NoError(t, err)
	assert.False(t, res.Duplicate)
	assert.NotEmpty(t, res.MessageID)

	// Envelope shape matches the spec.
	var env map[string]any
	require.NoError(t, json.Unmarshal(captured, &env))
	assert.Equal(t, service.QueueName, env["name"])
	assert.NotNil(t, env["message"])

	pub.AssertExpectations(t)
	dedup.AssertExpectations(t)
}

func TestSendMessage_Duplicate_ShortCircuits(t *testing.T) {
	pub := new(mockPublisher)
	dedup := new(mockDedup)

	dedup.On("IsDuplicate", mock.Anything, mock.Anything).Return(true, nil)

	svc := service.NewLineService(pub, dedup)
	res, err := svc.SendMessage(context.Background(), []byte(`{}`), validMessage())

	require.NoError(t, err)
	assert.True(t, res.Duplicate)
	assert.Empty(t, res.MessageID)
	pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendMessage_ValidationError_NoPublish(t *testing.T) {
	pub := new(mockPublisher)
	dedup := new(mockDedup)

	dedup.On("IsDuplicate", mock.Anything, mock.Anything).Return(false, nil)

	svc := service.NewLineService(pub, dedup)
	invalid := domain.LineMessage{To: "", Messages: nil}
	res, err := svc.SendMessage(context.Background(), []byte(`{}`), invalid)

	assert.ErrorIs(t, err, domain.ErrValidation)
	assert.False(t, res.Duplicate)
	pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendMessage_PublishError(t *testing.T) {
	pub := new(mockPublisher)
	dedup := new(mockDedup)

	dedup.On("IsDuplicate", mock.Anything, mock.Anything).Return(false, nil)
	pub.On("Publish", mock.Anything, service.QueueName, mock.Anything).
		Return(errors.New("broker down"))

	svc := service.NewLineService(pub, dedup)
	_, err := svc.SendMessage(context.Background(), []byte(`{}`), validMessage())

	require.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrValidation)
}

func TestSendMessage_DedupError(t *testing.T) {
	pub := new(mockPublisher)
	dedup := new(mockDedup)

	dedup.On("IsDuplicate", mock.Anything, mock.Anything).Return(false, errors.New("redis down"))

	svc := service.NewLineService(pub, dedup)
	_, err := svc.SendMessage(context.Background(), []byte(`{}`), validMessage())

	require.Error(t, err)
	pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}
