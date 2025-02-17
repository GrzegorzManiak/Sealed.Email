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
}


export {
    SendPlainEmail,
    SendEncryptedEmail,
    GetEmailList,
    
    type PlainEmail,
    type SignedPlainEmail,
    type PostEncryptedEmail,
    type EmailListFilters,
    type ComputedEncryptedInbox,
    type Inbox,
    type Email,
    type EmailListResponse
}