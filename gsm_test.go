package bird

import (
	"fmt"
	"math"
	"testing"
)

func TestAlignment(t *testing.T) {
	var data []byte
	entries := 3
	expectedLength := 3

	mask := ^(1 << 7)

	//fmt.Println(math.MaxUint8)
	//fmt.Println(uint8(math.MaxUint8 & mask))

	for i := 0; i < entries; i++ {
		entry := byte(math.MaxUint8 & mask)
		data = append(data, entry)
	}

	res := bit7Align(data)

	if length := len(res); length != expectedLength {
		t.Errorf("Alignment is of length %d, should be %d", length, expectedLength)
	}
	fmt.Println(res)
}
