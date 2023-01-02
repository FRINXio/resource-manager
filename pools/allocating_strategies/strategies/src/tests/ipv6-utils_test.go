package tests

import (
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"math/big"
	"reflect"
	"testing"
)

func TestIpv6IAddressToAmountAddresses(t *testing.T) {
	b, _ := src.Ipv6InetAton("dead::beef")
	if eq := reflect.DeepEqual(b.String(), "295986882420777848964380943247191621359"); !eq {
		t.Fatalf("different output of nil expected, got: %s", b.String())
	}
}

func TestIpv6AmountAddressesToStringAddress(t *testing.T) {
	input, _ := new(big.Int).SetString("295986882420777848964380943247191621359", 10)
	b:= src.Ipv6InetNtoa(input)
	if eq := reflect.DeepEqual(b, "dead::beef"); !eq {
		t.Fatalf("different output of nil expected, got: %s", b)
	}
}
