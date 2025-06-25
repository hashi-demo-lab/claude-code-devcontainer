# Environment Variable Management in Dev Containers

This project uses a multi-layered approach to manage environment variables in the development container:

## 1. Non-sensitive Variables in devcontainer.json

Standard, non-sensitive environment variables are defined directly in `devcontainer.json` under the `remoteEnv` property:

```json
"remoteEnv": {
  "NODE_OPTIONS": "--max-old-space-size=4096",
  "CLAUDE_CONFIG_DIR": "/home/node/.claude",
  "POWERLEVEL9K_DISABLE_GITSTATUS": "true"
}
```

These variables are version-controlled and shared across all developers.

## 2. Sensitive Variables in devcontainer.env

Sensitive variables (API keys, tokens, etc.) are stored in a separate `.devcontainer/devcontainer.env` file that is excluded from Git:

```
SECRET_TOKEN=your_secret_token
API_KEY=your_api_key
```

This file is loaded via the `--env-file` argument in `devcontainer.json`:

```json
"runArgs": [
  "--env-file", "${localWorkspaceFolder}/.devcontainer/devcontainer.env"
]
```

## 3. User-specific Variables (Optional)

For variables that are specific to individual developers or require frequent changes without container rebuilds, you can use `.devcontainer/.env.local`:

1. Copy `.devcontainer/.env.local.template` to `.devcontainer/.env.local`
2. Add your custom variables
3. These will be loaded by the `env-setup.sh` script during container creation

## Getting Started

When setting up this project for the first time:

1. Create your own copy of `devcontainer.env` based on team documentation
2. (Optional) Create `.env.local` if you need user-specific variables
3. Make sure both files are listed in your `.gitignore`
4. Rebuild your container

## Best Practices

- Only put non-sensitive variables in `devcontainer.json`
- Never commit sensitive credentials to Git
- Document required environment variables for new team members
- Use `.env.local` for temporary or user-specific overrides
