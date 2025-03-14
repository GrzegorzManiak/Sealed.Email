<script lang='ts'>
    // -- Svelte imports
    import { get } from 'svelte/store';
    import { onMount } from "svelte";
    import { type Writable, writable } from "svelte/store";

    // -- Third-party library imports
    import VirtualList from 'svelte-tiny-virtual-list';
    import InfiniteLoading from 'svelte-infinite-loading';

    // -- UI component imports
    import * as Resizable from '$shadcn/resizable';
    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
    import { Label } from '$shadcn/label';

    // -- Local UI components
    import { EmailCard } from "@/inbox/emailCard";
    import SelectedInboxHeader from "@/inbox/selectedInboxHeader.svelte";

    // -- Utilities and services
    import { cn } from "$lib/utils";
    import { InboxManager } from "$lib/inboxManager";
    import * as Stores from '$lib/stores';
    import * as API from '$api/lib';
	import {DomainService, EmailProvider, EmailStorage, Session} from "../../../api/lib";

    // -- Types
    import type { User } from "$lib/stores";
    import type { EmailMetadata } from "$api/lib/services/emailStorage";
	import {StorageService} from "$api/lib/services/storageServices";



	export let headerHeight: Writable<number>;
	export let stateManager: Writable<string>;

	const inboxManager = new InboxManager('1', 'Encrypted Inbox');
	const groupCounter = writable(0);
	const groupSelectStore = writable(new Set<string>());
	const selectedStore = writable<string>();
	const collapsed = writable(false);

	let storageService: StorageService | null = null;
	let domainService: DomainService | null = null;
	let emailService: EmailStorage | null = null;
	let emailProvider: EmailProvider | null = null;

	let data: Array<EmailMetadata> = [];
	let session: Session;
	let page = 0;

	async function reloadUser(user: User) {
        if (!user || !user.session || !user.isLoggedIn) return;
        session = API.Session.Deserialize(user.session);
        await session.DecryptKeys();
        if (!session || session instanceof Error) {
            console.info('No session');
            return;
        }
        console.log('Session loaded');
    }

	let reloadInProgress = false;
	async function reloadDomainService() {
        while (reloadInProgress) {
            console.log('Waiting for previous reload to finish...');
            await new Promise(resolve => setTimeout(resolve, 100));
        }
		reloadInProgress = true;

        if (!session) {
			reloadInProgress = false;
            return console.error('No session', session);
        }

		const currentDomain = get(Stores.selectedDomain);
		if (!currentDomain) {
            reloadInProgress = false;
			return console.error('No domain selected', currentDomain);
        }

		const domains = get(Stores.domains);
		const domain = domains[currentDomain.domainID];
		if (!domain) {
            reloadInProgress = false;
			return console.error('No domain found', domains);
        }

		storageService = new API.StorageServices.IndexedDBStorageService();
        domainService = await API.DomainService.Decrypt(session, domain.full);
		emailService = new API.EmailStorage(storageService, session);
		emailProvider = new API.EmailProvider(emailService, session);

		console.log(`All domain & email services loaded for ${domainService.Domain}`);
		reloadInProgress = false;

		const newEmails = await emailProvider.getEmails(domainService, {
			domainID: domainService.DomainID,
			order: 'desc',
			perPage: 10
		});

		console.log(newEmails);
    }

	// async function infiniteHandler({ detail: { complete, error } }) {
	// 	try {
	// 		console.log('Loading more emails...');
    //         if (!session || !domainService) {
    //             console.warn('No session or domain service', session, domainService);
    //             return error();
	// 		}
    //
	// 		try {
	// 			const newEmails = await emailProvider.getEmails(domainService, {
	// 				domainID: domainService.DomainID,
	// 				order: 'desc',
	// 				perPage: 10
	// 			});
    //
	// 			if (!newEmails) {
	// 				console.error('No emails found', newEmails);
	// 				return error();
	// 			}
    //
	// 			data = data.concat(newEmails);
	// 			console.log('Emails loaded:', data.length);
	// 			complete();
    //         }
    //
    //         catch (e) {
    //             console.error('Error loading emails', e);
    //             error();
    //         }
	// 	}
    //
	// 	catch (e) {
	// 		error();
	// 	}
	// }

	// Stores.user.subscribe(async (user) => reloadUser(user));
    // Stores.selectedDomain.subscribe(async () => reloadDomainService());

	onMount(async () => {
        const user = get(Stores.user);
		await reloadUser(user);
		await reloadDomainService();
    });
</script>

<Resizable.Pane
        maxSize={100}
        minSize={30}
        defaultSize={40}
        collapsedSize={0}
        collapsible={false}
        onCollapse={() => collapsed.set(true)}
        onExpand={() => collapsed.set(false)}
        class="flex-shrink-0 h-full">

    <span bind:offsetHeight={$headerHeight}>
        <SelectedInboxHeader {collapsed} {inboxManager}/>
    </span>

    <div class={cn({ 'hidden': $collapsed })}>
<!--        <VirtualList width="100%" height={600} itemCount={data.length} itemSize={50}>-->
<!--            <div slot="item" let:index let:style {style}>-->
<!--                Letter: {data[index]}, Row: #{index}-->
<!--            </div>-->

<!--            <div slot="footer">-->
<!--                <InfiniteLoading on:infinite={infiniteHandler} />-->
<!--            </div>-->
<!--        </VirtualList>-->
    </div>

<!--        &lt;!&ndash; Fixed height pane content &ndash;&gt;-->
<!--        <EmailCard-->
<!--                groupCounter={groupCounter}-->
<!--                selectedStore={selectedStore}-->
<!--                groupSelectStore={groupSelectStore}-->
<!--                data={{-->
<!--        id: '1',-->
<!--        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',-->

<!--        fromEmail: 'bob@dylan.com',-->
<!--        fromName: 'Bob Dylan',-->

<!--        toEmail: 'me@grzegorz.ie',-->
<!--        toName: 'Greg',-->

<!--        subject: 'Re: Welcome to Noise Email!',-->
<!--        body: 'Thank you for your email. I will get back to you as soon as possible.',-->

<!--        read: false,-->
<!--        pinned: false,-->
<!--        starred: false,-->

<!--        chain: [-->
<!--            {-->
<!--        id: '11',-->
<!--        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',-->

<!--        fromEmail: 'bob@dylan.com',-->
<!--        fromName: 'Bob Dylan',-->

<!--        toEmail: 'me@grzegorz.ie',-->
<!--        toName: 'Greg',-->

<!--        subject: 'Re: Welcome to Noise Email!',-->
<!--        body: 'Thank you for your email. I will get back to you as soon as possible.',-->

<!--        read: false,-->
<!--        pinned: false,-->
<!--        starred: false,-->

<!--        totalAttachments: 0-->
<!--    },{-->
<!--        id: '12',-->
<!--        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',-->

<!--        fromEmail: 'bob@dylan.com',-->
<!--        fromName: 'Bob Dylan',-->

<!--        toEmail: 'me@grzegorz.ie',-->
<!--        toName: 'Greg',-->

<!--        subject: 'Re: Welcome to Noise Email!',-->
<!--        body: 'Thank you for your email. I will get back to you as soon as possible.',-->

<!--        read: false,-->
<!--        pinned: false,-->
<!--        starred: false,-->


<!--        totalAttachments: 0-->
<!--    },{-->
<!--        id: '13',-->
<!--        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',-->

<!--        fromEmail: 'bob@dylan.com',-->
<!--        fromName: 'Bob Dylan',-->

<!--        toEmail: 'me@grzegorz.ie',-->
<!--        toName: 'Greg',-->

<!--        subject: 'Re: Welcome to Noise Email!',-->
<!--        body: 'Thank you for your email. I will get back to you as soon as possible.',-->

<!--        read: false,-->
<!--        pinned: false,-->
<!--        starred: false,-->

<!--        showcaseAttachment: {-->
<!--            id: '1',-->
<!--            filename: 'noise-email-welcome.pdf',-->
<!--            type: 'pdf',-->
<!--        },-->
<!--        totalAttachments: 5-->
<!--    }-->
<!--        ],-->

<!--        showcaseAttachment: {-->
<!--            id: '1',-->
<!--            filename: 'noise-email-welcome.pdf',-->
<!--            type: 'pdf',-->
<!--        },-->
<!--        totalAttachments: 5-->
<!--    }}/>-->
<!--        <EmailCard-->
<!--                groupCounter={groupCounter}-->
<!--                selectedStore={selectedStore}-->
<!--                groupSelectStore={groupSelectStore}-->
<!--                data={-->
<!--        {-->
<!--            id: '2f',-->
<!--            date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',-->

<!--            fromEmail: 'TOS@Noise.email',-->
<!--            fromName: 'Noise Email',-->

<!--            toEmail: 'Unknown',-->
<!--            toName: 'Unknown',-->

<!--            subject: 'Terms of Service Update',-->
<!--            body: 'We have updated our terms of service. Please review them at your earliest convenience.',-->

<!--            read: true,-->
<!--            pinned: true,-->
<!--            starred: false,-->

<!--            totalAttachments: 0-->
<!--            }-->
<!--    }/>-->

<!--        <p class="text-muted-foreground text-sm text-center p-4">No more emails to show.</p>-->
<!--    </div>-->

</Resizable.Pane>