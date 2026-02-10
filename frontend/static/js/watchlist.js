// Watchlist JavaScript

document.addEventListener('DOMContentLoaded', function() {
    loadWatchlist();
});

function loadWatchlist() {
    fetch('/api/user/watchlist')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayWatchlist(data.data);
            }
        })
        .catch(error => {
            console.error('Error loading watchlist:', error);
        });
}

function displayWatchlist(movies) {
    const grid = document.getElementById('watchlist-grid');
    const emptyMessage = document.getElementById('empty-message');

    if (!movies || movies.length === 0) {
        grid.style.display = 'none';
        emptyMessage.style.display = 'block';
        return;
    }

    grid.style.display = 'grid';
    emptyMessage.style.display = 'none';

    grid.innerHTML = movies.map(movie => `
        <div class="movie-card">
            <a href="/movie/${movie._id}">
                <img src="${movie.posterUrl}" alt="${movie.title}" class="movie-poster">
                <div class="movie-info">
                    <h3>${movie.title}</h3>
                    <p class="movie-year">${movie.year}</p>
                    <div class="movie-rating">
                        <span class="star">★</span>
                        <span>${movie.ranking.toFixed(1)}</span>
                    </div>
                    <div class="movie-genres">
                        ${movie.genre.slice(0, 2).map(genre => 
                            `<span class="genre-badge">${genre}</span>`
                        ).join('')}
                    </div>
                </div>
            </a>
            <div style="padding: 10px;">
                <button onclick="removeFromWatchlist('${movie._id}')" class="btn btn-secondary btn-block">
                    Remove
                </button>
            </div>
        </div>
    `).join('');
}

function removeFromWatchlist(movieId) {
    fetch(`/api/user/watchlist/${movieId}`, {
        method: 'DELETE',
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            loadWatchlist();
        } else {
            alert(data.error || 'Failed to remove from watchlist');
        }
    })
    .catch(error => {
        console.error('Error removing from watchlist:', error);
        alert('An error occurred. Please try again.');
    });
}