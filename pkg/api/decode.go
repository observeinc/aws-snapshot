package api

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func DecodeConfig(input interface{}, output interface{}) error {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      &output,
		ErrorUnused: true,
	})

	if err := decoder.Decode(input); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	return nil
}
