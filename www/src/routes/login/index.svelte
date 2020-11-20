<script>
  async function handleClick(event) {
    let response = await fetch("https://content.gophercon.com/auth/local", {
      method: "POST",
      headers: {
        "Content-Type": "application/json;charset=UTF-8",
      },
      body: JSON.stringify({
        identifier: "mail@bjk.fyi",
        password: "testing123",
      }),
    });

    let result = await response.json();
    console.log(result);
  }
</script>

<style>
</style>

<main>
  <h1>Log In</h1>
  <button on:click={handleClick}> log In </button>
</main>
