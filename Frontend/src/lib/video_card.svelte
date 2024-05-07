<script>
    export let video;
    export let query;

    function decorateDialogue(dialogue, query) {
        let queryWords = query.trim().split(/\b/);

        dialogue = '<span class="dialogue">' + dialogue + "</span>";

        for (let word of queryWords) {
            if (word.trim() !== "") {
                let regex = new RegExp("\\b" + word + "\\b", "g");

                dialogue = dialogue.replaceAll(
                    regex,
                    '<span style="color: #fc4c0a;">' + word + "</span>",
                );
            }
        }

        return dialogue;
    }

    function moveQueryContainingDialoguesToFront(timeDialogues, query) {
        let dialoguesContainedQuery = [];
        let timeDialoguesContainedQuery = [];
        let dialoguesNotContainedQuery = [];
        let timeDialoguesNotContainedQuery = [];

        for (let i = 0; i < timeDialogues.length; i++) {
            if (timeDialogues[i].Dialogue.includes(query)) {
                dialoguesContainedQuery.push(timeDialogues[i].Dialogue);
                timeDialoguesContainedQuery.push(timeDialogues[i]);
            } else {
                dialoguesNotContainedQuery.push(timeDialogues[i].Dialogue);
                timeDialoguesNotContainedQuery.push(timeDialogues[i]);
            }
        }

        let finalArray = Object.entries(
            timeDialoguesContainedQuery.concat(timeDialoguesNotContainedQuery),
        );

        return finalArray;
    }
</script>

<div id="video_card">
    <h3>{video.title}</h3>

    <div class="thum-dialogue">
        <a href="https://www.youtube.com/watch?v={video.id}" target="_blank">
            <img
                class="video_card_thumbnail"
                src="https://img.youtube.com/vi/{video.id}/0.jpg"
                alt=""
            />
        </a>

        <div class="dialogues">
            {#each moveQueryContainingDialoguesToFront(video.diloguesMap, query) as [d, t], index (index)}
                <a
                    class="dialogue-link"
                    target="_blank"
                    href="https://www.youtube.com/watch?v={video.id}&t={t[
                        'Time'
                    ].replaceAll('"', '')}"
                >
                    {@html decorateDialogue(t["Dialogue"], query)}
                </a>
            {/each}
        </div>
    </div>
</div>

<style>
    #video_card {
        background-color: #313030;
        display: flex;
        align-items: flex-start;
        padding: 1em;
        margin-bottom: 1em;
        border-radius: 0.3em;
        flex-direction: column;
    }

    .video_card_thumbnail {
        width: 15vw;
        border-radius: 0.5em;
    }

    .thum-dialogue {
        display: flex;
        align-items: center;
    }
    .dialogues {
        display: flex;
        flex-direction: column;
        padding-left: 2em;
        text-align: start;
        height: 12vw;
        contain: content;
        overflow-y: scroll;
    }

    a {
        color: #8b8b8b;
    }

    .dialogue-link {
        padding-top: 0.7em;
    }

    @media only screen and (max-width: 768px) {
        /* For mobile phones: */
        [id*="video_card"] {
            width: 79vw;
        }
    }

    @media only screen and (max-width: 768px) {
        /* For mobile phones: */
        [class*="dialogues"] {
            height: 25vh;
        }
    }
</style>
