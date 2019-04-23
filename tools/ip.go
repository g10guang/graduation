package tools

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func GetLocalIpStr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, v := range addrs {
		if ipNet, ok := v.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	panic(errors.New("Cannot get current IP"))
}

func GetLocalIpI32() int32 {
	ipStr := GetLocalIpStr()
	ipSlice := strings.Split(ipStr, ".")
	if len(ipSlice) != 4 {
		return 0
	}
	ipSliceI32 := make([]int, 4)
	for i, v := range ipSlice {
		ipSliceI32[i], _ = strconv.Atoi(v)
	}
	var ipI32 int32
	ipI32 = int32(ipSliceI32[0] << 24) + int32(ipSliceI32[1] << 16) + int32(ipSliceI32[2] << 8) + int32(ipSliceI32[3])
	return ipI32
}


func GetLocalIpInt() int {
	ipStr := GetLocalIpStr()
	ipSlice := strings.Split(ipStr, ".")
	if len(ipSlice) != 4 {
		return 0
	}
	ipSliceI32 := make([]int, 4)
	for i, v := range ipSlice {
		ipSliceI32[i], _ = strconv.Atoi(v)
	}
	var ipI32 int
	ipI32 = ipSliceI32[0] << 24 + ipSliceI32[1] << 16 + ipSliceI32[2] << 8 + ipSliceI32[3]
	return ipI32
}
