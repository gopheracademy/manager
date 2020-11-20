<script>
  export let segment;
</script>

<nav class="nav">
  <div class="nav-left">
    <a class="brand" href="/">GopherCon</a>
    <div class="tabs">
      <a
        class={segment == '2021' ? 'active' : ''}
        href="/gophercon/2021">Gophercon 2021</a>
    </div>
  </div>
  <div class="nav-right">
    <a class="button outline" href="/login">Log In</a>
  </div>
</nav>
