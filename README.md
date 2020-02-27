# hostess [![Linux Build Status](https://travis-ci.org/cbednarski/hostess.svg)](https://travis-ci.org/cbednarski/hostess) [![Windows Build Status](https://ci.appveyor.com/api/projects/status/wtxqb880b7v9dfgn/branch/master?svg=true)](https://ci.appveyor.com/project/cbednarski/hostess/branch/master) [![GoDoc](https://godoc.org/github.com/cbednarski/hostess?status.svg)](http://godoc.org/github.com/cbednarski/hostess)

An **idempotent** command-line utility for managing your `/etc/hosts`* file.

    hostess add local.example.com 127.0.0.1
    hostess add staging.example.com 10.0.2.16

Why? Because you edit `/etc/hosts` for development, testing, and debugging.
Because sometimes DNS doesn't work in production. And because editing
`/etc/hosts` by hand is a pain. Put hostess in your `Makefile` or deploy scripts
and call it a day.

\* And `C:\Windows\System32\drivers\etc\hosts` on Windows.

**Note: 0.4.0 has backwards incompatible changes in the API and CLI.** See
`CHANGELOG.md` for details.

## Installation

Download a [precompiled release](https://github.com/cbednarski/hostess/releases)
from GitHub, or build from source (with a [recent version of Go](https://golang.org/dl)):

    go get -u github.com/cbednarski/hostess

## Usage

Run `hostess` or `hostess -h` to see a full list of commands.

## Format

On unixes, hostess follows the format specified by `man hosts`, with one line
per IP address:

    127.0.0.1 localhost hostname2 hostname3
    127.0.1.1 machine.name
    # 10.10.20.30 some.host

On Windows, hostess writes each hostname on its own line.

    127.0.0.1 localhost
    127.0.0.1 hostname2
    127.0.0.1 hostname3

## Configuration

hostess may be configured via environment variables.

- `HOSTESS_FMT` may be set to `windows` or `unix` to override platform detection
  for the hosts file format. See Behavior, above, for details

- `HOSTESS_PATH` may be set to override platform detection for the location of
  the hosts file. By default this is `C:\Windows\System32\drivers\etc\hosts` on
  Windows and `/etc/hosts` everywhere else.

## IPv4 and IPv6

Your hosts file _may_ contain overlapping entries where the same hostname points
to both an IPv4 and IPv6 IP. In this case, hostess commands will apply to both
entries. Typically you won't have this kind of overlap and the default behavior
is OK. However, if you need to be more granular you can use `-4` or `-6` to
limit operations to entries associated with that type of IP.

## Contributing

I hope my software is useful, readable, fun to use, and helps you learn
something new. I maintain this software in my spare time. I rarely merge PRs
because I am both lazy and a snob. Bug reports, fixes, questions, and comments
are welcome but expect a delayed response. No refunds. 👻
