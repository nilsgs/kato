# kato shell integration
function kato
  if test "$argv[1]" = "nav"
    set dir (command kato $argv)
    and cd $dir
  else
    command kato $argv
  end
end

# kato alias loader — reads ~/.kato/aliases at shell startup
for line in (cat ~/.kato/aliases 2>/dev/null)
  string match -qr '^\s*(#|$)' $line; and continue
  set parts (string split -m1 = $line)
  alias $parts[1] $parts[2]
end
