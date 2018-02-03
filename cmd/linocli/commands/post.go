package commands

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/light-client/commands"
	txcmd "github.com/tendermint/light-client/commands/txs"

	btypes "github.com/lino-network/lino/types"
)

//-------------------------
// SendTx

// SendTxCmd is CLI command to send tokens between basecoin accounts
var PostTxCmd = &cobra.Command{
	Use:   "post",
	Short: "send a short post",
	RunE:  commands.RequireInit(doPostTx),
}

//nolint
const (
	FlagTitle      = "title"
	FlagContent    = "content"
	FlagParentAddr = "parentaddr"
	FlagParentSeq  = "parentseq"
)

func init() {
	flags := PostTxCmd.Flags()
	flags.String(FlagTitle, "", "Post Title")
	flags.String(FlagAddress, "", "Post author")
	flags.String(FlagContent, "", "Post content")
	flags.Int(FlagSequence, -1, "Sequence number for this post")
}

// runDemo is an example of how to make a tx
func doPostTx(cmd *cobra.Command, args []string) error {
	// load data from json or flags
	tx := new(btypes.PostTx)
	err := readPostTxFlags(tx)
	if err != nil {
		return err
	}

	// Wrap and add signer
	post := &PostTx{
		chainID: commands.GetChainID(),
		Tx:      tx,
	}
	post.AddSigner(txcmd.GetSigner())
	// Sign if needed and post.  This it the work-horse
	bres, err := txcmd.SignAndPostTx(post)
	if err != nil {
		return err
	}

	// Output result
	return txcmd.OutputTx(bres)
}

func readPostTxFlags(tx *btypes.PostTx) error {
	//parse the fee and amounts into coin types
	poster, err := hex.DecodeString(cmn.StripHex(FlagAddress))
	if err != nil {
		return errors.Wrap(err, "Invalid address")
	}
	parentAddr, _ := hex.DecodeString(cmn.StripHex(FlagParentAddr))
	if err == nil {
		// Set parent post
		tx.Parent = btypes.PostID(parentAddr, viper.GetInt(FlagParentSeq))
	}
	tx.Address = poster
	tx.Title = viper.GetString(FlagTitle)
	tx.Content = viper.GetString(FlagContent)
	tx.Sequence = viper.GetInt(FlagSequence)
	return nil
}
