Installation
============

This document describes how to install [actionlint](../docs).

> [!IMPORTANT]
> This repository is a fork of [rhysd/actionlint](https://github.com/rhysd/actionlint). The executable is still named `actionlint`,
> but the fork's packages use the name **`actionlint-zew1me`** (Homebrew cask, Docker image, release archives) and the Go module
> path is `github.com/zew1me/actionlint`. Third-party package managers listed in the sections below (Chocolatey, Scoop, winget,
> Arch Linux, Nix, Homebrew core, asdf, mise) distribute **upstream** actionlint, not this fork. Only the sections "Homebrew cask",
> "Prebuilt binaries", "Download script", "Docker image" and "Build from source" install this fork.

## Windows

### [Chocolatey](https://chocolatey.org/)

[`actionlint` package][chocolatey] is available in the community repository:

```powershell
choco install actionlint
```

### [Scoop](https://scoop.sh/)

[`actionlint` package][scoop] is available in the main bucket:

```powershell
scoop install actionlint
```

### [Winget](https://learn.microsoft.com/en-us/windows/package-manager/)

[`actionlint` package][winget] is available in the winget-pkgs repository:

```powershell
winget install actionlint
```

## Linux

### [Arch Linux](https://archlinux.org/)

[`actionlint` package][archlinux] is available in the official repository:

```sh
pacman -S actionlint
```

Alternatively actionlint is also available on [AUR][aur]. The packages can be installed via [`paru`][paru] command.

- [actionlint-bin](https://aur.archlinux.org/packages/actionlint-bin)
- [actionlint-git](https://aur.archlinux.org/packages/actionlint-git)

### [Nix](https://nixos.wiki/)

[`actionlint` package][nixpkgs] is available in the Nix ecosystem:

On NixOS:

```sh
nix-env -iA nixos.actionlint
```

On Non NixOS:

```sh
nix-env -iA nixpkgs.actionlint
```

## macOS

### [Homebrew][homebrew]

[`actionlint`][formula] formula is provided by Homebrew officially. Note that it installs **upstream** actionlint, not this fork.

```sh
brew install actionlint
```

This repository provides its own Homebrew cask named `actionlint-zew1me`, which is generated on new releases of this fork.
Tap the repository and install the cask with `--cask` option. The cask installs the `actionlint` executable, so it conflicts
with upstream's `actionlint` cask. Uninstall the upstream `actionlint` formula first, otherwise both link the same executable.

```sh
brew tap "zew1me/actionlint" "https://github.com/zew1me/actionlint"
brew install --cask actionlint-zew1me
```

> [!WARNING]
> Since the `actionlint` executable is unsigned, macOS displays a warning and tries to move it to the Trash. To allow it to run,
> go to 'Settings -> Privacy & Security' and grant the permission.

## Prebuilt binaries

Download an archive file from [the releases page][releases] for your platform, unarchive it and put the executable file to a
directory in `$PATH`.

Prebuilt binaries are built at each release of this fork by CI for the following OS and arch:

- macOS (x86_64, arm64)
- Linux (i386, x86_64, arm32, arm64)
- Windows (i386, x86_64, arm64)
- FreeBSD (i386, x86_64)

Note that the following targets are not tested since GitHub Actions doesn't support them:

- Linux i386, arm32
- Windows i386
- FreeBSD i386, x86_64

To install these binaries [`gh`][gh] command is useful. The following command is an example for x86_64 Linux.

```sh
gh release download --repo zew1me/actionlint --pattern '*_linux_amd64.tar.gz' vX.Y.Z
tar xf actionlint-zew1me_X.Y.Z_linux_amd64.tar.gz
./actionlint -version
```

Optionally you can verify the [attestation][attestations] of the downloaded artifact. This is highly recommended in terms of
security. Release archives of this fork are named `actionlint-zew1me_{version}_{os}_{arch}.{tar.gz|zip}` and contain the `actionlint`
executable.

```sh
gh attestation verify -R zew1me/actionlint actionlint-zew1me_X.Y.Z_linux_amd64.tar.gz
```

<a id="download-script"></a>
## Download script

To install `actionlint` executable with one command, [the download script](../scripts/download-actionlint.bash) is available.
It downloads the default version of actionlint from the releases of this fork (`actionlint.exe` on Windows and `actionlint` on other OSes) to the current
directory automatically. This is a recommended way if you install actionlint in some shell script.

```sh
bash <(curl https://raw.githubusercontent.com/zew1me/actionlint/main/scripts/download-actionlint.bash)
```

When you need to install specific version of actionlint, please give the version to the 1st command line argument. The following
example installs vX.Y.Z (replace it with a version released by this fork).

```sh
bash <(curl https://raw.githubusercontent.com/zew1me/actionlint/main/scripts/download-actionlint.bash) X.Y.Z
```

This script downloads `actionlint` (or `actionlint.exe` on Windows) binary to the current working directory. When you need to put
the downloaded binary to some other directory, please give the directory path to the 2nd command line argument. The following
example installs the latest version to `/usr/bin`.

```sh
bash <(curl https://raw.githubusercontent.com/zew1me/actionlint/main/scripts/download-actionlint.bash) latest /usr/bin
```

For the usage of actionlint on GitHub Actions, see [the usage document](usage.md#on-github-actions).

## Docker image

See [the usage document](./usage.md#docker) to know how to install and use the Docker image of this fork
(`ghcr.io/zew1me/actionlint-zew1me`).

## Cross-platform version managers

### asdf

You can install actionlint with the [asdf version manager][asdf] using the [asdf-actionlint][asdf-plugin] plugin, which
automates the process of installing (and switching between) various versions of GitHub release binaries. With asdf already
installed, run these commands to install actionlint:

```bash
# Add actionlint plugin
asdf plugin add actionlint

# Show all installable versions
asdf list-all actionlint

# Install specific version
asdf install actionlint latest

# Set a version globally (on your ~/.tool-versions file)
asdf global actionlint latest
```

### mise

You can install actionlint with the [mise-en-place][mise] which automates the process of installing (and switching
between) various versions of GitHub release binaries. With mise already installed, run these commands to install
actionlint:

```bash
# Show all installable versions
mise ls-remote actionlint

# Install specific version
mise install actionlint@latest

# Set a version globally (on your ~/.config/mise/config.toml file)
mise use -g actionlint@latest
```

## Build from source

Recent [Go][] toolchain is necessary to build actionlint from source. Last two major versions of Go are supported.

```sh
# Install the latest tagged version (when this fork has published a tagged release)
go install github.com/zew1me/actionlint/cmd/actionlint@latest

# Install the head of the main branch
go install github.com/zew1me/actionlint/cmd/actionlint@main
```

---

[Checks](checks.md) | [Usage](usage.md) | [Configuration](config.md) | [Go API](api.md) | [References](reference.md)

[formula]: https://formulae.brew.sh/formula/actionlint
[homebrew]: https://brew.sh/
[releases]: https://github.com/zew1me/actionlint/releases
[gh]: https://docs.github.com/en/github-cli/github-cli/about-github-cli
[attestations]: https://docs.github.com/en/actions/concepts/security/artifact-attestations
[Go]: https://golang.org/
[asdf]: https://asdf-vm.com/
[asdf-plugin]: https://github.com/crazy-matt/asdf-actionlint
[chocolatey]: https://community.chocolatey.org/packages/actionlint
[scoop]: https://scoop.sh/#/apps?q=actionlint&s=0&d=1&o=true
[winget]: https://github.com/microsoft/winget-pkgs/tree/master/manifests/r/rhysd/actionlint
[archlinux]: https://archlinux.org/packages/extra/x86_64/actionlint/
[aur]: https://aur.archlinux.org/
[paru]: https://github.com/Morganamilo/paru
[nixpkgs]: https://github.com/NixOS/nixpkgs/blob/master/pkgs/development/tools/analysis/actionlint/default.nix
[mise]: https://github.com/jdx/mise
