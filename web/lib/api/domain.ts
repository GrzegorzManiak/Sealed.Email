import { parseDomain, ParseResultType } from "parse-domain";
import Session from "../session/session";
import {CurrentCurve, Endpoints} from "../constants";
import {ClientError} from "../errors";
import {GetCurve} from "gowl-client-lib";
import {NewKey} from "../symetric";

function CleanDomain(domain: string): string {
    domain = domain.trim();
    domain = domain.toLowerCase();
    if (domain.startsWith("http://")) domain = domain.slice(7);
    if (domain.startsWith("https://")) domain = domain.slice(8);
    const domainParts = parseDomain(domain);
    if (domainParts.type === ParseResultType.Listed) return domainParts.hostname;
    throw new Error("Invalid domain");
}

async function AddDomainRequest(session: Session, domain: string, encRootKey: string): Promise<void> {
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

    if (!response.ok) {
        const errorText = await response.text();
        console.error("Failed to add domain:", errorText);
        throw new ClientError(
            'Failed to add domain',
            'Sorry, we were unable to add the domain to your account',
            'DOMAIN-ADD-FAIL'
        );
    }

    console.log("Domain added successfully", await response.json());
}

async function AddDomain(session: Session, domain: string): Promise<void> {
    domain = CleanDomain(domain);
    const domainKey = NewKey();
    const encRootKey = await session.EncryptKey(domainKey);

    console.log("Adding domain:", encRootKey);
    await AddDomainRequest(session, domain, encRootKey);
}

export {
    CleanDomain,
    AddDomain
}