<script lang='ts'>
	import * as Resizable from '$shadcn/resizable';
	import {EmailCard} from "@/inbox/emailCard";
	import { Button } from '$shadcn/button';
	import { Input } from '$shadcn/input';
	import { Label } from '$shadcn/label';
	import {type Writable, writable} from "svelte/store";
	import {InboxManager} from "$lib/inboxManager";
	import SelectedInboxHeader from "@/inbox/selectedInboxHeader.svelte";
	import {cn} from "$lib/utils";
	import {onMount} from "svelte";
	import VirtualList from 'svelte-tiny-virtual-list';
	import InfiniteLoading from 'svelte-infinite-loading';
	import type {EmailMetadata} from "../../../api/lib/services/emailStorage";
    import * as Stores from '$lib/stores';

	export let headerHeight: Writable<number>;
	export let stateManager: Writable<string>;

	const inboxManager = new InboxManager('1', 'Encrypted Inbox');
	const groupCounter = writable(0);
	const groupSelectStore = writable(new Set<string>());
	const selectedStore = writable<string>();
	const collapsed = writable(false);

	let data: Array<EmailMetadata> = [];
	let session: Session;
	let page = 0;

	Stores.user.subscribe((user) => {
		if (!user || !user.session || !user.isLoggedIn) return;
		session = user.session;
		console.log('Session loaded');
    });

	async function infiniteHandler({ detail: { complete, error } }) {
		try {
            if (!session) return;
            console.log('Loading more data', page);
			complete();
		}

		catch (e) {
			error();
		}
	}

	onMount(async () => {
        // data =
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
        <VirtualList width="100%" height={600} itemCount={data.length} itemSize={50}>
            <div slot="item" let:index let:style {style}>
                Letter: {data[index]}, Row: #{index}
            </div>

            <div slot="footer">
                <InfiniteLoading on:infinite={infiniteHandler} />
            </div>
        </VirtualList>
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