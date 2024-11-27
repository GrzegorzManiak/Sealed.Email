import { InboxManager } from '$lib/inboxManager';

type Attachment = {
    id: string;
    type: string;
    name: string;
};

type ChainEntry = {
    id: string;
    name: string;
    email: string;
    body: string;
    sender: boolean;
    date: string;
    read: boolean;
    showcaseAttachment?: Attachment;
    totalAttachments?: number;
};

type EmailCardData = {
    inbox: InboxManager;
    id: string;
    email: string;
    name: string;
    subject: string;
    body: string;
    avatar: string;
    date: string;
    read: boolean;
    pinned: boolean;
    favorite: boolean;
    contactID: string | number | null;
    showcaseAttachment?: Attachment;
    totalAttachments?: number;
    location: string;
    chain?: ChainEntry[];
};

export type {
    EmailCardData,
    ChainEntry,
    Attachment
}