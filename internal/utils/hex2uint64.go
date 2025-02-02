package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func HexToUint64(hexString string) (uint64, error) {
	hexString = strings.TrimSpace(hexString)
	hexString = strings.TrimPrefix(hexString, "0x")

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("convert hex string to uint64: %w", err)
	}

	return result, nil
}
