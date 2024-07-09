// Выход пользователя
document.getElementById('logout-button').addEventListener('click', function() {
    const token = localStorage.getItem('token');

    fetch('/logout?token=' + encodeURIComponent(token), {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => response.json())
    .then(data => {
        const responseElement = document.getElementById('logout-response');
        if (data.error) {
            responseElement.textContent = data.error;
        } else {
            responseElement.textContent = data.message;
            localStorage.removeItem('token');
            window.location.href = '/login'; // Перенаправление на страницу входа после выхода
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});