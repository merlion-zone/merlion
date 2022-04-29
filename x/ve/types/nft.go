package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

var VeNftClass = nft.Class{
	Id:          "veNFT",
	Name:        "veNFT",
	Symbol:      "veNFT",
	Description: "Merlion locks, can be used to boost gauge yields, vote on token emission, and receive bribes",
	Uri:         "",
}

func VeID(idNumber uint64) string {
	return fmt.Sprintf("ve-%d", idNumber)
}

func Uint64FromVeID(veID string) uint64 {
	splits := strings.Split(veID, "-")
	if len(splits) != 2 || splits[0] != "ve" {
		return EmptyVeID
	}
	id, err := strconv.ParseUint(splits[1], 10, 64)
	if err != nil || id == EmptyVeID {
		return EmptyVeID
	}
	return id
}

func VeNftUri(nftID string, balance sdk.Int, lockedEnd time.Time, value sdk.Int) string {
	output := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMinYMin meet" viewBox="0 0 350 350"><style>.base { fill: white; font-family: serif; font-size: 14px; }</style><rect width="100%%" height="100%%" fill="black" /><text x="10" y="20" class="base">token %s</text><text x="10" y="40" class="base">balanceOf %s</text><text x="10" y="60" class="base">locked_end %d</text><text x="10" y="80" class="base">value %s</text></svg>`, nftID, balance, lockedEnd.Unix(), value)

	var uri struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}
	uri.Name = fmt.Sprintf("lock #%s", nftID)
	uri.Description = VeNftClass.Description
	uri.Image = fmt.Sprintf("data:image/svg+xml;base64,%s", base64.URLEncoding.EncodeToString([]byte(output)))

	uriStr, err := json.Marshal(&uri)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("data:application/json;base64,%s", base64.URLEncoding.EncodeToString(uriStr))
}
