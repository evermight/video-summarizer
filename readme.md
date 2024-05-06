
# Installation

Copy `env.sample` as `.env`.  See the ChatGPT Setup and Google Setup on how to get the values necessary to fill out a `.env` file.

To get the `GOOGLE_API_TOKEN` needed in the `.env` file, you will need to run the `./oauth-step-1.sh` to generate a url to start the OAuth process.

You will eventually be led back to the `http://localhost...` as specified in `14.png` below.  Take the code from `code=[code]` (without the `code=`) and run `./oauth-step-2.sh -c "[code]"` (you may wish to quote the code).  You will find the API token in the `youtube/tokens.json`, which you can then use it to populate the `GOOGLE_API_TOKEN` in the `.env` file.

Start up the website with:

```
cp ~/.env ~/www/.env
docker-compose up --build -d
docker exec -it [name of the container] bash
go run main.go

```

Now visit `http://localhost:8181`.

You can make a production build of the `~/www/main.go` with `go build main.go` from within the docker container. More details to come.

## ChatGPT Setup

### ChatGPT API Key

Get your ChatGPT API key from https://platform.openai.com/api-keys as shown in the screenshot:
![step](/documentation/chatgpt.png)

## Google Setup

Follow screenshots below to get the necessary api keys, tokens and secrets to acces Google/YouTube APIs:

### Setup Project

![step](/documentation/00.png)
![step](/documentation/01.png)

### Setup API Key

![step](/documentation/02.png)
![step](/documentation/03.png)
![step](/documentation/04.png)
![step](/documentation/05.png)
![step](/documentation/06.png)
![step](/documentation/07.png)
![step](/documentation/08.png)
![step](/documentation/09.png)

### Setup OAuth Screen

![step](/documentation/10.png)
![step](/documentation/11.png)
![step](/documentation/12.png)

### Setup OAuth Client ID and Client Secret

![step](/documentation/13.png)
![step](/documentation/14.png)
![step](/documentation/15.png)
