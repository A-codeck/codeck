<script lang="ts">
    let arg = $props();
    let pokeName = arg;
    let pokemonUrl;

    async function searchPokemon(pokeName: string): Promise<string> {
        return new Promise((fulfil, reject) => {
            fetch(`http://localhost:8001/?pokename=${pokeName}`,{
            })
            .then(res => res.json())
            .then(data => fulfil(data.pokemonUrl))
            .catch(error => reject(error));
        });
    }

    pokemonUrl = searchPokemon(pokeName);
</script>

<h1> You are looking for {arg} </h1>
<h2> Here he is:  </h2>
<img src={pokemonUrl} alt="Imagem do Pokémon"/>


