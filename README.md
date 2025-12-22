# BRC-104 Cross-Language Signature Compatibility Test Suite

This test suite verifies that Python, TypeScript, and Go SDKs produce identical BRC-104 authentication signatures when given identical inputs.

## The Bug: TypeScript SDK Signature Data Preparation Error

### Problem Description

The TypeScript SDK had a critical bug in how it prepared signature data for BRC-104 authentication. This caused signature incompatibility between SDKs, preventing cross-implementation authentication.

### Root Cause

**Incorrect (TypeScript SDK - buggy):**
```typescript
// Concatenated base64 strings first, then decoded
data: Peer.base64ToBytes(message.initialNonce + sessionNonce)
```

**Correct (Go SDK - reference implementation):**
```go
// Decoded each nonce separately, then concatenated bytes
initialNonceBytes, _ := base64.StdEncoding.DecodeString(message.InitialNonce)
sessionNonceBytes, _ := base64.StdEncoding.DecodeString(session.SessionNonce)
sigData := append(initialNonceBytes, sessionNonceBytes...)
```

**Correct (Python SDK):**
```python
# Decoded each nonce separately, then concatenated bytes
initial_nonce_bytes = base64.b64decode(initial_nonce)
session_nonce_bytes = base64.b64decode(session_nonce)
return initial_nonce_bytes + session_nonce_bytes
```

### Why This Matters

The bug caused different byte arrays to be signed:
- **Buggy approach**: Decodes a concatenated base64 string (produces ~32-34 bytes)
- **Correct approach**: Concatenates individually decoded 32-byte arrays (produces 64 bytes)

This resulted in different signatures and failed authentication between TypeScript clients and Go/Python servers.

### Files Fixed

- `ts-sdk/src/auth/Peer.ts` (lines 528, 576-578)

### Before/After

**Before (buggy):**
```typescript
// Signing
data: Peer.base64ToBytes(message.initialNonce + sessionNonce)

// Verification
const dataToVerify = Peer.base64ToBytes(
  (peerSession.sessionNonce ?? '') + (message.initialNonce ?? '')
)
```

**After (fixed):**
```typescript
// Signing
data: [
  ...Peer.base64ToBytes(message.initialNonce),
  ...Peer.base64ToBytes(sessionNonce)
]

// Verification
const dataToVerify = [
  ...Peer.base64ToBytes(peerSession.sessionNonce ?? ''),
  ...Peer.base64ToBytes(message.initialNonce ?? '')
]
```

## Test Structure

```
cross-language-tests/
├── constants.py          # Python shared constants
├── constants.ts          # TypeScript shared constants
├── constants.go          # Go shared constants
├── test_brc104_signatures.py     # Python test suite
├── test_brc104_signatures.ts     # TypeScript test suite
├── test_brc104_signatures_test.go # Go test suite
├── package.json          # TypeScript dependencies
└── README.md            # This file
```

## Running Tests

### Prerequisites

- Python 3.8+
- Node.js 16+
- Go 1.19+
- pytest (for Python tests)
- jest (for TypeScript tests)

### Python Tests

```bash
cd cross-language-tests
python -m pytest test_brc104_signatures.py -v
```

### TypeScript Tests

```bash
cd cross-language-tests
npm install
npm test
```

### Go Tests

```bash
cd cross-language-tests
go test -v
```

## Test Results

All tests should pass after the fix:

- **Python**: 9/9 tests passing ✅
- **TypeScript**: 10/10 tests passing ✅
- **Go**: 9/9 tests passing ✅

## Test Coverage

### Signature Data Preparation
- Verifies correct nonce decoding and concatenation
- Tests both signing and verification data preparation
- Demonstrates the historical bug for documentation

### Cross-Language Compatibility
- Ensures identical constants across languages
- Validates nonce decoding correctness
- Tests key ID formatting

### Integration Tests
- Simulates complete BRC-104 authentication flows
- Tests both server signing and client verification paths

## Reference Implementation

The **Go SDK** is considered the authoritative reference implementation for BRC-104 authentication. All other SDKs should produce identical results when given identical inputs.

## BRC-104 Specification

This implementation follows [BRC-104: HTTP Transport for BRC-103 Mutual Authentication](https://github.com/bitcoin-sv/BRCs/blob/master/peer-to-peer/0104.md).

## Impact

This fix ensures:
- TypeScript clients can authenticate with Go servers
- TypeScript clients can authenticate with Python servers
- All SDKs maintain cryptographic compatibility
- Cross-implementation BRC-104 authentication works correctly

## Key Derivation Comparison Tests

The `test_key_derivation.*` files test key derivation functions across all three languages to identify any implementation differences that could cause signature verification failures.

### Running Key Derivation Tests

**Individual tests:**
```bash
# Python
python3 test_key_derivation.py

# TypeScript
npx ts-node test_key_derivation.ts

# Go
go run test_key_derivation_test.go constants.go
```

**Compare all languages:**
```bash
./compare_key_derivation.sh
```

### What the Tests Do

The key derivation tests:
1. Use identical root keys across all languages
2. Derive keys with the same inputs (protocol, keyID, counterparty)
3. Compare derived private keys
4. Compare derived public keys (with `forSelf=True` and `forSelf=False`)
5. Verify that public keys derived from private keys match `forSelf=True` results

### Expected Results

All three languages should produce:
- **Identical private keys** when given the same inputs
- **Identical public keys** for `forSelf=True` (should match public key from private key)
- **Identical public keys** for `forSelf=False` (different from `forSelf=True`)

If any differences are found, this indicates a bug in key derivation that could cause signature verification failures.
