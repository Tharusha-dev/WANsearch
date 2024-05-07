<script>
  import { fade } from "svelte/transition";
  let quotes_examples = [
    "apple claiming utterly implausible customer satisfaction",
    "it's a tax write off",
    "sending pictures of my balls and taint to the chinese government",
    "kim kardashian probably needs a butt helmet",
    "i need a freaking dash cam",
  ];
  var quote_example =
    quotes_examples[Math.floor(Math.random() * quotes_examples.length)];
</script>

<div class="tutorial-card" out:fade>
  <span
    >To search for a specific phrase, surround your search terms in quotation
    marks (faster). For example</span
  >
  <span class="example-quotes">"{quote_example}"</span>
  <span class="small-text"
    >Broader search term combinations (slower) might take some time to complete,
    please be patient</span
  >
</div>

<style>
  .tutorial-card {
    background-color: #313030;
    width: 35vh;
    height: 16vh;
    position: fixed;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    text-align: start;
    bottom: 0;
    left: 0;
    padding: 1.5vw;
    border-radius: 0 2vw 0 0;
  }

  .example-quotes {
    margin-top: 0.5em;
    color: #a3a3a3;
  }

  .small-text {
    margin-top: 0.5em;

    font-size: 0.8em;
  }

  @media only screen and (max-width: 1650px) {
    [class*="tutorial-card"] {
      font-size: 1vw;
    }
  }

  @media only screen and (max-width: 1300px) {
    [class*="tutorial-card"] {
      font-size: 12px;
    }
  }

  @media only screen and (max-width: 500px) {
    [class*="tutorial-card"] {
      margin-bottom: 10vh;
      border-radius: 0 2vw 2vw 0;
    }
  }
</style>
