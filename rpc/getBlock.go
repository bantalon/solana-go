package rpc

import (
	"context"
	"fmt"

	bin "github.com/dfuse-io/binary"
	"github.com/gagliardetto/solana-go"
)

type TransactionDetailsType string

const (
	TransactionDetailsFull       TransactionDetailsType = "full"
	TransactionDetailsSignatures TransactionDetailsType = "signatures"
	TransactionDetailsNone       TransactionDetailsType = "none"
)

type GetBlockOpts struct {
	// Encoding for each returned Transaction, either "json", "jsonParsed", "base58" (slow), "base64".
	// If parameter not provided, the default encoding is "json".
	// - "jsonParsed" encoding attempts to use program-specific instruction parsers to return
	//   more human-readable and explicit data in the transaction.message.instructions list.
	// - If "jsonParsed" is requested but a parser cannot be found, the instruction falls back
	//   to regular JSON encoding (accounts, data, and programIdIndex fields).
	//
	// This parameter is optional.
	Encoding solana.EncodingType

	// Level of transaction detail to return.
	// If parameter not provided, the default detail level is "full".
	//
	// This parameter is optional.
	TransactionDetails TransactionDetailsType

	// Whether to populate the rewards array.
	// If parameter not provided, the default includes rewards.
	//
	// This parameter is optional.
	Rewards *bool

	// "processed" is not supported.
	// If parameter not provided, the default is "finalized".
	//
	// This parameter is optional.
	Commitment CommitmentType
}

// GetBlock returns identity and transaction information about a confirmed block in the ledger.
//
// NEW: This method is only available in solana-core v1.7 or newer.
// Please use `getConfirmedBlock` for solana-core v1.6
func (cl *Client) GetBlock(
	ctx context.Context,
	slot uint64,
) (out *GetBlockResult, err error) {
	return cl.GetBlockWithOpts(
		ctx,
		slot,
		nil,
	)
}

// GetBlock returns identity and transaction information about a confirmed block in the ledger.
//
// NEW: This method is only available in solana-core v1.7 or newer.
// Please use `getConfirmedBlock` for solana-core v1.6
func (cl *Client) GetBlockWithOpts(
	ctx context.Context,
	slot uint64,
	opts *GetBlockOpts,
) (out *GetBlockResult, err error) {

	obj := M{
		"encoding": solana.EncodingJSON,
	}

	if opts != nil {
		if opts.TransactionDetails != "" {
			obj["transactionDetails"] = opts.TransactionDetails
		}
		if opts.Rewards != nil {
			obj["rewards"] = opts.Rewards
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Encoding != "" {
			if !solana.IsAnyOfEncodingType(
				opts.Encoding,
				// Valid encodings:
				solana.EncodingJSON,
				solana.EncodingJSONParsed,
				solana.EncodingBase58,
				solana.EncodingBase64,
			) {
				return nil, fmt.Errorf("provided encoding is not supported: %s", opts.Encoding)
			}
			obj["encoding"] = opts.Encoding
		}
	}

	params := []interface{}{slot, obj}

	err = cl.rpcClient.CallForInto(ctx, &out, "getBlock", params)

	if err != nil {
		return nil, err
	}
	if out == nil {
		// Block is not confirmed.
		return nil, ErrNotConfirmed
	}
	return
}

type GetBlockResult struct {
	// The blockhash of this block.
	Blockhash solana.Hash `json:"blockhash"`

	// The blockhash of this block's parent;
	// if the parent block is not available due to ledger cleanup,
	// this field will return "11111111111111111111111111111111".
	PreviousBlockhash solana.Hash `json:"previousBlockhash"`

	// The slot index of this block's parent.
	ParentSlot bin.Uint64 `json:"parentSlot"`

	// Present if "full" transaction details are requested.
	Transactions []TransactionWithMeta `json:"transactions"`

	// Present if "signatures" are requested for transaction details;
	// an array of signatures, corresponding to the transaction order in the block.
	Signatures []solana.Signature `json:"signatures"`

	// Present if rewards are requested.
	Rewards []BlockReward `json:"rewards"`

	// Estimated production time, as Unix timestamp (seconds since the Unix epoch).
	// Nil if not available.
	BlockTime *bin.Int64 `json:"blockTime"`

	// The number of blocks beneath this block.
	BlockHeight *bin.Uint64 `json:"blockHeight"`
}
