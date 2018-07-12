# miitomo-assetscraper

### An asset scraper for Miitomo

----

## Getting and using Miitomo Asset Scraper

1. Open the [releases page](https://github.com/JoshuaDoes/miitomo-assetscraper/releases) and download the latest appropriate version of Miitomo Asset Scraper to whichever folder on your computer you'd like to download Miitomo assets to
2. Run `miitomo-assetscraper [url]` where url = the Miitomo asset manifest URL
3. Enjoy!

## What does it do?

When you run Miitomo Asset Scraper with the asset manifest URL provided, it scans through the asset manifest and downloads all available assets, including running MD5 checksum validation on each downloaded asset and extracting any compressed assets.

----

## Building it yourself

In order to build Miitomo Asset Scraper locally, you must have already installed
a working Golang environment on your development system and installed the package
dependencies that Miitomo Asset Scraper relies on to function properly.

Miitomo Asset Scraper is currently built using Golang `1.10.2`.

### Dependencies

| Package Name |
| ------------ |
| [go-unarr](https://github.com/gen2brain/go-unarr) |

### Building

Simply run `go build` in this repo's directory once all dependencies are satisfied.

### Running Miitomo Asset Scraper

Finally, to run Miitomo Asset Scraper, simply type `./miitomo-assetscraper [url]` in your
terminal/shell or `.\miitomo-assetscraper.exe [url]` in your command prompt. If everything
goes well, you'll see the download progress in your terminal and all downloaded/extracted
Miitomo assets available once the download completes.

### Contributing notes

When pushing to your repo or submitting pull requests to this repo, it is highly
advised that you clean up the working directory to only contain `LICENSE`, `main.go`,
`README.md`, and the `.git` folder. A proper `.gitignore` will be written soon to
mitigate this requirement.

----

## Support
For help and support with Miitomo Asset Scraper, create an issue on the issues page. If you do not have a GitHub account, send me a message on Discord (@JoshuaDoes#1685) or join the [Kaeru Network Discord](https://discord.me/kaeru).

## License
The source code for Miitomo Asset Scraper is released under the MIT License. See LICENSE for more details.

## Donations
All donations are highly appreciated. They help motivate me to continue working on side projects like these, especially when it comes to something you may really want added!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/JoshuaDoes)