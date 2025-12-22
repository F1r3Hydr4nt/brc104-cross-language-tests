# Key Derivation Bug Report

## Comparison Results

### ✅ TypeScript vs Go: **MATCH** (Both Correct)

**TypeScript Results:**
- Private Key: `26bbbb010eab5d1ae2666c851c145dcb45f5e926d63a8e81a153d2b6ed8ab2fc`
- Public Key (from private): `024467ed1ba7c9d46ad972fe31c090be264e8892c8d12dcb5441d7c0c307ad99dc`
- Public Key (forSelf=true): `024467ed1ba7c9d46ad972fe31c090be264e8892c8d12dcb5441d7c0c307ad99dc` ✅
- Public Key (forSelf=false): `02528a60dee9a4124298e5db896185abf248962487d1da43bf414f7e9948fa6ab3`
- **Consistency**: Public key from private **MATCHES** forSelf=true ✅

**Go Results:**
- Private Key: `26bbbb010eab5d1ae2666c851c145dcb45f5e926d63a8e81a153d2b6ed8ab2fc` ✅ MATCHES
- Public Key (from private): `024467ed1ba7c9d46ad972fe31c090be264e8892c8d12dcb5441d7c0c307ad99dc` ✅ MATCHES
- Public Key (forSelf=true): `024467ed1ba7c9d46ad972fe31c090be264e8892c8d12dcb5441d7c0c307ad99dc` ✅ MATCHES
- Public Key (forSelf=false): `02528a60dee9a4124298e5db896185abf248962487d1da43bf414f7e9948fa6ab3` ✅ MATCHES
- **Consistency**: Public key from private **MATCHES** forSelf=true ✅

### ❌ Python: **BUG FOUND**

**Python Results:**
- Private Key: `26bbbb010eab5d1ae2666c851c145dcb45f5e926d63a8e81a153d2b6ed8ab2fc` ✅ MATCHES
- Public Key (from private): `024467ed1ba7c9d46ad972fe31c090be264e8892c8d12dcb5441d7c0c307ad99dc` ✅ MATCHES
- Public Key (forSelf=True): `02e0a7c5cde04b4404e7036b1c564a6b26bcc5f00802622c24864ad883a8efbccb` ❌ **WRONG**
- Public Key (forSelf=False): `02528a60dee9a4124298e5db896185abf248962487d1da43bf414f7e9948fa6ab3` ✅ MATCHES
- **Consistency**: Public key from private **DOES NOT MATCH** forSelf=True ❌

**Expected forSelf=True**: `024467ed1ba7c9d46ad972fe31c090be264e8892c8d12dcb5441d7c0c307ad99dc`
**Actual forSelf=True**: `02e0a7c5cde04b4404e7036b1c564a6b26bcc5f00802622c24864ad883a8efbccb`

## Root Cause Analysis

### Correct Implementation (Go/TypeScript)

**Go SDK:**
```go
if forSelf {
    privKey, err := kd.rootKey.DeriveChild(counterpartyKey, invoiceNumber)
    return privKey.PubKey(), nil
}
```

**TypeScript SDK:**
```typescript
if (forSelf) {
    return this.rootKey
        .deriveChild(counterparty, this.computeInvoiceNumber(protocolID, keyID))
        .toPublicKey()
}
```

**What they do:**
1. Derive the **private key** first: `rootKey.deriveChild(counterparty, invoice)`
2. Then get the **public key** from that derived private key: `derivedPriv.toPublicKey()`

### Buggy Implementation (Python)

**Python SDK (CURRENT - WRONG):**
```python
if for_self:
    # Determine counterparty pub used for tweak
    cp_pub = counterparty.to_public_key(self._root_public_key) if not for_self else self._root_public_key
    delta = self._branch_scalar(invoice_number, cp_pub)
    
    # forSelf=True: derived from root's perspective
    # tweaked public = root_pub + delta*G
    delta_point = curve_multiply(delta, curve.g)
    new_point = curve_add(self._root_public_key.point(), delta_point)
    return PublicKey(new_point)
```

**What it does (WRONG):**
1. When `for_self=True`, sets `cp_pub = self._root_public_key` (line 159)
2. Computes `delta = _branch_scalar(invoice_number, self._root_public_key)` - **WRONG**: should use counterparty's public key
3. Computes `root_pub + delta*G` directly - **WRONG**: should derive private key first, then get public key

**The Problem:**
- When `for_self=True`, Python uses `self._root_public_key` to compute the branch scalar, but it should use the **counterparty's public key**
- Python tries to compute the public key directly as `root_pub + delta*G`, but it should:
  1. Derive the private key: `root_priv + delta`
  2. Then get the public key from that private key

## Fix Required

### File: `py-sdk/bsv/wallet/key_deriver.py`

### Location: `derive_public_key` method, lines 141-173

### Current Code (WRONG):
```python
def derive_public_key(
    self,
    protocol: Protocol,
    key_id: str,
    counterparty: Counterparty,
    for_self: bool = False,
) -> PublicKey:
    invoice_number = self.compute_invoice_number(protocol, key_id)
    # Determine counterparty pub used for tweak
    cp_pub = counterparty.to_public_key(self._root_public_key) if not for_self else self._root_public_key
    delta = self._branch_scalar(invoice_number, cp_pub)
    
    if for_self:
        # forSelf=True: derived from root's perspective
        # tweaked public = root_pub + delta*G
        delta_point = curve_multiply(delta, curve.g)
        new_point = curve_add(self._root_public_key.point(), delta_point)
    else:
        # forSelf=False: derived from counterparty's perspective
        # tweaked public = cp_pub + delta*G
        delta_point = curve_multiply(delta, curve.g)
        new_point = curve_add(cp_pub.point(), delta_point)
    
    return PublicKey(new_point)
```

### Fixed Code (CORRECT):
```python
def derive_public_key(
    self,
    protocol: Protocol,
    key_id: str,
    counterparty: Counterparty,
    for_self: bool = False,
) -> PublicKey:
    """Derives a public key based on protocol ID, key ID, and counterparty.
    
    This implementation matches TypeScript/Go SDK behavior:
    - forSelf=True: rootKey.deriveChild(counterparty, invoice).toPublicKey()
    - forSelf=False: counterparty.deriveChild(rootKey, invoice)
    """
    invoice_number = self.compute_invoice_number(protocol, key_id)
    
    if for_self:
        # forSelf=True: Derive private key first, then get public key
        # This matches Go: privKey = rootKey.DeriveChild(counterparty, invoice); return privKey.PubKey()
        # This matches TS: rootKey.deriveChild(counterparty, invoice).toPublicKey()
        cp_pub = counterparty.to_public_key(self._root_public_key)
        delta = self._branch_scalar(invoice_number, cp_pub)
        derived_priv = PrivateKey((self._root_private_key.int() + delta) % CURVE_ORDER)
        return derived_priv.public_key()
    else:
        # forSelf=False: derived from counterparty's perspective
        # tweaked public = cp_pub + delta*G
        cp_pub = counterparty.to_public_key(self._root_public_key)
        delta = self._branch_scalar(invoice_number, cp_pub)
        delta_point = curve_multiply(delta, curve.g)
        new_point = curve_add(cp_pub.point(), delta_point)
        return PublicKey(new_point)
```

### Key Changes:
1. **Line 159**: Remove the ternary that sets `cp_pub = self._root_public_key` when `for_self=True`
2. **Lines 162-166**: Replace the direct public key computation with:
   - Use counterparty's public key to compute branch scalar
   - Derive the private key: `root_priv + delta`
   - Get the public key from the derived private key

## Impact

This bug causes signature verification failures when:
- Server uses `forSelf=True` to verify signatures
- Any code path that relies on `derive_public_key(forSelf=True)` matching `derive_private_key().public_key()`

The bug is currently masked because:
- Most verification uses `forSelf=False` (the default)
- The BRC-104 general message verification uses `forSelf=False` with `counterparty=sender_public_key`

However, if any code tries to use `forSelf=True`, it will fail.

