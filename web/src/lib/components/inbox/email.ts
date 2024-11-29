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

    totalAttachments: number;
    showcaseAttachment?: Attachment;

    chain?: Array<Email>;
}

const colors = {
    normal: "bg-[#020817]",
    chain: "bg-[#020611]",

    hovered: "bg-gray-950",
    selected: "bg-slate-900",
}

export type {
    Email,
    Attachment
}

export {
    colors
}