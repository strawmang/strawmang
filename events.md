# Events

## login

Client -> Server

required fields: `username`

## message

Client <-> Server

required fields C->S: `text`

required fields S-C: `text` `color`

## leave

Client -> Server

## vote

Client <-> Server

Still being worked out

## error

Client <- Server

required fields: `text`

## status

This is sent from the server to notify the user of successfull events.

Currently used as a reponse to login events

Client <- Server

required fields: `text`

## newtopic

Used by the server to notify clients of a new topic

Used by the client to request a new topic

Client <-> Server

required fields: Still being worked out
required fields S->C: "topic-id"?
