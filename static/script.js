var socket = null;

// Only scroll when the user is at the bottom of the messages list
const SCROLL_MESSAGES_DISTANCE = 100;

document.addEventListener("DOMContentLoaded", function (event) {
    socket = new WebSocket("ws://192.168.1.199:8000/api/v1/realtime/ws");

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
        let text = event.data;
        console.log(text);
        console.log("GOT MESSAGE FROM SERVER", event);
        createMessageElement(text, false);
        scrollToBottomIfAtEnd();
    }

    document.getElementById("message-form").addEventListener("submit", function (e) {
        e.preventDefault();
        sendmessage();
    });
})

function sendmessage() {
    console.log("Sending message")
    let messageInput = document.getElementById("message-input")
    console.log("Sending text ", messageInput.value)
    if (messageInput.value) {
        socket.send(messageInput.value);
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