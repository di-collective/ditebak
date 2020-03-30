<script>
  import { Route } from 'svelte-router-spa'
  import { FixedAdjust as Adjust } from "@smui/top-app-bar"

  import TopBar from '../../components/TopBar.svelte'

  export let currentRoute, params
</script>

<style>
.grid-container {
  display: grid;
  /* grid-template-columns: 250px auto; */
  grid-template-columns: auto;
  padding-left: 16px;
  padding-right: 16px;
}
/* .grid-item {
} */
</style>
  
<main>
  <TopBar></TopBar>

  <div use:Adjust class="grid-container">
    <div class="grid-item">
    </div>


    <div class="grid-item">
      <section>
        <Route {currentRoute}  {params} />
      </section>
    </div>

  </div>

</main>