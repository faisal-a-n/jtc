package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/faisal-a-n/jtc/models"
	"github.com/faisal-a-n/jtc/pkg/oauth2"
	"github.com/faisal-a-n/jtc/pkg/reddit"
	redditRepo "github.com/faisal-a-n/jtc/repotsitory/reddit"
	"google.golang.org/api/youtube/v3"
)

func DownloadFromReddit(sub string, post_type string, count int) (models.Post, string) {
	return reddit.GetVideo(redditRepo.Query{
		Sub:      sub,
		Awards:   0,
		Score:    0,
		PostType: post_type,
		Count:    count,
	})
}

func UploadToYoutube(post models.Post, video_dir string) {
	service := oauth2.Init()
	keywords := ""
	if len(post.Data.Title) > 25 {
		post.Data.Title = post.Data.Title[:25]
	}
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       fmt.Sprintf("#shorts r/%s %s", post.Data.Subreddit, post.Data.Title),
			Description: "#shorts Check out more videos",
			CategoryId:  "24",
		},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	call := service.Videos.Insert([]string{"snippet,status"}, upload)
	filename := video_dir
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}

	response, err := call.Media(file).Do()
	handleError(err, "")
	log.Printf("Upload successful! Video ID: %v\n", response.Id)

}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
