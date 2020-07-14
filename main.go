package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var dataPath string

type File struct {
	ID       uint `gorm:"primary_key"`
	Filename string
	Path     string `json:"-"`
}

// newFile
func newFile(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "index", nil)
		return
	}

	filename := c.PostForm("name")
	content := c.PostForm("content")
	fPath := filepath.Join(dataPath, filename)

	err := ioutil.WriteFile(fPath, []byte(content), os.ModePerm)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("create file: %v failed, err: %v", filename, err))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("create file: %v successfully", filename))
	fmt.Printf("create file: %v, successfully, filepath: %v \n", filename, fPath)

	f := File{Filename: filename, Path: fPath}

	db.Save(&f)
}

// editFile
func editFile(c *gin.Context) {
	id := c.Param("id")
	var file File
	err := db.First(&file, id).Error
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Sorry, File id: %v not found", id))
		return
	}
	b, rerr := ioutil.ReadFile(file.Path)
	if rerr != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Sorry, File: %v not found", file.Filename))
		return
	}
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "draft", gin.H{"filename": file.Filename, "content": string(b)})
		return
	}

	content := c.PostForm("content")

	err = ioutil.WriteFile(file.Path, []byte(content), os.ModePerm)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("update file: %v-%v failed, err: %v", id, file.Filename, err))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("update file: %v successfully", file.Filename))
}

// viewFileList
func viewFile(c *gin.Context) {
	var files []File
	db.Find(&files)
	c.HTML(http.StatusOK, "dashboard", files)
}

func downloadFile(c *gin.Context) {
	id := c.Param("id")
	var file File
	err := db.First(&file, id).Error
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Sorry, File id: %v not found", id))
		return
	}
	b, rerr := ioutil.ReadFile(file.Path)
	if rerr != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Sorry, File: %v not found", file.Filename))
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	_, _ = c.Writer.Write(b)
}

func setupDB(dbPath string) {
	var err error
	db, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("open sqlite:%v failed\nerr: %v\n", dbPath, err)
	}
	db.AutoMigrate(File{})
}

func setupText(dataPath string) {
	err := os.Mkdir(dataPath, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			log.Panicf("create directory:data failed, err: %v\n", err)
		}
	}
}

func setupRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromString("index", `
<!DOCTYPE html>
<html lang="en">
<head>
   <meta charset="UTF-8">
   <title>Create TxT File</title>
</head>
<body>
<div>
   <form id="myForm" action="/" method="post">
       Name:    <input name="name" type="text"> <br>
       Content: <textarea cols="33" rows="5" name="content"></textarea><br>
       <input type="submit" value="Save">
   </form>
</div>
<div>
   <a href="/text">VIEWALL</a>
</div>
</body>
</html>
`)
	r.AddFromString("draft", `
<!DOCTYPE html>
<html lang="en">
<head>
   <meta charset="UTF-8">
   <title>Edit TxT File</title>
</head>
<body>
<div>
   <form id="myForm" action="/" method="post">
       Name:    <input name="name" type="text" readonly value="{{.filename}}"><br>
       Content: <textarea cols="33" rows="5" name="content">{{.content}}</textarea><br>
       <input type="submit" value="Update">
   </form>
</div>
<div>
   <a href="/text">VIEWALL</a>
</div>
</body>
</html>
`)
	r.AddFromString("dashboard", `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>ALL Files</title>
</head>
<body>
{{range .}}
	{{ .Filename }}		<a href="/text/{{ .ID }}">Download</a></div>	<a href="/draft/{{.ID}}">Edit</a><br>
{{end}}
<br>
</body>
</html>
`)
	return r
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.HTMLRender = setupRender()
	r.GET("/", newFile)
	r.POST("/", newFile)
	r.GET("/text/:id", downloadFile)
	r.GET("/text", viewFile)
	r.GET("/draft/:id", editFile)
	r.POST("/draft/:id", editFile)
	return r
}

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Panicf("get binary absolute path failed, %v", err)
	}
	cwd := filepath.Dir(exe)

	dbPath := filepath.Join(cwd, "rc.db")
	_, err = os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("sqlite: %v not exist, try to create\n", dbPath)
			f, err := os.Create(dbPath)
			if err != nil {
				log.Panicf("create %v failed, err:%v", dbPath, err)
			}
			f.Close()
		} else {
			log.Panicf("err: %v\n", err)
		}
	}

	setupDB(dbPath)
	defer db.Close()

	dataPath = filepath.Join(cwd, "data")
	setupText(dataPath)

	r := setupRouter()
	log.Fatalf("serve failed, %v", r.Run(":8080"))
}
