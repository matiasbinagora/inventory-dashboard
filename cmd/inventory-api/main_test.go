package main

import "testing"

func TestDefaultAPIAddressIsLocalOnly(t *testing.T) {
	if defaultAPIAddress != "127.0.0.1:8080" {
		t.Fatalf("default API address = %q, want loopback", defaultAPIAddress)
	}
}

func TestEnvOrUsesFallbackWithoutConfiguration(t *testing.T) {
	t.Setenv("INVENTORY_API_ADDR", "")
	if got := envOr("INVENTORY_API_ADDR", defaultAPIAddress); got != defaultAPIAddress {
		t.Fatalf("envOr() = %q, want %q", got, defaultAPIAddress)
	}
}
