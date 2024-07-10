// chat.js

// chat.js

let ws;
let currentUsername;

function authenticateAndConnect(chatName, token) {
    const formData = new FormData();
    formData.append('token', token);

    fetch(`/auth_ws/${chatName}`, {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Authentication failed');
        }
        return response.json();
    })
    .then(data => {
        if (data.message === "Authenticated") {
            currentUsername = data.username; // Сохраняем имя пользователя
            connectWebSocket(chatName);
        } else {
            throw new Error('Authentication failed');
        }
    })
    .catch(error => {
        console.error("Error during authentication:", error);
    });
}

function connectWebSocket(chatName) {
    ws = new WebSocket(`ws://${window.location.host}/ws/${chatName}`);

    ws.onmessage = function(event) {
        const message = JSON.parse(event.data);
        const messagesDiv = document.getElementById("messages");
        const messageDiv = document.createElement("div");
        // Путь до иконки пользователя
        const userIconPath = `/static/images/${message.username}.png`
        messageDiv.classList.add("message");

        if (message.username === currentUsername) {
            messageDiv.classList.add("my-message");
        } else {
            messageDiv.classList.add("other-message");
        }
        
        const messageContent = `
        <div class="photo" style="background-image: url('${userIconPath}');"></div>
        <p class="text">${message.username}: ${message.text}</p>
        `;
        messageDiv.innerHTML = messageContent;

        messagesDiv.appendChild(messageDiv);
        messagesDiv.scrollTop = messagesDiv.scrollHeight; // Прокрутка вниз
    };

    ws.onerror = function(event) {
        console.error("WebSocket error:", event);
    };

    ws.onclose = function(event) {
        console.log("WebSocket connection closed:", event);
    };
}

function sendMessage(chatName) {
    const input = document.getElementById("message");
    const message = input.value;
    if (message.trim()) {
        ws.send(JSON.stringify({ chatUid: chatName, username: currentUsername, text: message }));
        input.value = '';
    }
}

function addNewChat(chatName) {
    const discussions = document.querySelector(".discussions");

    const newChatDiv = document.createElement("div");
    newChatDiv.classList.add("discussion");
    
    const chatContent = `
        <div class="photo" style="background-image: url('https://via.placeholder.com/40');"></div>
        <div class="desc-contact">
            <p class="name">${chatName}</p>
            <p class="message">New chat created</p>
        </div>
        <div class="timer">Just now</div>
    `;

    newChatDiv.innerHTML = chatContent;
    discussions.appendChild(newChatDiv);
    discussions.scrollTop = discussions.scrollHeight; // Прокрутка вниз
}

window.onload = function() {
    const chatName = document.body.getAttribute("data-chat-name");
    const token = localStorage.getItem('token');
    authenticateAndConnect(chatName, token);

    // Добавляем обработчик события для отправки сообщения по нажатию Enter
    document.getElementById("message").addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            sendMessage(chatName);
        }
    });

    // Example of creating a new chat
    // document.querySelector(".send").addEventListener("click", function() {
    //     const newChatName = "New Chat";
    //     addNewChat(newChatName);
    // });
};









/*
let ws;
let currentUsername;

function authenticateAndConnect(chatName, token) {
    const formData = new FormData();
    formData.append('token', token);

    fetch(`/auth_ws/${chatName}`, {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Authentication failed');
        }
        return response.json();
    })
    .then(data => {
        if (data.message === "Authenticated") {
            currentUsername = data.username; // Сохраняем имя пользователя
            connectWebSocket(chatName);
        } else {
            throw new Error('Authentication failed');
        }
    })
    .catch(error => {
        console.error("Error during authentication:", error);
    });
}

function connectWebSocket(chatName) {
    ws = new WebSocket(`ws://${window.location.host}/ws/${chatName}`);

    ws.onmessage = function(event) {
        const message = JSON.parse(event.data);
        const messagesDiv = document.getElementById("messages");
        const messageDiv = document.createElement("div");
        messageDiv.classList.add("message");
        
        const messageContent = `
            <div class="photo" style="background-image: url('https://via.placeholder.com/40');"></div>
            <p class="text">${message.username}: ${message.text}</p>
        `;
        messageDiv.innerHTML = messageContent;

        messagesDiv.appendChild(messageDiv);
        messagesDiv.scrollTop = messagesDiv.scrollHeight; // Прокрутка вниз
    };

    ws.onerror = function(event) {
        console.error("WebSocket error:", event);
    };

    ws.onclose = function(event) {
        console.log("WebSocket connection closed:", event);
    };
}

function sendMessage(chatName) {
    const input = document.getElementById("message");
    const message = input.value;
    if (message.trim()) {
        ws.send(JSON.stringify({ chatUid: chatName, username: currentUsername, text: message }));
        input.value = '';
    }
}

function addNewChat(chatName) {
    const discussions = document.querySelector(".discussions");

    const newChatDiv = document.createElement("div");
    newChatDiv.classList.add("discussion");
    
    const chatContent = `
        <div class="photo" style="background-image: url('https://via.placeholder.com/40');"></div>
        <div class="desc-contact">
            <p class="name">${chatName}</p>
            <p class="message">New chat created</p>
        </div>
        <div class="timer">Just now</div>
    `;

    newChatDiv.innerHTML = chatContent;
    discussions.appendChild(newChatDiv);
    discussions.scrollTop = discussions.scrollHeight; // Прокрутка вниз
}

window.onload = function() {
    const chatName = document.body.getAttribute("data-chat-name");
    const token = localStorage.getItem('token');
    authenticateAndConnect(chatName, token);

    // Example of creating a new chat
    // document.querySelector(".send").addEventListener("click", function() {
    //     const newChatName = "New Chat";
    //     addNewChat(newChatName);
    // });
};


*/













// let ws;
// let currentUsername;

// function authenticateAndConnect(chatName, token) {
//     const formData = new FormData();
//     formData.append('token', token);

//     fetch(`/auth_ws/${chatName}`, {
//         method: 'POST',
//         body: formData
//     })
//     .then(response => {
//         if (!response.ok) {
//             throw new Error('Authentication failed');
//         }
//         return response.json();
//     })
//     .then(data => {
//         if (data.message === "Authenticated") {
//             currentUsername = data.username; // Сохраняем имя пользователя
//             connectWebSocket(chatName);
//         } else {
//             throw new Error('Authentication failed');
//         }
//     })
//     .catch(error => {
//         console.error("Error during authentication:", error);
//     });
// }

// function connectWebSocket(chatName) {
//     ws = new WebSocket(`ws://${window.location.host}/ws/${chatName}`);

//     ws.onmessage = function(event) {
//         const message = JSON.parse(event.data);
//         const messagesDiv = document.getElementById("messages");
//         const messageDiv = document.createElement("div");
//         messageDiv.textContent = `${message.username}: ${message.text}`;
//         messagesDiv.appendChild(messageDiv);
//         messagesDiv.scrollTop = messagesDiv.scrollHeight; // Прокрутка вниз
//     };

//     ws.onerror = function(event) {
//         console.error("WebSocket error:", event);
//     };

//     ws.onclose = function(event) {
//         console.log("WebSocket connection closed:", event);
//     };
// }

// function sendMessage(chatName) {
//     const input = document.getElementById("message");
//     const message = input.value;
//     ws.send(JSON.stringify({ chatUid: chatName, username: currentUsername, text: message }));
//     input.value = '';
// }

// window.onload = function() {
//     const chatName = document.body.getAttribute("data-chat-name");
//     const token = localStorage.getItem('token');
//     authenticateAndConnect(chatName, token);
// };



// let ws;
// let currentUsername;

// function authenticateAndConnect(chatName, token) {
//     fetch(`/auth_ws/${chatName}`, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/x-www-form-urlencoded',
//         },
//         body: `token=${encodeURIComponent(token)}`
//     })
//     .then(response => response.json())
//     .then(data => {
//         if (data.message === "Authenticated") {
//             currentUsername = data.username; // Сохраняем имя пользователя
//             connectWebSocket(chatName);
//         } else {
//             console.error("Authentication failed:", data);
//         }
//     })
//     .catch(error => {
//         console.error("Error during authentication:", error);
//     });
// }

// function connectWebSocket(chatName) {
//     ws = new WebSocket(`ws://${window.location.host}/ws/${chatName}`);

//     ws.onmessage = function(event) {
//         const message = JSON.parse(event.data);
//         const messagesDiv = document.getElementById("messages");
//         const messageDiv = document.createElement("div");
//         messageDiv.textContent = `${message.username}: ${message.text}`;
//         messagesDiv.appendChild(messageDiv);
//         messagesDiv.scrollTop = messagesDiv.scrollHeight; // Прокрутка вниз
//     };
// }

// function sendMessage(chatName) {
//     const input = document.getElementById("message");
//     const message = input.value;
//     ws.send(JSON.stringify({ chatUid: chatName, username: currentUsername, text: message }));
//     input.value = '';
// }

// window.onload = function() {
//     const chatName = document.body.getAttribute("data-chat-name");
//     const token = document.body.getAttribute("data-token");
//     authenticateAndConnect(chatName, token);
// };





// let ws;

// function authenticateAndConnect(chatName, token) {
//     fetch(`/auth_ws/${chatName}`, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/x-www-form-urlencoded',
//         },
//         body: `token=${encodeURIComponent(token)}`
//     })
//     .then(response => response.json())
//     .then(data => {
//         if (data.message === "Authenticated") {
//             connectWebSocket(chatName);
//         } else {
//             console.error("Authentication failed:", data);
//         }
//     })
//     .catch(error => {
//         console.error("Error during authentication:", error);
//     });
// }

// function connectWebSocket(chatName) {
//     ws = new WebSocket(`ws://${window.location.host}/ws/${chatName}`);

//     ws.onopen = function() {
//         console.log("WebSocket connection opened");
//     };

//     ws.onmessage = function(event) {
//         const message = JSON.parse(event.data);
//         const messagesDiv = document.getElementById("messages");
//         const messageDiv = document.createElement("div");
//         messageDiv.textContent = `${message.username}: ${message.text}`;
//         messagesDiv.appendChild(messageDiv);
//         messagesDiv.scrollTop = messagesDiv.scrollHeight; // Прокрутка вниз
//     };

//     ws.onclose = function() {
//         console.log("WebSocket connection closed");
//     };

//     ws.onerror = function(error) {
//         console.error("WebSocket error observed:", error);
//     };
// }

// function sendMessage(chatName) {
//     const input = document.getElementById("message");
//     const message = input.value;
//     ws.send(JSON.stringify({ chatUid: chatName, username: "currentUser", text: message }));
//     input.value = '';
// }

// window.onload = function() {
//     const chatName = document.body.getAttribute("data-chat-name");
//     const token = localStorage.getItem('token'); // Получите токен из вашего механизма аутентификации
//     authenticateAndConnect(chatName, token);
// };









