package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

type Memo uint64

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Username types.AccountKey `json:"username"`
	Created  int64            `json:"created"`
	PostKey  crypto.PubKey    `json:"post_key"`
	OwnerKey crypto.PubKey    `json:"owner_key"`
	Address  sdk.Address      `json:"address"`
}

// AccountBank uses Address as the key instead of Username
type AccountBank struct {
	Address  sdk.Address      `json:"address"`
	Balance  types.Coin       `json:"balance"`
	Username types.AccountKey `json:"username"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence         int64 `json:"sequence"`
	LastActivity     int64 `json:"last_activity"`
	ActivityBurden   int64 `json:"activity_burden"`
	StakeRatio       int64 `json:"stake_ratio"`
	LastSRUpdateTime int64 `json:"last_stake_ratio_update_time"`
}

// AccountInfraConsumption records infra utility consumption
type AccountInfraConsumption struct {
	Storage   int64 `json:"storage"`
	Bandwidth int64 `json:"bandwidth"`
}

// record all meta info about this relation
type FollowerMeta struct {
	CreatedAt    int64            `json:"created_at"`
	FollowerName types.AccountKey `json:"follower_name"`
}

// record all meta info about this relation
type FollowingMeta struct {
	CreatedAt     int64            `json:"created_at"`
	FollowingName types.AccountKey `json:"following_name"`
}

// reward get from the inflation pool
type Reward struct {
	OriginalIncome types.Coin `json:"original_income"`
	ActualReward   types.Coin `json:"actual_reward"`
	UnclaimReward  types.Coin `json:"unclaim_reward"`
}