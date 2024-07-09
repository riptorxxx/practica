// Вход пользователя
document.getElementById('login-form').addEventListener('submit', function(event) {
    event.preventDefault();

    const login = document.getElementById('login').value;
    const password = document.getElementById('password').value;

    fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ login, password }),
    })
    .then(response => response.json())
    .then(data => {
        const responseElement = document.getElementById('login-response');
        if (data.error) {
            responseElement.textContent = data.error;
        } else {
            responseElement.textContent = "Login successful";
            // Сохраняем токен в localStorage для последующего использования
            localStorage.setItem('token', data.token);
            window.location.href = data.redirect; // Перенаправление на главную страницу после успешного входа
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});

function connectWebSocket() {
    const token = localStorage.getItem('token');
    if (!token) {
        console.error('No token found, please login first.');
        return;
    }

    // Создаем WebSocket соединение, передавая токен как параметр
    const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    socket.onopen = function(event) {
        console.log('WebSocket is open now.');
    };

    socket.onmessage = function(event) {
        const message = JSON.parse(event.data);
        console.log('Message from server:', message);
    };

    socket.onclose = function(event) {
        console.log('WebSocket is closed now.');
    };

    socket.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

// Вызываем функцию подключения WebSocket после успешного входа
if (window.location.pathname === '/main') {
    connectWebSocket();
}
