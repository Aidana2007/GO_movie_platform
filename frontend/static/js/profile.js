document.addEventListener('DOMContentLoaded', function() {
    loadStats();
    loadFriends();
});

function loadStats() {
    fetch('/api/user/watchlist')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const count = data.data ? data.data.length : 0;
                document.getElementById('watchlist-count').textContent = count;
            }
        })
        .catch(error => {
            console.error('Error loading watchlist:', error);
        });

    fetch('/api/user/friends')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const count = data.data ? data.data.length : 0;
                document.getElementById('friends-count').textContent = count;
            }
        })
        .catch(error => {
            console.error('Error loading friends:', error);
        });
}

function loadFriends() {
    fetch('/api/user/friends')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayFriends(data.data);
            }
        })
        .catch(error => {
            console.error('Error loading friends:', error);
        });
}

function displayFriends(friends) {
    const friendsList = document.getElementById('friends-list');

    if (!friends || friends.length === 0) {
        friendsList.innerHTML = '<p class="no-results">No friends yet. Visit the <a href="/users" style="color: var(--primary);">Users page</a> to find friends!</p>';
        return;
    }

    friendsList.innerHTML = friends.map(friend => `
        <div class="friend-item">
            <a href="/user/${friend._id}">${friend.username}</a>
            <button onclick="removeFriend('${friend._id}')" class="btn btn-secondary btn-small">Remove</button>
        </div>
    `).join('');
}

function removeFriend(friendId) {
    if (!confirm('Are you sure you want to remove this friend?')) {
        return;
    }

    fetch(`/api/user/friends/${friendId}`, {
        method: 'DELETE'
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                loadFriends();
                loadStats();
            } else {
                alert(data.error || 'Failed to remove friend');
            }
        })
        .catch(error => {
            console.error('Error removing friend:', error);
            alert('An error occurred. Please try again.');
        });
}