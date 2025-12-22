"""
Shared constants for cross-language BRC-104 signature testing.

These constants are shared across Python, TypeScript, and Go test suites
to ensure identical inputs produce identical outputs.
"""

import base64

# Private key in WIF format (fixed for testing)
TEST_PRIVATE_KEY_WIF = "L4B2postXdaP7TiUrUBYs53Fqzheu7WhSoQVPuY8qBdoBeEwbmZx"

# Test nonces (32 bytes each, base64 encoded)
# These are fixed values to ensure deterministic testing
INITIAL_NONCE_B64 = "QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUE="  # 32 bytes of 'A' characters
SESSION_NONCE_B64 = "QkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkI="  # 32 bytes of 'B' characters

# Decode to bytes for direct use
INITIAL_NONCE_BYTES = base64.b64decode(INITIAL_NONCE_B64)
SESSION_NONCE_BYTES = base64.b64decode(SESSION_NONCE_B64)

# Protocol parameters for BRC-104 authentication
PROTOCOL_ID = [2, "auth message signature"]  # Security level 2, protocol name

# Key ID format: "{initial_nonce} {session_nonce}"
def make_key_id(initial_nonce: str, session_nonce: str) -> str:
    return f"{initial_nonce} {session_nonce}"

# Test counterparty (identity key) - using a fixed public key for testing
# This should be derived from TEST_PRIVATE_KEY_WIF, but for simplicity we'll use a fixed value
TEST_COUNTERPARTY_KEY = "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"  # Fixed test key

# Expected signature data (what should be signed)
# For signing: initial_nonce_bytes + session_nonce_bytes
EXPECTED_SIG_DATA_SIGNING = b'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB'
# For verification: session_nonce_bytes + initial_nonce_bytes
EXPECTED_SIG_DATA_VERIFICATION = b'BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA'

# Test message for general message signing
TEST_MESSAGE = b"Hello, this is a test message for BRC-104 authentication"

# Timeout for async operations (in seconds)
TEST_TIMEOUT = 10
