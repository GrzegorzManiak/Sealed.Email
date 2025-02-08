<script lang="ts">
    export let className: string | undefined | null = undefined;
    export let text: string;
    import {cn} from "$lib/utils";
</script>

<p class={cn("text-muted-foreground text-xs font-medium leading-tight", className)}>{text}</p>
