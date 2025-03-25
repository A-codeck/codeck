<script lang="ts">
    let pokeName = $state('pikachu');
    let pokemonUrl = $state();

    async function searchPokemon(pokeName: string): Promise<string> {
        return new Promise((fulfil, reject) => {
            fetch(`http://localhost:8001/?pokename=${pokeName}`,{
            })
            .then(res => res.json())
            .then(data => fulfil(data.pokemonUrl))
            .catch(error => reject(error));
        });
    }
</script>

<input bind:value={pokeName} />
<button onclick={() => pokemonUrl = searchPokemon(pokeName)}>
    Buscar Pokémon
</button>

{#await pokemonUrl}
    <p>Estamos procurando seu amiguinho.</p>
{:then pokemonUrl}
    <img src={pokemonUrl} alt="Imagem do Pokémon"/>
{:catch error}
    <p>Seu amiguinho morreu.</p>
{/await}

