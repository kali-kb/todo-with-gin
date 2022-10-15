package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
)

//accessible model fields should always be capitalized
//add more fields
type Todo struct {
	ID int
	Task string
	CreatedAt time.Time
}

type User struct {
	ID int
	Username string
	Email string
	Password string
}


func main(){
	router := gin.Default()

    db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
    if err != nil {
    	panic("not connected")
    } else {
    	db.AutoMigrate(&Todo{})
    	db.AutoMigrate(&User{})
    	fmt.Println("database", db)
    }
    if 
	router.LoadHTMLGlob("templates/*")

	auth := func (c *gin.Context){
        cookie, err := c.Request.Cookie("user")
        if err != nil {
        	fmt.Println("cookie does'nt exist")
        	c.Redirect(http.StatusFound, "/user/login")
        } else {
        	c.Set("logged in", true)
        	c.Set("user", cookie)
        	c.Next()
        }

	}
	router.GET("/home", auth,  func(c *gin.Context){
		var todos []Todo
		userdata, loggedIn := c.MustGet("user"), c.MustGet("logged in").(bool)
		fmt.Println(userdata)
		db.Find(&todos)
		for _, data := range todos {
			fmt.Println(data)
			fmt.Println(data.Task)
		}
		fmt.Println("all todos", todos)
		c.HTML(http.StatusOK, "index.html", gin.H{"todolist": todos, "loggedIn": loggedIn})
	})

	todos := router.Group("/todos")
	{
		todos.GET("/:id", func(c *gin.Context){
			task_id := c.Param("id")
			var todo Todo
			fmt.Println(task_id)
			taskID, err := strconv.Atoi(task_id)
			if err != nil {
				log.Fatal("param not converted to int")
			}
			db.Where(&Todo{ID: taskID}).Find(&todo)
			// db.Model(Todo{ID: taskID}).First(&todo)
			// todos := db.Find(&Todo{ID: taskID})
			fmt.Println(todo.Task)
			c.HTML(http.StatusOK, "task.html", gin.H{"task": todo.Task, "created_at": todo.CreatedAt})
		})



		todos.POST("/add", func(c *gin.Context){
			t := c.PostForm("task")
			file, err := c.FormFile("img_data")
			os.Mkdir("media", os.ModePerm)
			if err != nil {
				panic(err)
			} else{
				fmt.Println("file", file.Filename)
				c.SaveUploadedFile(file, "media")
			}
			fmt.Println("filetype", reflect.TypeOf(t))
			todo := Todo{Task:t}
			db.Create(&todo)
			c.Redirect(http.StatusFound, c.Request.Referer())  //add "/home optionally"
		})

		todos.DELETE("/delete/all", func(c *gin.Context){
			db.Exec("DELETE FROM todo")
			c.Redirect(http.StatusFound, c.Request.Referer())
		})
	}




	//grouped

	user := router.Group("/user")
	{

		user.GET("/signup", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "signup.html", gin.H{})
		})

		user.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{})
		})

		user.POST("/login/get", func(ctx *gin.Context){
			email, password := ctx.PostForm("email"), ctx.PostForm("password")
			var user User
			db.Where(&User{Email:email, Password: password}).Find(&user)
			if user.Email == email && user.Password == password {
				ctx.SetCookie("user", password, 30000, "/", "localhost", false, true)
				ctx.Redirect(http.StatusFound, "/home")
			} else{
				fmt.Println("passed info not correct")
				ctx.Redirect(http.StatusFound, ctx.Request.Referer())
			}
			fmt.Println(user.Email, user.Password)
		})

		user.POST("/register/add", func(ctx *gin.Context) {
			email := ctx.PostForm("email")
			name := ctx.PostForm("name")
			password := ctx.PostForm("password")
			newuser := User{Username: name, Email: email, Password: password}
			db.Create(&newuser)
			ctx.SetCookie("user", password, 30000, "/", "localhost", false, true)
			ctx.Redirect(http.StatusFound, ctx.Request.Referer())
		})

	}


	router.Run()
}