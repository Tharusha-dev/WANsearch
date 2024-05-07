# SQLite Database structure
The latest version of the SQLite database file can be downloaded here.

All of these data were gathered using [yt-dlp](https://github.com/yt-dlp/yt-dlp), and computed using python with the help of [scikit-learn](https://scikit-learn.org)

## Tables 📊

 * **inverted_index**: Stores the mapping between words and the documents they appear in.
 * **word_time**: Stores words and the time stamp they appear in.
 * **time_dialogue**: Links timestamps to specific dialogue snippets within videos.
 * **all_dialogues**: Stores all dialogue text extracted from the videos.
 * **video_titles**: Stores the titles of the videos.
 * **term_positions**: Indicates the specific positions of words within documents.
 * **all_episodes_tfidf_count**: Stores tfidf, count and magnitude for each term in a document.

### inverted_index
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |3178
| term| TEXT |tax
| video_id| TEXT |z6alcVTRJMo,nVZTLLLkaXU,sEuVZcru_-A,3NjtRPnPsSc,td6zO4r2ogI,aBhIJQH-x7I,4QIMSsrJ8Vk,ILBz-SBVg08,w5boVxmH_Yc,c6pYG0olgEI, ...


### word_time
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |13240
| video_id| TEXT |3NjtRPnPsSc
| term| TEXT |tax
| times| TEXT |"578s","1211s"

### time_dialogue
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |72674
| video_id| TEXT |b6LnXwytBuA
| time| TEXT |1396s
| dialogue| TEXT |turn your nfts into a tax write-off that

### all_dialogues
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |95
| video_id| TEXT |z6alcVTRJMo
| words| TEXT |what is up everyone and welcome to the wan show we have a great topic lined up for you guys today luke and i both got some ...

### video_titles
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |326
| title| TEXT |4K YouTube Is Getting PAYWALLED - WAN Show October 7 2022
| video_id| TEXT |ltyntSIVsjA


### term_positions
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |537332
| video_id| TEXT |0Fx3DYIY-68
| terms| TEXT |tax
| positions| TEXT |262,270,7821,7825,7829,7895,7906,7914,7938, ...

### all_episodes_tfidf_count
| Column |  Data type|Example
|--|--|--|
| index| INTEGER |13460745
| video_id| TEXT |yhtJbq_C0V4
| term| TEXT |tax
| tfidf| REAL |0.09344002473490726
| count| INTEGER |54
| magnitude| REAL |0.00873103822246008




