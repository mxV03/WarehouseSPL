//go:build audit

package auditlog

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mxV03/wms/internal/features/audit"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
)

func Log(ctx context.Context, action, entity, entityRef, details string) {
	actor := strings.TrimSpace(os.Getenv("WMS_ACTOR"))
	if actor == "" {
		actor = "system"
	}

	client := clictx.AppCtx().Client()

	_ = audit.NewAuditService(client).Log(ctx, actor, action, entity, entityRef, details)
}

func Logf(ctx context.Context, action, entity, entityRef, format string, args ...any) {
	Log(ctx, action, entity, entityRef, fmt.Sprintf(format, args...))
}
