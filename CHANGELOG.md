# Change Log

## v0.5.2 (March 13, 2020)

Bug Fixes

- `hostess fmt -n` works properly again, and has more specific behavior:
- `hostess fmt` will replace duplicates without asking for help
- `hostess fmt -n` will *not* replace duplicates, and will exit with error if any are found (#41)
- `hostess fmt` with and without `-n` will exit with error if conflicting hostnames are found because hostess cannot fix the conflicts

## v0.5.1 (March 10, 2020)

Bug Fixes

- Format will no longer exit with an error when encountering a duplicate entry (#39)

## v0.5.0 (March 7, 2020)

Breaking changes

- Windows now has a platform-specific hosts format with one IP and hostname per line

## v0.4.1 (February 28, 2020)

Bug Fixes

- Fix hostfiles not saving on Windows #27

## v0.4.0 (February 28, 2020)

0.4.0 is a major refactor of the frontend (CLI) focused on simplifying the UI
and code, supporting newer Go tooling (i.e. go mod), and removing external
dependencies.

Breaking Changes

- Moved CLI to `github.com/cbednarski/hostess`. `go get` should now do what you probably wanted the first time.
- Moved library to `github.com/cbednarski/hostess/hostess`
- Many command aliases and flags have been removed
- `Hostlist.Enable` and `Hostlist.Disable` now return an `error` instead of `bool`. Check against `ErrHostnameNotFound`.
- Several functions will now return `ErrInvalidVersionArg` instead of panicking in that case

Improvements

- Removed `codegangsta/cli`
- Removed `aff` command
- Removed `del` command (use `rm` instead)
- Removed `list` command (use `ls` instead)
- Removed `fixed` command (just run `fmt`)
- Command `fix` renamed to `fmt`
- Removed `-s` and `-q` flags. Errors are now shown always. Redirect stderr if you don't want to see them.
- Removed `-f` from various commands. Use `fmt` instead.
- Added Go mod support
- Added AppVeyor for Windows builds
- Overhauled the Makefile for easier builds

## v0.3.0 (February 18, 2018)

Improvements

- Added `fixed` subcommand which checks whether the hosts file is already formatted by hostess

Bug Fixes

- Show an error when there is a parsing failure instead of silently truncating the hosts file
- Global flags between hostess and the subcommand are no longer ignored
- Binary should now display the correct version of the software

## v0.2.1 (May 17, 2016)

Bug Fixes

- Fix vendor path for `codegangsta/cli`

## v0.2.0 (May 10, 2016)

Improvements

- Vendor `codegangsta/cli` for more reliable builds

Bug Fixes

- Fix panic in `hostess ls` #14
- Fix incompatible API in CLI library #15

## v0.1.0 (June 6, 2015)

Initial release
