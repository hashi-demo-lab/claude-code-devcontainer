#!/bin/bash

# Environment variable setup script for dev containers
# This script supplements the environment variable handling in the devcontainer.json
#
# Current environment variable strategy:
# 1. Non-sensitive variables: Defined directly in devcontainer.json under remoteEnv
# 2. Sensitive variables: Stored in gitignored .devcontainer/devcontainer.env 
#    and loaded via --env-file in runArgs
# 3. User-specific variables: Can be loaded from .env.local using this script
#
# This script provides an additional method to load environment variables
# that can be updated without rebuilding the container

# Check if user-specific environment variables file exists
ENV_FILE="/workspace/.devcontainer/.env.local"
if [ -f "$ENV_FILE" ]; then
  echo "Loading user-specific environment variables from $ENV_FILE"
  # Load environment variables and add them to ~/.zshrc
  while IFS='=' read -r key value; do
    # Skip comments and empty lines
    [[ $key == \#* ]] && continue
    [[ -z "$key" ]] && continue
    
    # Remove quotes if present
    value=$(echo "$value" | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")
    
    # Add to shell configuration
    echo "export $key=\"$value\"" >> /home/node/.zshrc
  done < "$ENV_FILE"
  
  echo "User-specific environment variables have been set up"
else
  echo "No .env.local file found, using only devcontainer.json and devcontainer.env variables"
  echo "Create .devcontainer/.env.local for additional user-specific variables"
fi

# Note: This script is optional and supplements the primary environment variable
# configuration in devcontainer.json (remoteEnv) and devcontainer.env (--env-file)

# Other initialization tasks can be added below this line
