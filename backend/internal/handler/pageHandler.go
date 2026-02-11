package handler

import (
	"html/template"

	"github.com/Aidana2007/GO_movie_platform/internal/model"
	"github.com/Aidana2007/GO_movie_platform/internal/service"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PageHandler struct {
	movieService  *service.MovieService
	authService   *service.AuthService
	userService   *service.UserService
	friendService *service.FriendService
}

func NewPageHandler(movieService *service.MovieService, authService *service.AuthService, userService *service.UserService, friendService *service.FriendService) *PageHandler {
	return &PageHandler{
		movieService:  movieService,
		authService:   authService,
		userService:   userService,
		friendService: friendService,
	}
}

func (h *PageHandler) render(c *gin.Context, templateName string, data gin.H) {
	funcMap := template.FuncMap{
		"extractYouTubeEmbed": extractYouTubeEmbed,
	}

	tmpl, err := template.New("base.html").Funcs(funcMap).ParseFiles(
		"../../../frontend/templates/base.html",
		"../../../frontend/templates/"+templateName,
	)
	if err != nil {
		c.String(500, "Template error: "+err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		c.String(500, "Template execution error: "+err.Error())
	}
}

func (h *PageHandler) HomePage(c *gin.Context) {
	user := h.getCurrentUser(c)

	topRated, _ := h.movieService.GetTopRated(8)
	recent, _ := h.movieService.GetAllMovies()

	h.render(c, "home.html", gin.H{
		"User":         user,
		"Movies":       topRated,
		"RecentMovies": recent,
	})
}

func (h *PageHandler) MoviesPage(c *gin.Context) {
	user := h.getCurrentUser(c)

	query := c.Query("search")
	genre := c.Query("genre")
	sort := c.DefaultQuery("sort", "latest")
	minRatingStr := c.DefaultQuery("minRating", "0")

	minRating, _ := strconv.ParseFloat(minRatingStr, 64)

	var movies interface{}
	if query != "" || genre != "" || sort != "latest" || minRating > 0 {
		movies, _ = h.movieService.SearchMovies(query, genre, sort, minRating)
	} else {
		movies, _ = h.movieService.GetAllMovies()
	}

	h.render(c, "movies.html", gin.H{
		"User":      user,
		"Movies":    movies,
		"Search":    query,
		"Genre":     genre,
		"Sort":      sort,
		"MinRating": minRatingStr,
	})
}

func (h *PageHandler) MovieDetailsPage(c *gin.Context) {
	user := h.getCurrentUser(c)

	movieID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.String(400, "Invalid movie ID")
		return
	}

	movie, err := h.movieService.GetMovieByID(movieID)
	if err != nil {
		c.String(404, "Movie not found")
		return
	}

	data := gin.H{
		"User":  user,
		"Movie": movie,
	}

	h.render(c, "movie-details.html", data)
}

func (h *PageHandler) LoginPage(c *gin.Context) {
	h.render(c, "login.html", nil)
}

func (h *PageHandler) RegisterPage(c *gin.Context) {
	h.render(c, "register.html", nil)
}

func (h *PageHandler) WatchlistPage(c *gin.Context) {
	user := h.getCurrentUser(c)
	if user == nil {
		c.Redirect(303, "/login")
		return
	}

	h.render(c, "watchlist.html", gin.H{"User": user})
}

func (h *PageHandler) ProfilePage(c *gin.Context) {
	user := h.getCurrentUser(c)
	if user == nil {
		c.Redirect(303, "/login")
		return
	}

	h.render(c, "profile.html", gin.H{"User": user})
}

func (h *PageHandler) getCurrentUser(c *gin.Context) interface{} {
	token, err := c.Cookie("token")
	if err != nil {
		return nil
	}

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		return nil
	}

	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		return nil
	}

	return user
}

func extractYouTubeEmbed(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	if strings.Contains(u.Host, "youtu.be") {
		videoID := strings.TrimPrefix(u.Path, "/")
		return "https://www.youtube.com/embed/" + videoID
	}

	if strings.Contains(u.Host, "youtube.com") {
		videoID := u.Query().Get("v")
		if videoID != "" {
			return "https://www.youtube.com/embed/" + videoID
		}
	}

	return urlStr
}

func (h *PageHandler) UsersPage(c *gin.Context) {
	user := h.getCurrentUser(c)
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	search := c.Query("search")
	var searchResults interface{}

	if search != "" {
		searchResults, _ = h.userService.SearchUsers(search, 20)
	}

	h.render(c, "users.html", gin.H{
		"User":          user,
		"Search":        search,
		"SearchResults": searchResults,
	})
}

func (h *PageHandler) UserProfilePage(c *gin.Context) {
	user := h.getCurrentUser(c)

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Error": "User not found"})
		return
	}

	targetUser, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Error": "User not found"})
		return
	}

	friends, _ := h.userService.GetFriends(userID)

	var isFriend bool
	var hasPendingRequest bool
	if user != nil {
		isFriend, _ = h.userService.IsFriend(user.(*model.User).ID, userID)
		hasPendingRequest, _ = h.friendService.CheckPendingRequest(user.(*model.User).ID, userID)
	}

	c.HTML(http.StatusOK, "user-profile.html", gin.H{
		"User":              user,
		"TargetUser":        targetUser,
		"Friends":           friends,
		"IsFriend":          isFriend,
		"HasPendingRequest": hasPendingRequest,
	})
}

func init() {
}
