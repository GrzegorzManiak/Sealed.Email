import { parseDomain, ParseResultType } from "parse-domain";
import Session from "../session/session";
import {Endpoints} from "../constants";
import {ClientError, GenericError} from "../errors";
import {NewKey} from "../symetric";
import {EncodeToBase64} from "gowl-client-lib";

//
// -- Types
//

type DomainRefID = string;

type DomainDnsData = {
    mx: string;
    dkim: string;
    verification: string;
    identity: string;
    spf: string;
}

type AddDomainResponse = {
    domainID: DomainRefID;
    sentVerification: boolean;
    dns: DomainDnsData;
}

type RefreshDomainVerificationResponse = {
    sentVerification: boolean;
}

type DomainBrief = {
    domainID: DomainRefID;
    domain: string;
    verified: boolean;
    dateAdded: number;
    catchAll: boolean;
    version: number;
}

type DomainFull = {
    domainID: DomainRefID;
    dns: DomainDnsData;
    symmetricRootKey: string;
} & DomainBrief;

type DomainListResponse = {
    domains: DomainBrief[];
}

//
// -- Helpers
//

function CleanDomain(domain: string): string {
    domain = domain.trim();
    domain = domain.toLowerCase();
    if (domain.startsWith("http://")) domain = domain.slice(7);
    if (domain.startsWith("https://")) domain = domain.slice(8);
    const domainParts = parseDomain(domain);
    if (domainParts.type === ParseResultType.Listed) return domainParts.hostname;
    throw new Error("Invalid domain");
}

//
// -- Requests
//

async function AddDomainRequest(session: Session, domain: string, symmetricRootKey: string): Promise<AddDomainResponse> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.DOMAIN_ADD[0], {
        method: Endpoints.DOMAIN_ADD[1],
        body: JSON.stringify({
            domain,
            symmetricRootKey
        }),
        headers
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to add domain',
        'Sorry, we were unable to add the domain to your account',
        'DOMAIN-ADD-FAIL'
    ));

    return await response.json();
}

async function GetDomainRequest(session: Session, domainID: DomainRefID): Promise<DomainFull> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.DOMAIN_GET[0], {
        method: Endpoints.DOMAIN_GET[1],
        body: JSON.stringify({
            domainID
        }),
        headers
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to get domain',
        'Sorry, we were unable to get the domain from your account',
        'DOMAIN-GET-FAIL'
    ));

    return await response.json();
}

async function RefreshDomainVerificationRequest(session: Session, domainID: DomainRefID): Promise<RefreshDomainVerificationResponse> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.DOMAIN_REFRESH[0], {
        method: Endpoints.DOMAIN_REFRESH[1],
        body: JSON.stringify({
            domainID
        }),
        headers
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to refresh domain verification',
        'Sorry, we were unable to refresh the domain verification',
        'DOMAIN-REFRESH-FAIL'
    ));

    return await response.json();
}

async function DeleteDomainRequest(session: Session, domainID: DomainRefID): Promise<void> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.DOMAIN_DELETE[0], {
        method: Endpoints.DOMAIN_DELETE[1],
        body: JSON.stringify({
            domainID
        }),
        headers
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to delete domain',
        'Sorry, we were unable to delete the domain from your account',
        'DOMAIN-DELETE-FAIL'
    ));
}

async function GetDomains(session: Session, page: number, perPage: number): Promise<DomainListResponse> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.DOMAIN_LIST[0], {
        method: Endpoints.DOMAIN_LIST[1],
        body: JSON.stringify({
            pagination: {
                page,
                perPage
            }
        }),
        headers
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to get domain list',
        'Sorry, we were unable to get the domain list from your account',
        'DOMAIN-LIST-FAIL'
    ));

    return await response.json();
}

//
// -- Exports
//

async function AddDomain(session: Session, domain: string): Promise<AddDomainResponse> {
    domain = CleanDomain(domain);
    const domainKey = NewKey();
    const symmetricRootKey = await session.EncryptKey(domainKey);
    return await AddDomainRequest(session, domain, symmetricRootKey);
}

async function GetDomain(session: Session, domainID: DomainRefID): Promise<DomainFull> {
    return await GetDomainRequest(session, domainID);
}

async function RefreshDomainVerification(session: Session, domainID: DomainRefID): Promise<void> {
    await RefreshDomainVerificationRequest(session, domainID);
}

async function DeleteDomain(session: Session, domainID: DomainRefID): Promise<void> {
    await DeleteDomainRequest(session, domainID);
}

async function GetDomainList(session: Session, page: number, perPage: number): Promise<DomainListResponse> {
    return await GetDomains(session, page, perPage);
}

const Requests = {
    AddDomainRequest,
    RefreshDomainVerificationRequest,
    DeleteDomainRequest,
    GetDomains,
    GetDomainRequest
};

export {
    RefreshDomainVerification,
    GetDomainList,
    CleanDomain,
    DeleteDomain,
    AddDomain,
    GetDomain,

    Requests,

    type DomainBrief,
    type DomainFull,
    type DomainListResponse,
    type AddDomainResponse,
    type DomainRefID,
    type DomainDnsData,
    type RefreshDomainVerificationResponse
}