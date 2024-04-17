package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Cors Cross-domain middleware
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Request method
		method := c.Request.Method
		//request header
		origin := c.Request.Header.Get("Origin")
		// Declare the request header keys
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//This is to allow access to all domains
			c.Header("Access-Control-Allow-Origin", "*")
			//All cross-domain request methods supported by the server, in order to avoid multiple 'preflight' requests for browsing sub-requests.
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//Type of header
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//Allow cross-domain settings to return other subsections Cross-domain key settings for browsers to parse
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			//Cache request information in seconds
			c.Header("Access-Control-Max-Age", "172800")
			//Whether cross-domain requests need to carry cookie information Set to true by default
			c.Header("Access-Control-Allow-Credentials", "false")
			//Set the return format to json
			c.Set("content-type", "application/json")
		}

		//Release all OPTIONS methods
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// Processing requests
		c.Next() //  Processing requests
	}
}
