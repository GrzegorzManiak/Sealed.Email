<script lang='ts'>
    import {EmailCard} from '$local/emailCard';
    import {writable} from "svelte/store";
    import {InboxManager} from "$lib/inboxManager";
    const manger = new InboxManager('1', 'Inbox');

    import * as Resizable from '$shadcn/resizable';
    import {CardCompact} from "@/inbox/cardCompact/index.js";
</script>

<div class="max-h-screen h-screen w-screen flex flex-col">
    <Resizable.PaneGroup direction="horizontal" class="h-full">
        <Resizable.Pane class="flex-shrink-0 h-full">
            <!-- Fixed height pane content -->
            <CardCompact data={{
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

                showcaseAttachment: {
                    id: '1',
                    filename: 'noise-email-welcome.pdf',
                    type: 'pdf',
                },
                totalAttachments: 5
            }}/>
        </Resizable.Pane>

        <Resizable.Handle withHandle />

        <Resizable.Pane class="flex-grow h-full overflow-y-visible">
<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="iyggut1"-->
<!--                    email='no-reply@noise.email'-->
<!--                    name='Noise Email'-->
<!--                    subject='Welcome to Noise Email!'-->
<!--                    date='Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Welcome to Noise Email! We are excited to have you on board.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email'-->
<!--                    read={false}-->
<!--                    contactID='1'-->
<!--                    pinned={true}-->
<!--                    showcaseAttachment={{-->
<!--                            id: '1',-->
<!--                            name: 'noise-email-welcome.pdf',-->
<!--                            type: 'pdf',-->
<!--                        }}-->
<!--                    totalAttachments={5}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[-->
<!--                        { showcaseAttachment: { id: '1zxcv', name: 'noise-email-welcome.pdf', type: 'pdf' },totalAttachments: 5, read: false, date: 'Sat Sep 28 2024 21:21:53 GMT+0100 (Irish Standard Time)', sender: false, id: '11', name: 'Noise Email', email: 'no-reply@noise.email', body: 'Oops! We made a mistake, here is the' },-->
<!--                        { read: true, date: 'Sat Sep 28 2024 22:31:22 GMT+0100 (Irish Standard Time)', sender: true, id: '1mnb2', name: 'Greg', email: 'me@grzegorz.ie', body: 'Thank you for your email. I will get back to you as soon as possible.' },-->
<!--                    ]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="4fg7hdj"-->
<!--                    email='events@noise.email'-->
<!--                    name='Noise Events'-->
<!--                    subject='Upcoming Webinar: Privacy in the Digital Age'-->
<!--                    date='Sat Sep 28 2024 15:30:00 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Join us for an insightful webinar on privacy in the digital age. Register now!'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=events.noise.email'-->
<!--                    read={false}-->
<!--                    contactID='4'-->
<!--                    pinned={true}-->
<!--                    totalAttachments={2}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="56fsdfgd"-->
<!--                    email='alerts@noise.email'-->
<!--                    name='Noise Alerts'-->
<!--                    subject='Security Alert: Unusual Login Activity'-->
<!--                    date='Sat Sep 28 2024 14:00:00 GMT+0100 (Irish Standard Time)'-->
<!--                    body='We detected unusual login activity on your account. If this was not you, please reset your password immediately.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=alerts.noise.email'-->
<!--                    read={false}-->
<!--                    contactID='5'-->
<!--                    pinned={false}-->
<!--                    totalAttachments={1}-->
<!--                    favorite={true}-->
<!--                    location='inbox'-->
<!--                    chain={[]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="51sdfgasd5f"-->
<!--                    email='no-reply@noise.email'-->
<!--                    name='Noise Email'-->
<!--                    subject='Welcome to Noise Email!'-->
<!--                    date='Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Welcome to Noise Email! We are excited to have you on board.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email'-->
<!--                    read={true}-->
<!--                    contactID='1'-->
<!--                    pinned={false}-->
<!--                    totalAttachments={5}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[-->
<!--                            { showcaseAttachment: { id: '1adfs4', name: 'noise-email-welcome.pdf', type: 'pdf' },totalAttachments: 5, read: false, date: 'Sat Sep 28 2024 21:21:53 GMT+0100 (Irish Standard Time)', sender: false, id: '11nfg', name: 'Noise Email', email: 'no-reply@noise.email', body: 'Oops! We made a mistake, here is the' },-->
<!--                            { read: true, date: 'Sat Sep 28 2024 22:31:22 GMT+0100 (Irish Standard Time)', sender: true, id: '12afdsg', name: 'Greg', email: 'me@grzegorz.ie', body: 'Thank you for your email. I will get back to you as soon as possible.' },-->
<!--                        ]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="iyut5gg1"-->
<!--                    email='no-reply@noise.email'-->
<!--                    name='Noise Email'-->
<!--                    subject='Welcome to Noise Email!'-->
<!--                    date='Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Welcome to Noise Email! We are excited to have you on board.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email'-->
<!--                    read={false}-->
<!--                    contactID='1'-->
<!--                    pinned={true}-->
<!--                    showcaseAttachment={{-->
<!--                            id: '1',-->
<!--                            name: 'noise-email-welcome.pdf',-->
<!--                            type: 'pdf',-->
<!--                        }}-->
<!--                    totalAttachments={5}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[-->
<!--                            { showcaseAttachment: { id: '1zxcv', name: 'noise-email-welcome.pdf', type: 'pdf' },totalAttachments: 5, read: false, date: 'Sat Sep 28 2024 21:21:53 GMT+0100 (Irish Standard Time)', sender: false, id: '11', name: 'Noise Email', email: 'no-reply@noise.email', body: 'Oops! We made a mistake, here is the' },-->
<!--                            { read: true, date: 'Sat Sep 28 2024 22:31:22 GMT+0100 (Irish Standard Time)', sender: true, id: '1mnb2', name: 'Greg', email: 'me@grzegorz.ie', body: 'Thank you for your email. I will get back to you as soon as possible.' },-->
<!--                        ]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="4fghj1"-->
<!--                    email='events@noise.email'-->
<!--                    name='Noise Events'-->
<!--                    subject='Upcoming Webinar: Privacy in the Digital Age'-->
<!--                    date='Sat Sep 28 2024 15:30:00 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Join us for an insightful webinar on privacy in the digital age. Register now!'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=events.noise.email'-->
<!--                    read={false}-->
<!--                    contactID='4'-->
<!--                    pinned={true}-->
<!--                    totalAttachments={2}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="5fd1"-->
<!--                    email='alerts@noise.email'-->
<!--                    name='Noise Alerts'-->
<!--                    subject='Security Alert: Unusual Login Activity'-->
<!--                    date='Sat Sep 28 2024 14:00:00 GMT+0100 (Irish Standard Time)'-->
<!--                    body='We detected unusual login activity on your account. If this was not you, please reset your password immediately.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=alerts.noise.email'-->
<!--                    read={false}-->
<!--                    contactID='5'-->
<!--                    pinned={false}-->
<!--                    totalAttachments={1}-->
<!--                    favorite={true}-->
<!--                    location='inbox'-->
<!--                    chain={[]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="51asdf1"-->
<!--                    email='no-reply@noise.email'-->
<!--                    name='Noise Email'-->
<!--                    subject='Welcome to Noise Email!'-->
<!--                    date='Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Welcome to Noise Email! We are excited to have you on board.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email'-->
<!--                    read={true}-->
<!--                    contactID='1'-->
<!--                    pinned={false}-->
<!--                    totalAttachments={5}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[-->
<!--                            { showcaseAttachment: { id: '1adfs4', name: 'noise-email-welcome.pdf', type: 'pdf' },totalAttachments: 5, read: false, date: 'Sat Sep 28 2024 21:21:53 GMT+0100 (Irish Standard Time)', sender: false, id: '11nfg', name: 'Noise Email', email: 'no-reply@noise.email', body: 'Oops! We made a mistake, here is the' },-->
<!--                            { read: true, date: 'Sat Sep 28 2024 22:31:22 GMT+0100 (Irish Standard Time)', sender: true, id: '12afdsg', name: 'Greg', email: 'me@grzegorz.ie', body: 'Thank you for your email. I will get back to you as soon as possible.' },-->
<!--                        ]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="iyut1"-->
<!--                    email='no-reply@noise.email'-->
<!--                    name='Noise Email'-->
<!--                    subject='Welcome to Noise Email!'-->
<!--                    date='Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Welcome to Noise Email! We are excited to have you on board.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email'-->
<!--                    read={false}-->
<!--                    contactID='1'-->
<!--                    pinned={true}-->
<!--                    showcaseAttachment={{-->
<!--                            id: '1',-->
<!--                            name: 'noise-email-welcome.pdf',-->
<!--                            type: 'pdf',-->
<!--                        }}-->
<!--                    totalAttachments={5}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[-->
<!--                            { showcaseAttachment: { id: '1zxcv', name: 'noise-email-welcome.pdf', type: 'pdf' },totalAttachments: 5, read: false, date: 'Sat Sep 28 2024 21:21:53 GMT+0100 (Irish Standard Time)', sender: false, id: '11', name: 'Noise Email', email: 'no-reply@noise.email', body: 'Oops! We made a mistake, here is the' },-->
<!--                            { read: true, date: 'Sat Sep 28 2024 22:31:22 GMT+0100 (Irish Standard Time)', sender: true, id: '1mnb2', name: 'Greg', email: 'me@grzegorz.ie', body: 'Thank you for your email. I will get back to you as soon as possible.' },-->
<!--                        ]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="4fghj"-->
<!--                    email='events@noise.email'-->
<!--                    name='Noise Events'-->
<!--                    subject='Upcoming Webinar: Privacy in the Digital Age'-->
<!--                    date='Sat Sep 28 2024 15:30:00 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Join us for an insightful webinar on privacy in the digital age. Register now!'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=events.noise.email'-->
<!--                    read={false}-->
<!--                    contactID='4'-->
<!--                    pinned={true}-->
<!--                    totalAttachments={2}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="5fd"-->
<!--                    email='alerts@noise.email'-->
<!--                    name='Noise Alerts'-->
<!--                    subject='Security Alert: Unusual Login Activity'-->
<!--                    date='Sat Sep 28 2024 14:00:00 GMT+0100 (Irish Standard Time)'-->
<!--                    body='We detected unusual login activity on your account. If this was not you, please reset your password immediately.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=alerts.noise.email'-->
<!--                    read={false}-->
<!--                    contactID='5'-->
<!--                    pinned={false}-->
<!--                    totalAttachments={1}-->
<!--                    favorite={true}-->
<!--                    location='inbox'-->
<!--                    chain={[]}-->
<!--            />-->

<!--            <EmailCard-->
<!--                    inbox={manger}-->
<!--                    id="51asdf"-->
<!--                    email='no-reply@noise.email'-->
<!--                    name='Noise Email'-->
<!--                    subject='Welcome to Noise Email!'-->
<!--                    date='Sat Sep 28 2024 20:27:32 GMT+0100 (Irish Standard Time)'-->
<!--                    body='Welcome to Noise Email! We are excited to have you on board.'-->
<!--                    avatar='https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email'-->
<!--                    read={true}-->
<!--                    contactID='1'-->
<!--                    pinned={false}-->
<!--                    totalAttachments={5}-->
<!--                    favorite={false}-->
<!--                    location='inbox'-->
<!--                    chain={[-->
<!--                            { showcaseAttachment: { id: '1adfs4', name: 'noise-email-welcome.pdf', type: 'pdf' },totalAttachments: 5, read: false, date: 'Sat Sep 28 2024 21:21:53 GMT+0100 (Irish Standard Time)', sender: false, id: '11nfg', name: 'Noise Email', email: 'no-reply@noise.email', body: 'Oops! We made a mistake, here is the' },-->
<!--                            { read: true, date: 'Sat Sep 28 2024 22:31:22 GMT+0100 (Irish Standard Time)', sender: true, id: '12afdsg', name: 'Greg', email: 'me@grzegorz.ie', body: 'Thank you for your email. I will get back to you as soon as possible.' },-->
<!--                        ]}-->
<!--            />-->
        </Resizable.Pane>

<!--        <Resizable.Handle withHandle/>-->

<!--        <Resizable.Pane>Three</Resizable.Pane>-->
    </Resizable.PaneGroup>
</div>