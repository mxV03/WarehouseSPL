package cli

import (
	"context"
	"os"

	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func Run() error {
	ctx := context.Background()
	return registry.Dispatch(ctx, os.Args[1:])
}
