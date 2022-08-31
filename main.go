package main

import (
	app "github.com/faisal-a-n/jtc/pkg/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	post, dir := app.DownloadFromReddit("abruptchaos", "hot", 60)
	app.UploadToYoutube(post, dir)
}
