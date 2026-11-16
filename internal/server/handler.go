package server

import (
	"crypto/md5"
	"fmt"
	"html"
	"net/http"
	"time"
	"unicode/utf8"

	"tokens-bot/internal/database"
	"tokens-bot/internal/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db *database.DB
}

func New(db *database.DB) *Server {
	return &Server{db: db}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/", s.handleRoot)
}

func (s *Server) handleRoot(c *gin.Context) {
	tokens, err := s.db.GetTokens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch info"})
		return
	}

	svg := s.generateSVG(tokens)
	etag := fmt.Sprintf(`"%x"`, md5.Sum([]byte(svg)))

	c.Header("Content-Type", "image/svg+xml")
	c.Header("Cache-Control", "public, max-age=1800")
	c.Header("ETag", etag)
	c.Header("Expires", time.Now().Add(30*time.Minute).Format(http.TimeFormat))

	c.String(http.StatusOK, svg)
}

func (s *Server) generateSVG(tokens string) string {
	formattedTokens := utils.FormatTokensWithSeparator(tokens)
	fontSize := 75
	if n := utf8.RuneCountInString(formattedTokens); n > 10 {
		fontSize = max(75*10/n, 24)
	}

	safeTokens := html.EscapeString(formattedTokens)

	return fmt.Sprintf(`
<svg xmlns="http://www.w3.org/2000/svg" width="600" height="300" viewBox="0 0 600 300" role="img">
  <defs>
    <linearGradient id="background" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#FFFFFF"/>
      <stop offset="1" stop-color="#F0F0F0"/>
    </linearGradient>
    <linearGradient id="number" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#111111"/>
      <stop offset="1" stop-color="#444444"/>
    </linearGradient>
    <pattern id="grid" width="24" height="24" patternUnits="userSpaceOnUse">
      <path d="M 24 0 H 0 V 24" fill="none" stroke="#000000" stroke-opacity="0.08" />
    </pattern>
    <clipPath id="card">
      <rect x="9" y="9" width="582" height="282" rx="21"/>
    </clipPath>
  </defs>
  <g clip-path="url(#card)">
    <rect x="9" y="9" width="582" height="282" fill="url(#background)" />
    <rect x="9" y="9" width="582" height="282" fill="url(#grid)" />
    <rect x="42" y="40.5" width="33" height="33" rx="10.5" fill="#000000" fill-opacity="0.035" stroke="#000000" stroke-opacity="0.08" />
    <path d="M 60.75 47.25 L 51.75 58.5 H 58.5 L 56.25 66.75 L 66 54.75 H 59.25 Z" fill="#333333" />
    <text x="87" y="61.5" font-family="Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif" font-size="18" font-weight="600" letter-spacing="2.25" fill="#444">TOKENS WASTED</text>
    <text x="300" y="165" text-anchor="middle" dominant-baseline="middle" font-family="Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif" font-size="%[2]d" font-weight="800" letter-spacing="-1.5" fill="url(#number)" style="font-variant-numeric: tabular-nums;">%[1]s</text>
  </g>
  <rect x="9.5" y="9.5" width="581" height="281" rx="20.5" fill="none" stroke="#000000" stroke-opacity="0.1" />
</svg>
`, safeTokens, fontSize)
}
