function list_domains
  hostess list $argv | cut -d " " -f 1
end

complete -c hostess -x -n '__fish_use_subcommand' -a add -d 'Add or replace a hosts entry.'
complete -c hostess -x -n '__fish_use_subcommand' -a aff -d 'Add or replace a hosts entry in an off state.'
complete -c hostess -x -n '__fish_use_subcommand' -a has -d 'Exit code 0 if the domain is in your hostfile, 1 otherwise.'
complete -c hostess -f -n '__fish_use_subcommand' -a fix -d 'Rewrite your hosts file.'
complete -c hostess -x -n '__fish_use_subcommand' -a dump -d 'Dump your hostfile as JSON.'
complete -c hostess -r -n '__fish_use_subcommand' -a apply -d 'Add entries from a JSON file.'

complete -c hostess -f -n '__fish_use_subcommand' -a list -d 'List domains, target ips, and on/off status.'
complete -c hostess -f -n '__fish_use_subcommand' -a ls -d 'List domains, target ips, and on/off status.'
complete -c hostess -f -A -n '__fish_seen_subcommand_from list' -a off -d 'List only disabled domains.'
complete -c hostess -f -A -n '__fish_seen_subcommand_from ls' -a off -d 'List only disabled domains.'
complete -c hostess -f -A -n '__fish_seen_subcommand_from list' -a on -d 'List only enabled domains.'
complete -c hostess -f -A -n '__fish_seen_subcommand_from ls' -a on -d 'List only enabled domains.'

complete -c hostess -x -n '__fish_use_subcommand' -a del -d 'Remove a host entry.'
complete -c hostess -x -n '__fish_use_subcommand' -a rm -d 'Remove a host entry.'
complete -c hostess -f -A -n '__fish_seen_subcommand_from del' -a '(list_domains)'
complete -c hostess -f -A -n '__fish_seen_subcommand_from rm' -a '(list_domains)'

complete -c hostess -x -n '__fish_use_subcommand' -a off -d "Disable a domain (but don't remove it completely)."
complete -c hostess -f -A -n '__fish_seen_subcommand_from off' -a '(list_domains on)'

complete -c hostess -x -n '__fish_use_subcommand' -a on -d 'Re-enable a domain that was disabled.'
complete -c hostess -f -A -n '__fish_seen_subcommand_from on' -a '(list_domains off)'
