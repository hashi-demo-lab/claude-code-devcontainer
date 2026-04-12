FROM node:20

ARG TZ
ENV TZ="$TZ"
# Set versions as arguments
ARG TERRAFORM_VERSION="1.12.1"
ARG TERRAFORM_DOCS_VERSION="0.20.0"
ARG TFSEC_VERSION="1.28.13"
ARG TERRASCAN_VERSION="1.19.9"
ARG TFLINT_VERSION="0.48.0"
ARG TFLINT_AWS_RULESET_VERSION="0.23.1"
ARG TFLINT_AZURE_RULESET_VERSION="0.23.0"
ARG TFLINT_GCP_RULESET_VERSION="0.23.1"
ARG INFRACOST_VERSION="0.10.41"
ARG CHECKOV_VERSION="3.2.439"
ARG TERRAFORM_ALPHA="false"
# Install basic development tools and iptables/ipset
RUN apt update && apt install -y less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  python3 \
  python3-pip \
  python3-venv \
  python3-full \
  jq \
  curl

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share && \
  rm -rf /var/lib/apt/lists/*

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/node/.claude && \
  chown -R node:node /workspace /home/node/.claude

WORKDIR /workspace

RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/0.18.2/git-delta_0.18.2_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_0.18.2_${ARCH}.deb" && \
  rm "git-delta_0.18.2_${ARCH}.deb"

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Default powerline10k theme
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.0/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Install Claude
RUN npm install -g @anthropic-ai/claude-code

USER root
# Install uv system-wide before terraform tools (needed by terraform-tools.sh)
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/root/.local/bin:$PATH"

# Create wheel group if it doesn't exist and add node user to it
RUN groupadd -f wheel && usermod -aG wheel node

# Add the wheel group to sudoers with NOPASSWD
RUN echo "%wheel ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/wheel-group && \
    chmod 0440 /etc/sudoers.d/wheel-group

# Copy library scripts
COPY library-scripts/*.sh /tmp/library-scripts/
# COPY ALPHA library files (directory may be empty)
COPY library-scripts/alpha/ /tmp/library-scripts/alpha/
RUN chmod +x /tmp/library-scripts/*.sh

# Install cloud CLI tools
RUN /tmp/library-scripts/cloud-cli-tools.sh

RUN /tmp/library-scripts/terraform-tools.sh \
    "${TERRAFORM_VERSION}" \
    "${TERRAFORM_DOCS_VERSION}" \
    "${TFSEC_VERSION}" \
    "${TERRASCAN_VERSION}" \
    "${TFLINT_VERSION}" \
    "${TFLINT_AWS_RULESET_VERSION}" \
    "${TFLINT_AZURE_RULESET_VERSION}" \
    "${TFLINT_GCP_RULESET_VERSION}" \
    "${INFRACOST_VERSION}" \
    "${CHECKOV_VERSION}" \
    "${TERRAFORM_ALPHA}"


RUN mkdir -p /home/node/.config && chown -R node:node /home/node/.config
# Switch back to node user and install uv for them too
USER node
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

# Add uv tools to node user's PATH
RUN echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc

# Install pre-commit for the node user so they can use it
RUN ~/.local/bin/uv tool install pre-commit --with pre-commit-uv

# ── Go runtime ────────────────────────────────────────────────────────────────
USER root
ARG GO_VERSION=1.23.4
RUN ARCH=$(dpkg --print-architecture) && \
    case "$ARCH" in \
        amd64) GOARCH="amd64" ;; \
        arm64) GOARCH="arm64" ;; \
        *) echo "Unsupported arch: $ARCH" && exit 1 ;; \
    esac && \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz" \
        | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"
RUN mkdir -p /home/node/go && chown -R node:node /home/node/go

# Go tools (installed as node user so they land in /home/node/go/bin)
USER node
ENV GOPATH="/home/node/go"
ENV PATH="${GOPATH}/bin:${PATH}"
RUN echo 'export GOPATH="$HOME/go"' >> ~/.zshrc && \
    echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.zshrc
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4 && \
    go install github.com/securego/gosec/v2/cmd/gosec@latest

# ── Ansible tooling ───────────────────────────────────────────────────────────
# ansible-lint  : used at runtime by the validate_playbook MCP tool
# molecule      : used at runtime by the generate_test_cases MCP tool
USER root
RUN python3 -m venv /opt/ansible-venv && \
    /opt/ansible-venv/bin/pip install --upgrade pip && \
    /opt/ansible-venv/bin/pip install ansible-lint molecule "molecule-plugins[docker]" && \
    ln -sf /opt/ansible-venv/bin/ansible-lint /usr/local/bin/ansible-lint && \
    ln -sf /opt/ansible-venv/bin/molecule /usr/local/bin/molecule
USER node