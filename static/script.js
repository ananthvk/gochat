var socket = null;

var clientId = null;

const socketBase = "ws://localhost:8000"
const apiBase = "http://localhost:8000"

// Only scroll when the user is at the bottom of the messages list
const SCROLL_MESSAGES_DISTANCE = 100;

document.addEventListener("DOMContentLoaded", function (event) {
    socket = new WebSocket(`${socketBase}/api/v1/realtime/ws`);

    socket.onopen = () => {
        console.log("Connection established");
    };

    socket.onclose = event => {
        console.log("Closed connection: ", event);
    };

    socket.onerror = error => {
        console.log("Error occured: ", error);
    };

    socket.onmessage = event => {
        let message = event.data;
        const response = JSON.parse(message)
        if (response["type"] == "welcome") {
            console.log("connected to server: client id: ", response["payload"]["id"])
            clientId = response["payload"]["id"]
            document.getElementById("clientid").textContent = `Client ID: ${clientId}`
        } else if (response["type"] == "chat_message") {
            createMessageElement(response["payload"]["message"], false);
            scrollToBottomIfAtEnd();
        }
    }

    document.getElementById("message-form").addEventListener("submit", function (e) {
        e.preventDefault();
        sendmessage();
    });

    scrollToBottom();
})

document.getElementById("create-room-btn").addEventListener("click", async () => {
    const name = document.getElementById("roomname").value.trim();
    if (!name) return;
    let roomId;

    // try create
    let res = await fetch(`${apiBase}/api/v1/realtime/room`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name })
    });

    if (res.status === 201) {
        roomId = (await res.json()).id;
    } else {
        // fallback to get-by-name
        const getRes = await fetch(`${apiBase}/api/v1/realtime/room/by-name/${encodeURIComponent(name)}`);
        roomId = (await getRes.json()).id;
    }

    document.getElementById("roomid").value = roomId;

    // join room
    await fetch(`${apiBase}/api/v1/realtime/room/join`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ client_id: clientId, room_id: roomId })
    });
});

function sendmessage() {
    console.log("Sending message")
    let messageInput = document.getElementById("message-input")
    const room_id = document.getElementById("roomid").value;
    if (!room_id) {
        alert("Please specify room id")
        return
    }
    console.log("Sending text ", messageInput.value)
    if (messageInput.value) {
        const message = JSON.stringify(
            {
                "type": "chat_message",
                "payload": {
                    "room_id": room_id,
                    "message": messageInput.value
                }
            }
        )
        socket.send(message);
        createMessageElement(messageInput.value, true);
        messageInput.value = "";

        // Always scroll to bottom if the user is sending a message
        scrollToBottom();
    }
}

function scrollToBottomIfAtEnd() {
    const messagesDiv = document.getElementById('messages');

    const isNearBottom =
        messagesDiv.scrollTop + messagesDiv.clientHeight >= messagesDiv.scrollHeight - SCROLL_MESSAGES_DISTANCE;

    if (isNearBottom) {
        scrollToBottom();
    }
}

function scrollToBottom() {
    const messagesDiv = document.getElementById('messages');
    messagesDiv.scrollTo({
        top: messagesDiv.scrollHeight,
        behavior: 'smooth'
    });
}

function createMessageElement(messageContent, isSelf) {
    let newNode = document.createElement("div")
    newNode.classList.add("message", isSelf ? "message-self" : "message-other")
    newNode.textContent = messageContent
    document.getElementById("messages").appendChild(newNode)
}