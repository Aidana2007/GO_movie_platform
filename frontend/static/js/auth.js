function togglePassword(inputId) {
    const input = document.getElementById(inputId);
    if (input.type === 'password') {
        input.type = 'text';
    } else {
        input.type = 'password';
    }
}

function showModal(title, message, isSuccess = true) {
    const modal = document.getElementById('modal');
    const modalIcon = document.getElementById('modal-icon');
    const modalTitle = document.getElementById('modal-title');
    const modalMessage = document.getElementById('modal-message');

    modalTitle.textContent = title;
    modalMessage.textContent = message;
    modalIcon.className = 'modal-icon ' + (isSuccess ? 'success' : 'error');
    modalIcon.textContent = isSuccess ? '✓' : '✕';

    modal.classList.add('active');
}

function closeModal() {
    const modal = document.getElementById('modal');
    modal.classList.remove('active');

    const modalIcon = document.getElementById('modal-icon');
    const isSuccess = modalIcon.classList.contains('success');

    if (isSuccess && (window.location.pathname === '/login' || window.location.pathname === '/register')) {
        window.location.href = '/';
    }
}

function handleLogin(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    fetch('/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Login failed');
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showModal('Success!', 'You logged in successfully', true);
            } else {
                showModal('Error', data.error || 'Login failed', false);
            }
        })
        .catch(error => {
            showModal('Error', error.message || 'An error occurred. Please try again.', false);
            console.error('Login error:', error);
        });
}

function handleRegister(event) {
    event.preventDefault();

    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    fetch('/api/auth/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, email, password }),
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Registration failed');
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                return fetch('/api/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password }),
                });
            } else {
                throw new Error(data.error || 'Registration failed');
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Login after registration failed');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showModal('Welcome!', 'You registered successfully! Welcome to MovieLand', true);
            } else {
                throw new Error('Login after registration failed');
            }
        })
        .catch(error => {
            showModal('Error', error.message || 'An error occurred. Please try again.', false);
            console.error('Registration error:', error);
        });
}