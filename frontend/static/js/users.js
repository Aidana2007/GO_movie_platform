// Users JavaScript

document.addEventListener('DOMContentLoaded', function () {
    if (window.location.pathname === '/users') {
        loadFriendRequests();
        loadSentRequests();
    }
});

function sendFriendRequest(userId) {
    fetch('/api/user/friend-request', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ targetUserId: userId }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('Success', 'Friend request sent successfully!', 'success');
                loadSentRequests(); // Refresh sent requests
            } else {
                showMessage('Error', data.error || 'Failed to send friend request', 'error');
            }
        })
        .catch(error => {
            console.error('Error sending friend request:', error);
            showMessage('Error', 'An error occurred. Please try again.', 'error');
        });
}

function loadFriendRequests() {
    fetch('/api/user/friend-requests')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayFriendRequests(data.data);
            }
        })
        .catch(error => {
            console.error('Error loading friend requests:', error);
        });
}

function loadSentRequests() {
    fetch('/api/user/sent-requests')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displaySentRequests(data.data);
            }
        })
        .catch(error => {
            console.error('Error loading sent requests:', error);
        });
}

function displayFriendRequests(requests) {
    const container = document.getElementById('friend-requests');

    if (!requests || requests.length === 0) {
        container.innerHTML = '<p class="no-results">No pending friend requests</p>';
        return;
    }

    container.innerHTML = requests.map(request => `
        <div class="user-card friend-request-card">
            <div class="user-info">
                <h4>${request.username}</h4>
                <p>Wants to be your friend</p>
            </div>
            <div style="display: flex; gap: 10px;">
                <button onclick="acceptFriendRequest('${request.requestId}')" class="btn btn-primary btn-small">Accept</button>
                <button onclick="rejectFriendRequest('${request.requestId}')" class="btn btn-secondary btn-small">Reject</button>
            </div>
        </div>
    `).join('');
}

function displaySentRequests(requests) {
    const container = document.getElementById('sent-requests');

    if (!requests || requests.length === 0) {
        container.innerHTML = '<p class="no-results">No pending sent requests</p>';
        return;
    }

    container.innerHTML = requests.map(request => `
        <div class="user-card" style="background: #e7f3ff; border-left: 4px solid var(--primary);">
            <div class="user-info">
                <h4>${request.username}</h4>
                <p>Request pending</p>
            </div>
            <div style="display: flex; gap: 10px;">
                <a href="/user/${request.userId}" class="btn btn-secondary btn-small">View Profile</a>
                <button onclick="cancelFriendRequest('${request.requestId}', '${request.username}')" class="btn btn-secondary btn-small">Cancel</button>
            </div>
        </div>
    `).join('');
}

function acceptFriendRequest(requestId) {
    fetch(`/api/user/friend-request/${requestId}/accept`, {
        method: 'POST'
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('Success', 'Friend request accepted!', 'success');
                loadFriendRequests();
            } else {
                showMessage('Error', data.error || 'Failed to accept friend request', 'error');
            }
        })
        .catch(error => {
            console.error('Error accepting friend request:', error);
            showMessage('Error', 'An error occurred. Please try again.', 'error');
        });
}

function rejectFriendRequest(requestId) {
    showConfirm(
        'Reject Friend Request',
        'Are you sure you want to reject this friend request?',
        () => {
            fetch(`/api/user/friend-request/${requestId}/reject`, {
                method: 'POST'
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showMessage('Success', 'Friend request rejected', 'success');
                        loadFriendRequests();
                    } else {
                        showMessage('Error', data.error || 'Failed to reject friend request', 'error');
                    }
                })
                .catch(error => {
                    console.error('Error rejecting friend request:', error);
                    showMessage('Error', 'An error occurred. Please try again.', 'error');
                });
        }
    );
}

function cancelFriendRequest(requestId, username) {
    showConfirm(
        'Cancel Friend Request',
        `Are you sure you want to cancel your friend request to ${username}?`,
        () => {
            fetch(`/api/user/friend-request/${requestId}`, {
                method: 'DELETE'
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showMessage('Success', 'Friend request cancelled', 'success');
                        loadSentRequests();
                    } else {
                        showMessage('Error', data.error || 'Failed to cancel friend request', 'error');
                    }
                })
                .catch(error => {
                    console.error('Error cancelling friend request:', error);
                    showMessage('Error', 'An error occurred. Please try again.', 'error');
                });
        }
    );
}

function removeFriend(friendId) {
    showConfirm(
        'Remove Friend',
        'Are you sure you want to remove this friend?',
        () => {
            fetch(`/api/user/friends/${friendId}`, {
                method: 'DELETE'
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showMessage('Success', 'Friend removed', 'success');
                        window.location.reload();
                    } else {
                        showMessage('Error', data.error || 'Failed to remove friend', 'error');
                    }
                })
                .catch(error => {
                    console.error('Error removing friend:', error);
                    showMessage('Error', 'An error occurred. Please try again.', 'error');
                });
        }
    );
}