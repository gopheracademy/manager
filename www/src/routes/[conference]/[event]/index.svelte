<script context="module">
  export async function preload({ params }) {
    console.log(params.event);
    return { conference: params.conference, event: params.event };
  }
</script>

<script>
  import { onMount } from "svelte";
  import ConferenceService from "$components/client.gen.js";
  export let conference = {};
  export let event = {};

  console.log(conference);
  console.log(event);

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
