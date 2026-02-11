document.addEventListener('DOMContentLoaded', () => {
    // Modal logic
    const modalOverlay = document.getElementById('movie-modal');
    const addMovieBtn = document.getElementById('add-movie-btn');
    const closeBtn = document.querySelector('.close-modal');
    const movieForm = document.getElementById('movie-form');

    if (addMovieBtn) {
        addMovieBtn.addEventListener('click', () => {
            openModal();
        });
    }

    if (closeBtn) {
        closeBtn.addEventListener('click', () => {
            closeModal();
        });
    }

    modalOverlay.addEventListener('click', (e) => {
        if (e.target === modalOverlay) {
            closeModal();
        }
    });

    if (movieForm) {
        movieForm.addEventListener('submit', handleMovieSubmit);
    }
});

function openModal(movie = null) {
    const modalOverlay = document.getElementById('movie-modal');
    const modalTitle = document.getElementById('movie-modal-title');
    const movieIdInput = document.getElementById('movie-id');

    // Reset form
    document.getElementById('movie-form').reset();

    if (movie) {
        modalTitle.textContent = 'Edit Movie';
        movieIdInput.value = movie.ID;
        document.getElementById('movie-title').value = movie.title;
        document.getElementById('movie-description').value = movie.description;
        document.getElementById('movie-year').value = movie.year;
        document.getElementById('movie-director').value = movie.director;
        document.getElementById('movie-cast').value = movie.cast ? movie.cast.join(', ') : '';
        document.getElementById('movie-genre').value = movie.genre ? movie.genre.join(', ') : '';
        document.getElementById('movie-poster').value = movie.posterUrl;
        document.getElementById('movie-trailer').value = movie.trailerUrl;
    } else {
        modalTitle.textContent = 'Add Movie';
        movieIdInput.value = '';
    }

    modalOverlay.classList.add('active');
}

function closeModal() {
    const modalOverlay = document.getElementById('movie-modal');
    modalOverlay.classList.remove('active');
}

async function handleMovieSubmit(e) {
    e.preventDefault();

    const movieId = document.getElementById('movie-id').value;
    const isEdit = !!movieId;
    const url = isEdit ? `/api/auth/movies/${movieId}` : '/api/auth/movies'; // Note: check route prefixes
    // Actually routes are /api/movies (POST) and /api/movies/:id (PUT) but under auth group.
    // My routes in main.go:
    // adminMovies.POST("/movies", movieHandler.CreateMovie) -> /api/movies (Wait, the group is apiAuth -> adminMovies)
    // apiAuth group is /api/ ... wait.
    // main.go: api := r.Group("/api") -> apiAuth := api.Group("/") -> adminMovies := apiAuth.Group("/")
    // So it is /api/movies for POST and /api/movies/:id for PUT. Correct.

    const cast = document.getElementById('movie-cast').value.split(',').map(s => s.trim());
    const genre = document.getElementById('movie-genre').value.split(',').map(s => s.trim());

    const movieData = {
        title: document.getElementById('movie-title').value,
        description: document.getElementById('movie-description').value,
        year: parseInt(document.getElementById('movie-year').value),
        director: document.getElementById('movie-director').value,
        cast: cast,
        genre: genre,
        posterUrl: document.getElementById('movie-poster').value,
        trailerUrl: document.getElementById('movie-trailer').value
    };

    try {
        const response = await fetch(url, {
            method: isEdit ? 'PUT' : 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(movieData)
        });

        if (response.ok) {
            location.reload();
        } else {
            const data = await response.json();
            alert('Error: ' + (data.error || 'Failed to save movie'));
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred');
    }
}

async function openEditMovie(id) {
    try {
        const response = await fetch(`/api/movies/${id}`);
        if (response.ok) {
            const data = await response.json();
            // API returns { "movie": ... } or just movie object? 
            // MovieHandler.GetMovieByID returns user, movie in template but JSON usually just movie?
            // Checking MovieHandler.GetMovieByID in backend...
            // It calls s.movieRepo.FindByID(id) -> *model.Movie
            // Utils.RespondJSON calls json.NewEncoder(w).Encode(data). 
            // If handler uses c.JSON(200, movie) it sends movie object directly.
            // Handler: c.JSON(http.StatusOK, movie) -> Yes.

            openModal(data);
        } else {
            alert('Failed to fetch movie details');
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

async function deleteMovie(id) {
    if (!confirm('Are you sure you want to delete this movie?')) return;

    try {
        const response = await fetch(`/api/movies/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            // Remove row
            const row = document.querySelector(`tr[data-movie-id="${id}"]`);
            if (row) row.remove();
        } else {
            alert('Failed to delete movie');
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

async function deleteUser(id) {
    if (!confirm('Are you sure you want to delete this user? This cannot be undone.')) return;

    try {
        const response = await fetch(`/api/admin/users/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            const row = document.querySelector(`tr[data-user-id="${id}"]`);
            if (row) row.remove();
        } else {
            const data = await response.json();
            alert('Error: ' + (data.error || 'Failed to delete user'));
        }
    } catch (error) {
        console.error('Error:', error);
    }
}
