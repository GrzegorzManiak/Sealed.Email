type Attachment = {
    id: string;
    filename: string;
    size: number;
    type: string;
    data: string;
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

    totalAttachments: number;
    showcaseAttachment?: Attachment;

    chain?: Array<Email>;
}

export type {
    Email,
    Attachment
}