# Tipimate

![License](https://img.shields.io/github/license/steveiliop56/tipimate)
![Release](https://img.shields.io/github/v/release/steveiliop56/tipimate)
![Issues](https://img.shields.io/github/issues/steveiliop56/tipimate)

Tipimate is an extremely simple and lightweight tool that check for updates in your [Runtipi](https://github.com/runtipi/runtipi) server and notifies you through your favorite notification system.

> [!NOTE]
> Tipimate v2 and forward only supports runtipi instances from version v4 and above due to the API changes.

## Getting started

You can run tipimate with two ways, either docker or binary. If you chose binary, you can grab the latest binary from the [releases](https://github.com/steveiliop56/tipimate/releases) page, then `chmod +x tipimate` and finally you can run it with `./tipimate`. _assuming the binary is named tipimate_

Running with docker is also very easy, you just need to download the docker compose file from [here](./docker-compose.yml) and run tipimate with `docker compose up -d`. _make sure to change the environment variables accordingly_

If you prefer a docker run command, you can run it with:

```bash
docker run -t -d --name tipimate -v ./data:/data -e TIPIMATE_NOTIFICATION_URL=some_shoutrrr_url -e TIPIMATE_RUNTIPI_URL=your_runtipi_url -e TIPIMATE_JWT_SECRET=your_jwt_secret ghcr.io/steveiliop56/tipimate:v2
```

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

If you are interested in helping with the development feel free to create a pull request or an issue about a bug or a feature. All kinds of contributions are welcome!

## License

Tipimate is licensed under the GNU General Public License v3.0. TL;DR — You may copy, distribute and modify the software as long as you track changes/dates in source files. Any modifications to or software including (via compiler) GPL-licensed code must also be made available under the GPL along with build & install instructions.
