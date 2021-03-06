package wallet

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlindTransaction(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(NewWalletFromMnemonicOpts{
		SigningMnemonic:  strings.Split("quarter multiply swarm depth slice security flight glad arrow express worth legend wasp mobile anchor dinner mutual six sure wear section delay initial thank", " "),
		BlindingMnemonic: strings.Split("okay door hammer betray reason zero fiction rigid vivid scorpion thunder crucial focus riot cancel wear autumn rely kangaroo rug raven mystery ability stem", " "),
	})
	if err != nil {
		t.Fatal(err)
	}

	blindingKey1, _ := hex.DecodeString("030034486c2d456f451459f35f693a8c60016ffd07197545df1b5324fa85e56eb4")
	blindingKey2, _ := hex.DecodeString("026a01a023da02ac39c319a6517d961d9f8a198a2507372a9e87e697ab2c1e87b3")
	opts := BlindTransactionOpts{
		PsetBase64:         "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",
		OutputBlindingKeys: [][]byte{blindingKey1, blindingKey2},
	}
	blindedPsetBase64, err := wallet.BlindTransaction(opts)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(blindedPsetBase64) > 0)
}

func TestFailingBlindTransaction(t *testing.T) {
	wallet, err := newTestWallet()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		opts BlindTransactionOpts
		err  error
	}{
		{
			opts: BlindTransactionOpts{
				PsetBase64: "",
			},
			err: ErrNullPset,
		},
		{
			opts: BlindTransactionOpts{
				PsetBase64: "cHNldP8BAOoCAAAAAAGA5RCreFagpc3t/LtM7IaVNJsxhUECqpKZTyY+NPBknQAAAAAA/////wMBJbJRBw4pyhkEPPM8zXMk4t2rA+zErgted8T8Dlz2yVoBAAAAAAL68IAAGXapFDk5cIC1HvIsWb10aa+s/77sDaEuiKwBJbJRBw4pyhkEPPM8zXMk4t2rA+zErgted8T8Dlz2yVoBAAAAAAL67PwAGXapFGWb7bXT08erEtf4UyPDobbAYO++iKwBJbJRBw4pyhkEPPM8zXMk4t2rA+zErgted8T8Dlz2yVoBAAAAAAAAAfQAAAAAAAAAAAAAAA==",
			},
			err: ErrNullInputWitnessUtxo,
		},
	}

	for _, tt := range tests {
		_, err := wallet.BlindTransaction(tt.opts)
		assert.Equal(t, tt.err, err)
	}
}
