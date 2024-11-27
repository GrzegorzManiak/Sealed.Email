import { parseDomain, ParseResultType } from "parse-domain";
import Session from "../session/session";
import {CurrentCurve, Endpoints} from "../constants";
import {ClientError, GenericError} from "../errors";
import {GetCurve} from "gowl-client-lib";
import {NewKey} from "../symetric";

type DomainRefID = string;

function CleanDomain(domain: string): string {
    domain = domain.trim();
    domain = domain.toLowerCase();
    if (domain.startsWith("http://")) domain = domain.slice(7);
    if (domain.startsWith("https://")) domain = domain.slice(8);
    const domainParts = parseDomain(domain);
    if (domainParts.type === ParseResultType.Listed) return domainParts.hostname;
    throw new Error("Invalid domain");
}

type AddDomainResponse = {
    domainID: DomainRefID;
    sentVerification: boolean;
    dns: {
        mx: string;
        dkim: string;
        verification: string;
        identity: string;
        spf: string;
    }
}

async function AddDomainRequest(session: Session, domain: string, encRootKey: string): Promise<AddDomainResponse> {
    const headers = new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(Endpoints.DOMAIN_ADD[0], {
        method: Endpoints.DOMAIN_ADD[1],
        body: JSON.stringify({
            domain,
            encRootKey
        }),
        headers
    });

    if (!response.ok) throw GenericError.from_server_string(await response.text(), new ClientError(
        'Failed to add domain',
        'Sorry, we were unable to add the domain to your account',
        'DOMAIN-ADD-FAIL'
    ));

    return await response.json();
}

type RefreshDomainVerificationResponse = {
    sentVerification: boolean;
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

    if (!response.ok) throw GenericError.from_server_string(await response.text(), new ClientError(
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

    if (!response.ok) throw GenericError.from_server_string(await response.text(), new ClientError(
        'Failed to delete domain',
        'Sorry, we were unable to delete the domain from your account',
        'DOMAIN-DELETE-FAIL'
    ));
}

type Domain = {
    domainID: DomainRefID;
    domain: string;
    verified: boolean;
    dateAdded: number;
    catchAll: boolean;
    version: number;
}

type DomainListResponse = {
    domains: Domain[];
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

    if (!response.ok) throw GenericError.from_server_string(await response.text(), new ClientError(
        'Failed to get domain list',
        'Sorry, we were unable to get the domain list from your account',
        'DOMAIN-LIST-FAIL'
    ));

    return await response.json();
}

async function AddDomain(session: Session, domain: string): Promise<AddDomainResponse> {
    domain = CleanDomain(domain);
    const domainKey = NewKey();
    const encRootKey = await session.EncryptKey(domainKey);
    return await AddDomainRequest(session, domain, encRootKey);
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

export {
    RefreshDomainVerification,
    GetDomainList,
    CleanDomain,
    DeleteDomain,
    AddDomain,
    type Domain
}