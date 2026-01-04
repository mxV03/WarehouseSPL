//go:build audit

package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/auditevent"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "audit.list",
		Usage:       "audit.list [limit]",
		Group:       "Optional / Audit",
		Description: "List latest audit event",
		Run: func(ctx context.Context, args []string) error {
			limit := 50
			if len(args) > 1 {
				return fmt.Errorf("usage: audit.list [limit]")
			}
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}
			if limit <= 0 || limit > 500 {
				limit = 50
			}

			client := clictx.AppCtx().Client()
			evs, err := client.AuditEvent.Query().
				Order(ent.Desc(auditevent.FieldTs)).
				Limit(limit).
				All(ctx)
			if err != nil {
				return fmt.Errorf("audit list: %w", err)
			}
			if len(evs) == 0 {
				fmt.Println("no audit events")
				return nil
			}

			for _, e := range evs {
				fmt.Printf(
					"%s | actor=%s | %s | %s:%s | %s\n",
					e.Ts.Format(time.RFC3339),
					empty(e.Actor),
					e.Action,
					e.Entity,
					empty(e.EntityRef),
					empty(e.Details),
				)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "audit.filter",
		Usage:       "audit.filter [action=<a>] [entity=<e>] [actor=<u>] [limit=<n>]",
		Group:       "Optional / Audit",
		Description: "Filter audit events by key=value args.",
		Run: func(ctx context.Context, args []string) error {
			action := ""
			entity := ""
			actor := ""
			limit := 50

			for _, a := range args {
				a = strings.TrimSpace(a)
				if a == "" {
					continue
				}
				parts := strings.SplitN(a, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid arg %q (expected key=value)", a)
				}
				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])

				switch k {
				case "action":
					action = v
				case "entity":
					entity = v
				case "actor":
					actor = v
				case "limit":
					n, err := strconv.Atoi(v)
					if err != nil {
						return fmt.Errorf("limit must be an integer")
					}
					limit = n
				default:
					return fmt.Errorf("unknown key %q (allowed: action, entity, actor, limit)", k)
				}
			}

			if limit <= 0 || limit > 500 {
				limit = 50
			}

			client := clictx.AppCtx().Client()
			q := client.AuditEvent.Query()
			if action != "" {
				q = q.Where(auditevent.ActionEQ(action))
			}
			if entity != "" {
				q = q.Where(auditevent.EntityEQ(entity))
			}
			if actor != "" {
				q = q.Where(auditevent.ActorEQ(actor))
			}

			evs, err := q.Order(ent.Desc(auditevent.FieldTs)).Limit(limit).All(ctx)
			if err != nil {
				return fmt.Errorf("audit filter: %w", err)
			}
			if len(evs) == 0 {
				fmt.Println("no audit events")
				return nil
			}

			for _, e := range evs {
				fmt.Printf(
					"%s | actor=%s | %s | %s:%s | %s\n",
					e.Ts.Format(time.RFC3339),
					empty(e.Actor),
					e.Action,
					e.Entity,
					empty(e.EntityRef),
					empty(e.Details),
				)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "audit.tail",
		Usage:       "audit.tail",
		Group:       "Optional / Audit",
		Description: "Show last 20 audit events.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("usage: audit.tail")
			}
			client := clictx.AppCtx().Client()
			evs, err := client.AuditEvent.Query().
				Order(ent.Desc(auditevent.FieldTs)).
				Limit(20).
				All(ctx)
			if err != nil {
				return fmt.Errorf("audit tail: %w", err)
			}
			if len(evs) == 0 {
				fmt.Println("no audit events")
				return nil
			}
			for _, e := range evs {
				fmt.Printf("%s | %s | %s:%s | actor=%s\n",
					e.Ts.Format(time.RFC3339),
					e.Action,
					e.Entity,
					empty(e.EntityRef),
					empty(e.Actor),
				)
			}
			return nil
		},
	})
}

func empty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
