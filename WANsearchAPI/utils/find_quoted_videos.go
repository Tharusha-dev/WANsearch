package utils

import (
	// "database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"strings"
)

func QuotedVideos(query string) []Video {

	db := ConnectionToDB()

	var videoId string

	var videoIds []string

	var videos []Video

	query = strings.ReplaceAll(query, `"`, ``)

	stmt, err := db.Prepare("select video_id from all_dialogues where words like ?")

	if err != nil {
		log.Fatal(err)
	}

	documents, err := stmt.Query(fmt.Sprintf("%%%s%%", query))

	if err != nil {
		log.Fatal(err)
	}

	for documents.Next() {
		err = documents.Scan(&videoId)
		if err != nil {
			log.Fatal(err)
		}

		videoIds = append(videoIds, videoId)

	}

	timeStamps := fetchTimeStamps(db, queryWordsSlice(query), &videoIds)

	for _, video := range videoIds {

		dialogues := fetchDialogueFromTimeStamps(db, &video, timeStamps[video])

		title := fetchTitleFromId(&video, db)

		videos = append(videos, Video{TimeDialogues2: dialogues, Video_id: video, Title: title})

	}

	return videos

}
