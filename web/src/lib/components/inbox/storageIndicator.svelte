<script lang="ts">
    import {cn} from "$lib/utils";

    const GB = 1073741824;

    export let label: string = "10 GB";
    export let maxStorage: number = 10737418240;
    export let currentStorage: number = 3436500000;

    const colors = {
        "100": "bg-red-500",
        "95": "bg-red-400",
        "90": "bg-yellow-500",
        "85": "bg-yellow-400",
        "80": "bg-yellow-300",
        "75": "bg-green-200",
        "50": "bg-green-300",
        "25": "bg-green-400",
        "0": "bg-green-500"
    };

    $: gbUsed = currentStorage / GB;
    $: percentage = (currentStorage / maxStorage) * 100;
    $: color = Object.entries(colors).find(([key, value]) => percentage <= parseInt(key)) || ["100", "bg-red-500"];
</script>

<div class="flex flex-col items-start gap-1 w-full px-2">
    <div class="rounded-md relative h-1 inset-0 bg-background bg-opacity-30 w-full flex-grow overflow-clip">
        <div class={cn("absolute h-1 bg-primary rounded-r-md", color[1])} style="width: {percentage}%;"></div>
    </div>

    <div class="flex justify-between w-full gap-1 align-middle mt-1">
        <p class="text-muted-foreground text-sm text-center">
            {gbUsed.toFixed(2)} / {label} used
        </p>

        <a href="/settings" class="text-muted-foreground text-xs text-center underline hover:text-primary">
            Upgrade
        </a>
    </div>

</div>