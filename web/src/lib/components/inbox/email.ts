type Attachment = {
    id: string;
    filename: string;
    type: string;
}

type Email = {
    id: string;
    date: string;

    fromEmail: string;
    fromName: string;

    toEmail: string;
    toName: string;

    subject: string;
    body: string;

    read: boolean;
    starred: boolean;
    pinned: boolean;

    avatar?: string;
    totalAttachments: number;
    showcaseAttachment?: Attachment;

    chain?: Array<Email>;
}

const colors = {
    normal: "bg-background",
    chain: "bg-muted bg-opacity-20",
    hovered: "bg-muted bg-opacity-10",
    selected: "bg-zinc-900",
}

enum ChainGroupSelect {
    FULL,
    PARTIAL,
    NONE
}

export type {
    Email,
    Attachment
}

export {
    colors,
    ChainGroupSelect
}