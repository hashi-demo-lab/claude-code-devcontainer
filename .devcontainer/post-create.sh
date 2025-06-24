#!/bin/bash
set -e

echo "Setting up GitHub CLI extensions..."

# Check if GH_TOKEN is set (can be passed via devcontainer.json)
if [ -n "$GH_TOKEN" ]; then
    echo "Authenticating with GitHub using provided token..."
    echo "$GH_TOKEN" | gh auth login --with-token
    
    # Install Claude Code extension for GitHub
    echo "Installing Claude GitHub extension..."
    gh extension install anthropic/gh-claude || echo "Failed to install GitHub extension, may need manual installation"
else
    echo "GH_TOKEN not set. GitHub CLI extensions will need to be installed manually."
    echo "To install manually after authentication:"
    echo "  gh extension install anthropic/gh-claude"
fi

echo "Setup complete!"
