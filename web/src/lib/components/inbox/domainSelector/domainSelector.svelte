<script lang="ts">
	import Check from "lucide-svelte/icons/check";
	import ChevronsUpDown from "lucide-svelte/icons/chevrons-up-down";
	import {onMount, tick} from "svelte";
	import { Button } from "@/ui/button";
	import { cn } from "$lib/utils.js";
    import * as Command from "@/ui/command";
    import * as Popover from "@/ui/popover";
	import ChevronDown from "lucide-svelte/icons/chevron-down";
	import * as Avatar from "$shadcn/avatar";
    import * as Stores from '$lib/stores';
	import {setAllDomains} from "../../../stores";
	import * as API from "$api/lib";
	import type {DomainBrief} from "../../../../api/lib/api/domain";
	import {get, writable} from "svelte/store";

	interface DomainOption {
		label: string;
		value: string;
	}

	let open = false;
	let value = "";
	let searchQuery = "";
	let filteredDomains: DomainOption[] = [];

	$: domains = [] as Array<DomainOption>;
	$: storedSelected = get(Stores.selectedDomain)?.domainName ?? "No domain selected";
	$: selectedValue = domains.find((f) => f.value === value)?.label ?? storedSelected;
	$: allDomains = [] as Array<DomainOption>;
	$: filteredDomains = searchDomains(searchQuery, allDomains);

	function searchDomains(query: string, allDomains: DomainOption[]): DomainOption[] {
		if (!query) return allDomains;
		const searchTerm = query.toLowerCase();
		return allDomains.filter(domain => {
			const segments = domain.label.toLowerCase().split('.');
			for (let i = 0; i < segments.length; i++)
				if (segments[i].includes(searchTerm)) return true;
        });
	}

	Stores.domains.subscribe((incomingDomains) => {
		const newDomains: Array<DomainOption> = [];
		if (!incomingDomains) return;

		const keys = Object.keys(incomingDomains);
		keys.forEach((key) => {
			const domain = incomingDomains[key];
			if (!domain.brief) return;
			newDomains.push({
				label: domain.service.Domain,
				value: domain.brief.domainID,
			});
		});

		allDomains = newDomains;
	});

	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => document.getElementById(triggerId)?.focus());
	}

	function selectDomain(domainId: string) {
		for (let i = 0; i < domains.length; i++) {
			if (domains[i].value !== domainId) continue;
			console.log('Selected domain:', domains[i].label);
			Stores.selectedDomain.set({
                domainID: domains[i].value,
                domainName: domains[i].label,
            });
			return;
		}
    }

    onMount(async () => {
        await setAllDomains();
    });
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

    <Popover.Content class="w-[300px] p-0 my-1">
        <Command.Root>
            <Command.Input placeholder="Search domains..." bind:value={searchQuery}/>
            <Command.Empty>No domains found.</Command.Empty>
            <Command.Group>
                {#each filteredDomains as framework}
                    <Command.Item
                        value={framework.value}
                        onSelect={(currentValue) => {
                            value = currentValue;
                            selectDomain(currentValue);
                            closeAndFocusTrigger(ids.trigger);
                        }}>

                        <Check class={cn("mr-2 h-4 w-4", value !== framework.value && "text-transparent")}/>
                        {framework.label}
                    </Command.Item>
                {/each}
            </Command.Group>
        </Command.Root>
    </Popover.Content>
</Popover.Root>