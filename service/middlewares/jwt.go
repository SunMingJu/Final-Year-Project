package middlewares

import (
	"net/http"
	"simple-video-net/models/users"
	"simple-video-net/utils/jwt"
	ControllersCommon "simple-video-net/utils/response"
	Response "simple-video-net/utils/response"
	"simple-video-net/utils/validator"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// VerificationToken Carrying a token in the request header
func VerificationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		claim, err := jwt.ParseToken(token)
		if err != nil {
			ControllersCommon.NotLogin(c, "Token expiry")
			c.Abort()
			return
		}
		u := new(users.User)
		if !u.IsExistByField("id", claim.UserID) {
			//Without changing the user
			ControllersCommon.NotLogin(c, "user anomaly")
			c.Abort()
			return
		}
		c.Set("uid", u.ID)
		c.Set("currentUserName", u.Username)
		c.Next()
	}
}

// VerificationTokenAsParameter Carrying a token in the body parameter
func VerificationTokenAsParameter() gin.HandlerFunc {
	type qu struct {
		Token string `json:"token"`
	}
	return func(c *gin.Context) {
		req := new(qu)
		if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
			validator.CheckParams(c, err)
			return
		}
		token := req.Token
		claim, err := jwt.ParseToken(token)
		if err != nil {
			ControllersCommon.NotLogin(c, "Token expiry")
			c.Abort()
			return
		}
		u := new(users.User)
		if !u.IsExistByField("id", claim.UserID) {
			//Without changing the user
			ControllersCommon.NotLogin(c, "user anomaly")
			c.Abort()
			return
		}
		c.Set("uid", u.ID)
		c.Set("currentUserName", u.Username)
		c.Next()
	}
}

// VerificationTokenNotNecessary Carrying a token in the request header (not required)
func VerificationTokenNotNecessary() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if len(token) == 0 {
			//No authentication when users are not logged in
			c.Next()
		} else {
			//User access
			claim, err := jwt.ParseToken(token)
			if err != nil {
				c.Next()
			}
			u := new(users.User)
			if !u.IsExistByField("id", claim.UserID) {
				//Without changing the user
				ControllersCommon.NotLogin(c, "user anomaly")
				c.Abort()
				return
			}
			c.Set("uid", u.ID)
			c.Set("currentUserName", u.Username)
			c.Next()
		}
	}
}

func VerificationTokenAsSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Upgrade ws to return messages
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			http.NotFound(c.Writer, c.Request)
			c.Abort()
			return
		}
		token := c.Query("token")
		claim, err := jwt.ParseToken(token)
		if err != nil {
			Response.NotLoginWs(conn, "Token validation failure")
			_ = conn.Close()
			c.Abort()
			return
		}
		u := new(users.User)
		if !u.IsExistByField("id", claim.UserID) {
			//Without changing the user
			Response.NotLoginWs(conn, "user anomaly")
			_ = conn.Close()
			c.Abort()
			return
		}
		c.Set("uid", u.ID)
		c.Set("conn", conn)
		c.Set("currentUserName", u.Username)
		c.Next()
	}
}
