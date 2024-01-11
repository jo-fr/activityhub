# backend

This folder contains the api.

## endpoints

### Activity Pub

These are the endpoints that a relavant for the activity pub communication

#### GET Webfinger

```
/.well-known/webfinger?resource=acct:user@instanstance.example
```

Endpoint for user discovery

Response:

```json
200 OK
{
  "subject": "acct:user@instanstance.example",
  "links": [
    {
      "rel": "self",
      "type": "application/activity+json",
      "href": "https://instance.example/ap/user"
    },
    {
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html",
      "href": "https://activityhub.web.app/feed/user"
    }
  ]
}
```

#### GET User

```
/ap/exampleuser
```

Get data about user

Response:

```json
200 OK
{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://w3id.org/security/v1"
  ],
  "id": "https://instance.example/ap/exampleuser",
  "type": "Service",
  "following": "https://instance.example/ap/exampleuser/following",
  "followers": "https://instance.example/ap/exampleuser/followers",
  "preferredUsername": "exampleuser",
  "name": "DRINNIES ActivityHub Bot",
  "summary": "Description for example user account",
  "url": "https://instance.example/api/users/exampleuser/redirect",
  "published": "2024-01-09T10:58:26Z",
  "inbox": "https://instance.example/ap/exampleuser/inbox",
  "publicKey": {
    "id": "https://instance.example/ap/exampleuser#main-key",
    "owner": "https://instance.example/ap/exampleuser",
    "publicKeyPem": "-----BEGIN RSA PUBLIC KEY----- ... -----END RSA PUBLIC KEY-----\n"
  },
  "attachment": null
}
```

#### GET User followings

```
/ap/exampleuser/following
```

Get followings of user. This endpoint is just implemented to meet the requirements. Its always returning 0 because all users on this instance are bots and cant follow other users.

Response:

```json
200 OK
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "https://instance.example/ap/exampleuser/following",
  "type": "OrderedCollection",
  "totalItems": 0,
  "orderedItems": []
}
```

#### GET User followers

```
/ap/exampleuser/followers
```

Get followers of user.

Response:

```json
200 OK
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "https://instance.example/ap/exampleuser/followers",
  "type": "OrderedCollection",
  "totalItems": 1,
  "orderedItems": ["https://mastodon.instance/users/examplefollower"]
}
```

#### POST User inbox

```
/ap/exampleuser/inbox
```

Post to inbox of user. Request bodies are just stacically checked. The request get handeled asynconous
Request example:

```json
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "https://mastodon.green/b9154d43-5c4e-4b3f-a66f-2dec121622ce",
  "type": "Follow",
  "actor": "https://mastodon.green/users/user_that_wants_to_follow",
  "object": "https://instance.example/exampleuser"
}
```

Response:

```json
202 Accepted
```

### Web App endpoints

These are the endpoints that serve the web app

#### POST Add new feed

```
/.well-known/api/feeds
```

Add new feed. account automatically gets created.

Request example:

```json
{
  "feedURL": "https:/example.website/example/feeds"
}
```

Response:

```json
201 Created
{
    "createdAt": "2024-01-08 15:28:33.707272 +0000 UTC",
    "id": "a4baf330-55a9-45af-a50b-1441e17c835c",
    "name": "example feed",
    "type": "RSS",
    "feedURL": "https:/example.website/example/feed",
    "host": "https:/example.website/info",
    "author": "",
    "description": "Recent content on example feed",
    "imageURL": "",
    "account": {
        "createdAt": "2024-01-08 15:28:33.701409 +0000 UTC",
        "id": "f548e62b-20ab-4121-a3a5-8ab392996edd",
        "username": "examplefeed_activityhub",
        "name": "Example Feed Bot",
        "uri": "examplefeed_activityhub@example.instance"
    }
}
```

#### GET list feeds

```
//api/feeds
```

Add new feed. account automatically gets created.

Response:

```json
200 OK
{
    "total": 10,
    "items": [
        {
            "createdAt": "2024-01-08 15:28:33.707272 +0000 UTC",
            "id": "a4baf330-55a9-45af-a50b-1441e17c835c",
            "name": "example feed",
            "type": "RSS",
            "feedURL": "https:/example.website/example/feed",
            "host": "https:/example.website/info",
            "author": "",
            "description": "Recent content on example feed",
            "imageURL": "",
            "account": {
                "createdAt": "2024-01-08 15:28:33.701409 +0000 UTC",
                "id": "f548e62b-20ab-4121-a3a5-8ab392996edd",
                "username": "examplefeed_activityhub",
                "name": "Example Feed Bot",
                "uri": "examplefeed_activityhub@example.instance"
            }
        },
//...
    ]
}
```

#### GET feed by ID

```
/api/feeds/a4baf330-55a9-45af-a50b-1441e17c835c
```

Response:

```json
200 OK
{
    "createdAt": "2024-01-08 15:28:33.707272 +0000 UTC",
    "id": "a4baf330-55a9-45af-a50b-1441e17c835c",
    "name": "example feed",
    "type": "RSS",
    "feedURL": "https:/example.website/example/feed",
    "host": "https:/example.website/info",
    "author": "",
    "description": "Recent content on example feed",
    "imageURL": "",
    "account": {
        "createdAt": "2024-01-08 15:28:33.701409 +0000 UTC",
        "id": "f548e62b-20ab-4121-a3a5-8ab392996edd",
        "username": "examplefeed_activityhub",
        "name": "Example Feed Bot",
        "uri": "examplefeed_activityhub@example.instance"
    }
}
```

#### GET status from feed with feed ID

```
/api/feeds/a4baf330-55a9-45af-a50b-1441e17c835c/status?offset=0&limit=10
```

Response:

```json
200 OK
{
    "total": 10,
    "items": [
        {
            "id": "d2502991-fd7d-4a00-ab75-2a2d9647a063",
            "createdAt": "2024-01-11T15:17:46.946041Z",
            "content": "<p><strong>Example Status</strong><br/>This is an example status.</br><a href=\"https://www.example.de//example/resource\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">https://www.example.de//example/re...</a></p>",
            "accountID": "f548e62b-20ab-4121-a3a5-8ab392996edd"
        },
        //..
    ]
}
```

#### GET feed by username

```
/api/users/exampleusername
```

Add new feed. account automatically gets created.

Response:

```json
200 OK
{
    "createdAt": "2024-01-08 15:28:33.707272 +0000 UTC",
    "id": "a4baf330-55a9-45af-a50b-1441e17c835c",
    "name": "example feed",
    "type": "RSS",
    "feedURL": "https:/example.website/example/feed",
    "host": "https:/example.website/info",
    "author": "",
    "description": "Recent content on example feed",
    "imageURL": "",
    "account": {
        "createdAt": "2024-01-08 15:28:33.701409 +0000 UTC",
        "id": "f548e62b-20ab-4121-a3a5-8ab392996edd",
        "username": "examplefeed_activityhub",
        "name": "Example Feed Bot",
        "uri": "examplefeed_activityhub@example.instance"
    }
}
```

#### GET feed by username

```
/api/users/exampleusername/redirect
```

redirects to configured web app.

Response:

```json
308 Permanent Redirect
```
