# JTC
Script written in Go which downloads videos from Reddit and uploads it to your YouTube account.

## How does it work?
It uses your reddit account to get a token and then a video from a subreddit of your choice is downloaded with it's audio. 
Then the video and the audio are merged together using **ffmpeg** with exec module.

## Config
Config json files should be stored in a config folder present in the root of the project directory. 

**Note** 
Google oauth2 will ask for web login only once.
