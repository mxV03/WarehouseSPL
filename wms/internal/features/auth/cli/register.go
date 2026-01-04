//go:build auth

package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mxV03/wms/internal/features/auth"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "auth.user.add",
		Usage:       "auth.user.add <username> <role> <password>",
		Group:       "Optional / Auth",
		Description: "Create a user (role: Admin|Worker|ReadOnly)",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("usage: auth.user.add <username> <role> <password>")
			}

			svc := auth.NewAuthService(clictx.AppCtx().Client())

			// take care for declaring first user
			if _, err := svc.RequireRole(ctx, auth.RoleAdmin); err != nil {
				return err
			}

			dto, err := svc.AddUser(ctx, args[0], args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Printf("created user: %s role=%s active=%v\n", dto.Username, dto.Role, dto.Active)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "auth.user.list",
		Usage:       "auth.user.list [limit]",
		Group:       "Optional / Auth",
		Description: "List users.",
		Run: func(ctx context.Context, args []string) error {
			limit := 50
			if len(args) > 1 {
				return fmt.Errorf("usage: auth.user.list [limit]")
			}
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}

			svc := auth.NewAuthService(clictx.AppCtx().Client())
			if _, err := svc.RequireRole(ctx, auth.RoleAdmin); err != nil {
				return err
			}

			us, err := svc.ListUser(ctx, limit)
			if err != nil {
				return err
			}
			if len(us) == 0 {
				fmt.Println("no users")
				return nil
			}
			for _, u := range us {
				fmt.Printf("user: %s role=%s active=%v\n", u.Username, u.Role, u.Active)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "auth.user.disable",
		Usage:       "auth.user.disable <username>",
		Group:       "Optional / Auth",
		Description: "Disable a user.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: auth.user.disable <username>")
			}

			svc := auth.NewAuthService(clictx.AppCtx().Client())
			if _, err := svc.RequireRole(ctx, auth.RoleAdmin); err != nil {
				return err
			}

			if err := svc.DisableUser(ctx, args[0]); err != nil {
				return err
			}
			fmt.Printf("disable user: %s\n", args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "auth.whoami",
		Usage:       "auth.whoami",
		Group:       "Optional / Auth",
		Description: "Show current authenticated user.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("usage: auth.whoami")
			}

			svc := auth.NewAuthService(clictx.AppCtx().Client())
			u, p := auth.CredentialsFromEnv()
			pr, err := svc.Authenticate(ctx, u, p)
			if err != nil {
				return err
			}

			actor := strings.TrimSpace(u)
			if actor == "" {
				actor = "-"
			}
			fmt.Printf("user=%s role=%s\n", actor, pr.Role)
			return nil
		},
	})

}
