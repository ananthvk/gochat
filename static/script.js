var socket = null;

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
        console.log(text)
        console.log("GOT MESSAGE FROM SERVER", event)

        let newNode = document.createElement("div")
        newNode.classList.add("message", "message-other")
        newNode.textContent = text

        document.getElementById("messages").appendChild(newNode)

        // TODO: Only scroll when the user is already at the end of the div, i.e. don't scroll if the user is reading a message up
        const messagesDiv = document.getElementById('messages');
        messagesDiv.scrollTo({
            top: messagesDiv.scrollHeight,
            behavior: 'smooth'
        });
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
        socket.send(messageInput.value)

        let newNode = document.createElement("div")
        newNode.classList.add("message", "message-self")
        newNode.textContent = messageInput.value
        document.getElementById("messages").appendChild(newNode)

        messageInput.value = "";

        const messagesDiv = document.getElementById('messages');
        messagesDiv.scrollTo({
            top: messagesDiv.scrollHeight,
            behavior: 'smooth'
        });
    }
}
