/*
Copyright © 2024 zvirgilx
*/
package cmd

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zvirgilx/searxng-go/kernel/internal/complete"
	"github.com/zvirgilx/searxng-go/kernel/internal/locale"
	"github.com/zvirgilx/searxng-go/kernel/internal/search"
	"github.com/zvirgilx/searxng-go/kernel/internal/util"
	"github.com/zvirgilx/searxng-go/kernel/templates"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run a api server of searxng-go",
	Run: func(cmd *cobra.Command, args []string) {
		runapi()
	},
}

func init() {
	apiCmd.Flags().StringP("addr", "a", ":8888", "address to listen on")
	apiCmd.Flags().StringP("mode", "m", "debug", "gin mode(debug, release, test)")
	viper.BindPFlag("addr", apiCmd.Flags().Lookup("addr"))
	viper.BindPFlag("mode", apiCmd.Flags().Lookup("mode"))

	rootCmd.AddCommand(apiCmd)
}

func runapi() {
	gin.SetMode(viper.GetString("mode"))

	router := gin.Default()

	// allows all origins when debugging
	if viper.GetString("mode") == "debug" {
		router.Use(cors.Default())
	}

	tmpl := template.Must(template.New("").ParseFS(templates.Files, "*.tmpl"))
	router.SetHTMLTemplate(tmpl)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "searxng-go",
		})
	})

	api := router.Group("/api")
	api.GET("/search", func(c *gin.Context) {
		opts, err := search.VerifySearchOptions(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		r := search.Search(c, opts)
		c.JSON(http.StatusOK, gin.H{
			"query":        opts.Query,
			"results":      r.GetSortedData(),
			"suggestions":  util.SetToArray[string](r.Suggestions),
			"info_box":     r.InfoBox,
			"next_page_no": opts.PageNo + 1,
		})
	})
	api.GET("/complete", func(c *gin.Context) {
		q, ok := c.GetQuery("q")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		r := complete.Complete(c, q, locale.ParseAcceptLanguage(c.GetHeader("Accept-Language"), "en_US"))
		c.JSON(http.StatusOK, gin.H{
			"query":   q,
			"results": r,
		})
	})

	router.Run(viper.GetString("addr"))
}
