## API routes

| Status | Method |Path|Description|
|--------|--------|----|-----------|
| done   | GET    |`/api/v1/health` | Health check |
| done   | GET    |`/api/v1/realtime/ws` | Upgrades the connection to websocket protocol| 
| done   | POST   |`/api/v1/realtime/room` | Creates a new room with the given name, and returns the id of the created room|
| done   | POST   |`/api/v1/realtime/join` | Body must contain the client id & the room id, this action adds the client to the room|
| done   | GET    |`/api/v1/realtime/by-name/{name}` | Returns the room which has the given name, for now rooms have unique names|
| done   | GET    |`/api/v1/realtime/room` | Returns a list of all the active rooms|
| done   | POST   |`/api/v1/group` | Creates a new group & makes the creating user the admin of the room|
| done   | GET    |`/api/v1/group` | Return all the groups the user is a part of (max limit of 256 groups) |
| done   | GET    |`/api/v1/group/{id}` | Returns details of the group |
| done   | DELETE |`/api/v1/group/{id}` | Deletes the group, it's associated room (if any), and other data related to the room|
| done   | PATCH  |`/api/v1/group/{id}` | Update group details |
| todo   | POST   |`/api/v1/group/{id}/members` | Adds a member to the room, the `userId`specified in the body is added. If it's not specified, the current user is added|
| todo   | GET    |`/api/v1/group/{id}/members` | Returns a list of users in the group |
| done   | GET    |`/api/v1/group/{id}/message?before=<id>&limit=<n>` | Get messages in a group, implements cursor based pagination, n can range from 1 to 100, it returns all messages which have id strictly less than the specified id|
| done   | DELETE |`/api/v1/group/{id}/message/{id}` | Deletes a message|
| done   | GET    |`/api/v1/group/{id}/message/{id}` | Returns detailed info about a message (TODO: later add message delivery status, read etc here)|
| done   | POST   |`/api/v1/group/{id}/message` | Creates a new message under the group and returns the id of the created message|
| future | POST   |`/api/v1/auth/signup` | Creates a new user |
| future | POST   |`/api/v1/auth/login`  | TBD (JWT/Session/Token) |


Note: The `room` api endpoints will be removed in the future, after groups get implemented. After which rooms will only be exposed internally.