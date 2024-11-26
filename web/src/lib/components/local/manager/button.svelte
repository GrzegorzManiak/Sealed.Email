<script lang="ts">
    import { Button } from '$shadcn/button';
    import { ManagerPopup } from "$local/manager";

    type Props = {};
    type $$props = Props & { class?: string };
    let className: $$props["class"];
    const props = $$props as Props & { class?: string };

    let showingPopup = false;
    function SetPopup(visible: boolean) {
        console.log("Setting popup to: ", visible);
        showingPopup = visible;
    }
</script>

<Button
        variant="outline"
        class={className}
        on:click={() => SetPopup(!showingPopup)}
>
    Dots Settings
</Button>

<ManagerPopup bind:showing={showingPopup} />