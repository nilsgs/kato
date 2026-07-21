# kato shell integration
kato() {
  if [ "$1" = "nav" ]; then
    local dir
    dir=$(command kato "$@")
    [ -n "$dir" ] && cd "$dir"
  else
    command kato "$@"
  fi
}

# kato alias loader — reads ~/.kato/aliases at shell startup
_kato_load_aliases() {
  local file="$HOME/.kato/aliases"
  [ -f "$file" ] || return
  while IFS= read -r line; do
    [[ -z "$line" || "$line" == \#* ]] && continue
    alias "${line%%=*}"="${line#*=}"
  done < "$file"
}
_kato_load_aliases
