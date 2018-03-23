package post

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
)

func NewHandler(pm PostManager, am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreatePostMsg:
			return handleCreatePostMsg(ctx, pm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleCreatePostMsg(ctx sdk.Context, pm PostManager, am acc.AccountManager, msg CreatePostMsg) sdk.Result {
	account := acc.NewProxyAccount(msg.Author, &am)
	if !account.IsAccountExist(ctx) {
		return ErrPostCreateNonExistAuthor().Result()
	}
	post := NewProxyPost(msg.Author, msg.PostID, &pm)
	if post.IsPostExist(ctx) {
		return ErrPostExist().Result()
	}
	if err := post.CreatePost(ctx, &msg.PostInfo); err != nil {
		return err.Result()
	}
	if len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0 {
		parentPost := NewProxyPost(msg.ParentAuthor, msg.ParentPostID, &pm)
		if err := parentPost.AddComment(ctx, post.GetPostKey()); err != nil {
			return err.Result()
		}
		if err := parentPost.Apply(ctx); err != nil {
			return err.Result()
		}
	}
	if err := post.Apply(ctx); err != nil {
		return err.Result()
	}
	if err := account.UpdateLastActivity(ctx); err != nil {
		return err.Result()
	}
	if err := account.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
