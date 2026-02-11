
function logout() {
    if (typeof showConfirm === 'function') {
        showConfirm(
            'Logout',
            'Are you sure you want to logout?',
            () => {
                fetch('/api/auth/logout', {
                    method: 'POST',
                })
                    .then(() => {
                        window.location.href = '/';
                    })
                    .catch(error => {
                        console.error('Logout error:', error);
                    });
            }
        );
    } else {
        if (confirm('Are you sure you want to logout?')) {
            fetch('/api/auth/logout', {
                method: 'POST',
            })
                .then(() => {
                    window.location.href = '/';
                })
                .catch(error => {
                    console.error('Logout error:', error);
                });
        }
    }
}

function showError(elementId, message) {
    const errorElement = document.getElementById(elementId);
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.style.display = 'block';
    }
}

function showSuccess(elementId, message) {
    const successElement = document.getElementById(elementId);
    if (successElement) {
        successElement.textContent = message;
        successElement.style.display = 'block';
    }
}

function hideMessage(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        element.style.display = 'none';
    }
}