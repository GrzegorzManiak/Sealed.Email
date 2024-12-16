import {DomainService, GenericError, InboxService, Session} from "../index";
import {DomainRefID} from "./domain";
import {CurrentCurve, Endpoints} from "../constants";
import {ClientError} from "../errors";
import {NewKey} from "../symetric";
import {GetCurve} from "gowl-client-lib";
import {HandleRequest} from "./common";

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

const AddInbox = async (session: Session, domainService: DomainService, inboxName: string): Promise<InboxService> => {
    const inboxKeys = await domainService.CreateInboxKeys(inboxName);
    await HandleRequest<void>({
        session,
        body: inboxKeys,
        endpoint: Endpoints.INBOX_ADD,
        fallbackError: new ClientError(
            'Failed to add inbox',
            'Sorry, we were unable to add the inbox to your account',
            'INBOX-ADD-FAIL'
        ),
    });
    return InboxService.Decrypt(domainService, inboxKeys);
}


export {
    AddInbox,

    type InboxRefID,
    type InboxKeys
};