# gysmo 🚀

![License](https://shields.io/github/license/grosheth/gysmo?style=for-the-badge&color=blue)
[![GitHub Tag](https://img.shields.io/github/v/tag/grosheth/gysmo?include_prereleases&sort=semver&style=for-the-badge&color=white)](https://github.com/grosheth/gysmo/releases/tag/Alpha-0.1.0)
[![CodeFactor](https://www.codefactor.io/repository/github/grosheth/gysmo/badge/0.1.0?style=for-the-badge)](https://www.codefactor.io/repository/github/grosheth/gysmo/overview/0.1.0)

# ADD IMAGES HERE

**gysmo** is a fun and visually engaging CLI tool written in Go, designed to enhance your terminal experience. While it does provide system information, its primary purpose is to showcase your stylish terminal setup.

**Key Features:**
- **Highly Customizable:** Tailor gysmo to your preferences with user-configurable settings. The default JSON configuration is cool, but we encourage you to make it your own.
- **Show Off Your Setup:** Perfect for proudly displaying that you use Arch, by the way.
- **Inspired by the Best:** Draws inspiration from [nitch](https://github.com/ssleert/nitch) and [neofetch](https://github.com/dylanaraps/neofetch).

**DISCLAIMER:**
gysmo is not intended to be an accurate system information tool. There are many other tools available for that purpose.

## 📥 Installation

### Build From Source (RECOMMENDED)
If you are a true Linux user and you won't submit to anyone else automated processes (as is your right), you can always prove your superiority by building the project from source.
This also allows you to understand how the project works a little better.

```shell
git clone https://github.com/grosheth/gysmo.git
cd gysmo/src
# -o only renames the binary (default is src)
go build -o gysmo
mkdir ~/.config/gysmo
cp -R config ~/.config/gysmo/
cp -R ascii ~/.config/gysmo/
```

### Use the installation/update script (Why not take the easy way?)
The installation/update script basically does the following
1. Download the binary
2. Create the directory structure ~/.config/gysmo/... && ~/bin/gysmo
3. Copy the template config.json file IF it doesn't exist (don't want to overwrite your config in an update)
4. Copy the schema validator file
5. Copy a sample ASCII art IF it doesn't exist

```shell
wget installation.sh
```

### Home-manager (Nix users go brrrrrrr)
⚠️ **NOT OFFICIALLY DONE YET**: still to be done.
Package will be included in nixpkgs and home-manger eventually.

⚠️ **BUT**: Here is a way to do it in you nix config without an "official" home-manager module.
```nix
{ pkgs, ... }:
{

}
```

## ⚙️ Configuration
Like I said in the introduction, the default configuration is not meant to be used and although it can show off your system, it also can show anything you want.

Here are every section of the configuration file you can modify:

### File structure

<details>
  <summary>📝 Example Configuration</summary>

  ```json
  {
    "items": [
      {
        "text": "user",
        "keyword": "user",
        "icon": "",
        "value_color": "red",
        "text_color": "",
        "icon_color": "red",
        "value": "My user"
      },
      {
        "text": "shell",
        "keyword": "shell",
        "icon": "",
        "value_color": "yellow",
        "text_color": "",
        "icon_color": "yellow"
      }
    ],
    "ascii": {
      "path": "ascii/gysmo2",
      "colors": "",
      "enabled": true,
      "horizontal_padding": 0,
      "vertical_padding": 0,
      "position": "left"
    },
    "header": {
      "enabled": true,
      "text": "NixOS",
      "text_color": "purple",
      "line": true,
      "line_color": ""
    },
    "footer": {
      "enabled": true,
      "text": "gysmo",
      "text_color": "blue",
      "line": true,
      "line_color": ""
    },
    "general": {
      "menu_type": "box",
      "columns": false
    }
  }

  ```

</details>
<details>
  <summary>📝 general</summary>

  ```json
  "general": {
    "menu_type": "box",
    "columns": false
  }
  ```

</details>
<details>
  <summary>📝Items</summary>
  The items section is where you define what you want to show in your gysmo main menu. The following is an example configuration:

  ```json
  "items": [
    {
      "text": "user",
      "keyword": "user",
      "icon": "",
      "value_color": "red",
      "text_color": "",
      "icon_color": "red",
      "value": "My user"
    },
    {
      "text": "shell",
      "keyword": "shell",
      "icon": "",
      "value_color": "yellow",
      "text_color": "",
      "icon_color": "yellow"
    }
  ],
  ```

Here is a brief explanation of each option:

| Option       | Description                                                                 | Example Value       |
|--------------|-----------------------------------------------------------------------------|---------------------|
| `text`      | This is the value that will be shown in the middle of the menu.             | `"username"`        |
| `keyword`       | This is the system value gysmo will return. (does  not work with "value")              | `"user"`            |
| `icon`       | An icon to display alongside the item. Can also be text.                                     | `""`               |
| `value_color`| The color of the value text.                                                | `"purple"`          |
| `text_color` | The color of the item text.                                                 | `"green"`           |
| `icon_color`| The color of the icon.                                                      | `"red"`             |
| `value`      | A custom value to display for the item. (Does not work with keyword)                                    | `"Custom value"`    |

## Text

## Keywords for `keyword` Option
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
| `processes`            | Number of running processes                      | `"121"`|

## Icon
  ⚠️ **WARNING**: Icons is the most fragile part of gysmo. You technically can use text or multiple icons on one line but it's not that stable.

## Value



</details>

<details>
  <summary>🎨 ascii</summary>
I don't aim to keep millions of ASCII art in this repo.

Instead, I will open a discussion on the repo so people can share their ASCII art and configs.

I suggest you get the ASCII art you like from the following sources:

### Sources
- [asciiart.eu](https://www.asciiart.eu/)
- [ascii.co.uk](https://ascii.co.uk/art)

Here is an example of the ASCII configuration:
the ascii section is a simple dictionnary with the following options:

  ```json
    "ascii": {
      "path": "ascii/gysmo2",
      "colors": "",
      "enabled": true,
      "horizontal_padding": 0,
      "vertical_padding": 0,
      "position": "left"
    }
  ```

| Option       | Description                                                                 | Example Value       |
|--------------|-----------------------------------------------------------------------------|---------------------|
| `text`      | This is the value that will be shown in the middle of the menu.             | `"username"`        |
| `keyword`       | This is the system value gysmo will return. (does  not work with "value")              | `"user"`            |
| `icon`       | An icon to display alongside the item. Can also be text.                                     | `""`               |
| `value_color`| The color of the value text.                                                | `"purple"`          |
| `text_color` | The color of the item text.                                                 | `"green"`           |
| `icon_color`| The color of the icon.                                                      | `"red"`             |
| `value`      | A custom value to display for the item. (Does not work with keyword)                                    | `"Custom value"`    |

</details>

<details>
  <summary>📝 header</summary>

  ```json
  "header": {
    "enabled": true,
    "text": "NixOS",
    "text_color": "purple",
    "line": true,
    "line_color": ""
  },
  ```

| Option       | Description                                                                 | Example Value       |
|--------------|-----------------------------------------------------------------------------|---------------------|
| `text`      | This is the value that will be shown in the middle of the menu.             | `"username"`        |
| `keyword`       | This is the system value gysmo will return. (does  not work with "value")              | `"user"`            |
| `icon`       | An icon to display alongside the item. Can also be text.                                     | `""`               |
| `value_color`| The color of the value text.                                                | `"purple"`          |
| `text_color` | The color of the item text.                                                 | `"green"`           |
| `icon_color`| The color of the icon.                                                      | `"red"`             |
| `value`      | A custom value to display for the item. (Does not work with keyword)                                    | `"Custom value"`    |

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
  | Option       | Description                                                                 | Example Value       |
  |--------------|-----------------------------------------------------------------------------|---------------------|
  | `text`      | This is the value that will be shown in the middle of the menu.             | `"username"`        |
  | `keyword`       | This is the system value gysmo will return. (does  not work with "value")              | `"user"`            |
  | `icon`       | An icon to display alongside the item. Can also be text.                                     | `""`               |
  | `value_color`| The color of the value text.                                                | `"purple"`          |
  | `text_color` | The color of the item text.                                                 | `"green"`           |
  | `icon_color`| The color of the icon.                                                      | `"red"`             |
  | `value`      | A custom value to display for the item. (Does not work with keyword)                                    | `"Custom value"`    |

</details>

Here are some examples of what you can do with gysmo.
## Examples
- GitHub stats (stars, forks, issues, pull requests)
- Weather information
- System information

## Other Information

### 🎨Colors
You can specify any of these values in the color fields in the config to use the ANSI colors from you terminal.
```
Red
Green
Yellow
Blue
Purple
Cyan
White
```

if you wish to use any other colors, you can specify the RGB values in the following format:
```
#FFFFFF
```

## Json Validation
The configuration file is validated every time you run gysmo so you can be sure that your configuration is not missing anything. You can change this schema to your liking.


## 🛤️ ROADMAP

- [ ] Set ascii art in the background of the menu

- [ ] Add option to use images

- [ ] Easier integration of custom script to generate config files

## 📜 License
MIT LICENSE
