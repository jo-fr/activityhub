# ActivityHub - extending the mastodon feed with external content streams

ActivityHub merges diverse content streams (news, social media, videos) into a unified Mastodon feed. Using open protocols like Activity Pub and RSS, the service makes any RSS feed compatible. Users paste the feed URL, receive a user handle, and updates from the RSS feed are posted directly into their Mastodon feed.

## repo structure

```
.
├── app           # contains source files of web app
├── backend       # contains source files of backend
└── terraform     # contains infrastructure deployment files
```

## setup

After cloning the repo there a a few steps that you should run:

```
$ mv .env.example .env
$ cd app && npm install
```

## local testing

to run everything in docker run

```
$ docker compose up
```

- api is listening on `http://localhost:8080`
- web app ist listening on `http://localhost:5173`

You can also run the api outside of docker. This can be helpful in the developing procress. For that run the following commands:

```
$ docker compose up postgresql pub-sub-emulator
$ cd backend && make run-local
```

- api is listening on `http://localhost:8080`

If you want to start the web app outside of docker with a dev server

```
$ cd app
$ sdocker ru npm install
$ npm run dev
```

- web app ist listening on `http://localhost:5173`
