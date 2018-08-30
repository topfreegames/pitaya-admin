# Pitaya Admin

[Pitaya's](https://github.com/topfreegames/pitaya) games management and monitoring interface

# Getting Started

Setup dependencies (make sure you have dep installed)
```
make setup
```
Start pitaya admin
```
make run
```

# API

**Errors**
---
Whenever pitaya admin respond with an error, it will be in the following format:
`````
{"success":false, "message":[msg], "reason": [error]}
`````

**Client Request**
---
Establishes a websocket connection with pitaya admin that can be used to send requests to handlers and receives its outputs through a pitaya client connected to given server address.

* **URL**

    /request

* **Method**

    `GET`

* **URL Params**

    **Required:**

    `address=[pitaya server address]`

* **Data Params**

    After the websocket connection is set, pitaya admin will listen to requests in the following structure:

    `````
        Route       Pitaya route to handler
        Payload     The data that will be sent to the handler
        IsRequest   Whether this should be sent as a request or a notify
    `````

* **Response Example**

    After the websocket connection is set and the pitaya client connects to the given server address, any message sent to the client will be forwarded as a response. It will return an error if it fails to establish any of the connections.


**Documentation**
---
Returns server auto documentation given its type. Target server must implement pitaya's auto documentation remote.
* **URL**

    /docs

* **Method**

    `GET`

* **URL Params**

    **Optional:**

    `methodtype=[either "remote" or "handler"]`

    Specifies whether to get documentation only for remote's or handler's methods.

    `route=[pitaya route]`

    Specifies a single method route to get its documentation. Must have specified `methodtype` as well.

    **Required:**

    `type=[server type]`


* **Response Example**

   Auto documentation as implemented in remote or an error.




**Kick Users**
---
Send a kick packet to users in a given user list

* **URL**

    /user/kick

* **Method**

    `POST`

* **Data Params**

    A kick message struct:
    ````
        Uids          Array with users IDs
        FrontendType  The type of the frontend server the users are connected to
    ````


* **Response Example**

    `{"success" : true}`

    If the user were kicked or an error otherwise

**List Servers**
---
Returns a json array with servers information

* **URL**

    /servers

* **Method**

    `GET`

* **URL Params**

    **Optional:**

    `type=[server type]`

    If no server type is given, all servers are returned

* **Response Example**

    `{"success" : true, "response":[{"id":"740bcdba-1a9b-4790-9f07-ddf9f33a912e","type":"admin","metadata":{},"frontend":false}]}`

**RPC**
---
Sends a RPC to a pitaya server. Target server must implement pitaya's autodoc feature and a route that provides protobuff descriptors based on their filename.

* **URL**

    /rpc

* **Method**

    `POST`

* **Data Params**

    A RPC message struct:
    `````
        Route          Remote's route
        FrontendType   Remote's type
        ServerID       Remote's server ID (can be empty)
        Meta           Remote's args data serialized
        ProtoSource    Proto source file name
    `````

    Since pitaya's RPC methods takes protos as an argument and returns protos as a response, the proto source code filename is used to create dynamic messages and send it to a given method specified by its route and type. A server id can be specified if one wants the method to be sent to given server.

* **Response Example**

    `{"success" : true, "response":{"Msg":"hi im a rpc msg"}}`

    If the RPC was sent successfully, response will contain the response returned by the method, otherwise, it will be an error.


**Send Push**
---
Sends a push notification to users in a given user list

* **URL**

    /user/push

* **Method**

    `POST`

* **Data Params**

    A push message struct:
    ````
        Uids          Array with users IDs
        Route         Route to push
        Message       Notification message
        FrontendType  The type of the frontend server the users are connected to
    ````


* **Response Example**

    `{"success" : true}`

    If the message was sent or an error otherwise.
