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

const colors = {
    normal: "bg-[#020817]",

    hovered: "bg-gray-950",
    selected: "bg-slate-900",

    chainHovered: "#f8f8f8",
    chainSelected: "#e8e8e8"
}

export type {
    Email,
    Attachment
}

export {
    colors
}