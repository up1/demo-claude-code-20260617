package service

import (
	"context"
	"fmt"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type inboxService struct {
	repo ports.InboxRepository
}

// NewInboxService wires the inbox use-case logic to its repository port.
func NewInboxService(repo ports.InboxRepository) ports.InboxService {
	return &inboxService{repo: repo}
}

// ListMessages validates and normalizes the request, queries the repository,
// and maps (items, total) into a paginated result.
func (s *inboxService) ListMessages(ctx context.Context, req ports.ListInboxInput) (*ports.ListInboxResult, error) {
	if req.Channel != "" && !domain.ValidChannel(req.Channel) {
		return nil, fmt.Errorf("%w: channel must be one of facebook, line, instagram", domain.ErrValidation)
	}
	if req.Status != "" && !domain.ValidStatus(req.Status) {
		return nil, fmt.Errorf("%w: status must be one of pending, replied", domain.ErrValidation)
	}
	if req.From != nil && req.To != nil && req.From.After(*req.To) {
		return nil, fmt.Errorf("%w: from must be less than or equal to to", domain.ErrValidation)
	}

	page := req.Page
	if page < 1 {
		page = defaultPage
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	filter := ports.InboxFilter{
		Channel: req.Channel,
		Status:  req.Status,
		Search:  req.Search,
		From:    req.From,
		To:      req.To,
		Offset:  (page - 1) * pageSize,
		Limit:   pageSize,
	}

	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &ports.ListInboxResult{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}
