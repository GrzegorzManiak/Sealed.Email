<script lang="ts">
	import Check from "lucide-svelte/icons/check";
	import ChevronsUpDown from "lucide-svelte/icons/chevrons-up-down";
	import { tick } from "svelte";
	import { Button } from "@/ui/button";
	import { cn } from "$lib/utils.js";
    import * as Command from "@/ui/command";
    import * as Popover from "@/ui/popover";
	import ChevronDown from "lucide-svelte/icons/chevron-down";
	import * as Avatar from "$shadcn/avatar";


	const frameworks = [
		{
			value: "sveltekit",
			label: "SvelteKit"
		},
		{
			value: "next.js",
			label: "Next.js"
		},
		{
			value: "nuxt.js",
			label: "Nuxt.js"
		},
		{
			value: "remix",
			label: "Remix"
		},
		{
			value: "astro",
			label: "Astro"
		}
	];

	let open = false;
	let value = "";
	const avatar = `https://api.dicebear.com/9.x/lorelei/svg?seed=sdfg&options[mood][]=happy`;
	$: selectedValue =
		frameworks.find((f) => f.value === value)?.label ?? "beta.noise.email";

	// We want to refocus the trigger button when the user selects
	// an item from the list so users can continue navigating the
	// rest of the form with the keyboard.
	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}
</script>

<Popover.Root bind:open let:ids>
    <Popover.Trigger asChild let:builder>
        <Button
            builders={[builder]}
            variant="outline"
            role="combobox"
            aria-expanded={open}
            class="w-full justify-between bg-transparent"
        >
            {selectedValue}
            <ChevronDown size="16" class="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
    </Popover.Trigger>


    <Popover.Content class="w-[200px] p-0 my-1">
        <Command.Root>
            <Command.Input placeholder="Search framework..." />
            <Command.Empty>No framework found.</Command.Empty>
            <Command.Group>
                {#each frameworks as framework}
                    <Command.Item
                            value={framework.value}
                            onSelect={(currentValue) => {
                               value = currentValue;
                               closeAndFocusTrigger(ids.trigger);
                              }}
                                            >
                                                <Check
                                                        class={cn(
                                "mr-2 h-4 w-4",
                                value !== framework.value && "text-transparent"
                              )}
                        />
                        {framework.label}
                    </Command.Item>
                {/each}
            </Command.Group>
        </Command.Root>
    </Popover.Content>
</Popover.Root>