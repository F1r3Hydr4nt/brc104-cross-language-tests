package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/bsv-blockchain/go-sdk/wallet"
)

type KeyDerivationResult struct {
	PrivateKeyHex            string `json:"private_key_hex"`
	PublicKeyForSelfHex      string `json:"public_key_for_self_hex"`
	PublicKeyNotForSelfHex   string `json:"public_key_not_for_self_hex"`
	PublicKeyFromPrivateHex  string `json:"public_key_from_private_hex"`
}

type TestOutput struct {
	Language    string                 `json:"language"`
	Protocol    map[string]interface{} `json:"protocol"`
	KeyID       string                 `json:"key_id"`
	Counterparty string                `json:"counterparty"`
	Results     KeyDerivationResult    `json:"results"`
}

func deriveKeysGo(protocol wallet.Protocol, keyID string, counterpartyHex string, forSelf bool) KeyDerivationResult {
	// Create root key from WIF
	rootKey, err := ec.PrivateKeyFromWif(TestPrivateKeyWIF)
	if err != nil {
		panic(fmt.Sprintf("Failed to create root key: %v", err))
	}

	// Create key deriver
	keyDeriver := wallet.NewKeyDeriver(rootKey)

	// Create counterparty public key
	counterpartyPub, err := ec.PublicKeyFromString(counterpartyHex)
	if err != nil {
		panic(fmt.Sprintf("Failed to create counterparty key: %v", err))
	}
	counterparty := wallet.Counterparty{
		Type:         wallet.CounterpartyTypeOther,
		Counterparty: counterpartyPub,
	}

	// Derive private key
	derivedPriv, err := keyDeriver.DerivePrivateKey(protocol, keyID, counterparty)
	if err != nil {
		panic(fmt.Sprintf("Failed to derive private key: %v", err))
	}

	// Derive public keys
	derivedPubForSelf, err := keyDeriver.DerivePublicKey(protocol, keyID, counterparty, true)
	if err != nil {
		panic(fmt.Sprintf("Failed to derive public key (forSelf=true): %v", err))
	}

	derivedPubNotForSelf, err := keyDeriver.DerivePublicKey(protocol, keyID, counterparty, false)
	if err != nil {
		panic(fmt.Sprintf("Failed to derive public key (forSelf=false): %v", err))
	}

	return KeyDerivationResult{
		PrivateKeyHex:            hex.EncodeToString(derivedPriv.Serialize()),
		PublicKeyForSelfHex:      derivedPubForSelf.ToDERHex(),
		PublicKeyNotForSelfHex:   derivedPubNotForSelf.ToDERHex(),
		PublicKeyFromPrivateHex:  derivedPriv.PubKey().ToDERHex(),
	}
}

func TestKeyDerivation() {
	// Test parameters matching BRC-104 general message scenario
	protocol := wallet.Protocol{
		SecurityLevel: wallet.SecurityLevelEveryAppAndCounterparty, // 2
		Protocol:      "auth message signature",
	}
	keyID := MakeKeyId(InitialNonceB64, SessionNonceB64)
	counterpartyHex := TestCounterpartyKey

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("GO KEY DERIVATION TEST")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Protocol: security_level=%d, protocol='%s'\n", protocol.SecurityLevel, protocol.Protocol)
	fmt.Printf("Key ID: '%s'\n", keyID)
	fmt.Printf("Counterparty: %s\n", counterpartyHex)
	fmt.Println()

	// Derive keys
	results := deriveKeysGo(protocol, keyID, counterpartyHex, false)

	fmt.Println("Derived Keys:")
	fmt.Printf("  Private Key (hex): %s\n", results.PrivateKeyHex)
	fmt.Printf("  Public Key (from private): %s\n", results.PublicKeyFromPrivateHex)
	fmt.Printf("  Public Key (forSelf=true): %s\n", results.PublicKeyForSelfHex)
	fmt.Printf("  Public Key (forSelf=false): %s\n", results.PublicKeyNotForSelfHex)
	fmt.Println()

	// Output JSON for comparison
	output := TestOutput{
		Language: "go",
		Protocol: map[string]interface{}{
			"security_level": protocol.SecurityLevel,
			"protocol":       protocol.Protocol,
		},
		KeyID:        keyID,
		Counterparty: counterpartyHex,
		Results:      results,
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}

	fmt.Println("JSON Output (for comparison):")
	fmt.Println(string(jsonOutput))
	fmt.Println()

	// Verify consistency
	fmt.Println("Consistency Checks:")
	fmt.Printf("  Public key from private matches forSelf=true: %v\n",
		results.PublicKeyFromPrivateHex == results.PublicKeyForSelfHex)
	fmt.Printf("  Public key from private matches forSelf=false: %v\n",
		results.PublicKeyFromPrivateHex == results.PublicKeyNotForSelfHex)
	fmt.Println()
}

func main() {
	TestKeyDerivation()
}

