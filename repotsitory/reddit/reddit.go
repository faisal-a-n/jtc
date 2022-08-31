package reddit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faisal-a-n/jtc/models"
	helper "github.com/faisal-a-n/jtc/pkg/helpers"
	oauthHandler "github.com/faisal-a-n/jtc/pkg/oauth2"
	"golang.org/x/oauth2"
)

//Get an http client with a bearer token attached to it
//TODO Use env
func GetClient(ctx context.Context) http.Client {
	config := &oauth2.Config{
		ClientID:     "dDNzrSm6hYxu2_7p6ocJPQ",
		ClientSecret: "VZB7AZGVDDvr2Rvv02KXn_THd4sx2g",
		Endpoint: oauth2.Endpoint{
			TokenURL:  "https://www.reddit.com/api/v1/access_token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
	tokenDir := "./config/reddit_secret.json"
	token, err := oauthHandler.TokenFromFile(tokenDir)
	if token.Expiry.Before(time.Now()) {
		token, err = config.PasswordCredentialsToken(ctx, os.Getenv("username"), os.Getenv("password"))
		if err != nil {
			log.Fatal("Couldn't get a new token", err)
		}
		oauthHandler.SaveToken(tokenDir, token)
	}
	return *config.Client(ctx, token)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

/*
	Reddit API calls to get data
*/
type Query struct {
	Sub      string
	Score    int
	Awards   int
	PostType string
	Count    int
}

const ERR_AUDIO_UNAVAILABLE = "AUDIO_UNAVAILABLE"
const ERR_VIDEO_UNAVAILABLE = "VIDEO_UNAVAILABLE"

const Post_RISING = "rising"
const Post_HOT = "hot"
const Post_POPULAR = "popular"

const baseUrl = "https://oauth.reddit.com/"

func url(endpoint string) string {
	return fmt.Sprint(baseUrl, endpoint)
}

func makeReq(method string, endpoint string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprint(baseUrl, endpoint), body)
	req.Header.Add("User-Agent", "jtc:v0.1 (by /u/rxmen_)")
	return req
}

func GetPost(client *http.Client, query Query) (string, string, string, models.Post, error) {
	res, err := client.Do(makeReq("GET",
		fmt.Sprintf("r/%s/%s?limit=%d", query.Sub, query.PostType, query.Count),
		nil))
	handleError(err)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	redditResponse := &models.RedditResponse{}
	json.Unmarshal(body, redditResponse)

	best_score, postIndex := 0, 0
	if len(redditResponse.Data.Children) == 0 {
		log.Fatal("No posts available")
	}
	for i, post := range redditResponse.Data.Children {
		if !post.Data.IsVideo {
			continue
		}
		if best_score < post.Data.Score {
			best_score = post.Data.Score
			postIndex = i
		}
	}
	post := redditResponse.Data.Children[postIndex]
	url := post.Data.SecureMedia.RedditVideo.FallbackURL

	video_media := fmt.Sprint("./videos/", getFileName(url))
	video_media, _ = filepath.Abs(video_media)
	audio_media := fmt.Sprint(video_media, "_audio.mp4")
	post_title, _ := filepath.Abs(fmt.Sprint("./videos/", strings.Map(func(r rune) rune {
		if r == '.' {
			return -1
		}
		if r == ' ' {
			return '_'
		}
		return r
	}, post.Data.Title), ".mp4"))

	err = helper.DownloadFile(fmt.Sprint(video_media), url)
	if err != nil {
		log.Fatal("There was an error downloading file ", url, err)
	}

	log.Printf("Downloading media for https://reddit.com/%s\n", post.Data.Permalink)
	err = helper.DownloadFile(fmt.Sprint(audio_media),
		fmt.Sprint(post.Data.URL, "/DASH_audio.mp4"))
	if err != nil {
		fmt.Printf("There was an error downloading file %s. %s\nSkipping audio\n", fmt.Sprint(post.Data.URL, "/DASH_audio.mp4"), err)
		return post_title, video_media, audio_media, post, errors.New(ERR_AUDIO_UNAVAILABLE)
	}

	return post_title, video_media, audio_media, post, nil
}

func getFileName(url string) string {
	name := strings.Replace(url, "?source=fallback", "", -1)
	name = name[strings.LastIndex(name, "/")+1:]
	return name
}
