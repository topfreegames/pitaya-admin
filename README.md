# Pitaya Admin [![Build Status][1]][2] [![GoDoc][3]][4] [![Go Report Card][5]][6] [![MIT licensed][7]][8]

[1]: https://travis-ci.org/topfreegames/pitaya-admin.svg?branch=master
[2]: https://travis-ci.org/topfreegames/pitaya-admin
[3]: https://godoc.org/github.com/topfreegames/pitaya-admin?status.svg
[4]: https://godoc.org/github.com/topfreegames/pitaya-admin
[5]: https://goreportcard.com/badge/github.com/topfreegames/pitaya-admin
[6]: https://goreportcard.com/report/github.com/topfreegames/pitaya-admin
[7]: https://img.shields.io/badge/license-MIT-blue.svg
[8]: LICENSE

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

## **Errors**

Whenever pitaya admin responds with an error, it will be in the following format:

```
{"success":false, "message":[msg], "reason": [error]}
```

## **Client Request**

Establishes a websocket connection with pitaya admin that can be used to send requests to handlers and receive theirs outputs through a pitaya client connected to given server address.

- **URL**

  /request

- **Method**

  `GET`

- **URL Params**

  **Required:**

  `address=[pitaya server address]`

- **Data Params**

  After the websocket connection is set, pitaya admin will listen to requests in the following structure:

  ```
      Route       Pitaya route to handler
      Payload     The data that will be sent to the handler
      IsRequest   Whether this should be sent as a request or a notify
  ```

- **Response Example**

  After the websocket connection is set and the pitaya client connects to the given server address, any message sent to the client will be forwarded as a response. It will return an error if it fails to establish any of the connections.

## **Documentation**

Returns server auto documentation given its type. Target server must implement pitaya's auto documentation remote.

- **URL**

  /docs

- **Method**

  `GET`

- **URL Params**

  **Optional:**

  `methodtype=[either "remote" or "handler"]`

  Specifies whether to get documentation only for remote's or handler's methods.

  `route=[pitaya route]`

  Specifies a single method route to get its documentation. Must have specified `methodtype` as well.

  `getProtos=["1"]`

  Specifies if the documentation will have protobuff message names when set to 1.

  **Required:**

  `type=[server type]`

* **Response Example**

  Auto documentation as implemented in remote or an error.

* **Requirements**

  As stated above, in order for documentation routes to work, you will need to implement a remote for pitaya's auto documentation. A simple example is presented below.

  ```golang
      func (c *ConnectorRemote) Docs(ctx context.Context, flag *protos.DocMsg) (*protos.Doc, error) {
          d, err := pitaya.Documentation(flag.GetGetProtos())

          if err != nil {
              return nil, err
          }
          doc, err := json.Marshal(d)

          if err != nil {
              return nil, err
          }

          return &protos.Doc{Doc: string(doc)}, nil
  }
  ```

  The boolean inside `flag` passed to `pitaya.Documentation` specifies if the documentation will have protobuff message names. Note that in order to pitaya admin RPC to work, documentation **must** have protos name.

## **Kick Users**

Send a kick packet to users in a given user list, breaking the users connections to the pitaya server

- **URL**

  /user/kick

- **Method**

  `POST`

- **Data Params**

  A kick message struct:

  ```
      Uids          Array with users IDs
      FrontendType  The type of the frontend server the users are connected to
  ```

* **Response Example**

  `{"success" : true}`

  If the user were kicked or an error otherwise

## **List Servers**

Returns a json array with servers information

- **URL**

  /servers

- **Method**

  `GET`

- **URL Params**

  **Optional:**

  `type=[server type]`

  If no server type is given, all servers are returned

- **Response Example**

  `{"success" : true, "response":[{"id":"740bcdba-1a9b-4790-9f07-ddf9f33a912e","type":"admin","metadata":{},"frontend":false}]}`

## **RPC**

Sends a RPC to a pitaya server. Target server must implement pitaya's autodoc feature and a remote that provides protobuff descriptors using reflection.

- **URL**

  /rpc

- **Method**

  `POST`

- **Data Params**

  A RPC message struct:

  ```
      Route          Remote's route
      FrontendType   Remote's type
      ServerID       Remote's server ID (can be empty)
      Meta           Remote's args data serialized
  ```

  Since pitaya's RPC methods take protos as an argument and return protos as a response, pitaya admin uses reflection alongside with pitaya's autodoc feature to build dynamic messages and send them to the server

- **Response Example**

  `{"success" : true, "response":{"Msg":"hi im a rpc msg"}}`

  If the RPC was sent successfully, response will contain the response returned by the method, otherwise, it will be an error.

- **Requirements**

  As stated, in order to RPC work properly, you will need a remote that provides protobuff descriptors via reflection **and** documentation's remote. This allows pitaya-admin to work with RPCs without any proto specification at compile time. An remote example is presented below

  ```golang
  func (c *ConnectorRemote) Docs(ctx context.Context, flag *protos.DocMsg) (*protos.Doc, error) {
  		d, err := pitaya.Documentation(flag.GetGetProtos())

  		if err != nil {
  			return nil, err
  		}
  		doc, err := json.Marshal(d)

  		if err != nil {
  			return nil, err
  		}

  		return &protos.Doc{Doc: string(doc)}, nil
  }

  func (c *ConnectorRemote) Proto(ctx context.Context, name *protos.ProtoName) (*protos.ProtoDescriptor, error) {
       protoName := name.Name
       protoReflectTypePointer := proto.MessageType(protoName)
       protoReflectType := protoReflectTypePointer.Elem()
       protoValue := reflect.New(protoReflectType)
       descriptorMethod, ok := protoReflectTypePointer.MethodByName("Descriptor")

       if !ok {
           return nil, errors.New("failed to get proto descriptor")
       }

       descriptorValue := descriptorMethod.Func.Call([]reflect.Value{protoValue})
       protoDescriptor := descriptorValue[0].Bytes()
       return &protos.ProtoDescriptor{
           Desc: protoDescriptor,
       }, nil
    }
    ```

    It is important to note that ``proto`` refers to [golang/proto](https://github.com/golang/protobuf/tree/master/proto).

## **Send Push**

Sends a push notification to users in a given user list

- **URL**

  /user/push

- **Method**

  `POST`

- **Data Params**

  A push message struct:

  ```
      Uids          Array with users IDs
      Route         Route to push
      Message       Notification message
      FrontendType  The type of the frontend server the users are connected to
  ```

* **Response Example**

  `{"success" : true}`

  If the message was sent or an error otherwise.

# License

[MIT License](./LICENSE)
