package utils

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/client"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"gopkg.in/yaml.v2"
)

func ToPrintObj(ctx client.Context, toPrint interface{}) ([]byte, error) {
	var (
		out []byte
		err error
	)

	switch ctx.GetOutputFormat() {
	case "text":
		out, err = yaml.Marshal(&toPrint)

	case "json":
		if ctx.Indent() {
			out, err = ctx.Codec().MarshalJSONIndent(toPrint, "", "  ")
		} else {
			if canBePretty, ok := toPrint.(types.Prettifier); ok {
				out, err = canBePretty.PrettifyJSON(ctx.Codec())
			} else {
				out, err = ctx.Codec().MarshalJSON(toPrint)
			}
		}
	}

	return out, err
}

// PrintOutput prints output while respecting output and indent flags
func PrintOutput(ctx client.Context, toPrint interface{}) error {
	out, err := ToPrintObj(ctx, toPrint)

	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
