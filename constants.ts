/**
 * Shared constants for cross-language BRC-104 signature testing.
 *
 * These constants are shared across Python, TypeScript, and Go test suites
 * to ensure identical inputs produce identical outputs.
 */

// Private key in WIF format (fixed for testing)
export const TEST_PRIVATE_KEY_WIF = "L4B2postXdaP7TiUrUBYs53Fqzheu7WhSoQVPuY8qBdoBeEwbmZx";

// Test nonces (32 bytes each, base64 encoded)
// These are fixed values to ensure deterministic testing
export const INITIAL_NONCE_B64 = "QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUE="; // 32 bytes of 'A' characters
export const SESSION_NONCE_B64 = "QkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkI="; // 32 bytes of 'B' characters

// Decode to byte arrays for direct use
export const INITIAL_NONCE_BYTES: number[] = [65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65];
export const SESSION_NONCE_BYTES: number[] = [66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66];

// Protocol parameters for BRC-104 authentication
export const PROTOCOL_ID = [2, "auth message signature"]; // Security level 2, protocol name

// Key ID format: "{initial_nonce} {session_nonce}"
export function makeKeyId(initialNonce: string, sessionNonce: string): string {
  return `${initialNonce} ${sessionNonce}`;
}

// Test counterparty (identity key) - using a fixed public key for testing
// This should be derived from TEST_PRIVATE_KEY_WIF, but for simplicity we'll use a fixed value
export const TEST_COUNTERPARTY_KEY = "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"; // Fixed test key

// Expected signature data (what should be signed)
// For signing: initial_nonce_bytes + session_nonce_bytes
export const EXPECTED_SIG_DATA_SIGNING: number[] = [65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66];
// For verification: session_nonce_bytes + initial_nonce_bytes
export const EXPECTED_SIG_DATA_VERIFICATION: number[] = [66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65,65];

// Test message for general message signing
export const TEST_MESSAGE: number[] = Array.from(Buffer.from("Hello, this is a test message for BRC-104 authentication", 'utf8'));

// Timeout for async operations (in milliseconds)
export const TEST_TIMEOUT = 10000;