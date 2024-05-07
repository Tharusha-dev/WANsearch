# WANsearch 🔍

**A fan-made search engine for the WAN show podcast by LMG.** 

##  Frontend
Written in [Svelte](https://svelte.dev/).
 

## Backend
Written in Golang using [Gin](https://gin-gonic.com/).


## Database

Using a SQLite database and the [driver by mattn.](https://github.com/mattn/go-sqlite3)
[More info](https://github.com/Tharusha-dev/WANsearch/blob/main/WANsearchAPI/db/README.md).

## How it works

The algorithm currently uses these techniques to rank documents. (in this order)

- Inverted Index
- Cosine similarity
-  Word proximity
- Title weight (If the title includes a query term it is rated higher)

When retrieving relevent documents,

And uses a simple `LIKE` sql statement in [all_dialogues](https://github.com/Tharusha-dev/WANsearch/blob/main/WANsearchAPI/db/README.md) table for quoted search. 

## Infrastucture

Frontend is hosted in Cloudflare pages.
The API is running in EC2. Uses Cloudflare proxy.

