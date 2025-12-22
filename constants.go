package main

// Shared constants for cross-language BRC-104 signature testing.
//
// These constants are shared across Python, TypeScript, and Go test suites
// to ensure identical inputs produce identical outputs.

var (
	// Private key in WIF format (fixed for testing)
	TestPrivateKeyWIF = "L4B2postXdaP7TiUrUBYs53Fqzheu7WhSoQVPuY8qBdoBeEwbmZx"

	// Test nonces (32 bytes each, base64 encoded)
	// These are fixed values to ensure deterministic testing
	InitialNonceB64 = "QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUE=" // 32 bytes of 'A' characters
	SessionNonceB64 = "QkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkI=" // 32 bytes of 'B' characters
)

// Decode nonces to bytes for direct use
var (
	InitialNonceBytes = []byte{65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65}
	SessionNonceBytes = []byte{66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66}
)

// Protocol parameters for BRC-104 authentication
var ProtocolID = []interface{}{2, "auth message signature"} // Security level 2, protocol name

// MakeKeyId creates key ID in format: "{initial_nonce} {session_nonce}"
func MakeKeyId(initialNonce, sessionNonce string) string {
	return initialNonce + " " + sessionNonce
}

// Test counterparty (identity key) - using a fixed public key for testing
// This should be derived from TestPrivateKeyWIF, but for simplicity we'll use a fixed value
const TestCounterpartyKey = "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798" // Fixed test key

// Expected signature data (what should be signed)
// For signing: initial_nonce_bytes + session_nonce_bytes
var ExpectedSigDataSigning = []byte{
	65,65,65,65,65,65,65,65,
	65,65,65,65,65,65,65,65,
	65,65,65,65,65,65,65,65,
	65,65,65,65,65,65,65,65,
	66,66,66,66,66,66,66,66,
	66,66,66,66,66,66,66,66,
	66,66,66,66,66,66,66,66,
	66,66,66,66,66,66,66,66,
}

// For verification: session_nonce_bytes + initial_nonce_bytes
var ExpectedSigDataVerification = []byte{
	66,66,66,66,66,66,66,66,
	66,66,66,66,66,66,66,66,
	66,66,66,66,66,66,66,66,
	66,66,66,66,66,66,66,66,
	65,65,65,65,65,65,65,65,
	65,65,65,65,65,65,65,65,
	65,65,65,65,65,65,65,65,
	65,65,65,65,65,65,65,65,
}

// Test message for general message signing
var TestMessage = []byte("Hello, this is a test message for BRC-104 authentication")

// Timeout for operations (in seconds)
const TestTimeout = 10