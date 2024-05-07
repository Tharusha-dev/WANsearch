package utils

import (
	"database/sql"

	"fmt"
	"log"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Video struct {
	Video_id       string
	TimeDialogues2 []TimeDialogue
	Title          string
}

type TimeDialogue struct {
	Time     string
	Dialogue string
}

func Videos(query string) []Video {

	db := ConnectionToDB()

	defer db.Close()

	if isQuotedSearch(query) {

		return quotedVideos(query, db)

	}

	query = sanitizeQuery(query)
	query_words := queryWordsSlice(query)

	query_words_tf := computeTF(query_words)

	documentsToSearch := fetchDocumentsMatchingTerms(query_words, db)

	documentsMatchingAllTerms := intersection(documentsToSearch)

	videos := matchingVideos(&documentsMatchingAllTerms, query_words_tf, query_words, db)

	return videos
}

func checkError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}

func isQuotedSearch(query string) bool {

	query = strings.ToLower(strings.TrimSuffix(query, " "))
	characters := strings.Split(query, "")

	if characters[0] == `"` && characters[len(characters)-1] == `"` {
		return true

	} else {
		return false
	}

}

func quotedVideos(query string, db *sql.DB) []Video {

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

func sanitizeQuery(query string) string {
	// defer timer("sanitizeQuery")()

	regForPunctuations, _ := regexp.Compile("[^a-zA-Z0-9]+")

	return strings.ToLower(regForPunctuations.ReplaceAllString(query, " "))

}

func cosineSimilarity(queryVector map[string]float64, documentVector map[string][]float64, sqrtQueryMagnitude float64) float64 {
	// defer timer("cosineSimilarity")()

	var dotProduct, magnitudeV2 float64

	for k, v := range queryVector {

		if v2Value, ok := documentVector[k]; ok {
			dotProduct += v * v2Value[0]
		}
	}

	for _, v := range documentVector {
		magnitudeV2 += v[1]
	}
	return dotProduct / (sqrtQueryMagnitude * math.Sqrt(magnitudeV2))
}

func cosineSimilarityDocs(tfs *sql.Rows, docsToSearch []string, queryVector map[string]float64) map[string]float64 {

	var queryMagnitude float64

	for _, v := range queryVector {
		queryMagnitude += v * v
	}

	sqrtQueryMag := math.Sqrt(queryMagnitude)

	defer tfs.Close()
	// defer timer("cosineSimilarityDocs")()

	documentsWithValue := make(map[string]float64, len(docsToSearch))
	documentVector := make(map[string][]float64)

	for j := 0; j < len(docsToSearch); j++ {

		var doc string
		var term_ string
		var tfidf float64
		var docTemp string
		var magnitude float64

		for {
			if tfs.Next() {
				err := tfs.Scan(&doc, &term_, &tfidf, &magnitude)

				checkError(err)

				if docTemp != "" && doc != docTemp {

					break
				}

				documentVector[term_] = []float64{tfidf, magnitude}

				docTemp = doc

			} else {
				break
			}
		}

		score := cosineSimilarity(queryVector, documentVector, sqrtQueryMag)

		documentsWithValue[doc] = score

	}

	return documentsWithValue

}

func computeCosineSimilarityBatch(tfs *sql.Rows, docsToSearch []string, queryVector map[string]float64) map[string]float64 {

	var queryMagnitude float64
	for _, v := range queryVector {
		queryMagnitude += v * v
	}
	sqrtQueryMag := math.Sqrt(queryMagnitude)

	defer tfs.Close()
	// defer timer("cosineSimilarityDocs")()

	documentsWithValue := make(map[string]float64, len(docsToSearch))
	batchSize := 100
	batchData := make(map[string]map[string][]float64)

	for {

		for i := 0; i < batchSize && tfs.Next(); i++ {
			var doc string
			var term_ string
			var tfidf float64
			var magnitude float64

			err := tfs.Scan(&doc, &term_, &tfidf, &magnitude)
			checkError(err)

			// Better way to implement Document vector, while grouping document vectors in to a batch,
			// TODO: messure performance difference
			if innerMap, ok := batchData[doc]; ok {
				innerMap[term_] = []float64{tfidf, magnitude}
			} else {
				batchData[doc] = map[string][]float64{term_: {tfidf, magnitude}}
			}
		}

		// Process batch when full || no more results
		if len(batchData) > 0 {
			for docID, docVec := range batchData {
				score := cosineSimilarity(queryVector, docVec, sqrtQueryMag)
				documentsWithValue[docID] = score
			}
			batchData = make(map[string]map[string][]float64) // Reset
		}

		if !tfs.Next() {
			break
		}
	}

	return documentsWithValue
}

func queryWordsSlice(sentence string) []string {
	// defer timer("queryWordsSlice")()

	return strings.Split(sentence, " ")

}

func computeTF(query_words []string) map[string]float64 {
	// defer timer("computeTF")()

	// query_words := getQueryWords(sentence)

	query_words_tf := make(map[string]float64)

	for _, word := range query_words {
		_, present := query_words_tf[word]

		if present {
			query_words_tf[word] += 1
		} else {
			query_words_tf[word] = 1

		}
	}

	return query_words_tf
}

func fetchInvertedIndex(term *string, db *sql.DB) []string {
	// defer timer("fetchInvertedIndex")()

	var docs string

	stmt, err := db.Prepare("select video_id from inverted_index where term = ?")

	if err != nil {
		log.Fatal(err)
	}

	documents, err := stmt.Query(*term)

	if err != nil {
		log.Fatal(err)
	}

	defer documents.Close()

	for documents.Next() {

		err = documents.Scan(&docs)
		if err != nil {
			log.Fatal(err)
		}
	}

	return strings.Split(docs, ",")
}

func fetchDocumentsMatchingTerms(query_words []string, db *sql.DB) [][]string {
	// defer timer("fetchDocumentsMatchingTerms")()

	var documentSlices [][]string

	for _, word := range query_words {
		documnts := fetchInvertedIndex(&word, db)

		if len(documnts) > 0 && documnts[0] != "" {
			documentSlices = append(documentSlices, documnts)

		}
	}

	return documentSlices

}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func fetchTitleFromId(document *string, db *sql.DB) string {
	// defer timer("getTitlesFromId")()

	var title string

	stmt, err := db.Prepare("SELECT title FROM video_titles WHERE vid_id=?")

	if err != nil {
		log.Fatal(err)
	}

	tfs, err := stmt.Query(*document)

	if err != nil {
		log.Fatal(err)
	}

	for tfs.Next() {

		err = tfs.Scan(&title)
		if err != nil {
			log.Fatal(err)
		}

	}
	return title
}

func sanitizeTitle(title string) string {
	// defer timer("sanitizeTitle")()

	title = sanitizeQuery(title)

	return title

}

func fetchTitle(documents []string, db *sql.DB) map[string]string {
	// defer timer("fetchTitle")()

	var title string
	var vid_id string

	titlesVideoMap := make(map[string]string, len(documents))

	placeholders, args := prepareForInStatement(documents)

	stmt, err := db.Prepare("SELECT vid_id,title FROM video_titles WHERE vid_id in" + placeholders)

	if err != nil {
		log.Fatal(err)
	}

	tfs, err := stmt.Query(args...)
	checkError(err)

	defer tfs.Close()

	if err != nil {
		log.Fatal(err)
	}

	for tfs.Next() {

		err = tfs.Scan(&vid_id, &title)
		if err != nil {
			log.Fatal(err)
		}

		titlesVideoMap[vid_id] = title

	}
	return titlesVideoMap
}

func titleWeight(queryVector []string, title string) float64 {
	// defer timer("titleWeight")()

	c := 0

	wordsInTitle := strings.Split(title, " ")

	for _, queryWord := range queryVector {

		if slices.Contains(wordsInTitle, queryWord) {
			c += 1
		}

	}

	return float64(c)

}

func prepareForInStatement(items []string) (string, []interface{}) {
	// defer timer("PrepareForInStatement")()

	lenItems := len(items)

	args := make([]interface{}, lenItems)
	for i, id := range items {
		args[i] = id
	}

	var builder strings.Builder

	builder.WriteString("(")

	if lenItems > 0 {
		builder.WriteString("?")

		for i := 0; i < lenItems-1; i++ {
			builder.WriteString(",")
			builder.WriteString("?")
		}
	}

	builder.WriteString(")")

	return builder.String(), args

}

func matchingVideos(docsToSearch *[]string, queryVector map[string]float64, queryWords []string, db *sql.DB) []Video {

	video_id_placeholders, video_id_args := prepareForInStatement(*docsToSearch)
	terms_placeholders, terms_args := prepareForInStatement(queryWords)

	stmt, err := db.Prepare("Select video_id, term, tfidf,magnitude From all_episodes_tfidf_count Where video_id IN  " + video_id_placeholders + "and term in " + terms_placeholders + "order by video_id")

	checkError(err)

	allArgs := append(video_id_args, terms_args...)

	tfs, err := stmt.Query(allArgs...)

	checkError(err)
	defer tfs.Close()

	checkError(err)

	vidTitlesMap := fetchTitle(*docsToSearch, db)

	documentsWithValue := computeCosineSimilarityBatch(tfs, *docsToSearch, queryVector)

	documentsWithValueSortedMap, documentsWithValueSorted := sortDocumentsByValue(documentsWithValue)

	documentsWithValueFar := make(map[string]float64)

	if len(documentsWithValueSorted) > 10 {

		documentsWithValueSorted = documentsWithValueSorted[:10]
	}

	videoTermPositions := fetchTermPositions(db, documentsWithValueSorted, queryWords)

	for _, doc := range documentsWithValueSorted {

		score := documentsWithValueSortedMap[doc]

		ds := distanceScore(db, videoTermPositions, doc, queryWords)

		if ds > score {
			score = 0
		} else {
			score -= ds

		}

		tw := titleWeight(queryWords, sanitizeTitle(vidTitlesMap[doc]))

		score += tw

		documentsWithValueFar[doc] = score

	}

	_, documentsWithValueFarSorted := sortDocumentsByValue(documentsWithValueFar)

	var finalVideos []Video

	if len(documentsWithValueFarSorted) > 6 {

		documentsWithValueFarSorted = documentsWithValueFarSorted[:6]
	}

	documentsWithTimesMap := fetchTimeStamps(db, queryWords, docsToSearch) //doc : timestamps of all terms in query within doc group by term

	for _, video := range documentsWithValueFarSorted {

		dialoguesWithWordInDocument := fetchDialogueFromTimeStamps(db, &video, documentsWithTimesMap[video])

		finalVideos = append(finalVideos, Video{Video_id: video, TimeDialogues2: dialoguesWithWordInDocument, Title: vidTitlesMap[video]})

	}

	return finalVideos

}

func intersection(slices [][]string) []string {
	// defer timer("intersection")()

	commonDocsMap := make(map[string]int, 0)
	for _, slice := range slices {
		sliceMap := make(map[string]struct{}, 0)
		for _, doc := range slice {
			sliceMap[doc] = struct{}{}
		}
		for doc := range sliceMap {
			commonDocsMap[doc]++
		}
	}
	commonDocs := make([]string, 0)
	for doc, count := range commonDocsMap {
		if count == len(slices) {
			commonDocs = append(commonDocs, doc)
		}
	}
	return commonDocs
}

func fetchTermPositions(db *sql.DB, documents []string, queryWords []string) map[string]map[string][]int {

	// defer timer("fetchTermPositions")()

	var term string
	var positions string
	var videoId string

	terms_placeholders, terms_args := prepareForInStatement(queryWords)

	video_id_placeholders, video_id_args := prepareForInStatement(documents)

	stmt, err := db.Prepare(`select video_id,terms,positions from term_positions Where video_id in ` + video_id_placeholders + `and terms in ` + terms_placeholders)

	checkError(err)
	allArgs := append(video_id_args, terms_args...)

	tfs, err := stmt.Query(allArgs...)
	checkError(err)
	defer tfs.Close()

	checkError(err)

	result := make(map[string]map[string][]int)

	for tfs.Next() {
		err = tfs.Scan(&videoId, &term, &positions)
		checkError(err)
		if _, ok := result[videoId]; !ok {
			result[videoId] = make(map[string][]int)
		}
		result[videoId][term] = commaSeparatedToIntSlice(positions)
	}

	return result
}

func fetchLengthOfDocument(db *sql.DB, document string) int {

	var length string

	stmt, err := db.Prepare(`SELECT positions FROM term_positions WHERE video_id= ? ORDER BY "index" DESC LIMIT 1;`)
	checkError(err)

	tfs, err := stmt.Query(document)

	for tfs.Next() {
		err = tfs.Scan(&length)
	}
	n, err := strconv.Atoi(length)

	return n

}

func distanceScore(db *sql.DB, videoTermPositions map[string]map[string][]int, document string, queryWords []string) float64 {
	// defer timer("distanceScore")()

	lenOfDoc := fetchLengthOfDocument(db, document)

	var score float64
	numOfWords := len(queryWords) - 1

	for j := 0; j < numOfWords; j++ {

		posList1 := videoTermPositions[document][queryWords[j]]
		posList2 := videoTermPositions[document][queryWords[j+1]]

		c := findClosestFromTwoSlicec(posList1, posList2)
		score += float64(c)

	}

	return score / float64(lenOfDoc)

}

func findClosestFromTwoSlicec(list1 []int, list2 []int) int {
	// probably not the best way to do this

	// defer timer("getClosestFromTwoSlicec")()

	close := 100000 // large number

	if len(list1) == 0 {
		return close
	}

	startNum := list1[0]

	for _, n2 := range list2 {
		for _, n1 := range list1 {

			c := n2 - n1

			if c < 0 || n2 < startNum {
				break
			}

			if c < close {
				close = c
			}

		}
	}

	return close
}

func commaSeparatedToIntSlice(str string) []int {
	// defer timer("commaSeparatedToIntSlice")()
	//
	var result []int

	result = make([]int, 0, len(str)+1)

	var n int
	var err error
	var start int

	for i, char := range str {
		if char == ',' {
			n, err = strconv.Atoi(str[start:i])
			checkError(err)

			result = append(result, n)
			start = i + 1
		}
	}

	// Handle trailing comma or empty string
	if start < len(str) {
		n, err = strconv.Atoi(str[start:])
		checkError(err)
		result = append(result, n)
	}

	return result
}

func sortDocumentsByValue(originalMap map[string]float64) (map[string]float64, []string) {
	// defer timer("sortDocumentsByValue")()

	documents := make([]string, 0, len(originalMap))

	documentsWithSimilarity := make(map[string]float64, len(originalMap))

	for document := range originalMap {
		documents = append(documents, document)
	}
	sort.Slice(documents, func(i, j int) bool { return originalMap[documents[i]] > originalMap[documents[j]] }) //this gives a orders slice of video_ids ("documents")

	for _, document := range documents {
		documentsWithSimilarity[document] = originalMap[document]

	}

	return documentsWithSimilarity, documents
}

func fetchDialogueFromTimeStamps(db *sql.DB, document *string, timeStamps string) []TimeDialogue {
	// defer timer("fetchDialogueFromTimeStamps")()

	var dialogues []TimeDialogue
	var dialogue string
	var time string

	tfs, err := db.Query(fmt.Sprintf(`select time,dialogue from time_dialogue where time in (%s) and video_id="%s"`, timeStamps, *document))

	checkError(err)

	for tfs.Next() {
		err = tfs.Scan(&time, &dialogue)

		checkError(err)

		dialogues = append(dialogues, TimeDialogue{Dialogue: dialogue, Time: time})

	}

	defer tfs.Close()

	return dialogues

}

func fetchTimeStamps(db *sql.DB, words []string, documents *[]string) map[string]string {
	// defer timer("fetchTimeStamps")()

	//
	var timeStamps string
	var doc string

	documentsTimeStampsMap := make(map[string]string, len(*documents))

	term_placeholders, term_args := prepareForInStatement(words)
	video_id_placeholders, video_id_args := prepareForInStatement(*documents)

	stmt, err := db.Prepare("select video_id,GROUP_CONCAT(times, ',') AS times from word_time where term in" + term_placeholders + " and video_id in " + video_id_placeholders + "GROUP BY video_id")

	checkError(err)

	allArgs := append(term_args, video_id_args...)

	tfs, err := stmt.Query(allArgs...)

	checkError(err)

	for tfs.Next() {
		err = tfs.Scan(&doc, &timeStamps)
		if err != nil {
			log.Fatal(err)
		}

		documentsTimeStampsMap[doc] = timeStamps

	}

	return documentsTimeStampsMap

}
