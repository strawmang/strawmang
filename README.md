# strawmang

Strawmang is a simple chat application.  

There can be up to three chats running at any time.  If there are three chats
you can vote to kill a thread.  If the majority of people connected vote to
kill the chat disappears and a new one can be created.

# API/endpoints

None of this is final yet and I'm just using this as a base to figure out how
the data structures will be orginized 

All endpoints and websocket can optioannly return an error field that will
indicate if any errors are returned

`/chats`

Returns a JSON response listing the available chat endpoints.  Expected response:
``` JSON
{
  "chats": [
  {
    "id": 5,
      "topic": "Tea vs coffee",
      "started": 1440290981,
      "ends": 1440294201,
      "url": "/chat/5",
  },
  {
    "id": 6,
    "topic": "Gentoo vs Arch",
    "started": 1440290981,
    "ends": 1440294201,
    "url": "/chat/6",
  },
  {
    "id": 7,
    "topic": "Javascript vs Coffeescript",
    "started": 1440290981,
    "ends": 1440294201,
    "url": "/chat/7",
  },
  ],
}
```

`/users`

Returns information about the set of online users
``` JSON
{
  "online": 12,
}

```


