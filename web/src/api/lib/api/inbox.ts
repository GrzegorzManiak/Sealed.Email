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

type InboxShort = {
    inboxID: InboxRefID;
    inboxName: string;
    dateAdded: number;
    version: number;
}

type InboxListResponse = {
    inboxes: InboxShort[];
}

//
// -- Requests
//

const AddInbox = async (session: Session, domainService: DomainService, inboxName: string): Promise<InboxService> => {
    const inboxKeys = await domainService.CreateInboxKeys(inboxName);
    await HandleRequest<void>({
        session,
        body: {
            ...inboxKeys,
            domainID: domainService.DomainID,
        },
        endpoint: Endpoints.INBOX_ADD,
        fallbackError: new ClientError(
            'Failed to add inbox',
            'Sorry, we were unable to add the inbox to your account',
            'INBOX-ADD-FAIL'
        ),
    });
    return InboxService.Decrypt(domainService, inboxKeys);
}

const ListInboxes = async (session: Session, domainID: DomainRefID, page: number, perPage: number): Promise<InboxListResponse> => HandleRequest<InboxListResponse>({
    session,
    query: { domainID, page, perPage },
    endpoint: Endpoints.INBOX_LIST,
    fallbackError: new ClientError(
        'Failed to list inboxes',
        'Sorry, we were unable to list the inboxes from your account',
        'INBOX-LIST-FAIL'
    ),
});

export {
    AddInbox,
    ListInboxes,

    type InboxRefID,
    type InboxKeys
};