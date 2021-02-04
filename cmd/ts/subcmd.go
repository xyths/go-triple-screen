package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xyths/go-triple-screen/cmd/utils"
	"github.com/xyths/go-triple-screen/triple"
	"github.com/xyths/hs"
)

func tripleAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := triple.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	t, err := triple.NewTrader(cfg)
	if err != nil {
		return err
	}
	defer t.Close(ctx.Context)
	if err := t.Init(ctx.Context); err != nil {
		return err
	}

	if err := t.Start(ctx.Context); err != nil {
		return err
	}
	<-ctx.Done()
	if err := t.Stop(ctx.Context); err != nil {
		return err
	}
	return nil
}

func print(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := triple.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	t, err := triple.NewTrader(cfg)
	if err != nil {
		return err
	}
	defer t.Close(ctx.Context)
	if err := t.Init(ctx.Context); err != nil {
		return err
	}

	if err := t.Print(ctx.Context); err != nil {
		return err
	}
	return nil
}
func clear(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := triple.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	t, err := triple.NewTrader(cfg)
	if err != nil {
		return err
	}
	defer t.Close(ctx.Context)
	if err := t.Init(ctx.Context); err != nil {
		return err
	}

	if err := t.Clear(ctx.Context); err != nil {
		return err
	}
	return nil
}
