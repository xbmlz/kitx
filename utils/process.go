package utils

import (
	"context"
	"log"
	"os"

	"github.com/sourcegraph/run"
)

func ExecuteCommand(cmd string) error {
	ctx := context.Background()
	err := run.Cmd(ctx, cmd).Run().Stream(os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}
