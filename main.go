package main

import (
	"context"
	"fmt"
	timestampcmds "github.com/ProtoconNet/mitum-timestamp/cmds"
	"os"

	"github.com/ProtoconNet/mitum-minic/cmds"
	launchcmd "github.com/ProtoconNet/mitum2/launch/cmd"

	credentialcmds "github.com/ProtoconNet/mitum-credential/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	nftcmds "github.com/ProtoconNet/mitum-nft/v2/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
)

var (
	Version   = "v0.0.1"
	BuildTime = "-"
	GitBranch = "master"
	GitCommit = "-"
	version   util.Version
)

//revive:disable:nested-structs
var CLI struct { //nolint:govet //...
	launch.BaseFlags
	Init      currencycmds.INITCommand `cmd:"" help:"init node"`
	Run       cmds.RunCommand          `cmd:"" help:"run node"`
	Storage   launchcmd.Storage        `cmd:""`
	Operation struct {
		Currency   currencycmds.CurrencyCommand     `cmd:"" help:"currency operation"`
		Suffrage   currencycmds.SuffrageCommand     `cmd:"" help:"suffrage operation"`
		NFT        nftcmds.NFTCommand               `cmd:"" help:"nft operation"`
		Credential credentialcmds.CredentialCommand `cmd:"" help:"credential operation"`
		Timestamp  timestampcmds.TimestampCommand   `cmd:"" help:"timestamp operation"`
	} `cmd:"" help:"create operation"`
	Network struct {
		Client cmds.NetworkClientCommand `cmd:"" help:"network client"`
	} `cmd:"" help:"network"`
	Key struct {
		New     currencycmds.KeyNewCommand     `cmd:"" help:"generate new key"`
		Address currencycmds.KeyAddressCommand `cmd:"" help:"generate address from key"`
		Load    launchcmd.KeyLoadCommand       `cmd:"" help:"load key"`
		Sign    launchcmd.KeySignCommand       `cmd:"" help:"sign"`
	} `cmd:"" help:"key"`
	Handover launchcmd.HandoverCommands `cmd:""`
	Version  struct{}                   `cmd:"" help:"version"`
}

var flagDefaults = kong.Vars{
	"log_out":                           "stderr",
	"log_format":                        "terminal",
	"log_level":                         "debug",
	"log_force_color":                   "false",
	"design_uri":                        launch.DefaultDesignURI,
	"create_account_threshold":          "100",
	"create_contract_account_threshold": "100",
	"safe_threshold":                    base.SafeThreshold.String(),
	"network_id":                        "mitum",
}

func main() {
	kctx := kong.Parse(&CLI, flagDefaults)

	if err := checkVersion(); err != nil {
		kctx.FatalIfErrorf(err)
	}

	if kctx.Command() == "version" {
		showVersion()
		return
	}

	pctx := context.Background()
	pctx = context.WithValue(pctx, launch.VersionContextKey, version)
	pctx = context.WithValue(pctx, launch.FlagsContextKey, CLI.BaseFlags)
	pctx = context.WithValue(pctx, launch.KongContextContextKey, kctx)

	pss := launch.DefaultMainPS()

	switch i, err := pss.Run(pctx); {
	case err != nil:
		kctx.FatalIfErrorf(err)
	default:
		pctx = i
		kctx = kong.Parse(&CLI, kong.BindTo(pctx, (*context.Context)(nil)), flagDefaults)
	}

	var log *logging.Logging
	if err := util.LoadFromContextOK(pctx, launch.LoggingContextKey, &log); err != nil {
		kctx.FatalIfErrorf(err)
	}

	log.Log().Debug().Interface("flags", os.Args).Msg("flags")
	log.Log().Debug().Interface("main_process", pss.Verbose()).Msg("processed")

	if err := func() error {
		defer log.Log().Debug().Msg("stopped")
		return errors.WithStack(kctx.Run(pctx))
	}(); err != nil {
		log.Log().Error().Err(err).Msg("stopped by error")
		kctx.FatalIfErrorf(err)
	}
}

func checkVersion() error {
	if len(Version) < 1 {
		return errors.Errorf("empty version")
	}

	v, err := util.ParseVersion(Version)
	if err != nil {
		return err
	}

	if err := v.IsValid(nil); err != nil {
		return err
	}

	version = v

	return nil
}

func showVersion() {
	_, _ = fmt.Fprintf(os.Stdout, `version: %s
 branch: %s
 commit: %s
  build: %s
`, version, GitBranch, GitCommit, BuildTime)
}
