// Movie Details JavaScript

document.addEventListener('DOMContentLoaded', function () {
    const pathParts = window.location.pathname.split('/');
    const movieId = pathParts[pathParts.length - 1];

    if (movieId) {
        loadReviews(movieId);
        checkWatchlistStatus(movieId);
    }
});

function checkWatchlistStatus(movieId) {
    fetch('/api/user/watchlist')
        .then(response => response.json())
        .then(data => {
            if (data.success && data.data) {
                const isInWatchlist = data.data.some(movie => movie._id === movieId);
                const btn = document.getElementById('watchlist-btn');
                if (btn) {
                    if (isInWatchlist) {
                        btn.textContent = 'Remove from Watchlist';
                        btn.classList.remove('btn-secondary');
                        btn.classList.add('btn-primary');
                    } else {
                        btn.textContent = 'Add to Watchlist';
                        btn.classList.remove('btn-primary');
                        btn.classList.add('btn-secondary');
                    }
                }
            }
        })
        .catch(error => {
            console.error('Error checking watchlist:', error);
        });
}

function loadReviews(movieId) {
    fetch(`/api/movies/${movieId}/reviews`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayReviews(data.data);
            }
        })
        .catch(error => {
            console.error('Error loading reviews:', error);
        });
}

function displayReviews(reviews) {
    const reviewsList = document.getElementById('reviews-list');
    const reviewCountDisplay = document.getElementById('review-count-display');
    const currentUserIdInput = document.getElementById('current-user-id');
    const currentUserId = currentUserIdInput ? currentUserIdInput.value : null;
    const currentUserRoleInput = document.getElementById('current-user-role');
    const currentUserRole = currentUserRoleInput ? currentUserRoleInput.value : null;
    const addReviewSection = document.getElementById('add-review-section');

    if (!reviews || reviews.length === 0) {
        reviewsList.innerHTML = '<p class="no-results">No reviews yet. Be the first to review!</p>';
        if (reviewCountDisplay) {
            reviewCountDisplay.textContent = '(0 reviews)';
        }
        if (addReviewSection) addReviewSection.style.display = 'block';
        return;
    }

    if (reviewCountDisplay) {
        reviewCountDisplay.textContent = `(${reviews.length} review${reviews.length !== 1 ? 's' : ''})`;
    }

    let userHasReviewed = false;
    if (currentUserId) {
        userHasReviewed = reviews.some(review => review.user === currentUserId);
    }

    if (addReviewSection) {
        addReviewSection.style.display = userHasReviewed ? 'none' : 'block';
    }

    reviewsList.innerHTML = reviews.map(review => {
        const isAuthor = currentUserId && review.user === currentUserId;
        const isAdmin = currentUserRole === 'admin';

        let deleteButton = '';
        if (isAuthor) {
            deleteButton = `<button class="btn-text delete" onclick="deleteReview('${review._id}')">Delete</button>`;
        } else if (isAdmin) {
            deleteButton = `<button class="btn-text delete" onclick="deleteReviewAdmin('${review._id}')">Delete (Admin)</button>`;
        }

        const editButton = isAuthor ?
            `<button class="btn-text edit" onclick="startEditReview('${review._id}', ${review.rating}, '${review.comment.replace(/'/g, "\\'")}')">Edit</button>` : '';

        return `
        <div class="review-card" id="review-${review._id}">
            <div class="review-header">
                <div class="review-meta">
                    <span class="review-user">${review.username}</span>
                    <span class="review-date">${new Date(review.createdAt).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}</span>
                </div>
                <div class="review-rating-badge">
                    ${review.rating}/10
                </div>
            </div>
            <p class="review-comment">${review.comment}</p>
            <div class="review-actions">
                ${editButton}
                ${deleteButton}
            </div>
        </div>
    `}).join('');
}

function submitReview(event, movieId) {
    event.preventDefault();

    const rating = parseInt(document.getElementById('rating').value);
    const comment = document.getElementById('comment').value;

    fetch(`/api/movies/${movieId}/reviews`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ rating, comment }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById('review-form').reset();
                loadReviews(movieId);
                showMessage('Success', 'Review submitted successfully!', 'success');
                setTimeout(() => {
                    window.location.reload();
                }, 1500);
            } else {
                showMessage('Error', data.error || 'Failed to submit review', 'error');
            }
        })
        .catch(error => {
            console.error('Error submitting review:', error);
            showMessage('Error', 'An error occurred. Please try again.', 'error');
        });
}

let currentEditId = null;

function startEditReview(reviewId, rating, comment) {
    currentEditId = reviewId;
    const editSection = document.getElementById('edit-review-section');
    const addSection = document.getElementById('add-review-section');

    if (addSection) addSection.style.display = 'none';
    if (editSection) {
        editSection.style.display = 'block';
        document.getElementById('edit-rating').value = rating;
        document.getElementById('edit-comment').value = comment;

        editSection.scrollIntoView({ behavior: 'smooth' });

        const editForm = document.getElementById('edit-review-form');
        editForm.onsubmit = (e) => updateReview(e);
    }
}

function cancelEdit() {
    currentEditId = null;
    const editSection = document.getElementById('edit-review-section');
    const addSection = document.getElementById('add-review-section');

    if (editSection) editSection.style.display = 'none';
}

function updateReview(event) {
    event.preventDefault();

    if (!currentEditId) return;

    const rating = parseInt(document.getElementById('edit-rating').value);
    const comment = document.getElementById('edit-comment').value;

    fetch(`/api/reviews/${currentEditId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ rating, comment }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                cancelEdit();
                const pathParts = window.location.pathname.split('/');
                const movieId = pathParts[pathParts.length - 1];
                loadReviews(movieId);
                showMessage('Success', 'Review updated successfully!', 'success');
                setTimeout(() => {
                    window.location.reload();
                }, 1500);
            } else {
                showMessage('Error', data.error || 'Failed to update review', 'error');
            }
        })
        .catch(error => {
            console.error('Error updating review:', error);
            showMessage('Error', 'An error occurred. Please try again.', 'error');
        });
}

function deleteReview(reviewId) {
    showConfirm('Delete Review', 'Are you sure you want to delete this review?', () => {
        fetch(`/api/reviews/${reviewId}`, {
            method: 'DELETE'
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    const pathParts = window.location.pathname.split('/');
                    const movieId = pathParts[pathParts.length - 1];
                    loadReviews(movieId);
                    showMessage('Success', 'Review deleted successfully', 'success');
                    setTimeout(() => {
                        window.location.reload();
                    }, 1500);
                } else {
                    showMessage('Error', data.error || 'Failed to delete review', 'error');
                }
            })
            .catch(error => {
                console.error('Error deleting review:', error);
                showMessage('Error', 'An error occurred. Please try again.', 'error');
            });
    });
}

function deleteReviewAdmin(reviewId) {
    showConfirm('Delete Review (Admin)', 'Are you sure you want to delete this review?', () => {
        fetch(`/api/admin/reviews/${reviewId}`, {
            method: 'DELETE'
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    const pathParts = window.location.pathname.split('/');
                    const movieId = pathParts[pathParts.length - 1];
                    loadReviews(movieId);
                    showMessage('Success', 'Review deleted successfully', 'success');
                    setTimeout(() => {
                        window.location.reload();
                    }, 1500);
                } else {
                    showMessage('Error', data.error || 'Failed to delete review', 'error');
                }
            })
            .catch(error => {
                console.error('Error deleting review:', error);
                showMessage('Error', 'An error occurred. Please try again.', 'error');
            });
    });
}

function toggleWatchlist(movieId) {
    const btn = document.getElementById('watchlist-btn');
    const isInWatchlist = btn.textContent.includes('Remove');

    const url = `/api/user/watchlist/${movieId}`;
    const method = isInWatchlist ? 'DELETE' : 'POST';

    fetch(url, { method })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to update watchlist');
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                if (isInWatchlist) {
                    btn.textContent = 'Add to Watchlist';
                    btn.classList.remove('btn-primary');
                    btn.classList.add('btn-secondary');
                    showMessage('Success', 'Removed from watchlist', 'success');
                } else {
                    btn.textContent = 'Remove from Watchlist';
                    btn.classList.remove('btn-secondary');
                    btn.classList.add('btn-primary');
                    showMessage('Success', 'Added to watchlist!', 'success');
                }
            } else {
                showMessage('Error', data.error || 'Failed to update watchlist', 'error');
            }
        })
        .catch(error => {
            console.error('Error toggling watchlist:', error);
            showMessage('Error', error.message || 'An error occurred. Please try again.', 'error');
        });
}