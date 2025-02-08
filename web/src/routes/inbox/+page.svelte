<script lang='ts'>
	import * as Resizable from '$shadcn/resizable';
	import {EmailCard} from "@/inbox/emailCard";
	import { Button } from '$shadcn/button';
	import { Input } from '$shadcn/input';
	import { Label } from '$shadcn/label';

	import {writable} from "svelte/store";
    import {InboxManager} from "$lib/inboxManager";
	import {SquarePen} from "lucide-svelte";
    const manger = new InboxManager('1', 'Inbox');

    const groupCounter = writable(0);
    const groupSelectStore = writable(new Set<string>());
    const selectedStore = writable<string>();
</script>

<Resizable.Pane class="flex-shrink-0 h-full">
    <div class="flex flex-col items-start justify-end p-4 py-4 gap-5 border-b h-36">
        <h1 class="text-2xl font-bold">Encrypted Inbox</h1>

        <div class="flex flex-row gap-2 w-full">
            <Input id="search" type="text" placeholder="Search" class="w-full rounded-md bg-primary-foreground bg-opacity-40 text-foreground"/>
            <Button variant="ghost" class="bg-foreground rounded-md ">
                <SquarePen class="h-5 w-5 text-muted"/>
            </Button>

        </div>
    </div>

    <!-- Fixed height pane content -->
    <EmailCard
            groupCounter={groupCounter}
            selectedStore={selectedStore}
            groupSelectStore={groupSelectStore}
            data={{
        id: '1',
        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',

        fromEmail: 'bob@dylan.com',
        fromName: 'Bob Dylan',

        toEmail: 'me@grzegorz.ie',
        toName: 'Greg',

        subject: 'Re: Welcome to Noise Email!',
        body: 'Thank you for your email. I will get back to you as soon as possible.',

        read: false,
        pinned: false,
        starred: false,

        chain: [
            {
        id: '11',
        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',

        fromEmail: 'bob@dylan.com',
        fromName: 'Bob Dylan',

        toEmail: 'me@grzegorz.ie',
        toName: 'Greg',

        subject: 'Re: Welcome to Noise Email!',
        body: 'Thank you for your email. I will get back to you as soon as possible.',

        read: false,
        pinned: false,
        starred: false,

        totalAttachments: 0
    },{
        id: '12',
        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',

        fromEmail: 'bob@dylan.com',
        fromName: 'Bob Dylan',

        toEmail: 'me@grzegorz.ie',
        toName: 'Greg',

        subject: 'Re: Welcome to Noise Email!',
        body: 'Thank you for your email. I will get back to you as soon as possible.',

        read: false,
        pinned: false,
        starred: false,


        totalAttachments: 0
    },{
        id: '13',
        date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',

        fromEmail: 'bob@dylan.com',
        fromName: 'Bob Dylan',

        toEmail: 'me@grzegorz.ie',
        toName: 'Greg',

        subject: 'Re: Welcome to Noise Email!',
        body: 'Thank you for your email. I will get back to you as soon as possible.',

        read: false,
        pinned: false,
        starred: false,

        showcaseAttachment: {
            id: '1',
            filename: 'noise-email-welcome.pdf',
            type: 'pdf',
        },
        totalAttachments: 5
    }
        ],

        showcaseAttachment: {
            id: '1',
            filename: 'noise-email-welcome.pdf',
            type: 'pdf',
        },
        totalAttachments: 5
    }}/>
    <EmailCard
            groupCounter={groupCounter}
            selectedStore={selectedStore}
            groupSelectStore={groupSelectStore}
            data={
        {
            id: '2f',
            date: 'Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)',

            fromEmail: 'TOS@Noise.email',
            fromName: 'Noise Email',

            toEmail: 'Unknown',
            toName: 'Unknown',

            subject: 'Terms of Service Update',
            body: 'We have updated our terms of service. Please review them at your earliest convenience.',

            read: true,
            pinned: true,
            starred: false,

            totalAttachments: 0
            }
    }/>
</Resizable.Pane>

<Resizable.Handle withHandle />

<Resizable.Pane class="flex-grow h-full overflow-y-visible">
</Resizable.Pane>

