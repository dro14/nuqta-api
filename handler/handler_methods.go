package handler

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.SetTrustedProxies([]string{"127.0.0.1"})

	h.engine.GET("", h.root)

	group := h.engine.Group("/download")
	group.GET("/:referrer", h.download)

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

var pattern = regexp.MustCompile(`iPhone; CPU iPhone OS (\d+_\d+(_\d+)?) like Mac OS X`)

func (h *Handler) download(c *gin.Context) {
	userAgent := c.Request.UserAgent()
	referrer := c.Param("referrer")
	referrer = strings.TrimSpace(referrer)

	if strings.HasPrefix(referrer, "0x") {
		_, err := strconv.ParseInt(referrer[2:], 16, 64)
		if err == nil {
			osVersion := ""
			match := pattern.FindStringSubmatch(userAgent)
			if match != nil {
				osVersion = strings.ReplaceAll(match[1], "_", ".")
			}
			h.data.SetReferrer(c.Request.Context(), c.ClientIP(), osVersion, referrer)
		} else {
			referrer = ""
		}
	} else {
		referrer = ""
	}

	var url string
	if strings.Contains(userAgent, "Mac OS X") {
		url = "https://apps.apple.com/us/app/nuqta/id6743650655"
	} else {
		url = "https://play.google.com/store/apps/details?id=uz.chuqurtech.nuqta"
		if referrer != "" {
			url += "&referrer=" + referrer
		}
	}

	c.Redirect(http.StatusFound, url)
}
