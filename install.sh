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

KATO_DIR="$HOME/.kato"
INIT_SH="$KATO_DIR/init.sh"
INIT_FISH="$KATO_DIR/init.fish"

# Copy the kato-owned init scripts from the repo to ~/.kato/.
write_init_sh() {
    mkdir -p "$KATO_DIR"
    cp "$REPO_DIR/src/internal/shell/scripts/init.sh" "$INIT_SH"
    chmod 640 "$INIT_SH"
    echo "Wrote $INIT_SH"
}

write_init_fish() {
    mkdir -p "$KATO_DIR"
    cp "$REPO_DIR/src/internal/shell/scripts/init.fish" "$INIT_FISH"
    chmod 640 "$INIT_FISH"
    echo "Wrote $INIT_FISH"
}

# Inject (or replace) the kato block in a profile using BEGIN/END markers.
# Safe to run on every install — replaces the block in-place if present.
inject_kato_block() {
    local profile="$1"
    local init_script="$2"
    local block
    block="$(printf '# BEGIN kato\nsource "%s"\n# END kato' "$init_script")"

    mkdir -p "$(dirname "$profile")"
    if [ -f "$profile" ] && grep -q '# BEGIN kato' "$profile"; then
        python3 - "$profile" "$block" << 'PY'
import sys, re
path, block = sys.argv[1], sys.argv[2]
text = open(path).read()
text = re.sub(r'# BEGIN kato\n.*?# END kato', block, text, flags=re.DOTALL)
open(path, 'w').write(text)
PY
        echo "Updated kato block in $profile"
    else
        printf '\n%s\n' "$block" >> "$profile"
        echo "Added kato block to $profile"
    fi
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

# Write init scripts and wire them into shell profiles.
case "$shell_name" in
    zsh)
        write_init_sh
        inject_kato_block "$HOME/.zshrc" "$INIT_SH"
        ;;
    bash)
        write_init_sh
        if [ -f "$HOME/.bash_profile" ]; then
            inject_kato_block "$HOME/.bash_profile" "$INIT_SH"
        else
            inject_kato_block "$HOME/.bashrc" "$INIT_SH"
        fi
        ;;
    fish)
        write_init_fish
        fish_conf="$HOME/.config/fish/conf.d/kato.fish"
        inject_kato_block "$fish_conf" "$INIT_FISH"
        ;;
esac

echo "Done. Run 'kato --help' to get started."
