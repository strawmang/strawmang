# strawmang

Strawmang is a simple chat application.  

There can be three debates where people can anonymously or not debate about a topic.
People can at anytime start a vote to kill the current topic where it will get deleted
and possibly replaced by another topic. Topics can also be pruned after a certain amount
of time.

The exact details haven't been completely established yet so it's still up in the air.

# API/endpoints

None of this is final yet and I'm just using this as a base to figure out how
the data structures will be organized 

All endpoints and websocket can optionally return an error field that will
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


