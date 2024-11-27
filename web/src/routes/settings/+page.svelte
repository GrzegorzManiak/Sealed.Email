<script lang='ts'>
    import {onMount} from "svelte";
    import {redirect} from "$lib/redirect";

    onMount(() => redirect('/settings/account'));
</script>
