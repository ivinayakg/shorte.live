# shorte.live

- An open-source url-shortner.

## Info

- Build with:
  - Golang
  - Mux (Gorilla)
  - MongoDB
  - Redis
  - React
  - Typescript
  - TailwindCSS
  - Shadcn/UI

## Setup Locally

1. Clone the project locally with
   ```
   git clone https://github.com/ivinayakg/shorte.live.git
   ```
2. Setup env files
   - for frontend env variables, run the following command from the repo root, in the terminal.
     ```
     cp client/sample.env client/.env
     ```
   - for backend env variables, run the following command from the repo root, in the terminal.
     ```
     cp api/sample.env api/.env
     ```
3. Setup Google Oauth keys

   - Currently the service only have one primary means of authentication, and that is with Google OAuth. [Read in detail here](https://developers.google.com/identity/protocols/oauth2)
   - You need to obtain your secret key and client id to make this work.
     1. Login to your google console [here](https://console.cloud.google.com).
     2. Create a new project [here](https://console.cloud.google.com/projectcreate).
     3. Open the API dashboard, and go to `OAuth consent screen` [here](https://console.cloud.google.com/apis/credentials/consent). Once there configure your screen, read more [here](https://developers.google.com/workspace/guides/configure-oauth-consent)
     4. Now navigate to `credentials` and create a new one.
        - add `http://localhost:3100` into `Authorized JavaScript origins`
        - add `http://localhost:3100/user/google/callback` into `Authorized redirect URIs`
        - Read more [here](https://developers.google.com/identity/protocols/oauth2/web-server#creatingcred)
     5. Once done copy your `client id` and your `client secret`, and paste in your `api/.env` file for the vairable `GOOGLE_OAUTH_CLIENT_ID, GOOGLE_OAUTH_CLIENT_SECRET` and done.

4. Setup Local DB and Local Redis

   - You need `Docker Desktop` installed for this on your system. [check here](https://docs.docker.com/desktop/)
   - Once verified, open your `Docker Desktop` and keep it running in the background.
   - Run the follwing command from the repo terminal
     ```
     docker-compose -f config/compose.yml up
     ```
   - If you want it to be running in the background add the `-d` flag.
     ```
     docker-compose -f config/compose.yml up -d
     ```

5. Test
   - Open `client` directory in your terminal and run `yarn dev`
   - Open `api` directory in your terminal and run `go run .`
   - and done

## Deploy

- Coming soon...

## Contribution

- Hey, thank you for wanting to contribute to this project. This project is really special to me, so my sincere gratitude to you for interacting with this project to whatever extent you have. Thank you.
- You can start with contributing from [here]()
