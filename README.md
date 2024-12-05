# TipiMate

TipiMate is an extremely simple tool that periodically checks for app updates in your [Runtipi](https://github.com/runtipi/runtipi) instance and then sends notifications to your prefered Discord channel/server. It is super fast, lightweight and only ~20mb in size.

> [!NOTE]
> TipiMate only supports runtipi instances from v3.7.0 and above since this version added the ability to use an API to communicate with the server.

> [!WARNING]
> TipiMate is still in early stages of development so issues are to be expected. If you encounter any please create an issue so I can fix them as soon as possible.

## Roadmap

This project is still in early stages of development so it only includes the basic features but I am planning to add the following ones too.

- [ ] Multiple instances
- [x] Configuration using environment variables (docker)
- [ ] Support for other notifications service
- [x] Project rename
- [x] Possibly a CLI check mode like [Cup](https://github.com/sergi0g/cup)
- [ ] Check for main Runtipi version

## Getting started

You can run tipimate with two ways, docker or binary. If you chose binary, you can grab the latest binary from the [releases](https://github.com/steveiliop56/tipimate/releases) page, then `chmod +x tipimate` and finally you can run it with `./tipimate`. *assuming the binary is named tipimate*

Running with docker is also very easy, you just need to download the docker compose file from [here](./docker-compose.yml) and run tipimate with `docker compose up -d`. *make sure to change the environment variables accordingly*

If you prefer docker run command you can run it with

```bash
docker run -t -d --name tipimate -v ./data:/data -e DISCORD=your-discord-url -e RUNTIPI=your-runtipi-url -e JWT_SECRET=your-jwt-secret ghcr.io/steveiliop56/tipimate:latest
```

> [!TIP]
> You can set the `--runtipi-internal` flag or the `RUNTIPI_INTERNAL` environment variable to something like `http://localhost` if TipiMate is running on the same server as your Runtipi server and then set the `--runtipi` flag or `RUNIPI` to the public URL of your instance e.g. `https://runtipi.mydomain.com` so TipiMate can both connect directly to Runtipi and show the correct URL on Discord. 

## Building

To build the project you need to have Go and Git installed. 

You firstly have to clone the repository with

```bash
git clone https://github.com/steveiliop56/tipimate
cd tipimate
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
go build
```

If everything succeeds you should have a binary named `tipimate`.

> [!NOTE]
> You can also build for other operating systems/architectures using `GOOS=windows` and `GOARCH=arm64`.

> [!NOTE]
> You can also run a "development" docker compose file by copying the `.env.example` file to `.env`, changing your environment variables and running `docker compose -f docker-compose.dev.yml up --build --force-recreate`. With this way you can test your changes in the docker image too.

## Contributing

This project is still in early stages of development so bugs are to be expected. If you are interested in helping with my terrible go skills, feel free to create a pull request. Any help is appreciated!

## License

TipiMate is licensed under the GNU General Public License v3.0. TL;DR — You may copy, distribute and modify the software as long as you track changes/dates in source files. Any modifications to or software including (via compiler) GPL-licensed code must also be made available under the GPL along with build & install instructions.
