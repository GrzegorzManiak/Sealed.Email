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


export {
    SendPlainEmail,
    SendEncryptedEmail,
    
    type PlainEmail,
    type SignedPlainEmail,
    type PostEncryptedEmail,
    type ComputedEncryptedInbox,
    type Inbox
}