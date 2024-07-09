// Регистрация пользователя
document.getElementById('register-form').addEventListener('submit', function(event) {
    event.preventDefault();

    const login = document.getElementById('login').value;
    const email = document.getElementById('email').value;
    const phone = document.getElementById('phone').value;
    const password = document.getElementById('password').value;

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ login, email, phone, password }),
    })
    .then(response => response.json())
    .then(data => {
        const responseElement = document.getElementById('register-response');
        if (data.error) {
            responseElement.textContent = data.error;
        } else {
            responseElement.textContent = data.message;
            window.location.href = data.redirect; // Перенаправление на главную страницу после регистрации
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});