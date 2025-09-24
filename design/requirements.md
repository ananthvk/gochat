# Real Time Chat Application

The goal of this project is to build a real time chat application that supports notifications, realtime message delivery, and retrieval of past chats.

## Requirements

- Realtime message delivery
- Login/Authentication
- Direct chats
- Group chats
- Persistence of chats
- Read receipts (status of messages)
- Notifications
- Different types of messages (media/text/video/etc)
- Admin/member roles
- Invite links
- Profile view permissions
- Cache messages on client side so that it does not need to be fetched repeatedly

Future enhancements
- End to end encryption
- Login/signup through qr codes
- Mobile app
- Typing/online indicators

## Tech Stack

- Golang (backend)
- Typescript + React (frontend)
- Postgresql (DB) _Note: Find out if SQL is actually a good choice for a chat app_
- Redis (caching)

## Assumptions

- A monolithic application, i.e. only a single server handles all websockets

## Architecture

The application will be split into four modules:
- Realtime service (handles realtime delivery of events & messages)
- Authentication service (handles signup/login/authorization/user)
- Message store service (handles retrieval of past chats and filters)
- Notification service (push notification/in app notification/etc)

## Data Model

- User (user_id, username, email, password, signup_date, role, status)
- Profile (name, profile pic, about, last_seen, user_id)
- Message (message_id, messagetype, sender_id, room_id, timestamp, content, content_type)
- MessageStatus (message_id, user_id, status, updated_at)
- Room (room_id, created_at, name, description, icon, created_by_user_id)
- RoomMembers (room_id, user_id, role, join_date)
- InviteLink (room_id, link_id, use_count, expired, created_at)

## MVP

The user should be able to send messages, if the recipients are online, they receive the messages in realtime. Otherwise the messages are stored, and once the receiver comes online, the message is delivered. The user should be able to communicate with another user, and should be able to join and create groups

## Flows

### Send Message Flow

#### System Flow

1. Client sends a message, with {message content, content type, room id, auth token} to the realtime service
2. The controller verifies if the client is authenticated, and they are allowed to send message to the room
3. Then the controller saves the message to the persistence layer, it receives a {message_id} in return
4. Then the realtime service fans out, sending the message event to currently connected clients in the same room
5. The clients acknowledge the delivery, after which the message status is recorded in the database

#### Implementation Notes

- If there is an error when saving the message to the database, the server sends back an error to the sender. The sender must back off and try sending the message after some time
- The server should only mark the message as delivered once the client acknowledges
- The server should not send back the message to client who sent it

#### Tasks

- [x] Define a MessageSent websocket event
- [ ] Implement token validation
- [x] Implement broadcasting of message event to connected clients in room
- [ ] Implement persistence of message
- [ ] Implement handling of message status (delivered/read)

### Receive Message Flow

Suppose user A has sent a message M to a room 123, which has members X, Y, Z

#### System Flow

1. The service first checks if the user is still a member of the room
2. If the user is still a member of the room, the service checks if the member is currently online
3. If the user is online, the client pushes the message  {message content, message id, room id} to the client
4. The client acknowledges the message
5. Once the client acknowledges the message, the message is marked as delivered in the database
6. If the user is offline, the server waits for the user to become online
7. All undelivered messages after the last acknowledged timestamp are sent to the client in batches once the client reconnects

#### Implementation Notes

- Should not silently fail, the client MUST send an acknowledgement (`Delivered` message)
- If the client is offline, and there are undelivered messages for the client, but later the client is removed from the room (say by an admin), the client should still get all messages until the removal point

#### Tasks

- [x] Implement `chat_message` web socket event message structure
- [ ] Implement `delivered` web socket event structure
- [ ] Implement receiving undelivered messages

### Past Messages Flow

#### System Flow

1. The user clicks on a room (chat/group), and scrolls up to view history
2. Older messages are not yet downloaded, so the client requests it from the history service
3. The history services validates if the token is valid
4. The history service sends messages in batches for a particular room

#### Implementation Notes

- The client should not be able to fetch too many messages at once (say max 50-100 messages in one scroll/ depending upon screen size or timestamp)
- The client should be able to see history before joining (like telegram)