import {DomainRefID} from "./domain";
import {Session} from "../index";
import {HandleRequest} from "./common";
import {Endpoints} from "../constants";
import {ClientError} from "../errors";

type EncryptedInbox = {
    displayName: string;
    emailHash: string;
    publicKey: string;
}

type ComputedEncryptedInbox = {
    encryptedEmailKey: string;
    nonce: string;
} & EncryptedInbox;

type EncryptedEmail = {
    domainID: DomainRefID;
    from: EncryptedInbox;
    inReplyTo: string;

    to: EncryptedInbox;
    cc: EncryptedInbox[];
    bcc: EncryptedInbox[];

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

export {
    SendPlainEmail,
    
    type PlainEmail,
    type SignedPlainEmail,
    type EncryptedEmail,
    type EncryptedInbox,
    type ComputedEncryptedInbox,
    type Inbox
}