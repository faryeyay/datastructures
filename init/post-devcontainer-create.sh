#!/bin/bash

info() {
    echo "[info] " $@
}

fatal() {
    echo "[fatal] " $@
    exit 1
}

warn() {
    echo "[warn] " $@
}

install_dev_dependecies() {
    info "Installing Dev Dependencies"
    sudo apt install -y binutils build-essentials gnupg sshpass gnupg2 pinentry-curses openssl unzip zip curl wget lsb-release software-properties-common apt-transport-https ca-certificates jq unzip zsh >/dev/null
}

install_git() {
    info "Installing Git"
    sudo apt install -y git >/dev/null
}

configure_git() {
    info "Allowing for the container directory to be safe."
    git config --global --add safe.directory /workspaces/$DEVPOD_WORKSPACE_ID
}

install_and_configure_git() {
    install_git
    configure_git
}

install_aws() {
    info "Installing AWS CLI"
    curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install
    rm awscliv2.zip
    rm -rf aws
}

configure_aws() {
    info "Configuring AWS"
    mkdir -p ~/.aws/
    cp init/aws.config ~/.aws/config 
}

install_and_configure_aws() {
    install_aws
    configure_aws
}

install_kubectl() {
    info "Installing Kubernetes CLI"
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/arm64/kubectl"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
}

install_vault_cli() {
    info "Installing Vault CLI"
    wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
    sudo apt update -y && sudo apt install -y vault
    echo 'export VAULT_ADDR="https://vault.corp.hadrian-automation.com:8200"' >> ~/.bashrc
    cat init/vault.envrc >> $HOME/.envrc
}

install_kubectl_krew() {
    info "Installing Krew for Kubectl"
(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)
    echo 'export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"' >> ~/.bashrc
}

install_kubectl_kubens() {
    info "Installing Kubens"
    # The path is needed for Krew
    export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
    kubectl krew install ns
}

install_kubectl_kubectx() {
    info "Installing Kubectx"
    # The path is needed for Krew
    export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
    kubectl krew install ctx
}

install_homebrew() {
    info "Installing Homebrew"
    NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    echo >> /home/$USER/.bashrc
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/$USER/.bashrc

    ## Do the same for zshrc
    echo >> /home/$USER/.zshrc
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/$USER/.zshrc

    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
}

install_kubelogin() {
    info "Installing Kubelogin"
    brew install int128/kubelogin/kubelogin
}

install_k9s() {
    info "Installing K9s"
    brew install k9s
}

# Used to ensure that we're able to leverage Golang based tooling
add_golang_path() {
    echo 'export PATH="$PATH:$HOME/go/bin"' >> /home/$USER/.bashrc
    echo 'export PATH="$PATH:$HOME/go/bin"' >> /home/$USER/.zshrc
}

install_git_crypt() {
    info "Setting up git-crypt"
    brew install git-crypt
}

init_git_crypt() {
    git-crypt init
}

setup_dependencies() {
    echo "Setting up dependencies..."
    install_dev_dependencies
    install_and_configure_git
    install_and_configure_aws
    install_homebrew
    install_kubectl
    install_vault_cli
    install_kubectl_krew
    install_kubectl_kubens
    install_kubectl_kubectx
    install_kubelogin
    install_k9s
    install_git_crypt
    init_git_crypt
}

setup_dependencies
add_golang_path