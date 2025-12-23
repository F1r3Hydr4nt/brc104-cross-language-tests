/**
 * Cross-language key derivation comparison test.
 * 
 * This test compares key derivation functions across Python, TypeScript, and Go
 * to identify any differences in implementation that could cause signature verification failures.
 */

import { KeyDeriver } from '../ts-sdk/src/wallet/KeyDeriver.js';
import { PrivateKey, PublicKey } from '../ts-sdk/src/primitives/index.js';
import {
  TEST_PRIVATE_KEY_WIF,
  INITIAL_NONCE_B64,
  SESSION_NONCE_B64,
  makeKeyId,
  TEST_COUNTERPARTY_KEY,
  PROTOCOL_ID
} from './constants';

function deriveKeysTypeScript(
  protocolID: [number, string],
  keyID: string,
  counterpartyHex: string,
  forSelf: boolean = false
): {
  privateKeyHex: string;
  publicKeyForSelfHex: string;
  publicKeyNotForSelfHex: string;
  publicKeyFromPrivateHex: string;
} {
  // Create root key from WIF
  const rootKey = PrivateKey.fromWif(TEST_PRIVATE_KEY_WIF);
  
  // Create key deriver
  const keyDeriver = new KeyDeriver(rootKey);
  
  // Create counterparty public key
  const counterpartyPub = PublicKey.fromString(counterpartyHex);
  
  // Derive private key
  const derivedPriv = keyDeriver.derivePrivateKey(protocolID as any, keyID, counterpartyPub);
  
  // Derive public keys
  const derivedPubForSelf = keyDeriver.derivePublicKey(protocolID as any, keyID, counterpartyPub, true);
  const derivedPubNotForSelf = keyDeriver.derivePublicKey(protocolID as any, keyID, counterpartyPub, false);
  
  return {
    privateKeyHex: derivedPriv.toHex(),
    publicKeyForSelfHex: derivedPubForSelf.toString(),
    publicKeyNotForSelfHex: derivedPubNotForSelf.toString(),
    publicKeyFromPrivateHex: derivedPriv.toPublicKey().toString()
  };
}

function testKeyDerivation(): void {
  // Test parameters matching BRC-104 general message scenario
  const protocolID = PROTOCOL_ID as [number, string];
  const keyID = makeKeyId(INITIAL_NONCE_B64, SESSION_NONCE_B64);
  const counterpartyHex = TEST_COUNTERPARTY_KEY;
  
  console.log('='.repeat(80));
  console.log('TYPESCRIPT KEY DERIVATION TEST');
  console.log('='.repeat(80));
  console.log(`Protocol: security_level=${protocolID[0]}, protocol='${protocolID[1]}'`);
  console.log(`Key ID: '${keyID}'`);
  console.log(`Counterparty: ${counterpartyHex}`);
  console.log();
  
  // Derive keys
  const results = deriveKeysTypeScript(protocolID, keyID, counterpartyHex);
  
  console.log('Derived Keys:');
  console.log(`  Private Key (hex): ${results.privateKeyHex}`);
  console.log(`  Public Key (from private): ${results.publicKeyFromPrivateHex}`);
  console.log(`  Public Key (forSelf=true): ${results.publicKeyForSelfHex}`);
  console.log(`  Public Key (forSelf=false): ${results.publicKeyNotForSelfHex}`);
  console.log();
  
  // Output JSON for comparison
  const output = {
    language: 'typescript',
    protocol: {
      security_level: protocolID[0],
      protocol: protocolID[1]
    },
    key_id: keyID,
    counterparty: counterpartyHex,
    results: results
  };
  
  console.log('JSON Output (for comparison):');
  console.log(JSON.stringify(output, null, 2));
  console.log();
  
  // Verify consistency
  console.log('Consistency Checks:');
  console.log(`  Public key from private matches forSelf=true: ${results.publicKeyFromPrivateHex === results.publicKeyForSelfHex}`);
  console.log(`  Public key from private matches forSelf=false: ${results.publicKeyFromPrivateHex === results.publicKeyNotForSelfHex}`);
  console.log();
}

// Run test if executed directly
// In ESM, check if this file is being run directly (not imported)
import { fileURLToPath } from 'url';
import { resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const isMainModule = process.argv[1] && resolve(process.argv[1]) === __filename;

if (isMainModule) {
  testKeyDerivation();
}

export { testKeyDerivation, deriveKeysTypeScript };

