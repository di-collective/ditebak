<script>
  import { Route } from 'svelte-router-spa'
  import { FixedAdjust as Adjust } from "@smui/top-app-bar"

  import TopBar from 'components/TopBar.svelte'
  import BottomBar from 'components/BottomBar.svelte'

  export let currentRoute, params
</script>

<style>
main {
  min-height: 75vh;
}
.grid-container {
  display: grid;
  grid-template-columns: auto;
  padding-left: 16px;
  padding-right: 16px;
}
</style>
  
<main>
  <TopBar></TopBar>

  <div use:Adjust class="grid-container mobile-wd">
    <div></div>

    <div>
      <section>
        <Route {currentRoute}  {params} />
      </section>
    </div>
  </div>
</main>

<BottomBar>
</BottomBar>