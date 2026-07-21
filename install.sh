#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="$HOME/.kato/bin"
REPO_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION=$(cat "$REPO_DIR/VERSION" | tr -d '\r\n')
COMMIT=$(git -C "$REPO_DIR" rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "Building kato v${VERSION}+${COMMIT}..."
cd "$REPO_DIR/src"
go build -ldflags "-s -w -X kato/cmd.version=${VERSION} -X kato/cmd.commit=${COMMIT}" -o "$REPO_DIR/kato" .

echo "Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
rm -f "$INSTALL_DIR/kg"
mv "$REPO_DIR/kato" "$INSTALL_DIR/kato"
chmod +x "$INSTALL_DIR/kato"

add_to_path() {
    local profile="$1"
    if [ -f "$profile" ] && grep -q '.kato/bin' "$profile"; then
        return
    fi
    echo '' >> "$profile"
    echo '# Kato CLI' >> "$profile"
    echo 'export PATH="$HOME/.kato/bin:$PATH"' >> "$profile"
    echo "Added to $profile"
}

add_shell_integration() {
    local profile="$1"
    if [ -f "$profile" ] && grep -q '# kato shell integration' "$profile"; then
        return
    fi
    cat >> "$profile" << 'EOF'

# kato shell integration — enables kato nav to change directory
kato() {
  if [ "$1" = "nav" ]; then
    local dir
    dir=$(command kato "$@")
    [ -n "$dir" ] && cd "$dir"
  else
    command kato "$@"
  fi
}
EOF
    echo "Added kato shell integration to $profile"
}

shell_name="$(basename "${SHELL:-/bin/bash}")"

if echo "$PATH" | tr ':' '\n' | grep -q "$INSTALL_DIR"; then
    echo "PATH already contains $INSTALL_DIR"
else
    case "$shell_name" in
        zsh)  add_to_path "$HOME/.zshrc" ;;
        bash)
            if [ -f "$HOME/.bash_profile" ]; then
                add_to_path "$HOME/.bash_profile"
            else
                add_to_path "$HOME/.bashrc"
            fi
            ;;
        fish)
            fish_conf="$HOME/.config/fish/conf.d/kato.fish"
            if [ ! -f "$fish_conf" ]; then
                mkdir -p "$(dirname "$fish_conf")"
                echo 'fish_add_path $HOME/.kato/bin' > "$fish_conf"
                echo "Added to $fish_conf"
            fi
            ;;
        *)
            echo "Unknown shell '$shell_name'. Add $INSTALL_DIR to your PATH manually."
            ;;
    esac
    echo "Restart your shell or run: export PATH=\"$INSTALL_DIR:\$PATH\""
fi

# Install shell integration (idempotent — skipped if already present).
case "$shell_name" in
    zsh)  add_shell_integration "$HOME/.zshrc" ;;
    bash)
        if [ -f "$HOME/.bash_profile" ]; then
            add_shell_integration "$HOME/.bash_profile"
        else
            add_shell_integration "$HOME/.bashrc"
        fi
        ;;
    fish)
        fish_fn="$HOME/.config/fish/functions/kato.fish"
        if [ ! -f "$fish_fn" ] || ! grep -q '# kato shell integration' "$fish_fn"; then
            mkdir -p "$(dirname "$fish_fn")"
            cat > "$fish_fn" << 'EOF'
# kato shell integration — enables kato nav to change directory
function kato
  if test "$argv[1]" = "nav"
    set dir (command kato $argv)
    and cd $dir
  else
    command kato $argv
  end
end
EOF
            echo "Added kato shell integration to $fish_fn"
        fi
        ;;
esac

echo "Done. Run 'kato --help' to get started."
