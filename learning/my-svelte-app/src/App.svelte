<script>
  import { Router, Route, Link } from "svelte-routing";
  import Home from "./Home.svelte";
  import PokeQuery from "./PokeQuery.svelte";
</script>

<Router>
  <!-- The main page at "/" -->
  <Route path="/">
      <Home />
  </Route>

  <!-- Dynamic route: "/:arg" where "arg" is a parameter -->
  <Route path="/:arg" let:params>
    <PokeQuery arg={params.arg} />
  </Route>
</Router>

