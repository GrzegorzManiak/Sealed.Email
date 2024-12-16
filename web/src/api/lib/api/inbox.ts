import {DomainService, GenericError, InboxService, Session} from "../index";
import {DomainRefID} from "./domain";
import {CurrentCurve, Endpoints} from "../constants";
import {ClientError} from "../errors";
import {NewKey} from "../symetric";
import {GetCurve} from "gowl-client-lib";

//
// -- Types
//

type InboxRefID = string;

type InboxKeys = {
    emailHash: string;
    symmetricRootKey: string;
    asymmetricPrivateKey: string;
    asymmetricPublicKey: string;
    encryptedInboxName: string;
}

//
// -- Requests
//

async function AddInboxRequest(session: Session, domainID: DomainRefID, inbox: InboxKeys): Promise<void> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.INBOX_ADD[0], {
        method: Endpoints.INBOX_ADD[1],
        body: JSON.stringify({
            domainID,
            ...inbox
        }),
        headers
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to add inbox',
        'Sorry, we were unable to add the inbox to your account',
        'INBOX-ADD-FAIL'
    ));
}

//
// -- Exports
//

async function AddInbox(session: Session, domainService: DomainService, inboxName: string): Promise<InboxService> {
    const inboxKeys = await domainService.CreateInboxKeys(inboxName);
    await AddInboxRequest(session, domainService.DomainID, inboxKeys);
    return InboxService.Decrypt(domainService, inboxKeys);
}

const Requests = {
    AddInboxRequest
}

export {
    Requests,

    AddInbox,

    type InboxRefID,
    type InboxKeys
};