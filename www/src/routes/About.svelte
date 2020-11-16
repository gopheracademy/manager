<script>
  import { onMount } from "svelte";
  import ConferenceService from "../client.gen.js";

  let conference = {};
  onMount(async () => {
    var conferenceService = new ConferenceService();

    const response = await conferenceService.get({});
    console.log(response);
    conference = response.conference;
  });
</script>

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  }

  main {
    text-align: center;
    padding: 1em;
    margin: 0 auto;
  }

  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 4rem;
    font-weight: 100;
    line-height: 1.1;
    margin: 4rem auto;
    max-width: 14rem;
  }

  p {
    max-width: 14rem;
    margin: 2rem auto;
    line-height: 1.35;
  }

  @media (min-width: 480px) {
    h1 {
      max-width: none;
    }

    p {
      max-width: none;
    }
  }
</style>

<main>
  <h1>Hello world!</h1>

  {#await conference then value}
    <!-- promise was fulfilled -->
    <p>Conference: {value.name}</p>
  {:catch error}
    <!-- promise was rejected -->
    <p>Something went wrong: {error.message}</p>
  {/await}
  <p>
    Visit the
    <a href="https://svelte.dev">svelte.dev</a>
    to learn how to build Svelte apps.
  </p>
</main>
