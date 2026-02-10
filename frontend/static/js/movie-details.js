// Movie Details JavaScript

// Load reviews when page loads
document.addEventListener('DOMContentLoaded', function() {
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
    
    if (!reviews || reviews.length === 0) {
        reviewsList.innerHTML = '<p class="no-results">No reviews yet. Be the first to review!</p>';
        return;
    }

    reviewsList.innerHTML = reviews.map(review => `
        <div class="review-card">
            <div class="review-header">
                <span class="review-user">${review.username}</span>
                <span class="review-date">${new Date(review.createdAt).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}</span>
            </div>
            <div class="review-rating">
                ${'★'.repeat(review.rating)}${'☆'.repeat(10 - review.rating)}
                <span>${review.rating}/10</span>
            </div>
            <p class="review-comment">${review.comment}</p>
        </div>
    `).join('');
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
            // Reset form
            document.getElementById('review-form').reset();
            // Reload reviews
            loadReviews(movieId);
            // Show success message
            alert('Review submitted successfully!');
            // Reload page to update rating
            setTimeout(() => {
                window.location.reload();
            }, 1000);
        } else {
            alert(data.error || 'Failed to submit review');
        }
    })
    .catch(error => {
        console.error('Error submitting review:', error);
        alert('An error occurred. Please try again.');
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
                    alert('Removed from watchlist');
                } else {
                    btn.textContent = 'Remove from Watchlist';
                    btn.classList.remove('btn-secondary');
                    btn.classList.add('btn-primary');
                    alert('Added to watchlist');
                }
            } else {
                alert(data.error || 'Failed to update watchlist');
            }
        })
        .catch(error => {
            console.error('Error toggling watchlist:', error);
            alert(error.message || 'An error occurred. Please try again.');
        });
}