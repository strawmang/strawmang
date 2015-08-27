# Events

## login

Client -> Server

required fields: `username`

## message

Client <-> Server

required fields C->S: `text` `topic-id`

required fields S-C: `text` `color` `topic-id`

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

Just getting it all stubbed out.  Eventually we will have more paramaters to change the 
behavior of the topic

Used by the server to notify clients of a new topic

Used by the client to request a new topic

Client <-> Server

required fields: C->S: `option-a` `option-b`
required fields S->C: `topic-id`?
