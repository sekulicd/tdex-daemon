package trade

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/tdex-network/tdex-daemon/pkg/bufferutil"
	"github.com/tdex-network/tdex-daemon/pkg/explorer"
	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/go-elements/payment"
	"github.com/vulpemventures/go-elements/pset"
	"github.com/vulpemventures/go-elements/transaction"
)

func NewSwapTx(
	unspents []explorer.Utxo,
	inAsset string,
	inAmount uint64,
	outAsset string,
	outAmount uint64,
	outScript []byte,
) (string, error) {
	ptx, err := pset.New([]*transaction.TxInput{}, []*transaction.TxOutput{}, 2, 0)
	if err != nil {
		return "", err
	}

	selectedUnspents, change, err := explorer.SelectUnspents(
		unspents,
		inAmount,
		inAsset,
	)
	if err != nil {
		return "", err
	}

	updater, _ := pset.NewUpdater(ptx)

	for _, in := range selectedUnspents {
		input, witnessUtxo, _ := in.Parse()
		updater.AddInput(input)
		err := updater.AddInWitnessUtxo(witnessUtxo, len(ptx.Inputs)-1)
		if err != nil {
			return "", err
		}
	}

	output, err := newTxOutput(outAsset, outAmount, outScript)
	if err != nil {
		return "", err
	}
	updater.AddOutput(output)

	if change > 0 {
		changeOutput, err := newTxOutput(inAsset, change, outScript)
		if err != nil {
			return "", err
		}
		updater.AddOutput(changeOutput)
	}

	return ptx.ToBase64()
}

type Wallet struct {
	privateKey         *btcec.PrivateKey
	blindingPrivateKey *btcec.PrivateKey
	network            *network.Network
}

func NewRandomWallet(net *network.Network) (*Wallet, error) {
	prvkey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	blindPrvkey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return &Wallet{prvkey, blindPrvkey, net}, nil
}

func NewWalletFromKey(privateKey, blindingKey []byte, net *network.Network) *Wallet {
	prvkey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)
	blindPrvkey, _ := btcec.PrivKeyFromBytes(btcec.S256(), blindingKey)

	return &Wallet{prvkey, blindPrvkey, net}
}

func (w *Wallet) Address() string {
	p2wpkh := payment.FromPublicKey(w.privateKey.PubKey(), w.network, w.blindingPrivateKey.PubKey())
	ctAddress, _ := p2wpkh.ConfidentialWitnessPubKeyHash()
	return ctAddress
}

func (w *Wallet) Script() ([]byte, []byte) {
	p2wpkh := payment.FromPublicKey(w.privateKey.PubKey(), w.network, w.blindingPrivateKey.PubKey())
	return p2wpkh.Script, p2wpkh.WitnessScript
}

func (w *Wallet) Sign(psetBase64 string) (string, error) {
	ptx, err := pset.NewPsetFromBase64(psetBase64)
	if err != nil {
		return "", err
	}
	updater, err := pset.NewUpdater(ptx)
	if err != nil {
		return "", err
	}

	for i, in := range ptx.Inputs {
		script, witnessScript := w.Script()
		if bytes.Equal(in.WitnessUtxo.Script, witnessScript) {
			hashForSignature := ptx.UnsignedTx.HashForWitnessV0(
				i,
				script,
				in.WitnessUtxo.Value,
				txscript.SigHashAll,
			)

			signature, err := w.privateKey.Sign(hashForSignature[:])
			if err != nil {
				return "", err
			}

			if !signature.Verify(hashForSignature[:], w.privateKey.PubKey()) {
				return "", fmt.Errorf(
					"signature verification failed for input %d",
					i,
				)
			}

			sigWithSigHashType := append(signature.Serialize(), byte(txscript.SigHashAll))
			_, err = updater.Sign(
				i,
				sigWithSigHashType,
				w.privateKey.PubKey().SerializeCompressed(),
				nil,
				nil,
			)
		}
	}

	return ptx.ToBase64()
}

func (w *Wallet) BlindingKey() []byte {
	return w.blindingPrivateKey.Serialize()
}

func newTxOutput(assetHex string, amount uint64, script []byte) (*transaction.TxOutput, error) {
	asset, err := bufferutil.AssetHashToBytes(assetHex)
	if err != nil {
		return nil, err
	}
	value, err := bufferutil.ValueToBytes(amount)
	if err != nil {
		return nil, err
	}
	return transaction.NewTxOutput(asset, value, script), nil
}
