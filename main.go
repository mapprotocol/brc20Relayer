package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/mapprotocol/brc20Relayer/config"
	"github.com/mapprotocol/brc20Relayer/resource/db"
	"github.com/mapprotocol/brc20Relayer/startstore"
)

var (
	app     = cli.NewApp()
	Version = "0.0.1"
)

var cliFlags = []cli.Flag{
	config.VerbosityFlag,
	config.ConfigFileFlag,
	config.BlockStorePathFlag,
}

func init() {
	app.Action = run
	app.Version = Version
	app.Name = "brc20Relayer"

	app.Flags = append(app.Flags, cliFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error("run failed", "error", err.Error())
		os.Exit(1)
	}
}

func Initialize(cfg *config.Config) error {
	return nil
}

func run(ctx *cli.Context) error {
	if err := setup(ctx); err != nil {
		return err
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	db.Init(cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.Name)

	start, err := startstore.ReadLatestStart(ctx.String(config.BlockStorePathFlag.Name))
	if err != nil {
		return err
	}
	cfg.StartNumber = start

	go func() {
		// poll
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		log.Info("received the exit signal, ready to exit...")
	}

	exit()
	return nil
}

func setup(ctx *cli.Context) error {
	if err := startLogger(ctx); err != nil {
		return err
	}
	log.Info("setup scaner...")
	return nil
}

func exit() {
	log.Info("exit scaner...")
}

func startLogger(ctx *cli.Context) error {
	var lvl log.Lvl
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(true)))

	if lvlToInt, err := strconv.Atoi(ctx.String(config.VerbosityFlag.Name)); err == nil {
		lvl = log.Lvl(lvlToInt)
	} else if lvl, err = log.LvlFromString(ctx.String(config.VerbosityFlag.Name)); err != nil {
		return err
	}
	glogger.Verbosity(lvl)
	log.Root().SetHandler(glogger)

	return nil
}
