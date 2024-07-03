<script>
  import VideoCard from "./video_card.svelte";
  import TutorialCard from "./tutorial_card.svelte";
  import CountsCard from "./counts_card.svelte";

  let query = "";
  let responseGetting = false;
  let responseGot = false;
  let searchDone = true;
  let isNullResult = false;
  let isQuotedSearch = false;
  let fdata;
  let videosJson;
  let countsJson;


  let videos = [];

  async function getDocument() {
    searchDone = false;
    responseGetting = true;
    responseGot = false;


    const response = await fetch(`https://api.wansearch.xyz/search?q=${query}`);
    if (response.ok) {
      const jsonData = await response.json();
      fdata = jsonData;
      videosJson = jsonData.Videos
      
      countsJson = jsonData.Counts

      if (countsJson ==  null){
      
          isQuotedSearch = true
      }
      

      if (videosJson == null ){
        isNullResult = true;
       
      responseGot = true;

      }else{
        isNullResult = false
        
        videos = videosJson.map((video) => ({
        id: video.Video_id,
        title: video.Title,
        
        diloguesMap : video.TimeDialogues2
      }));

   

      }
    
    
      
      responseGot = true;
    } else {
      throw new Error(`Error fetching data: ${response.status}`);
    }
  }
</script>

<div class="search-bar">
  <svg
    width="2em"
    height="2em"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M15.7955 15.8111L21 21M18 10.5C18 14.6421 14.6421 18 10.5 18C6.35786 18 3 14.6421 3 10.5C3 6.35786 6.35786 3 10.5 3C14.6421 3 18 6.35786 18 10.5Z"
      stroke="white"
      stroke-width="1"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
  <input bind:value={query} type="search" name="" class="main_input" />
</div>

<br />
<button  class="search-button" on:click={getDocument}> Search </button>

<br />

{#if responseGetting}
  {#if !responseGot}
    <div class="loading-bar">
      <div class="loader"></div>
    </div>
  {/if}

  {#if !isNullResult}
  {#if responseGot}

  {#if !isQuotedSearch}
  <CountsCard counts={countsJson}/>
  {/if}
  {/if}
  {#each videos as video}
      
    <VideoCard video={video} query={query} />

  {/each}

  {:else}

  <span>No result found</span>
  {/if}


{:else}
  <p></p>
{/if}


{#if searchDone}
<TutorialCard />
{/if}

<style>
  .main_input {
   
    background-color: #313030;

    height: 5vh;
    width: 30vw;
    border: 0;
    font-size: 1.2em;
    padding-left: 2em;
  }

  .loader {
    border: 0.2em solid #f3f3f3; /* Light grey */
    border-top: 0.2em solid white; /* Blue */
    border-radius: 50%;
    width: 2em;
    height: 2em;
    animation: spin 2s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .loading-bar {
    display: flex;
  justify-content: center;
  padding-top: 2em;
  }

  .search-bar {
   
    background-color: #313030;
    display: flex;
    border: #242424 solid 0.2em;
    align-items: center;
    padding: 0.2em;
    padding-right: 1.2em;
    padding-left: 1em;
    

    border-radius: 2em;
  }

  .search-button {
    margin-bottom: 1.1em;
  }



</style>
