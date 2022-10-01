package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Param (path) example
	r.GET("/path/:param", func(c *gin.Context) {
		param := c.Param("param")
		c.String(200, "param: %s", param)
	})

	type Path struct {
		ParentID string `uri:"parentId"`
		ChildID  string `uri:"childId"`
	}
	// BindUri (path) example
	r.GET("/parent/:parentId/children/:childId", func(c *gin.Context) {
		path := Path{}
		if err := c.BindUri(&path); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		c.JSON(200, path)
	})

	type Body struct {
		Name string `json:"name" binding:"required"`
		Desc string `json:"desc"`
	}
	// BindJSON (body) example
	r.POST("/body", func(c *gin.Context) {
		body := Body{}
		if err := c.BindJSON(&body); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		c.JSON(200, body)
	})

	type Query struct {
		Test     string    `form:"test" binding:"required,min=3,max=10"`
		Earliest time.Time `form:"earliest"`
		Latest   time.Time `form:"latest" binding:"omitempty,gtfield=Earliest,lteNow"`
	}
	var lteNow validator.Func = func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			if date.After(time.Now()) {
				return false
			}
		}
		return true
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lteNow", lteNow); err != nil {
			return nil
		}
	}

	// BindQuery (form) example
	r.GET("/query", func(c *gin.Context) {
		q := Query{}
		if err := c.BindQuery(&q); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		if q.Latest.IsZero() {
			q.Latest = time.Now().UTC().Truncate(time.Second)
		}
		c.JSON(200, q)
	})

	type ManyParams struct {
		Path  string `uri:"id"`
		Query string `form:"q"`
		Body  string `json:"data" binding:"required"`
	}
	// Bind (all) example
	r.PUT("/all/:id", func(c *gin.Context) {
		params := ManyParams{}
		fmt.Printf("params: %v\n", params)
		if err := c.ShouldBindUri(&params); err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("params: %v\n", params)
		if err := c.ShouldBindQuery(&params); err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("params: %v\n", params)
		if err := c.BindJSON(&params); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("params: %v\n", params)
		c.JSON(200, params)
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
