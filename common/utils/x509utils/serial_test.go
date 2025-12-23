package x509utils

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestUnmarshalJSON_Decimal(t *testing.T) {
	var b BigInt
	data := []byte("\"12345678901234567890\"")
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("decimal unmarshal failed: %v", err)
	}
	if b.Int == nil || b.String() != "12345678901234567890" {
		t.Fatalf("unexpected value: %v", b.Int)
	}
}

func TestUnmarshalJSON_HexPrefixed(t *testing.T) {
	var b BigInt
	data := []byte("\"0xAABBCC\"")
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("hex-prefixed unmarshal failed: %v", err)
	}
	expected := new(big.Int)
	expected.SetString("AABBCC", 16)
	if b.Int == nil || b.Cmp(expected) != 0 {
		t.Fatalf("unexpected hex value: got %v want %v", b.Int, expected)
	}
}

func TestUnmarshalJSON_PlainHexRejected(t *testing.T) {
	var b BigInt
	data := []byte("\"AABBCC\"")
	if err := json.Unmarshal(data, &b); err == nil {
		t.Fatalf("plain hex should be rejected but succeeded: %v", b.Int)
	}
}

func TestUnmarshalJSON_NumberUnquoted(t *testing.T) {
	var b BigInt
	data := []byte("12345")
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("number unquoted unmarshal failed: %v", err)
	}
	if b.Int == nil || b.String() != "12345" {
		t.Fatalf("unexpected value for unquoted number: %v", b.Int)
	}
}

func TestUnmarshalJSON_Null(t *testing.T) {
	var b BigInt
	data := []byte("null")
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("null unmarshal failed: %v", err)
	}
	if b.Int != nil {
		t.Fatalf("expected nil Int for null, got %v", b.Int)
	}
}
