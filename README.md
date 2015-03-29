# hostess [![](https://travis-ci.org/cbednarski/hostess.svg)](https://travis-ci.org/cbednarski/hostess)

An idempotent command-line utility for managing your `/etc/hosts` file.

## Using Hostess

### Download and Install

Download a [precompiled release](https://github.com/cbednarski/hostess/releases) from GitHub.

### Usage

    hostess add domain ip   # Add or replace a hosts entry for this domain pointing to this IP
    hostess aff domain ip   # Add or replace a hosts entry in an off state
    hostess del domain      # (alias rm) Remove a domain from your hosts file
    hostess has domain      # exit code 0 if the domain is in your hostfile, 1 otherwise
    hostess off domain      # Disable a domain (but don't remove it completely), exit 1 if entry is missing
    hostess on domain       # Re-enable a domain that was disabled, exit 1 if entry is missing
    hostess list            # (alias ls) List domains, target ips, and on/off status
    hostess fix             # Rewrite your hosts file; use -n to dry run
    hostess dump            # Dump your hostfile as json
    hostess apply           # Add entries from a json file

    Flags

    -n   # Dry run. Show what will happen but don't do it; output to stdout
    -f   # Forcibly rewrite the hostfile, even if there are errors or conflicts
    -4   # Limit operation to ipv4 entries
    -6   # Limit operation to ipv6 entries

hostess may mangle your hosts file. In general it will probably look like this, with domains pointing at the same IP grouped together and disabled domains commented out.

    127.0.0.1 localhost hostname2 hostname3
    127.0.1.1 machine.name
    # 10.10.20.30 some.host

### IPv4 and IPv6

Your hosts file *can* contain overlapping entries where the same hostname points to both an IPv4 and IPv6 IP. In this case, hostess commands will apply to both entries. Typically you won't have this kind of overlap and the default behavior is OK. However, if you need to be more granular you can use `-4` or `-6` to limit operations to entries associated with that type of IP.

### Disclaimer

hostess uses readme-driven-development and may not actually do any of the things listed above. When in doubt, pass the `-n` flag to try hostess without changing your system.

## Developing Hostess

### Configuration

By default, hostess will read / write to `/etc/hosts`. You can use the `HOSTESS_PATH` environment variable to provide an alternate path (for testing).

### Building from Source

To build from source you'll need to have go 1.4+

#### Install with go get

    go get github.com/cbednarski/hostess/cmd/hostess

#### Install from source

    git clone https://github.com/cbednarski/hostess
    cd hostess
    make
    make install
