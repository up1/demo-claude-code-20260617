package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api/internal/core/domain"
	"api/internal/core/ports"
	"api/internal/core/service"
)

type mockInboxRepo struct{ mock.Mock }

func (m *mockInboxRepo) List(ctx context.Context, filter ports.InboxFilter) ([]domain.InboxMessage, int64, error) {
	args := m.Called(ctx, filter)
	var items []domain.InboxMessage
	if args.Get(0) != nil {
		items = args.Get(0).([]domain.InboxMessage)
	}
	return items, args.Get(1).(int64), args.Error(2)
}

func TestListMessages_DefaultsApplied(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)
	ctx := context.Background()

	repo.On("List", ctx, ports.InboxFilter{Offset: 0, Limit: 20}).
		Return([]domain.InboxMessage{}, int64(0), nil)

	res, err := svc.ListMessages(ctx, ports.ListInboxInput{})

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Page)
	assert.Equal(t, 20, res.PageSize)
	assert.Equal(t, int64(0), res.TotalItems)
	assert.Equal(t, 0, res.TotalPages)
	repo.AssertExpectations(t)
}

func TestListMessages_ClampsPageSizeToMax(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)
	ctx := context.Background()

	repo.On("List", ctx, ports.InboxFilter{Offset: 0, Limit: 100}).
		Return([]domain.InboxMessage{}, int64(0), nil)

	res, err := svc.ListMessages(ctx, ports.ListInboxInput{Page: 1, PageSize: 500})

	assert.NoError(t, err)
	assert.Equal(t, 100, res.PageSize)
	repo.AssertExpectations(t)
}

func TestListMessages_OffsetMathAndPassThrough(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)
	ctx := context.Background()

	from := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC)

	expected := ports.InboxFilter{
		Channel: "line",
		Status:  "pending",
		Search:  "tracking",
		From:    &from,
		To:      &to,
		Offset:  10, // (page 2 - 1) * 10
		Limit:   10,
	}
	repo.On("List", ctx, expected).Return([]domain.InboxMessage{}, int64(0), nil)

	_, err := svc.ListMessages(ctx, ports.ListInboxInput{
		Channel:  "line",
		Status:   "pending",
		Search:   "tracking",
		From:     &from,
		To:       &to,
		Page:     2,
		PageSize: 10,
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListMessages_TotalPagesCeil(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)
	ctx := context.Background()

	repo.On("List", ctx, mock.Anything).
		Return([]domain.InboxMessage{}, int64(128), nil)

	res, err := svc.ListMessages(ctx, ports.ListInboxInput{Page: 1, PageSize: 20})

	assert.NoError(t, err)
	assert.Equal(t, int64(128), res.TotalItems)
	assert.Equal(t, 7, res.TotalPages) // ceil(128/20)
}

func TestListMessages_InvalidChannel(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)

	res, err := svc.ListMessages(context.Background(), ports.ListInboxInput{Channel: "tiktok"})

	assert.Nil(t, res)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "List")
}

func TestListMessages_InvalidStatus(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)

	res, err := svc.ListMessages(context.Background(), ports.ListInboxInput{Status: "archived"})

	assert.Nil(t, res)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "List")
}

func TestListMessages_InvalidDateRange(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := service.NewInboxService(repo)

	from := time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)

	res, err := svc.ListMessages(context.Background(), ports.ListInboxInput{From: &from, To: &to})

	assert.Nil(t, res)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "List")
}
