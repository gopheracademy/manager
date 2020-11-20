<script context="module">
  export async function preload({ params }) {
    console.log(params.conference);
    return { conference: params.conference };
  }
</script>

<script>
  import { onMount } from "svelte";
  import ConferenceService from "$components/client.gen.js";
  export let conference = {};
  console.log(conference);

  onMount(async () => {
    var conferenceService = new ConferenceService();

    const response = await conferenceService.getBySlug({ slug: conference });
    console.log(response);
    conference = response.conference;
  });
</script>

<style>
</style>

{conference.name}
{#each conference.events as event}<a href={event.slug}>{event.name}</a>{/each}
