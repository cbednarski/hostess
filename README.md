# hostess [![](https://travis-ci.org/cbednarski/hostess.svg)](https://travis-ci.org/cbednarski/hostess)

An idempotent command-line utility for managing your `/etc/hosts` file.

## Usage

    hostess add domain ip   # Add or change a hosts entry for this domain pointing to this IP
    hostess add -off domain ip  # Add in a disabled state (if it already exists, disable it)
    hostess del domain      # Remove a domain from your hosts file
    hostess has domain      # exit code 0 or 1 if the domain is in your hostfile
    hostess off domain      # Disable a domain (but don't remove it completely)
    hostess on domain       # Re-enable a domain that was disabled
    hostess ls              # List domains, target ips, and on/off status
    hostess fix             # Read your hosts file and show warnings if there are bumps
    hostess dump            # Dump your hostfile as json
    hostess apply           # Add entries from a json file

    Flags

    -n   # Dry run. Show what will happen but don't do it; output to stdout
    -f   # Forcibly rewrite the hostfile, even if there are errors or conflicts

hostess may mangle your hosts file. In general it will probably look like this, with domains pointing at the same IP grouped together and disabled domains commented out.

    127.0.0.1 localhost hostname2 hostname3
    127.0.1.1 machine.name
    # 10.10.20.30 some.host

## Installation

Grab a [release](https://github.com/cbednarski/hostess/releases) or download the code and run `make install` (building probably requires go 1.4).

## Configuration

By default, hostess will read / write to `/etc/hosts`. You can use the `HOSTESS_FILE` environment variable to provide an alternate path (for testing).

## Disclaimer

hostess uses readme-driven-development and may not actually do any of the things listed above. When in doubt, pass the `-n` flag to try hostess without changing your system.
