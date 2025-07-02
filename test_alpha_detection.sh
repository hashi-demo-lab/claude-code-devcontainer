#!/bin/bash

# Test script to verify alpha detection logic
TERRAFORM_ALPHA=${1:-"false"}

echo "Testing alpha detection with TERRAFORM_ALPHA=$TERRAFORM_ALPHA"

if [ "$TERRAFORM_ALPHA" = "true" ]; then
    echo "Alpha mode detected"
    
    # Get the directory where this script is located
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
    ALPHA_DIR="${SCRIPT_DIR}/.devcontainer/library-scripts/alpha"
    
    echo "Looking in directory: $ALPHA_DIR"
    
    # Find terraform alpha binary
    if [ -f "${ALPHA_DIR}/terraform_"*".zip" ]; then
        TERRAFORM_ALPHA_ZIP=$(ls "${ALPHA_DIR}"/terraform_*.zip | head -1)
        echo "Found Terraform alpha binary: $(basename $TERRAFORM_ALPHA_ZIP)"
    else
        echo "Error: No Terraform alpha binary found in ${ALPHA_DIR}"
    fi
    
    # Find tfpolicy alpha binary
    if [ -f "${ALPHA_DIR}/tfpolicy_"*".zip" ]; then
        TFPOLICY_ALPHA_ZIP=$(ls "${ALPHA_DIR}"/tfpolicy_*.zip | head -1)
        echo "Found TFPolicy alpha binary: $(basename $TFPOLICY_ALPHA_ZIP)"
    else
        echo "Warning: No TFPolicy alpha binary found in ${ALPHA_DIR}"
    fi
else
    echo "Regular mode - would download from HashiCorp releases"
fi