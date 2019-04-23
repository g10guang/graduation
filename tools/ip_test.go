package tools

import (
	"fmt"
	"testing"
)

func TestGetCurrentIp(t *testing.T) {
	ipStr := GetLocalIpStr()
	fmt.Println(ipStr)
	ipI32 := GetLocalIpI32()
	fmt.Println(ipI32)
}

func TestGetLogID(t *testing.T) {
	fmt.Println(MakeLogID())
}

