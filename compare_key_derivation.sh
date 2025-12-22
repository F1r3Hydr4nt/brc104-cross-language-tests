#!/bin/bash
# Compare key derivation results across Python, TypeScript, and Go

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=========================================="
echo "Cross-Language Key Derivation Comparison"
echo "=========================================="
echo ""

# Run Python test
echo "Running Python test..."
python3 test_key_derivation.py > python_output.json 2>&1 || {
    echo "ERROR: Python test failed"
    cat python_output.json
    exit 1
}

# Run TypeScript test
echo "Running TypeScript test..."
cd "$SCRIPT_DIR"
npx ts-node test_key_derivation.ts > typescript_output.json 2>&1 || {
    echo "ERROR: TypeScript test failed"
    cat typescript_output.json
    exit 1
}

# Run Go test
echo "Running Go test..."
cd "$SCRIPT_DIR"
go run test_key_derivation_test.go constants.go > go_output.json 2>&1 || {
    echo "ERROR: Go test failed"
    cat go_output.json
    exit 1
}

echo ""
echo "=========================================="
echo "Extracting JSON results..."
echo "=========================================="
echo ""

# Extract JSON from outputs (they may have console.log mixed in)
python3 << 'PYTHON_SCRIPT'
import json
import re
import sys

def extract_json(filename):
    with open(filename, 'r') as f:
        content = f.read()
    
    # Try to find JSON object in the output
    # Look for lines between "JSON Output" and end of file or next section
    json_match = re.search(r'JSON Output.*?\n(\{.*\})', content, re.DOTALL)
    if json_match:
        json_str = json_match.group(1)
        try:
            return json.loads(json_str)
        except:
            pass
    
    # Fallback: try to parse entire file as JSON
    try:
        return json.loads(content)
    except:
        pass
    
    return None

python_data = extract_json('python_output.json')
typescript_data = extract_json('typescript_output.json')
go_data = extract_json('go_output.json')

if not python_data or not typescript_data or not go_data:
    print("ERROR: Failed to extract JSON from outputs")
    sys.exit(1)

print("Python Results:")
print(f"  Private Key: {python_data['results']['private_key_hex']}")
print(f"  Public Key (forSelf=True): {python_data['results']['public_key_for_self_hex']}")
print(f"  Public Key (forSelf=False): {python_data['results']['public_key_not_for_self_hex']}")
print()

print("TypeScript Results:")
print(f"  Private Key: {typescript_data['results']['privateKeyHex']}")
print(f"  Public Key (forSelf=true): {typescript_data['results']['publicKeyForSelfHex']}")
print(f"  Public Key (forSelf=false): {typescript_data['results']['publicKeyNotForSelfHex']}")
print()

print("Go Results:")
print(f"  Private Key: {go_data['results']['private_key_hex']}")
print(f"  Public Key (forSelf=true): {go_data['results']['public_key_for_self_hex']}")
print(f"  Public Key (forSelf=false): {go_data['results']['public_key_not_for_self_hex']}")
print()

print("=" * 80)
print("COMPARISON RESULTS")
print("=" * 80)
print()

# Compare private keys
priv_match = (python_data['results']['private_key_hex'] == 
              typescript_data['results']['privateKeyHex'] == 
              go_data['results']['private_key_hex'])
print(f"Private Keys Match: {priv_match}")
if not priv_match:
    print(f"  Python:    {python_data['results']['private_key_hex']}")
    print(f"  TypeScript: {typescript_data['results']['privateKeyHex']}")
    print(f"  Go:        {go_data['results']['private_key_hex']}")

# Compare public keys (forSelf=True)
pub_self_match = (python_data['results']['public_key_for_self_hex'] == 
                   typescript_data['results']['publicKeyForSelfHex'] == 
                   go_data['results']['public_key_for_self_hex'])
print(f"Public Keys (forSelf=True) Match: {pub_self_match}")
if not pub_self_match:
    print(f"  Python:    {python_data['results']['public_key_for_self_hex']}")
    print(f"  TypeScript: {typescript_data['results']['publicKeyForSelfHex']}")
    print(f"  Go:        {go_data['results']['public_key_for_self_hex']}")

# Compare public keys (forSelf=False)
pub_not_self_match = (python_data['results']['public_key_not_for_self_hex'] == 
                       typescript_data['results']['publicKeyNotForSelfHex'] == 
                       go_data['results']['public_key_not_for_self_hex'])
print(f"Public Keys (forSelf=False) Match: {pub_not_self_match}")
if not pub_not_self_match:
    print(f"  Python:    {python_data['results']['public_key_not_for_self_hex']}")
    print(f"  TypeScript: {typescript_data['results']['publicKeyNotForSelfHex']}")
    print(f"  Go:        {go_data['results']['public_key_not_for_self_hex']}")

print()
if priv_match and pub_self_match and pub_not_self_match:
    print("✅ ALL TESTS PASSED - Key derivation is consistent across languages!")
    sys.exit(0)
else:
    print("❌ TESTS FAILED - Key derivation differs between languages")
    sys.exit(1)
PYTHON_SCRIPT

