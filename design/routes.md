## API routes

| Status | Method |Path|Description|
|--------|--------|----|-----------|
| todo   | GET    |`/api/v1/health` | Health check |
| done   | GET    |`/api/v1/realtime/ws` | Upgrades the connection to websocket protocol| 
| done   | POST   |`/api/v1/realtime/room` | Creates a new room with the given name, and returns the id of the created room|
| done   | POST   |`/api/v1/realtime/join` | Body must contain the client id & the room id, this action adds the client to the room|
| done   | GET    |`/api/v1/realtime/by-name/{name}` | Returns the room which has the given name, for now rooms have unique names|
| done   | GET    |`/api/v1/realtime/room` | Returns a list of all the active rooms|
| todo   | POST   |`/api/v1/groups` | Creates a new group & makes the creating user the admin of the room|
| todo   | GET    |`/api/v1/groups/{id}` | Returns details of the group |
| todo   | DELETE |`/api/v1/groups/{id}` | Deletes the group, it's associated room (if any), and other data related to the room|
| todo   | PATCH  |`/api/v1/groups/{id}` | Update group details |
| todo   | POST   |`/api/v1/groups/{id}/members` | Adds a member to the room, the `userId`specified in the body is added. If it's not specified, the current user is added|
| todo   | GET    |`/api/v1/groups/{id}/members` | Returns a list of users in the group |
| future | POST   |`/api/v1/auth/signup` | Creates a new user |
| future | POST   |`/api/v1/auth/login`  | TBD (JWT/Session/Token) |


Note: The `room` api endpoints will be removed in the future, after groups get implemented. After which rooms will only be exposed internally.