function openCreateChatPopup() {
    document.getElementById("createChatPopup").style.display = "block";
}

function closeCreateChatPopup() {
    document.getElementById("createChatPopup").style.display = "none";
}

function createChat() {
    const name = document.getElementById("chatName").value;
    const lcc = document.getElementById("chatLcc").value;
    const cypher = document.getElementById("chatCypher").value;
    const userToken = localStorage.getItem('token'); // Предполагаем, что токен сохраняется в localStorage после авторизации

    console.log("User Token", userToken);
    
    fetch('/main/createChat', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({name: name, lifetime: lcc, cypher: cypher, token: userToken})
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === "Chat created successfully") {
            // Add miniaturized chat to the main page
            const chatDiv = document.createElement("div");
            chatDiv.innerHTML = `<a href="/chat/${name}">${name}</a>`;
            document.getElementById("chats").appendChild(chatDiv);
            closeCreateChatPopup();
        } else {
            alert(data.error);
        }
    });
}

// Убедитесь, что функции определены в глобальной области видимости
window.openCreateChatPopup = openCreateChatPopup;
window.closeCreateChatPopup = closeCreateChatPopup;
window.createChat = createChat;