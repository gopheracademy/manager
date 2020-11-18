<script context="module">
  export const prerender = true;
  export async function preload({ params }) {
    const res = await this.fetch(
      `http://content.44ec396da6024f08b75a.eastus2.aksapp.io/home`
    );
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
<img
  src="http://content.44ec396da6024f08b75a.eastus2.aksapp.io{homepage.head_image.url}" />
<div>{homepage.hero_text}</div>
