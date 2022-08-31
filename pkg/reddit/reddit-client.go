package reddit

import (
	"context"
	"log"

	"github.com/faisal-a-n/jtc/models"
	helper "github.com/faisal-a-n/jtc/pkg/helpers"
	redditRepo "github.com/faisal-a-n/jtc/repotsitory/reddit"
)

func GetVideo(query redditRepo.Query) (models.Post, string) {
	redditClient := redditRepo.GetClient(context.Background())
	title, video_media, audio_media, post, err := redditRepo.GetPost(&redditClient, query)
	if err != nil {
		return post, video_media
	}
	err = helper.AddAudioToVideo(video_media, audio_media, title)
	if err != nil {
		log.Fatal("Couldn't add audio to video", err)
	}
	return post, title
}
