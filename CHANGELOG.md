# Changelog

## v0.2.1 (2025-06-01)

### Fix
- Fixes done to implement gysmo into Nixpkgs.
- Correcting sysinfo tests so they succeed whenever the value returned by the functions are strings. It does not need to be more specific than this.
- moving EnsureConfigFileExists to main.go instead of having it in the LoadConfig function. This should help mocking the tests.
- removing extras folder from this repo.
- changing default build name

## v0.2.0 (2025-05-07)

### Feature
- Changing Path for ascii file so it searches in /ascii by default
- Adding a version flag

## v0.1.5 (2025-03-29)

### Fix
- Adding a release to Nixpkgs

## v0.1.4 (2025-03-25)

### Fix
- Fix in installation script

## v0.1.3 (2025-03-21)

### Fix
- Fixes in the installation script + fun times testing automated releases

## v0.1.2 (2025-03-21)

### Fix
- Tests to fix the binary build

## v0.1.1 (2025-03-17)

### Fix
- Correcting default ascii art in the config
- adding colors to installation script so it's readable.

## v0.1.0 (2025-03-16)

### Feature
- First Release

## Alpha-0.1.0 (2025-02-12)

### Feature
- Adding a first pre-release
