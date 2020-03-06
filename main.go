package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// MDConfig -
var MDConfig *Config

// ReadConfig -
func ReadConfig(configPath string) error {
	MDConfig = new(Config)
	err := MDConfig.ParseConfig(configPath)
	return err
}

// MdRender -
func MdRender(mdfile []byte, buf *bytes.Buffer) error {
	md := goldmark.New(goldmark.WithExtensions(extension.Table))

	if err := md.Convert(mdfile, buf); err != nil {
		return err
	}
	return nil
}

// MdReadFile -
func MdReadFile(path string) ([]byte, error) {
	fullpath := fmt.Sprintf("%s%s", MDConfig.MarkdownPath, path)
	print(fullpath)
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		return nil, err
	}
	return ioutil.ReadFile(fullpath)
}

// WikiPageRender -
func WikiPageRender(c *gin.Context) {
	var buf bytes.Buffer
	var mdfile []byte
	var err error

	mdpath := c.Param("mdpath")

	if mdfile, err = MdReadFile(mdpath); err != nil {
		c.HTML(http.StatusNotFound, "PageNotFound.tmpl", "")
	}

	err = MdRender(mdfile, &buf)
	if err != nil {
		print(err)
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"mdbody": buf.Bytes()})
}

// SafeHTMLForTemplate -
func SafeHTMLForTemplate(b []byte) template.HTML {
	return template.HTML(b)
}

func main() {
	err := ReadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(MDConfig)
	r := gin.Default()
	//r.LoadHTMLGlob("./template/*")
	//html := template.Must(template.ParseGlob(fmt.Sprintf("%s/*", MDConfig.TemplatePath)))
	//html.Funcs(template.FuncMap{"safeHTML": SafeHTMLForTemplate})
	html := template.Must(template.New("").Funcs(template.FuncMap{"safeHTML": SafeHTMLForTemplate}).ParseGlob(fmt.Sprintf("%s/*", MDConfig.TemplatePath)))
	r.SetHTMLTemplate(html)
	/*
		if MDConfig.BindFavicon {
			r.StaticFile("/favicon.ico", MDConfig.FaviconPath)

		}
	*/
	r.GET("/*mdpath", WikiPageRender)
	r.Run() // listen and serve on 0.0.0.0:8080
}
