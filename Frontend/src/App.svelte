<script>
  import Input from "./lib/Input.svelte";
  import Footer from "./lib/footer.svelte";
</script>

<main>
  <h1>WANsearch</h1>

  <Input />
</main>

<Footer />

<div class="github-logo">
  <a
  href="https://github.com/Tharusha-dev/WANsearch"
  target="_blank"
  aria-label="github"
>
  <img
    class="github-icon"
    src="github-mark-white.svg"
    alt="github-icon"
    srcset=""
  />
</a>
</div>

<div class="buy-me-coffee">
  <a
    href="https://buymeacoffee.com/tharushaj"
    target="_blank"
    aria-label="Support me"
  >
    <img
      class="buy-me-coffee-black"
      src="buymeacoffee-black.png"
      alt="buy-me-coffee-logo"
      srcset=""
    />
  </a>
</div>

<style>
  h1 {
    color: #fc4c0a;
  }

  .github-logo{
    position: fixed; /* Positioning relative to the viewport */
    top: 0; /* 0px from the bottom */
    right: 0; 
  }

  .github-icon {
    height: 4vh;

  }

  .buy-me-coffee {
    /* height: 5vh;
    width: 5vw; */
    position: fixed; /* Positioning relative to the viewport */
    bottom: 0; /* 0px from the bottom */
    right: 0; /* 0px from the right */
  }

  .buy-me-coffee-black {
    height: 6vh;
  }
</style>
