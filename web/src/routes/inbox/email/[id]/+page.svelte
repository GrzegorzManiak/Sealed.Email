<script lang="ts">
	import {page} from '$app/stores';
	import {onMount} from "svelte";
	import PostalMime, {type Email} from 'postal-mime';
	import {writable, type Writable} from "svelte/store";
	import {RandomHEXColor} from "$lib/common";
	import * as Avatar from "$shadcn/avatar";

	const devUrl = '/test/email.txt';

	export let headerHeight: Writable<number> = writable(192 / 2);
	let id: string = '';
	let parsedEmail: Email | null = null;
	let color = RandomHEXColor();
	let iframe: HTMLIFrameElement;

	async function getEmailData(): Promise<string> {
        const res = await fetch(devUrl);
        const emailData = await res.text();
		console.log('Fetched email data');
		return emailData;
    }

	async function parseEmailData(emailData: string) {
		parsedEmail = await PostalMime.parse(emailData);
    }

	async function loadEmailData() {
		const emailData = await getEmailData();
		await parseEmailData(emailData);

		if (iframe && parsedEmail) {
			iframe.srcdoc = (parsedEmail.html || parsedEmail.text) || '';
        }
	}

	onMount(() => {
		id = ''

		page.subscribe((value) => {
			console.log('Page changed:', value);
			if (id === value.params.id) return;
			id = value.params.id;
			loadEmailData();
		});

		headerHeight.subscribe((value) => {
			console.log('Header height changed:', value);
		});
    });
</script>

<div class="bg-primary-foreground border-b bg-opacity-40" style="height: {$headerHeight - 1}px; max-height: {$headerHeight - 1}px;">
    {#if parsedEmail}

        <div class="flex flex-col items-start gap-2 p-2">
            <div class="flex flex-row items-center gap-2">
                <Avatar.Root class="transition-colors duration-200 grid grid-cols-1 grid-rows-1 relative w-10 h-10" style="background-color: {color}">
                    <Avatar.Image src="https://api.dicebear.com/9.x/lorelei/svg?seed={parsedEmail.from.address}&options[mood][]=happy" class="rounded-full w-10 h-10" />
                    <Avatar.Fallback> {parsedEmail.from.name[0]} </Avatar.Fallback>
                </Avatar.Root>

                <div class="flex flex-col items-start justify-start w-full">
                    <p class="truncate font-bold">{parsedEmail.from.name} <span class="truncate font-normal text-sm text-muted-foreground">&lt;{parsedEmail.from.address}&gt;</span></p>
                    {#if parsedEmail.to}
                        <p class="truncate text-sm text-muted-foreground">
                            <span class="font-bold">To: </span>

                            {parsedEmail.to.map((to) => {
                                if (to.name) return `${to.name} <${to.address}>`;
                                return to.address;
                            }).join(', ')}</p>
                    {/if}

                    {#if parsedEmail.cc}
                        <p class="truncate text-sm text-muted-foreground">
                            <span class="font-bold">Cc: </span>

                            {parsedEmail.cc.map((cc) => {
                                if (cc.name) return `${cc.name} <${cc.address}>`;
                                return cc.address;
                            }).join(', ')}</p>
                    {/if}
                </div>
            </div>

            <h2>
                {parsedEmail.subject}
            </h2>
        </div>

        {#if parsedEmail.attachments.length > 0}
            <div class="flex flex-col items-start gap-2 p-2">
                <h3>Attachments</h3>
                <div class="flex flex-row items-center gap-2">
                    {#each parsedEmail.attachments as attachment}
                        <div class="flex flex-row items-center gap-2">
                            <Avatar.Root class="transition-colors duration-200 grid grid-cols-1 grid-rows-1 relative w-10 h-10" style="background-color: {color}">
                                <Avatar.Image src="https://api.dicebear.com/9.x/lorelei/svg?seed={attachment.filename}&options[mood][]=happy" class="rounded-full w-10 h-10" />
                                <Avatar.Fallback> {attachment.filename} </Avatar.Fallback>
                            </Avatar.Root>
                            <p>{attachment.filename}</p>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}

    {:else}
        <p>Loading...</p>
    {/if}
</div>

{#if parsedEmail}
    <div class="p-2">
        <div class="prose max-w-full" style="overflow-wrap: break-word;">
            <iframe
                    bind:this={iframe}
                    sandbox=""
                    class="w-full rounded-md"
                    style="height: calc(100vh - {$headerHeight - 8}px);" />
        </div>
    </div>
{/if}
