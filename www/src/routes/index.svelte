<script>
  import Counter from "$components/Counter.svelte";
  import { onMount } from "svelte";
  let homepage = {};
  let error = null;

  onMount(async () => {
    const parseJSON = (resp) => (resp.json ? resp.json() : resp);
    const checkStatus = (resp) => {
      if (resp.status >= 200 && resp.status < 300) {
        return resp;
      }
      return parseJSON(resp).then((resp) => {
        throw resp;
      });
    };
    const headers = {
      "Content-Type": "application/json",
    };

    try {
      const res = await fetch("http://localhost:1337/home", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(checkStatus)
        .then(parseJSON);
      homepage = res;
      console.log(homepage);
    } catch (e) {
      error = e;
    }
  });
</script>

<style>
</style>

{#if error !== null}{error}{:else}{homepage.title}{/if}
