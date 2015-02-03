# hostess

A command-line utility for managing your `/etc/hosts` file.

## Usage

    hostess add domain ip   # Add or change a hosts entry for this domain pointing to this IP
    hostess del domain      # Remove a domain from your hosts file
    hostess has domain      # exit code 0 or 1 depending on whether the domain is in your hosts file
    hostess off domain      # Disable a domain (but don't remove it completely)
    hostess on domain       # Re-enable a domain that was disabled
    hostess ls              # List domains, target ips, and on/off status

hostess may mangle your hosts file. In general it will probably look like this, with domains pointing at the same IP grouped together and disabled domains commented out.

    127.0.0.1 localhost hostname2 hostname3
    127.0.1.1 machine.name
    # 10.10.20.30 some.host

## Installation

Grab a [release](https://github.com/cbednarski/hostess/releases) or download the code and run `make && make install` (building probably requires go 1.4).

## Configuration

By default, hostess will read / write to `/etc/hosts`. You can use the `HOSTESS_FILE` environment variable to provide an alternate path (for testing).

## Disclaimer

hostess reserves the right to sort, parse, and validate as it sees fit (or not at all) and may be ruthlessly hostile towards comments, whitespace, and other things that robots don't care for. hostess may include different default entries depending on OS. hostess uses readme-driven-development and may not actually do any of the things listed above. You have been warned.
