#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="filepass"
BIN_DIR="$HOME/.local/bin"
BASHRC="$HOME/.bashrc"
DIST_DIR="$(pwd)/dist"
BUILD_OUTPUT="$DIST_DIR/$BINARY_NAME"
INSTALL_TARGET="$BIN_DIR/$BINARY_NAME"

# ── Colours ────────────────────────────────────────────────────────────────────
BOLD='\033[1m'
DIM='\033[2m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
RESET='\033[0m'

step()    { echo -e "\n${BOLD}${CYAN}▶  $1${RESET}"; }
ok()      { echo -e "   ${GREEN}✔${RESET}  $1"; }
skip()    { echo -e "   ${DIM}–  $1${RESET}"; }
warn()    { echo -e "   ${YELLOW}⚠  $1${RESET}"; }
divider() { echo -e "${DIM}   ──────────────────────────────────────${RESET}"; }

# ── Header ─────────────────────────────────────────────────────────────────────
echo -e "\n${BOLD}  scripts-organizer — build${RESET}"
divider

# ── 1. Build Go binary ─────────────────────────────────────────────────────────
step "Building Go binary"

mkdir -p "$DIST_DIR"
go build -o "$BUILD_OUTPUT"

ok "Build complete → dist/$BINARY_NAME"

# ── 2. Ensure ~/.local/bin exists ──────────────────────────────────────────────
step "Checking $BIN_DIR"
mkdir -p "$BIN_DIR"
ok "$BIN_DIR exists"

# ── 3. Add ~/.local/bin to PATH in ~/.bashrc if not already present ────────────
step "Checking PATH in $BASHRC"
PATH_EXPORT='export PATH="$HOME/.local/bin:$PATH"'
PATH_MARKER="# scripts-organizer: bin_dir"

touch "$BASHRC"

if grep -qE '(^|:)[^#]*\.local/bin([^a-zA-Z0-9_]|$)' "$BASHRC"; then
    skip "$BIN_DIR already declared in $BASHRC"
else
    echo "" >> "$BASHRC"
    echo "$PATH_MARKER" >> "$BASHRC"
    echo "$PATH_EXPORT" >> "$BASHRC"
    ok "Added $BIN_DIR to PATH in $BASHRC"
    warn "Restart your shell or run: source $BASHRC"
fi

# ── 4. Move binary into ~/.local/bin ───────────────────────────────────────────
step "Installing binary"

if [[ -f "$INSTALL_TARGET" ]]; then
    warn "Overwriting existing binary at $INSTALL_TARGET"
fi

mv "$BUILD_OUTPUT" "$INSTALL_TARGET"
chmod +x "$INSTALL_TARGET"

ok "Installed → $INSTALL_TARGET"

# ── Done ───────────────────────────────────────────────────────────────────────
divider
echo -e "\n${BOLD}${GREEN}  ✔  Done.${RESET}  Run: ${BOLD}$BINARY_NAME${RESET}\n"
