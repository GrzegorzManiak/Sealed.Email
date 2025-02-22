import {DomainRefID} from "./domain";
import {Session} from "../index";
import {HandleRequest} from "./common";
import {Endpoints} from "../constants";
import {ClientError} from "../errors";

type ComputedEncryptedInbox = {
    displayName: string;
    emailHash: string;
    publicKey: string;
    encryptedEmailKey: string
}

type PostEncryptedEmail = {
    domainID: DomainRefID;
    from: ComputedEncryptedInbox;
    inReplyTo: string;
    references: string[];

    to: ComputedEncryptedInbox;
    cc: ComputedEncryptedInbox[];
    bcc: ComputedEncryptedInbox[];

    subject: string;
    body: string;
}

type Inbox = {
    displayName: string;
    email: string;
}

type PlainEmail = {
    domainID: DomainRefID;
    from: Inbox;
    inReplyTo: string;

    to: Inbox;
    cc: Inbox[];
    bcc: Inbox[];

    subject: string;
    body: string;
}

type SignedPlainEmail = {
    signature: string;
    nonce: string;
} & PlainEmail;

type EmailListFilters = {
    domainID: DomainRefID;
    page?: number;
    perPage?: number;
    order?: 'asc' | 'desc';
    read?: 'all' | 'read' | 'unread';
    spam?: 'all' | 'only' | 'none';
    folders?: string[];
    encrypted?: 'all' | 'original' | 'post';
    sent?: 'all' | 'in' | 'out';
    pinned?: 'all' | 'only' | 'none';
}

type EmailModifiers = {
    read?: 'read' | 'unread' | 'unchanged';
    folder?: string;
    spam?: 'true' | 'false' | 'unchanged';
    pinned?: 'true' | 'false' | 'unchanged';
}

type Email = {
    emailID: string;
    receivedAt: number;
    bucketPath: string;
    read: boolean;
    folder: string;
    to: string;
    spam: boolean;
    sent: boolean;
    accessKey: string;
    expiration: number;
}

type EmailListResponse = {
    emails: Email[];
    total: number;
}

const SendPlainEmail = async (session: Session, email: SignedPlainEmail): Promise<void> => {
    await HandleRequest<void>({
        session,
        body: email,
        endpoint: Endpoints.EMAIL_SEND_PLAIN,
        fallbackError: new ClientError(
            'Failed to send email',
            'Sorry, we were unable to send the email',
            'EMAIL-PLAIN-SEND-FAIL'
        ),
    });
};

const SendEncryptedEmail = async (session: Session, email: PostEncryptedEmail, signature: string): Promise<void> => {
    await HandleRequest<void>({
        session,
        body: { ...email, signature },
        endpoint: Endpoints.EMAIL_SEND_ENCRYPTED,
        fallbackError: new ClientError(
            'Failed to send email',
            'Sorry, we were unable to send the email',
            'EMAIL-ENCRYPTED-SEND-FAIL'
        ),
    });
};

const GetEmailList = async (session: Session, filters: EmailListFilters): Promise<EmailListResponse> => {
    const defaultFilters = {
        page: 0,
        perPage: 10,
        order: 'asc',
        read: 'all',
        spam: 'all',
        folders: [],
        encrypted: 'all',
        sent: 'all',
        pinned: 'all',
        ...filters,
    }

    return await HandleRequest<EmailListResponse>({
        session,
        query: { ...defaultFilters, ...filters },
        endpoint: Endpoints.EMAIL_LIST,
        fallbackError: new ClientError(
            'Failed to get email list',
            'Sorry, we were unable to get the email list',
            'EMAIL-LIST-FAIL'
        ),
    });
};

const GetEmail = async (session: Session, domainID: DomainRefID, bucketPath: string): Promise<Email> => {
    return await HandleRequest<Email>({
        session,
        query: { domainID, bucketPath },
        endpoint: Endpoints.EMAIL_GET,
        fallbackError: new ClientError(
            'Failed to get email',
            'Sorry, we were unable to get the email',
            'EMAIL-GET-FAIL'
        ),
    });
};

const GetEmailData = async (session: Session, domainID: DomainRefID, email: Email): Promise<string> => {
    return await HandleRequest<Response>({
        session,
        query: {
            domainID,
            bucketPath: email.bucketPath,
            accessKey: email.accessKey,
            expiration: email.expiration,
        },
        parse: false,
        endpoint: Endpoints.EMAIL_DATA,
        fallbackError: new ClientError(
            'Failed to get email data',
            'Sorry, we were unable to get the email data',
            'EMAIL-DATA-FAIL'
        ),
    }).then(response => response.text());
}

const ModifyEmails = async (session: Session, domainID: DomainRefID, emails: Email[], modifiers: EmailModifiers): Promise<void> => {
    const send = (emailIds: string[]) => HandleRequest<void>({
        session,
        body: { emailIds, ...modifiers },
        endpoint: Endpoints.EMAIL_MODIFY,
        fallbackError: new ClientError(
            'Failed to modify emails',
            'Sorry, we were unable to modify the emails',
            'EMAIL-MODIFY-FAIL'
        ),
    });

    const chunkSize = 100;
    for (let i = 0; i < emails.length; i += chunkSize)
        await send(emails.slice(i, i + chunkSize).map(email => email.emailID));
}

export {
    SendPlainEmail,
    SendEncryptedEmail,
    GetEmailList,
    GetEmail,
    GetEmailData,
    ModifyEmails,
    
    type PlainEmail,
    type SignedPlainEmail,
    type PostEncryptedEmail,
    type EmailListFilters,
    type ComputedEncryptedInbox,
    type Inbox,
    type Email,
    type EmailListResponse,
    type EmailModifiers,
}