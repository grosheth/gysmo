# gysmo 🚀

**gysmo** is a fun and visually engaging CLI tool written in Go, designed to enhance your terminal experience. While it does provide system information, its primary purpose is to showcase your stylish terminal setup.

**Key Features:**
- **Highly Customizable:** Tailor gysmo to your preferences with user-configurable settings. The default JSON configuration is cool, but we encourage you to make it your own.
- **Show Off Your Setup:** Perfect for proudly displaying that you use Arch, by the way.
- **Inspired by the Best:** Draws inspiration from [nitch](https://github.com/ssleert/nitch) and [neofetch](https://github.com/dylanaraps/neofetch).

**DISCLAIMER:**
gysmo is not intended to be an accurate system information tool. There are many other tools available for that purpose.

## 📥 Installation

### Download the Binary
```shell
wget setup.sh
```

### Home-manager
The program is not included in nixpkgs, I might do it later.
There is still a way to generate your gysmo config with home-manager.

```nix
{ pkgs, ... }:
{

}
```

### Build From Source
```shell
git clone https://github.com/yourusername/gysmo.git
cd gysmo
go build
```

## ⚙️ Configuration
Like I said in the introduction, the default configuration is not meant to be used and although it can show off your system, it also can show anything you want.

Here are every section of the configuration file you can modify:
<details>
  <summary>📝Items</summary>
  The items section is where you define what you want to show in your gysmo main menu. The following is an example configuration:

  ```json
  "items": [
    {
      "alias": "user",
      "name": "user",
      "icon": "",
      "value_color": "purple",
      "text_color": "",
      "image_color": "red",
      "value": "Custom value"
    },
    {
      "alias": "name-alias",
      "name": "name",
      "icon": "",
      "value_color": "purple",
      "text_color": "green",
      "image_color": "red"
    }
  ]
  ```

Here is a brief explanation of each option:

| Option       | Description                                                                 | Example Value       |
|--------------|-----------------------------------------------------------------------------|---------------------|
| `alias`      | This is the value that will be shown in the middle of the menu.             | `"username"`        |
| `name`       | This is the system value gysmo will return.                                 | `"user"`            |
| `icon`       | An icon to display alongside the item.                                      | `""`               |
| `value_color`| The color of the value text.                                                | `"purple"`          |
| `text_color` | The color of the item text.                                                 | `"green"`           |
| `image_color`| The color of the icon.                                                      | `"red"`             |
| `value`      | A custom value to display for the item. (will replace the value returned by the keyword)                                    | `"Custom value"`    |

## Keywords for `name` Option
Some values of /etc/os-release are not available on some distros, look at [os-release](https://github.com/which-distro/os-release) to get an idea.

| Keyword                | Description                                      | Example Value            |
|------------------------|--------------------------------------------------|--------------------------|
| `os_ansi_color`        | ANSI color of the /etc/os-release                             | `"osRelease.ANSI_COLOR"`   |
| `os_pretty_name`       | Pretty name of the /etc/os-release                            | `"NixOS 25.05 (Warbler)"`  |
| `os_bug_report_url`    | Bug report URL of the /etc/os-release                         | `"https://github.com/NixOS/nixpkgs/issues"`|
| `os_build_id`          | Build ID of the /etc/os-release                               | `"25.05.20250204.799ba5b"`     |
| `os_cpe_name`          | CPE name of the /etc/os-release                               | `"cpe:/o:nixos:nixos:25.05"`     |
| `os_default_hostname`  | Default hostname of the /etc/os-release                       | `"nixos"`|
| `os_documentation_url` | Documentation URL of the /etc/os-release                      | `"https://nixos.org/learn.html"`|
| `os_home_url`          | Home URL of the /etc/os-release                               | `"https://nixos.org/"`     |
| `os_id`                | ID of the /etc/os-release                                     | `"nixos"`           |
| `os_id_like`           | ID like of the /etc/os-release                                | `"osRelease.ID_LIKE"`      |
| `os_image_id`          | Image ID of the /etc/os-release                               | `"osRelease.IMAGE_ID"`     |
| `os_image_version`     | Image version of the /etc/os-release                          | `"osRelease.IMAGE_VERSION"`|
| `os_version`           | Version of the /etc/os-release                                | `"25.05 (Warbler)"`      |
| `os_logo`              | Logo of the /etc/os-release                                   | `"nix-snowflake"`         |
| `os_name`              | Name of the /etc/os-release                                   | `"NixOS"`         |
| `os_support_url`       | Support URL of the /etc/os-release                            | `"https://nixos.org/community.html"`  |
| `os_variant`           | Variant of the /etc/os-release                                | `"osRelease.VARIANT"`      |
| `os_variant_id`        | Variant ID of the /etc/os-release                             | `"osRelease.VARIANT_ID"`   |
| `os_vendor_name`       | Vendor name of the /etc/os-release                            | `"NixOS"`  |
| `os_vendor_url`        | Vendor URL of the /etc/os-release                             | `"https://nixos.org/"`   |
| `os_version_codename`  | Version codename of the /etc/os-release                       | `"warbler"`|
| `os_version_id`        | Version ID of the /etc/os-release                             | `"25.05"`   |
| `user`                 | Username of the current user                     | `"user"`          |
| `hostname`             | Hostname of the system                           | `"hostname"`          |
| `kernel`               | Kernel version of the system                     | `"6.6.75"`     |
| `shell`                | Default shell of the user                        | `"zsh"`             |
| `uptime`               | System uptime                                    | `"19:44:53"`            |
| `dm`                   | Desktop manager                                  | `"KDE"`    |
| `gpu`                  | GPU information                                  | `"GPU Info"`           |
| `cpu`                  | CPU information                                  | `"CPU Info"`           |
| `ram`                  | RAM information                                  | `"RAM Info"`           |
| `drive`                | Drive information                                | `"Drive Info"`         |
| `gpu %`                | GPU usage percentage                             | `"GPU Usage"`          |
| `cpu %`                | CPU usage percentage                             | `"CPU Usage"`          |
| `ram %`                | RAM usage percentage                             | `"RAM Usage"`          |
| `drive %`              | Drive usage percentage                           | `"Drive Usage"`        |
| `term`                 | Terminal information                             | `"ghostty"`          |
| `processes`            | Number of running processes                      | `"121"`|
| `wm`            | Window Manager                     | `"none+bpswm"`|
| ``            | Number of running processes                      | `"121"`|
| `processes`            | Number of running processes                      | `"121"`|
| `processes`            | Number of running processes                      | `"121"`|
| `processes`            | Number of running processes                      | `"121"`|
| `processes`            | Number of running processes                      | `"121"`|
| `processes`            | Number of running processes                      | `"121"`|

  ⚠️ **WARNING**: If you are using custom images, the borders might get messed up. Using icons from [Nerd Fonts](https://www.nerdfonts.com/) seems to work fine.
</details>

<details>
  <summary>🎨 ascii</summary>
I don't aim to keep millions of ASCII art in this repo.

I might open an Issue so people can share their ASCII art.

Instead, I suggest you get the ASCII art you like from the following sources:

### Sources
- [asciiart.eu](https://www.asciiart.eu/) (all the gysmo ASCII art were made there)
- [asciiflow](https://asciiflow.com/#/)


Here is an example of the ASCII configuration:
the ascii section is a simple dictionnary with the following options:

  ```json
  "ascii": {
    "path": "ascii/gysmo4",
    "colors": "red",
    "enabled": true,
    "padding": 2
  },
  ```

</details>

<details>
  <summary>📝 header</summary>

  ```json
  "header": {
    "enabled": true,
    "text": "NixOS",
    "text_color": "white",
    "line": true,
    "line_color": "red"
  },
  ```
</details>

<details>
  <summary>📝 footer</summary>

  ```json
  "footer": {
    "enabled": true,
    "text": "gysmo",
    "text_color": "white",
    "line": true,
    "line_color": "red"
  },
  ```
</details>

<details>
  <summary>📝 general</summary>

  ```json
  "general": {
    "padding": 2,
    "left_padding": 200,
    "menu_type": "box"
  }
  ```

</details>

Here are some examples of what you can do with gysmo.
## Examples
- GitHub stats (stars, forks, issues, pull requests)
- Weather information
- System information

## 🤝 Contributing

### Versioning
A github workflow is in place to automatically bump the version of the project. I will merge the PRs when I am ready.
To version this project, I use [Semantic Versioning](https://semver.org/).
