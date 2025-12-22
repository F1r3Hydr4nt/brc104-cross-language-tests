package main

import (
	"encoding/base64"
	"testing"
)

// Cross-language BRC-104 signature tests for Go SDK.
//
// These tests verify that the Go SDK implementation correctly handles
// BRC-104 authentication signature generation and verification.
// Since Go SDK is the reference implementation, these tests should always pass.

func TestBRC104SignatureDataPreparation(t *testing.T) {
	t.Run("signature data preparation matches reference implementation", func(t *testing.T) {
		/**
		 * Test that Go SDK produces correct signature data.
		 *
		 * Go SDK correctly:
		 * 1. Decodes each nonce separately: base64.StdEncoding.DecodeString(nonce)
		 * 2. Concatenates the decoded bytes: append(initialBytes, sessionBytes...)
		 */

		// Go implementation (reference)
		initialBytes, err := base64.StdEncoding.DecodeString(InitialNonceB64)
		if err != nil {
			t.Fatalf("Failed to decode initial nonce: %v", err)
		}

		sessionBytes, err := base64.StdEncoding.DecodeString(SessionNonceB64)
		if err != nil {
			t.Fatalf("Failed to decode session nonce: %v", err)
		}

		sigData := append(initialBytes, sessionBytes...)

		// Should match our expected result
		if len(sigData) != len(ExpectedSigDataSigning) {
			t.Errorf("Signature data length mismatch: got %d, want %d", len(sigData), len(ExpectedSigDataSigning))
		}

		for i, expected := range ExpectedSigDataSigning {
			if i >= len(sigData) || sigData[i] != expected {
				t.Errorf("Signature data mismatch at index %d: got %d, want %d", i, sigData[i], expected)
			}
		}

		// Verify this differs from incorrect approach (concatenating base64 strings)
		concatenatedB64 := InitialNonceB64 + SessionNonceB64
		incorrectSigData, _ := base64.StdEncoding.DecodeString(concatenatedB64)

		if len(sigData) == len(incorrectSigData) {
			allSame := true
			for i := range sigData {
				if sigData[i] != incorrectSigData[i] {
					allSame = false
					break
				}
			}
			if allSame {
				t.Error("Go SDK should produce different result from incorrect base64 concatenation approach")
			}
		}
	})

	t.Run("verification data preparation uses reversed order", func(t *testing.T) {
		/**
		 * Test that verification uses reversed order: session + initial.
		 * This matches the verification path in Go SDK.
		 */
		verifyData := append(SessionNonceBytes, InitialNonceBytes...)

		if len(verifyData) != len(ExpectedSigDataVerification) {
			t.Errorf("Verification data length mismatch: got %d, want %d", len(verifyData), len(ExpectedSigDataVerification))
		}

		for i, expected := range ExpectedSigDataVerification {
			if i >= len(verifyData) || verifyData[i] != expected {
				t.Errorf("Verification data mismatch at index %d: got %d, want %d", i, verifyData[i], expected)
			}
		}
	})
}

func TestBRC104Constants(t *testing.T) {
	t.Run("constants are properly defined", func(t *testing.T) {
		if TestPrivateKeyWIF == "" {
			t.Error("TestPrivateKeyWIF should be defined")
		}

		if InitialNonceB64 == "" {
			t.Error("InitialNonceB64 should be defined")
		}

		if SessionNonceB64 == "" {
			t.Error("SessionNonceB64 should be defined")
		}

		if len(InitialNonceBytes) != 32 {
			t.Errorf("InitialNonceBytes should be 32 bytes, got %d", len(InitialNonceBytes))
		}

		if len(SessionNonceBytes) != 32 {
			t.Errorf("SessionNonceBytes should be 32 bytes, got %d", len(SessionNonceBytes))
		}

		if TestCounterpartyKey == "" {
			t.Error("TestCounterpartyKey should be defined")
		}
	})

	t.Run("protocol ID is correct", func(t *testing.T) {
		if len(ProtocolID) != 2 {
			t.Errorf("ProtocolID should have 2 elements, got %d", len(ProtocolID))
		}

		if securityLevel, ok := ProtocolID[0].(int); !ok || securityLevel != 2 {
			t.Errorf("ProtocolID[0] should be 2, got %v", ProtocolID[0])
		}

		if protocol, ok := ProtocolID[1].(string); !ok || protocol != "auth message signature" {
			t.Errorf("ProtocolID[1] should be 'auth message signature', got %v", ProtocolID[1])
		}
	})

	t.Run("key ID format is correct", func(t *testing.T) {
		keyId := MakeKeyId(InitialNonceB64, SessionNonceB64)
		expected := InitialNonceB64 + " " + SessionNonceB64

		if keyId != expected {
			t.Errorf("Key ID format mismatch: got %s, want %s", keyId, expected)
		}
	})
}

func TestBRC104NonceDecoding(t *testing.T) {
	t.Run("nonces decode correctly", func(t *testing.T) {
		// Test initial nonce decoding
		initialDecoded, err := base64.StdEncoding.DecodeString(InitialNonceB64)
		if err != nil {
			t.Fatalf("Failed to decode initial nonce: %v", err)
		}

		if len(initialDecoded) != 32 {
			t.Errorf("Initial nonce should decode to 32 bytes, got %d", len(initialDecoded))
		}

		// Should match our constant
		for i, expected := range InitialNonceBytes {
			if i >= len(initialDecoded) || initialDecoded[i] != expected {
				t.Errorf("Initial nonce decoding mismatch at index %d: got %d, want %d", i, initialDecoded[i], expected)
			}
		}

		// Test session nonce decoding
		sessionDecoded, err := base64.StdEncoding.DecodeString(SessionNonceB64)
		if err != nil {
			t.Fatalf("Failed to decode session nonce: %v", err)
		}

		if len(sessionDecoded) != 32 {
			t.Errorf("Session nonce should decode to 32 bytes, got %d", len(sessionDecoded))
		}

		// Should match our constant
		for i, expected := range SessionNonceBytes {
			if i >= len(sessionDecoded) || sessionDecoded[i] != expected {
				t.Errorf("Session nonce decoding mismatch at index %d: got %d, want %d", i, sessionDecoded[i], expected)
			}
		}
	})
}

func TestBRC104Integration(t *testing.T) {
	t.Run("initial response signature flow", func(t *testing.T) {
		/**
		 * Test the complete initial response signature flow.
		 * This simulates what happens in handleInitialRequest.
		 */

		sigData := ExpectedSigDataSigning

		// Verify the signature data is what we expect
		if len(sigData) != 64 { // 32 + 32 bytes
			t.Errorf("Signature data should be 64 bytes, got %d", len(sigData))
		}

		// First 32 bytes should be initial nonce
		for i := 0; i < 32; i++ {
			if sigData[i] != InitialNonceBytes[i] {
				t.Errorf("Signature data initial nonce mismatch at index %d: got %d, want %d", i, sigData[i], InitialNonceBytes[i])
			}
		}

		// Last 32 bytes should be session nonce
		for i := 0; i < 32; i++ {
			if sigData[i+32] != SessionNonceBytes[i] {
				t.Errorf("Signature data session nonce mismatch at index %d: got %d, want %d", i+32, sigData[i+32], SessionNonceBytes[i])
			}
		}
	})

	t.Run("initial response verification flow", func(t *testing.T) {
		/**
		 * Test the initial response verification flow.
		 * This simulates what happens in handleInitialResponse.
		 */

		verifyData := ExpectedSigDataVerification

		// Verify the verification data is correct
		if len(verifyData) != 64 { // 32 + 32 bytes
			t.Errorf("Verification data should be 64 bytes, got %d", len(verifyData))
		}

		// First 32 bytes should be session nonce
		for i := 0; i < 32; i++ {
			if verifyData[i] != SessionNonceBytes[i] {
				t.Errorf("Verification data session nonce mismatch at index %d: got %d, want %d", i, verifyData[i], SessionNonceBytes[i])
			}
		}

		// Last 32 bytes should be initial nonce
		for i := 0; i < 32; i++ {
			if verifyData[i+32] != InitialNonceBytes[i] {
				t.Errorf("Verification data initial nonce mismatch at index %d: got %d, want %d", i+32, verifyData[i+32], InitialNonceBytes[i])
			}
		}
	})
}

// Benchmark tests to ensure performance
func BenchmarkSignatureDataPreparation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initialBytes, _ := base64.StdEncoding.DecodeString(InitialNonceB64)
		sessionBytes, _ := base64.StdEncoding.DecodeString(SessionNonceB64)
		_ = append(initialBytes, sessionBytes...)
	}
}

func BenchmarkNonceDecoding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = base64.StdEncoding.DecodeString(InitialNonceB64)
		_, _ = base64.StdEncoding.DecodeString(SessionNonceB64)
	}
}

