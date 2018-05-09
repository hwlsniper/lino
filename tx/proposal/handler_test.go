package proposal

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	c4600 = types.Coin{4600 * types.Decimals}
)

func TestProposalBasic(t *testing.T) {
	ctx, am, pm, _, _, gm := setupTest(t, 0)
	handler := NewHandler(am, pm, gm)
	pm.InitGenesis(ctx)

	rat := sdk.Rat{Denom: 10, Num: 5}
	para := param.GlobalAllocationParam{
		ContentCreatorAllocation: rat,
	}
	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c4600)

	// let user1 create a proposal
	msg := NewChangeGlobalAllocationMsg("user1", para)
	resultPass := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, resultPass)

	// invalid create
	invalidMsg := NewChangeGlobalAllocationMsg("wqdkqwndkqwd", para)
	resultInvalid := handler(ctx, invalidMsg)
	assert.Equal(t, ErrUsernameNotFound().Result(), resultInvalid)

	result2 := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result2)

	proposal, _ := pm.storage.GetProposal(ctx, proposalID1)
	p := proposal.(*model.ChangeGlobalAllocationParamProposal)
	assert.Equal(t, true, p.Description.ContentCreatorAllocation.Equal(rat))

	// check proposal list is correct
	lst, _ := pm.storage.GetProposalList(ctx)
	assert.Equal(t, 2, len(lst.OngoingProposal))
	assert.Equal(t, proposalID1, lst.OngoingProposal[0])
	assert.Equal(t, proposalID2, lst.OngoingProposal[1])

	// test delete proposal
	pm.storage.DeleteProposal(ctx, proposalID2)
	_, err := pm.storage.GetProposal(ctx, proposalID2)
	assert.Equal(t, model.ErrGetProposal(), err)

}