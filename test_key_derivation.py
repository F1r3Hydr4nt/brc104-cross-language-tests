"""
Cross-language key derivation comparison test.

This test compares key derivation functions across Python, TypeScript, and Go
to identify any differences in implementation that could cause signature verification failures.
"""

import sys
import os
import json
import base64

# Add py-sdk to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'py-sdk'))

from bsv.wallet.key_deriver import KeyDeriver, Protocol, Counterparty, CounterpartyType
from bsv.keys import PrivateKey, PublicKey
from test_constants import (
    TEST_PRIVATE_KEY_WIF,
    INITIAL_NONCE_B64,
    SESSION_NONCE_B64,
    make_key_id,
    TEST_COUNTERPARTY_KEY
)


def derive_keys_python(protocol, key_id, counterparty_hex, for_self=False):
    """Derive keys using Python SDK."""
    # Create root key from WIF
    root_key = PrivateKey(TEST_PRIVATE_KEY_WIF)
    
    # Create key deriver
    key_deriver = KeyDeriver(root_key)
    
    # Create counterparty public key
    counterparty_pub = PublicKey(counterparty_hex)
    counterparty = Counterparty(CounterpartyType.OTHER, counterparty_pub)
    
    # Derive private key
    derived_priv = key_deriver.derive_private_key(protocol, key_id, counterparty)
    
    # Derive public keys
    derived_pub_for_self = key_deriver.derive_public_key(protocol, key_id, counterparty, for_self=True)
    derived_pub_not_for_self = key_deriver.derive_public_key(protocol, key_id, counterparty, for_self=False)
    
    return {
        'private_key_hex': derived_priv.hex(),
        'public_key_for_self_hex': derived_pub_for_self.hex(),
        'public_key_not_for_self_hex': derived_pub_not_for_self.hex(),
        'public_key_from_private_hex': derived_priv.public_key().hex()
    }


def test_key_derivation():
    """Test key derivation and output results for comparison."""
    # Test parameters matching BRC-104 general message scenario
    protocol = Protocol(security_level=2, protocol="auth message signature")
    key_id = make_key_id(INITIAL_NONCE_B64, SESSION_NONCE_B64)
    counterparty_hex = TEST_COUNTERPARTY_KEY
    
    print("=" * 80)
    print("PYTHON KEY DERIVATION TEST")
    print("=" * 80)
    print(f"Protocol: security_level={protocol.security_level}, protocol='{protocol.protocol}'")
    print(f"Key ID: '{key_id}'")
    print(f"Counterparty: {counterparty_hex}")
    print()
    
    # Derive keys
    results = derive_keys_python(protocol, key_id, counterparty_hex)
    
    print("Derived Keys:")
    print(f"  Private Key (hex): {results['private_key_hex']}")
    print(f"  Public Key (from private): {results['public_key_from_private_hex']}")
    print(f"  Public Key (forSelf=True): {results['public_key_for_self_hex']}")
    print(f"  Public Key (forSelf=False): {results['public_key_not_for_self_hex']}")
    print()
    
    # Output JSON for comparison
    output = {
        'language': 'python',
        'protocol': {
            'security_level': protocol.security_level,
            'protocol': protocol.protocol
        },
        'key_id': key_id,
        'counterparty': counterparty_hex,
        'results': results
    }
    
    print("JSON Output (for comparison):")
    print(json.dumps(output, indent=2))
    print()
    
    # Verify consistency
    print("Consistency Checks:")
    print(f"  Public key from private matches forSelf=True: {results['public_key_from_private_hex'] == results['public_key_for_self_hex']}")
    print(f"  Public key from private matches forSelf=False: {results['public_key_from_private_hex'] == results['public_key_not_for_self_hex']}")
    print()
    
    return output


if __name__ == '__main__':
    test_key_derivation()

