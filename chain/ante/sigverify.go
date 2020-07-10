package ante

import (
	"bytes"
	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Verify all signatures for a tx and return an error if any are invalid. Note,
// the SigVerificationDecorator decorator will not get executed on ReCheck.
//
// CONTRACT: Pubkeys are set in context for all signers before this decorator runs
// CONTRACT: Tx must implement SigVerifiableTx interface
type SigVerificationDecorator struct {
	ak keeper.AccountKeeper
}

func NewSigVerificationDecorator(ak keeper.AccountKeeper) SigVerificationDecorator {
	return SigVerificationDecorator{
		ak: ak,
	}
}

func (svd SigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// no need to verify signatures on recheck tx
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}
	sigTx, ok := tx.(SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type for SigVerifiableTx")
	}

	stdTx, ok := tx.(types.StdTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type for stdTx")
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	sigs := sigTx.GetSignatures()

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	signerAddrs := sigTx.GetSigners()

	// check that signer length and signature length are the same
	if len(sigs) != len(signerAddrs) {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
	}

	for i, sig := range sigs {
		seq, num, err := svd.ak.GetAuthSequence(ctx, signerAddrs[i])
		if err != nil {
			return ctx, err
		}

		// retrieve signBytes of tx
		signBytes := txutil.GetSignBytes(ctx, &stdTx, num, seq)

		// retrieve pubkey
		pubKey := sig.PubKey
		if !simulate && pubKey == nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey on account is not set")
		}

		// verify signature
		if !simulate && !pubKey.VerifyBytes(signBytes, sig.Signature) {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signature verification failed; verify correct account sequence and chain-id")
		}
	}

	return next(ctx, tx, simulate)
}

type IncrementSequenceDecorator struct {
	ak keeper.AccountKeeper
}

func NewIncrementSequenceDecorator(ak keeper.AccountKeeper) IncrementSequenceDecorator {
	return IncrementSequenceDecorator{
		ak: ak,
	}
}

func (isd IncrementSequenceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// no need to increment sequence on CheckTx or RecheckTx
	if ctx.IsCheckTx() && !simulate {
		return next(ctx, tx, simulate)
	}

	sigTx, ok := tx.(SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type for SigVerifiableTx in inc seq")
	}

	// increment sequence of all signers
	for _, addr := range sigTx.GetSigners() {
		isd.ak.IncAuthSequence(ctx, addr)
	}

	return next(ctx, tx, simulate)
}

// SetPubKeyDecorator sets PubKeys in context for any signer which does not already have pubkey set
// PubKeys must be set in context for all signers before any other sigverify decorators run
// CONTRACT: Tx must implement SigVerifiableTx interface
type SetPubKeyDecorator struct {
	ak keeper.AccountKeeper
}

func NewSetPubKeyDecorator(ak keeper.AccountKeeper) SetPubKeyDecorator {
	return SetPubKeyDecorator{
		ak: ak,
	}
}

func (spkd SetPubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	signatures := sigTx.GetSignatures()
	signers := sigTx.GetSigners()

	for i, sig := range signatures {
		// PublicKey was omitted from slice since it has already been set in context
		if sig.PubKey == nil {
			if !simulate {
				continue
			}
			// sig.PubKey = simSecp256k1Pubkey
		}
		// Only make check if simulate=false
		if !simulate && !bytes.Equal(sig.PubKey.Address(), signers[i]) {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey,
				"pubKey does not match signer address %s with signer index: %d", signers[i], i)
		}

		spkd.ak.SetPubKey(ctx, signers[i], sig.PubKey)
	}

	return next(ctx, tx, simulate)
}
