package handler

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dro14/nuqta-service/utils/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.SetTrustedProxies([]string{"127.0.0.1"})

	h.engine.GET("", h.root)

	h.engine.DELETE("/type/:name", h.deleteType)

	group := h.engine.Group("/download")
	group.GET("", h.download)
	group.GET("/:referrer", h.download)

	authorized := h.engine.Group("")
	authorized.Use(h.authMiddleware)

	authorized.GET("/ping", h.ping)
	authorized.GET("/update/:after", h.getUpdate)

	group = authorized.Group("/profile")
	group.POST("", h.createProfile)
	group.GET("", h.getProfile)
	group.PUT("", h.editProfile)
	group.DELETE("", h.deleteProfile)
	group.OPTIONS("", h.isAvailable)

	group = authorized.Group("/user")
	group.GET("", h.getUserList)
	group.GET("/username", h.getUserByUsername)

	group = authorized.Group("/post")
	group.POST("", h.createPost)
	group.GET("", h.getPostList)
	group.PUT("", h.editPost)
	group.PATCH("", h.hidePost)
	group.PATCH("/report", h.reportPost)
	group.DELETE("", h.deletePost)

	group = authorized.Group("/edge")
	group.POST("", h.createEdge)
	group.DELETE("", h.deleteEdge)

	group = authorized.Group("/chat")
	group.POST("", h.createChat)
	group.GET("/:type", h.getMessages)
	group.POST("/private/type", h.typePrivate)
	group.POST("/private/new", h.createPrivate)
	group.PATCH("/private/view", h.viewPrivate)
	group.PATCH("/private/like", h.likePrivate)
	group.PATCH("/private/unlike", h.unlikePrivate)
	group.PUT("/private", h.editPrivate)
	group.DELETE("/private/remove", h.removePrivate)
	group.DELETE("/private/delete", h.deletePrivate)
	group.POST("/yordamchi/:provider", h.createYordamchi)
	group.PUT("/yordamchi/:provider", h.editYordamchi)

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

	ctx := c.Request.Context()
	firebaseUid, err := h.firebase.VerifyIdToken(ctx, idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, failure(err))
		c.Abort()
		return
	}

	uid, _ := h.data.GetUidByFirebaseUid(ctx, firebaseUid)
	c.Set("firebase_uid", firebaseUid)
	c.Set("uid", uid)
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

func (h *Handler) deleteType(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	err := h.data.DeleteType(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
