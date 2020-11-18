<script>
  export let segment;
</script>

<nav class="nav">
  <div class="nav-left">
    <a class="brand" href="/">GopherCon</a>
    <div class="tabs">
      <a
        class={segment == 'sponsors' ? 'active' : ''}
        href="/sponsors">Sponsors</a>
      <a class={segment == 'agenda' ? 'active' : ''} href="/agenda">Agenda</a>
      <a
        class={segment == 'speakers' ? 'active' : ''}
        href="/speakers">Speakers</a>

      <a class={segment == 'coc' ? 'active' : ''} href="/coc">Code of Conduct</a>
    </div>
  </div>
  <div class="nav-right">
    <a class="button outline" href="/login">Log In</a>
  </div>
</nav>
