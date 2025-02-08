<script lang='ts'>
	import * as Resizable from '$shadcn/resizable';
	import {EmailCard} from "@/inbox/emailCard";
	import { Button } from '$shadcn/button';
	import { Input } from '$shadcn/input';
	import { Label } from '$shadcn/label';

	import {writable} from "svelte/store";
    import {InboxManager} from "$lib/inboxManager";
	import {Globe, Inbox, Mails, SquarePen} from "lucide-svelte";
	import {Checkbox} from "@/ui/checkbox";
	import ChevronLeft from "lucide-svelte/icons/chevron-left";
    const manger = new InboxManager('1', 'Inbox');

    const groupCounter = writable(0);
    const groupSelectStore = writable(new Set<string>());
    const selectedStore = writable<string>();
</script>

<Resizable.Pane class="flex-shrink-0 h-full">
    <div class="h-[181px] max-h-[181px] overflow-hidden">
        <div class="flex flex-col items-start justify-between align-bottom p-4 gap-2 border-b bg-primary-foreground bg-opacity-40 h-full">

            <div class="flex flex-col gap-1">
                <h1 class="text-2xl font-bold leading-tight p-0">Encrypted Inbox</h1>
                <div class="flex flex-row gap-2">

                    <div class="flex flex-row gap-1 items-center">
                        <Mails class="text-muted-foreground" size={15}/>
                        <p class="text-sm text-muted-foreground leading-tight">
                            1,400 Messages,
                        </p>
                    </div>

                    <div class="flex flex-row gap-1 items-center">
                        <Inbox class="text-muted-foreground" size={15}/>
                        <p class="text-sm text-muted-foreground leading-tight">
                            216 Unread
                        </p>
                    </div>

                </div>
            </div>

            <div class="flex flex-row gap-2 w-full pt-5">
                <Input id="search" type="text" placeholder="Search" class="w-full rounded-md bg-primary-foreground bg-opacity-40 text-foreground"/>
                <Button variant="ghost" class="bg-foreground rounded-md ">
                    <SquarePen class="h-5 w-5 text-muted"/>
                </Button>
            </div>

            <div class="flex justify-between align-middle items-center gap-2 w-full">
                <div class="flex flex-row gap-1 items-center">
                    <Checkbox id="selectAll" class="text-sm text-muted-foreground border-"/>
                    <Label for="selectAll" class="sr-only">Select All</Label>
                    <p class="text-sm text-muted-foreground leading-tight">
                        Select All
                    </p>
                </div>

                <div class="flex flex-row items-center">
                    <p class="text-sm text-muted-foreground leading-tight mr-1">
                        Showing 1-50 of 1,400
                    </p>

                    <Button variant="ghost" class="w-5 h-5 p-0">
                        <ChevronLeft class="text-muted-foreground"/>
                    </Button>

                    <Button variant="ghost" class="w-5 h-5 p-0">
                        <ChevronLeft class="text-muted-foreground transform rotate-180"/>
                    </Button>
                </div>
            </div>
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

    <p class="text-muted-foreground text-sm text-center p-4">No more emails to show.</p>

</Resizable.Pane>

<Resizable.Handle withHandle />

<Resizable.Pane class="flex-grow h-full overflow-y-visible">
</Resizable.Pane>

