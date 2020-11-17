<script context="module">
  export const prerender = true;
</script>

<script>
  export let segment;
  import Nav from "$components/Nav.svelte";
</script>

<style>
  main {
    position: relative;
    max-width: 56em;
    padding: 2em;
    margin: 0 auto;
    box-sizing: border-box;
  }
</style>

<Nav {segment} />

<main>
  <slot />
</main>
