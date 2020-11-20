<script context="module">
  export const prerender = true;
  export async function preload({ params }) {
    const res = await this.fetch(`https://content.gophercon.com/home`);
    const homepage = await res.json();
    return { homepage };
  }
</script>

<script>
  export let homepage;
  let error = null;
</script>

<style>
</style>

<h1>{homepage.title}</h1>
<img src="https://content.gophercon.com{homepage.head_image.url}" />
<div>{homepage.hero_text}</div>
