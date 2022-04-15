package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

type RegisterBackingProposalRequest struct {
	BaseReq     rest.BaseReq            `json:"base_req" yaml:"base_req"`
	Title       string                  `json:"title" yaml:"title"`
	Description string                  `json:"description" yaml:"description"`
	Deposit     sdk.Coins               `json:"deposit" yaml:"deposit"`
	RiskParams  types.BackingRiskParams `json:"risk_params" yaml:"risk_params"`
}

type RegisterCollateralProposalRequest struct {
	BaseReq     rest.BaseReq               `json:"base_req" yaml:"base_req"`
	Title       string                     `json:"title" yaml:"title"`
	Description string                     `json:"description" yaml:"description"`
	Deposit     sdk.Coins                  `json:"deposit" yaml:"deposit"`
	RiskParams  types.CollateralRiskParams `json:"risk_params" yaml:"risk_params"`
}

type SetBackingProposalRequest struct {
	BaseReq     rest.BaseReq            `json:"base_req" yaml:"base_req"`
	Title       string                  `json:"title" yaml:"title"`
	Description string                  `json:"description" yaml:"description"`
	Deposit     sdk.Coins               `json:"deposit" yaml:"deposit"`
	RiskParams  types.BackingRiskParams `json:"risk_params" yaml:"risk_params"`
}

type SetCollateralProposalRequest struct {
	BaseReq     rest.BaseReq               `json:"base_req" yaml:"base_req"`
	Title       string                     `json:"title" yaml:"title"`
	Description string                     `json:"description" yaml:"description"`
	Deposit     sdk.Coins                  `json:"deposit" yaml:"deposit"`
	RiskParams  types.CollateralRiskParams `json:"risk_params" yaml:"risk_params"`
}

type BatchSetBackingProposalRequest struct {
	BaseReq     rest.BaseReq              `json:"base_req" yaml:"base_req"`
	Title       string                    `json:"title" yaml:"title"`
	Description string                    `json:"description" yaml:"description"`
	Deposit     sdk.Coins                 `json:"deposit" yaml:"deposit"`
	RiskParams  []types.BackingRiskParams `json:"risk_params" yaml:"risk_params"`
}

type BatchSetCollateralProposalRequest struct {
	BaseReq     rest.BaseReq                 `json:"base_req" yaml:"base_req"`
	Title       string                       `json:"title" yaml:"title"`
	Description string                       `json:"description" yaml:"description"`
	Deposit     sdk.Coins                    `json:"deposit" yaml:"deposit"`
	RiskParams  []types.CollateralRiskParams `json:"risk_params" yaml:"risk_params"`
}

func RegisterBackingProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
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

			content := &types.RegisterBackingProposal{
				Title:       req.Title,
				Description: req.Description,
				RiskParams:  req.RiskParams,
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

func RegisterCollateralProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req RegisterCollateralProposalRequest

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

			content := &types.RegisterCollateralProposal{
				Title:       req.Title,
				Description: req.Description,
				RiskParams:  req.RiskParams,
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

func SetBackingProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req SetBackingProposalRequest

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

			content := &types.SetBackingRiskParamsProposal{
				Title:       req.Title,
				Description: req.Description,
				RiskParams:  req.RiskParams,
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

func SetCollateralProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req SetCollateralProposalRequest

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

			content := &types.SetCollateralRiskParamsProposal{
				Title:       req.Title,
				Description: req.Description,
				RiskParams:  req.RiskParams,
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

func BatchSetBackingProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req BatchSetBackingProposalRequest

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

			content := &types.BatchSetBackingRiskParamsProposal{
				Title:       req.Title,
				Description: req.Description,
				RiskParams:  req.RiskParams,
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

func BatchSetCollateralProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req BatchSetCollateralProposalRequest

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

			content := &types.BatchSetCollateralRiskParamsProposal{
				Title:       req.Title,
				Description: req.Description,
				RiskParams:  req.RiskParams,
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
