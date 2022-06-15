package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/merlion-zone/merlion/x/oracle/client/cli"
	"github.com/merlion-zone/merlion/x/oracle/client/rest"
)

var (
	RegisterTargetProposalHandler = govclient.NewProposalHandler(cli.NewRegisterTargetProposalCmd, rest.RegisterTargetProposalRESTHandler)
)
