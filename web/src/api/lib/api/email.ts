import {DomainRefID} from "./domain";
import {Session} from "../index";
import {HandleRequest} from "./common";
import {Endpoints} from "../constants";
import {ClientError} from "../errors";

type Inbox = {
    displayName: string;
    email: string;
};

type PlainEmail = {
    domainID: DomainRefID;
    from: Inbox;
    inReplyTo: string;

    to: Inbox;
    cc: Inbox[];
    bcc: Inbox[];

    subject: string;
    body: string;

    signature: string;
    nonce: string;
}

const SendPlainEmail = async (session: Session, email: PlainEmail): Promise<void> => {
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
}