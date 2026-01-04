//go:build audit

package audit

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/wms/ent"
)

type AuditService struct {
	client *ent.Client
}

func NewAuditService(client *ent.Client) *AuditService {
	return &AuditService{
		client: client,
	}
}

func (s *AuditService) Log(ctx context.Context, actor, action, entity, entityRef, details string) error {
	actor = strings.TrimSpace(actor)
	action = strings.TrimSpace(action)
	entity = strings.TrimSpace(entity)
	entityRef = strings.TrimSpace(entityRef)
	details = strings.TrimSpace(details)

	if actor == "" {
		actor = "system"
	}
	if action == "" || entity == "" {
		return fmt.Errorf("audit: action and entity must not be empty")
	}

	_, err := s.client.AuditEvent.Create().
		SetActor(actor).
		SetAction(action).
		SetEntity(entity).
		SetEntityRef(entityRef).
		SetDetails(details).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}
