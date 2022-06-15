package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/merlion-zone/merlion/x/oracle/types"
)

type RegisterBackingProposalRequest struct {
	BaseReq      rest.BaseReq       `json:"base_req" yaml:"base_req"`
	Title        string             `json:"title" yaml:"title"`
	Description  string             `json:"description" yaml:"description"`
	Deposit      sdk.Coins          `json:"deposit" yaml:"deposit"`
	TargetParams types.TargetParams `json:"target_params" yaml:"target_params"`
}

func RegisterTargetProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req RegisterBackingProposalRequest

			if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
				return
			}

			req.BaseReq = req.BaseReq.Sanitize()
			if !req.BaseReq.ValidateBasic(w) {
				return
			}

			from, err := sdk.AccAddressFromBech32(req.BaseReq.From)
			if rest.CheckBadRequestError(w, err) {
				return
			}

			content := &types.RegisterTargetProposal{
				Title:        req.Title,
				Description:  req.Description,
				TargetParams: req.TargetParams,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, from)
			if rest.CheckBadRequestError(w, err) {
				return
			}

			if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
				return
			}

			tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
		},
	}
}
