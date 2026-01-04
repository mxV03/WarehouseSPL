//go:build !audit

package auditlog

import "context"

func Log(ctx context.Context, action, entity, entityRef, details string)              {}
func Logf(ctx context.Context, action, entity, entityRef, format string, args ...any) {}
