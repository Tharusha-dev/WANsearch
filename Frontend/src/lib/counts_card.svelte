<script>
    export let counts;

    console.log(counts);

    //     for (const [key, value] of Object.entries(object1)) {
    //   console.log(`${key}: ${value}`);
    // }
</script>

<div class="counts-card">
    {#each Object.entries(counts) as [term, count]}
        <span>
            <span id="term">{term}&nbsp;</span>:<span id="count">&nbsp;{count}</span>
        </span>
        
    {/each}
    <span id="count-samllprint"><a href="https://github.com/Tharusha-dev/searchYard/blob/main/searchYardAPI/stop_words.txt">These</a> words, and words cencensored by YouTube substitles are missing from the counts and/or search</span>
</div>

<style>
    .counts-card {
        display: flex;
        flex-direction: column;
        align-items: flex-start;

        background-color: #222;

        padding: 1em;
        margin-bottom: 1em;
        border-radius: 0.3em;
    }

    #term {
        color: white;
        font-weight: bold;
    }

    #count-samllprint{
        font-size: 0.75em;
    }

    #count-samllprint a{
       color: white
    }

</style>
