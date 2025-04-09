package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.root)

	group := h.engine.Group("/client")
	group.GET("", h.getClientInfo)

	group = h.engine.Group("/schema")
	group.GET("", h.getSchema)
	group.PUT("", h.updateSchema)
	group.DELETE("", h.deleteSchema)
	group.DELETE("/:predicate", h.deletePredicate)

	authorized := h.engine.Group("")
	authorized.Use(h.authMiddleware)

	group = authorized.Group("/profile")
	group.POST("", h.createProfile)
	group.GET("", h.getProfile)
	group.PUT("", h.updateProfile)
	group.OPTIONS("", h.isAvailable)

	group = authorized.Group("/user")
	group.GET("", h.getUserList)
	group.GET("/username", h.getUserByUsername)

	group = authorized.Group("/post")
	group.POST("", h.createPost)
	group.GET("", h.getPostList)
	group.DELETE("", h.deletePost)

	group = authorized.Group("/edge")
	group.POST("", h.createEdge)
	group.DELETE("", h.deleteEdge)

	group = authorized.Group("/yordamchi")
	group.POST("", h.createResponse)

	return h.engine.Run(":" + port)
}

func (h *Handler) root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (h *Handler) authMiddleware(c *gin.Context) {
	idToken := c.GetHeader("my-id-token")
	if idToken == "" {
		c.JSON(http.StatusUnauthorized, failure(e.ErrNoAuthHeader))
		c.Abort()
		return
	}

	firebaseUid, err := h.auth.VerifyIdToken(c.Request.Context(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, failure(err))
		c.Abort()
		return
	}

	c.Set("firebase_uid", firebaseUid)
}

func (h *Handler) getClientInfo(c *gin.Context) {
	ua := c.Request.UserAgent()
	os, version := ExtractOSAndVersion(ua)
	if version != "" {
		os += " " + version
	}
	c.JSON(http.StatusOK, gin.H{
		"ip": c.ClientIP(),
		"os": os,
		"ua": ua,
	})
}

// ExtractOSAndVersion extracts the operating system and its version from a user-agent string.
// It returns two strings: the OS name and the version. If either cannot be determined,
// an empty string is returned for that component.
func ExtractOSAndVersion(ua string) (string, string) {
	// Extract the first parenthesized section
	osPart := extractOSPart(ua)
	if osPart == "" {
		return "", ""
	}

	// Parse the OS and version from the extracted section
	return parseOS(osPart)
}

// extractOSPart extracts the content within the first set of parentheses in the user-agent string.
func extractOSPart(ua string) string {
	start := strings.Index(ua, "(")
	if start == -1 {
		return ""
	}
	end := strings.Index(ua[start:], ")")
	if end == -1 {
		return ""
	}
	return ua[start+1 : start+end]
}

// parseOS determines the OS name and version from the extracted OS part.
func parseOS(osPart string) (string, string) {
	// Windows
	if strings.Contains(osPart, "Windows") {
		re := regexp.MustCompile(`Windows NT (\d+\.\d+)`)
		if match := re.FindStringSubmatch(osPart); match != nil {
			return "Windows", match[1] // e.g., "10.0"
		}
		return "Windows", ""
	}

	// macOS
	if strings.Contains(osPart, "Mac OS X") {
		re := regexp.MustCompile(`Mac OS X (\d+_\d+(_\d+)?)`)
		if match := re.FindStringSubmatch(osPart); match != nil {
			version := strings.Replace(match[1], "_", ".", -1) // e.g., "10_15_7" -> "10.15.7"
			return "Mac OS X", version
		}
		return "Mac OS X", ""
	}

	// Android
	if strings.Contains(osPart, "Android") {
		re := regexp.MustCompile(`Android (\d+(\.\d+)?)`)
		if match := re.FindStringSubmatch(osPart); match != nil {
			return "Android", match[1] // e.g., "10" or "4.4"
		}
		return "Android", ""
	}

	// iOS (iPhone or iPad)
	if strings.Contains(osPart, "iPhone") || strings.Contains(osPart, "iPad") {
		re := regexp.MustCompile(`OS (\d+_\d+(_\d+)?)`)
		if match := re.FindStringSubmatch(osPart); match != nil {
			version := strings.Replace(match[1], "_", ".", -1) // e.g., "14_4" -> "14.4"
			return "iOS", version
		}
		return "iOS", ""
	}

	// Linux
	if strings.Contains(osPart, "Linux") {
		return "Linux", "" // Version often not specified
	}

	// Unknown OS
	return "", ""
}
