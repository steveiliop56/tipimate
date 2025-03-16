# Tipimate

![License](https://img.shields.io/github/license/steveiliop56/tipimate)
![Release](https://img.shields.io/github/v/release/steveiliop56/tipimate)
![Commit activity](https://img.shields.io/github/commit-activity/w/steveiliop56/tipimate)
![Issues](https://img.shields.io/github/issues/steveiliop56/tipimate)

Tipimate is an extremely simple and lightweight tool that check for updates in your [Runtipi](https://github.com/runtipi/runtipi) server and notifies you in your favorite notification system.

> [!NOTE]
> Tipimate only supports runtipi instances from v3.7.0 and above since this version added the ability to use an API to communicate with the server.

> [!WARNING]
> Tipimate is still in early stages of development so issues are to be expected. If you encounter any please create an issue so I can fix them as soon as possible.

## Getting started

You can run tipimate with two ways, either docker or binary. If you chose binary, you can grab the latest binary from the [releases](https://github.com/steveiliop56/tipimate/releases) page, then `chmod +x tipimate` and finally you can run it with `./tipimate`. _assuming the binary is named tipimate_

Running with docker is also very easy, you just need to download the docker compose file from [here](./docker-compose.yml) and run tipimate with `docker compose up -d`. _make sure to change the environment variables accordingly_

If you prefer a docker run command, you can run it with:

```bash
docker run -t -d --name tipimate -v ./data:/data -e NOTIFY_URL=your-discord-url -e RUNTIPI=your-runtipi-url -e JWT_SECRET=your-jwt-secret ghcr.io/steveiliop56/tipimate:v1
```

> [!TIP]
> You can set the `--runtipi-internal` flag or the `RUNTIPI_INTERNAL` environment variable to something like `http://localhost` if tipimate is running on the same server as your runtipi server and then set the `--runtipi` flag or `RUNIPI` environment variable to the public URL of your instance e.g. `https://runtipi.mydomain.com` so tipimate can both connect directly to runtipi and show the correct URL on Discord.

## Building

To build the project you need to have Go and Git installed.

You firstly have to clone the repository with:

```bash
git clone https://github.com/steveiliop56/tipimate
cd tipimate
```

Then install dependencies:

```bash
go mod tidy
```

And finally run it with:

```bash
go run .
```

Or build it with:

```bash
go build
```

If everything succeeds you should have a binary named `tipimate`.

> [!NOTE]
> You can also build for other operating systems/architectures using `GOOS=windows` and `GOARCH=arm64`.

> [!NOTE]
> You can also run a "development" docker compose file by copying the `.env.example` file to `.env`, changing your environment variables and running `docker compose -f docker-compose.dev.yml up --build`. With this way you can test your changes in the docker image too.

## Contributing

This project is still in early stages of development so bugs are to be expected. If you are interested in helping with the development feel free to create a pull request or an issue about a bug or a feature.

## License

TipiMate is licensed under the GNU General Public License v3.0. TL;DR — You may copy, distribute and modify the software as long as you track changes/dates in source files. Any modifications to or software including (via compiler) GPL-licensed code must also be made available under the GPL along with build & install instructions.
