<script lang="ts">
    import { onMount } from 'svelte';
    import { redirectIfLoggedIn } from "$lib/redirect";
    onMount(() => redirectIfLoggedIn('/authentication/login'));
</script>