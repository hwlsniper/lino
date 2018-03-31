package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	FlagDonator = "donator"
	FlagAmount  = "amount"
)

// SendTxCommand will create a send tx and sign it with the given key
func DonateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate",
		Short: "donate to a post",
		RunE:  sendDonateTx(cdc),
	}
	cmd.Flags().String(FlagDonator, "", "donator of this transaction")
	cmd.Flags().String(FlagAuthor, "", "author of the target post")
	cmd.Flags().String(FlagPostID, "", "post id of the target post")
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	return cmd
}

// send register transaction to the blockchain
func sendDonateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		username := viper.GetString(FlagDonator)
		author := viper.GetString(FlagAuthor)
		postID := viper.GetString(FlagPostID)

		amount, err := sdk.ParseCoins(viper.GetString(FlagAmount))
		if err != nil {
			return err
		}

		msg := post.NewDonateMsg(acc.AccountKey(username), amount, acc.AccountKey(author), postID)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(username, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}