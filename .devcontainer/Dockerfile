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
  jq

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

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
# Copy library scripts
COPY library-scripts/*.sh /tmp/library-scripts/
RUN chmod +x /tmp/library-scripts/*.sh

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
    "${CHECKOV_VERSION}"

USER node

# # Copy and set up firewall script
# COPY init-firewall.sh /usr/local/bin/
# USER root
# RUN chmod +x /usr/local/bin/init-firewall.sh && \
#   echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
#   chmod 0440 /etc/sudoers.d/node-firewall
# USER node
