# Architecture of application

The application is divided into two layers, the realtime delivery layer, and the persistence layer. In the realtime layer, all objects are in memory and temporary. They can be recreated when necessary.

Realtime layer: Hub, Room, Client

Persistence layer: User, Group

A User can have many connections to the server (say from different devices), each of these connections is called a `Client`. A `Room` is mapped to a `Group`, and rooms are only used for realtime message and events delivery.

To simplify the implementation, WS is used only for online status, typing indicator, and receiving new messages. Sending a message uses the REST endpoint instead of using websocket.