#!/usr/bin/env bash
#!/usr/bin/env bash

set -e

LOCAL_BIN="$HOME/.local/bin"
EXPORT_LINE='export PATH="$HOME/.local/bin:$PATH"'

# Files to check (macOS + Linux)
FILES=(
  "$HOME/.bashrc"
  "$HOME/.bash_profile"
  "$HOME/.zshrc"
)

echo "Ensuring ~/.local/bin exists..."
mkdir -p "$LOCAL_BIN"

add_to_file() {
  local file="$1"

  # Create file if it doesn't exist
  [ -f "$file" ] || touch "$file"

  if grep -qxF "$EXPORT_LINE" "$file"; then
    echo "✓ PATH already set in $(basename "$file")"
  else
    echo "→ Adding PATH to $(basename "$file")"
    {
      echo ""
      echo "# Add local bin to PATH"
      echo "$EXPORT_LINE"
    } >> "$file"
  fi
}

echo "Updating shell config files..."

for file in "${FILES[@]}"; do
  add_to_file "$file"
done

echo ""
echo "Done!"
echo "Restart your shell or run:"
echo "  source ~/.zshrc   # for zsh"
echo "  source ~/.bashrc  # for bash"
