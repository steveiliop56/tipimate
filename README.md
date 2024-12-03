# TipiCord

TipiCord is an extremely simple tool that periodically checks for app updates in your [Runtipi](https://github.com/runtipi/runtipi) instance and then sends notifications to your prefered Discord channel/server. It is super fast, lightweight and only ~20mb in size.

> [!NOTE]
> TipiCord only supports runtipi instances from v3.7.0 and above since this version added the ability to use an API to communicate with the server.

> [!WARNING]
> TipiCord is still in early stages of development so issues are to be expected. If you encounter any please create an issue so I can fix them as soon as possible.

## Roadmap

This project is still in early stages of development so it only includes the basic features but I am planning to add the following ones too.

- [ ] Multiple instances
- [x] Configuration using environment variables (docker)
- [ ] Support for other notifications services (project rename)
- [x] Possibly a CLI check mode like [Cup](https://github.com/sergi0g/cup)

## Getting started

You can run tipicord with two ways, docker or binary. If you chose binary, you can grab the latest binary from the [releases](https://github.com/steveiliop56/tipicord/releases) page, then `chmod +x tipicord` and finally you can run it with `./tipicord`. *assuming the binary is named tipicord*

Running with docker is also very easy, you just need to download the docker compose file from [here](./docker-compose.yml) and run tipicord with `docker compose up -d`. *make sure to change the environment variables accordingly*

If you prefer docker run command you can run it with

```bash
docker run -t -d --name tipicord -v ./data:/data -e DISCORD=your-discord-url -e RUNTIPI=your-runtipi-url -e JWT_SECRET=your-jwt-secret ghcr.io/steveiliop56/tipicord:latest
```

## Building

To build the project you need to have Go and Git installed. 

You firstly have to clone the repository with

```bash
git clone https://github.com/steveiliop56/tipicord
cd tipicord
```

Then install dependencies

```bash
go mod tidy
```

And finally run it with

```bash
go run .
```

Or build it with

```bash
got build
```

If everything succeeds you should have a binary named `tipicord`.

> [!NOTE]
> You can also build for other operating systems/architectures using `GOOS=windows` and `GOARCH=arm64`.

> [!NOTE]
> You can also run a "development" docker compose file by copying the `.env.example` file to `.env`, changing your environment variables and running `docker compose -f docker-compose.dev.yml up --build --force-recreate`. With this way you can test your changes in the docker image too.

## Contributing

This project is still in early stages of development so bugs are to be expected. If you are interested in helping with my terrible go skills, feel free to create a pull request. Any help is appreciated!

## License

TipiCord is licensed under the GNU General Public License v3.0. TL;DR — You may copy, distribute and modify the software as long as you track changes/dates in source files. Any modifications to or software including (via compiler) GPL-licensed code must also be made available under the GPL along with build & install instructions.
