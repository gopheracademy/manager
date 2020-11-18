<script>
  import { onMount } from "svelte";
  import ConferenceService from "$components/client.gen.js";

  let conference = {};
  onMount(async () => {
    var conferenceService = new ConferenceService();

    const response = await conferenceService.get({});
    console.log(response);
    conference = response.conference;
  });
</script>

<style>
</style>

<main>
  <h1>Admin</h1>

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
