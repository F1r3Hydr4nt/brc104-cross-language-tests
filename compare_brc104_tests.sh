#!/bin/bash
# Compare BRC-104 signature test results across Python, TypeScript, and Go

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=========================================="
echo "BRC-104 Cross-Language Test Comparison"
echo "=========================================="
echo ""

# Run Python tests
echo "Running Python BRC-104 tests..."
PYTHON_PASSED=0
PYTHON_FAILED=0
if python3 -m pytest test_brc104_signatures.py -v > python_brc104_output.txt 2>&1; then
    # Extract from "9 passed in 0.02s" format
    PYTHON_PASSED=$(grep "passed" python_brc104_output.txt | grep -oE '[0-9]+ passed' | grep -oE '[0-9]+' | head -1 || echo "0")
    PYTHON_FAILED=$(grep "failed" python_brc104_output.txt | grep -oE '[0-9]+ failed' | grep -oE '[0-9]+' | head -1 || echo "0")
    PYTHON_STATUS="✅ PASSED"
    PYTHON_EXIT=0
else
    PYTHON_PASSED=$(grep "passed" python_brc104_output.txt | grep -oE '[0-9]+ passed' | grep -oE '[0-9]+' | head -1 || echo "0")
    PYTHON_FAILED=$(grep "failed" python_brc104_output.txt | grep -oE '[0-9]+ failed' | grep -oE '[0-9]+' | head -1 || echo "0")
    PYTHON_STATUS="❌ FAILED"
    PYTHON_EXIT=1
fi
echo "Python: $PYTHON_STATUS ($PYTHON_PASSED passed, $PYTHON_FAILED failed)"
echo ""

# Run TypeScript tests
echo "Running TypeScript BRC-104 tests..."
TYPESCRIPT_PASSED=0
TYPESCRIPT_FAILED=0
if npx jest test_brc104_signatures.ts > typescript_brc104_output.txt 2>&1; then
    # Extract from "Tests: 10 passed, 10 total" format
    TYPESCRIPT_PASSED=$(grep "Tests:" typescript_brc104_output.txt | grep -oE '[0-9]+ passed' | grep -oE '[0-9]+' | head -1 || echo "0")
    TYPESCRIPT_FAILED=$(grep "Tests:" typescript_brc104_output.txt | grep -oE '[0-9]+ failed' | grep -oE '[0-9]+' | head -1 || echo "0")
    TYPESCRIPT_STATUS="✅ PASSED"
    TYPESCRIPT_EXIT=0
else
    TYPESCRIPT_PASSED=$(grep "Tests:" typescript_brc104_output.txt | grep -oE '[0-9]+ passed' | grep -oE '[0-9]+' | head -1 || echo "0")
    TYPESCRIPT_FAILED=$(grep "Tests:" typescript_brc104_output.txt | grep -oE '[0-9]+ failed' | grep -oE '[0-9]+' | head -1 || echo "0")
    TYPESCRIPT_STATUS="❌ FAILED"
    TYPESCRIPT_EXIT=1
fi
echo "TypeScript: $TYPESCRIPT_STATUS ($TYPESCRIPT_PASSED passed, $TYPESCRIPT_FAILED failed)"
echo ""

# Run Go tests
echo "Running Go BRC-104 tests..."
GO_PASSED=0
GO_FAILED=0
if go test -v -run TestBRC104 > go_brc104_output.txt 2>&1; then
    # Count individual test passes (--- PASS: lines)
    GO_PASSED=$(grep -c -- "--- PASS:" go_brc104_output.txt || echo "0")
    GO_FAILED=$(grep -c -- "--- FAIL:" go_brc104_output.txt || echo "0")
    GO_STATUS="✅ PASSED"
    GO_EXIT=0
else
    GO_PASSED=$(grep -c -- "--- PASS:" go_brc104_output.txt || echo "0")
    GO_FAILED=$(grep -c -- "--- FAIL:" go_brc104_output.txt || echo "0")
    GO_STATUS="❌ FAILED"
    GO_EXIT=1
fi
echo "Go: $GO_STATUS ($GO_PASSED passed, $GO_FAILED failed)"
echo ""

echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo ""
echo "Python:    $PYTHON_STATUS"
echo "TypeScript: $TYPESCRIPT_STATUS"
echo "Go:        $GO_STATUS"
echo ""

# Show detailed results if any failed
if [ $PYTHON_EXIT -ne 0 ] || [ $TYPESCRIPT_EXIT -ne 0 ] || [ $GO_EXIT -ne 0 ]; then
    echo "=========================================="
    echo "Failed Test Details"
    echo "=========================================="
    echo ""
    
    if [ $PYTHON_EXIT -ne 0 ]; then
        echo "--- Python Test Failures ---"
        grep -A 5 "FAILED" python_brc104_output.txt | head -20 || echo "No detailed failure info"
        echo ""
    fi
    
    if [ $TYPESCRIPT_EXIT -ne 0 ]; then
        echo "--- TypeScript Test Failures ---"
        grep -A 5 "FAIL\|●" typescript_brc104_output.txt | head -20 || echo "No detailed failure info"
        echo ""
    fi
    
    if [ $GO_EXIT -ne 0 ]; then
        echo "--- Go Test Failures ---"
        grep -A 5 "FAIL:" go_brc104_output.txt | head -20 || echo "No detailed failure info"
        echo ""
    fi
fi

# Final result
echo "=========================================="
if [ $PYTHON_EXIT -eq 0 ] && [ $TYPESCRIPT_EXIT -eq 0 ] && [ $GO_EXIT -eq 0 ]; then
    echo "✅ ALL TESTS PASSED - BRC-104 signatures are compatible across languages!"
    exit 0
else
    echo "❌ SOME TESTS FAILED - Check output files for details:"
    echo "   - python_brc104_output.txt"
    echo "   - typescript_brc104_output.txt"
    echo "   - go_brc104_output.txt"
    exit 1
fi

