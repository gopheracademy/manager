# Content Management 

Content management is powered by [strapi](https://strapi.io), which provides a graphical 
user interface for managing content.

## Content Types    

Strapi provides for Collection Types and Single types.  A Single type is a piece of content
that can only exist in one place.  `Home` is a good example.  The structure of the content
on the home page is different from any other, so it's a Single type.

Collection Types are groups of pages that share the same structure.  Each page in a collection
has the same fields. Most interior pages on the site will be collection pages.  They'll share attributes
like "title", "content", "slug", etc.

## Fetching Content

Content is available via the Strapi API through REST type calls that return JSON data. In the svelte pages, 
content is loaded in a "module" context, which happens during the page build(server side), as opposed to a regular script
block that would happen on page load client-side.

```javascript
<script context="module">
  export const prerender = true;
  export async function preload({ params }) {
    const res = await this.fetch(`https://content.gophercon.com/home`);
    const homepage = await res.json();
    return { homepage };
  }
</script>
```
The content in the object `homepage` is available to the client-side script block:
```html
<script>
  export let homepage;
  let error = null;
</script>

<style>
</style>

<h1>{homepage.title}</h1>
<img src="https://content.gophercon.com{homepage.head_image.url}" />
<div>{homepage.hero_text}</div>
```