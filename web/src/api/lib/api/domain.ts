import { parseDomain, ParseResultType } from "parse-domain";
import Session from "../session/session";
import {Endpoints} from "../constants";
import {ClientError} from "../errors";
import {NewKey} from "../symetric";
import {HandleRequest} from "./common";
import * as Asym from "../asymmetric";
import * as Sym from "../symetric";
import { UrlSafeBase64Encode } from "../common";


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

type DomainBrief = {
    domainID: DomainRefID;
    domain: string;
    verified: boolean;
    dateAdded: number;
    catchAll: boolean;
    version: number;
    symmetricRootKey: string;
    publicKey: string;
    encryptedPrivateKey: string;
}

type DomainFull = {
    domainID: DomainRefID;
    dns: DomainDnsData;
} & DomainBrief;

type DomainListResponse = {
    domains: DomainBrief[];
    count: number;
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

const AddDomain = async (session: Session, domain: string): Promise<AddDomainResponse> => {
    domain = CleanDomain(domain);
    const domainKey = NewKey();
    const symmetricRootKey = await session.EncryptKey(domainKey);

    const key = Asym.GenerateKeyPair();
    const encryptedKey = await Sym.Encrypt(UrlSafeBase64Encode(key.priv), domainKey);
    const compressedKey = Sym.Compress(encryptedKey);
    const signature = await Asym.SignData(domain, key.priv);

    return HandleRequest<AddDomainResponse>({
        session,
        body: {
            domain,
            symmetricRootKey,
            publicKey: UrlSafeBase64Encode(key.pub),
            encryptedPrivateKey: UrlSafeBase64Encode(compressedKey),
            proofOfPossession: signature,
        },
        endpoint: Endpoints.DOMAIN_ADD,
        fallbackError: new ClientError(
            'Failed to add domain',
            'Sorry, we were unable to add the domain to your account',
            'DOMAIN-ADD-FAIL'
        ),
    });
}

const GetDomain = async (session: Session, domainID: DomainRefID): Promise<DomainFull> => HandleRequest<DomainFull>({
    session,
    query: { domainID },
    endpoint: Endpoints.DOMAIN_GET,
    fallbackError: new ClientError(
        'Failed to get domain',
        'Sorry, we were unable to get the domain from your account',
        'DOMAIN-GET-FAIL'
    ),
});

const RefreshDomainVerification = async (session: Session, domainID: DomainRefID): Promise<void> => HandleRequest<void>({
    session,
    body: { domainID },
    endpoint: Endpoints.DOMAIN_REFRESH,
    fallbackError: new ClientError(
        'Failed to refresh domain verification',
        'Sorry, we were unable to refresh the domain verification',
        'DOMAIN-REFRESH-FAIL'
    ),
});

const DeleteDomain = async (session: Session, domainID: DomainRefID): Promise<void> => HandleRequest<void>({
    session,
    body: { domainID },
    endpoint: Endpoints.DOMAIN_DELETE,
    fallbackError: new ClientError(
        'Failed to delete domain',
        'Sorry, we were unable to delete the domain from your account',
        'DOMAIN-DELETE-FAIL'
    ),
});

const GetDomainList = async (session: Session, page: number, perPage: number): Promise<DomainListResponse> => HandleRequest<DomainListResponse>({
    session,
    query: { page, perPage  },
    endpoint: Endpoints.DOMAIN_LIST,
    fallbackError: new ClientError(
        'Failed to get domain list',
        'Sorry, we were unable to get the domain list from your account',
        'DOMAIN-LIST-FAIL'
    ),
});



export {
    RefreshDomainVerification,
    GetDomainList,
    CleanDomain,
    DeleteDomain,
    AddDomain,
    GetDomain,

    type DomainBrief,
    type DomainFull,
    type DomainListResponse,
    type AddDomainResponse,
    type DomainRefID,
    type DomainDnsData,
}