import { parseDomain, ParseResultType } from "parse-domain";
import Session from "../session/session";

function CleanDomain(domain: string): string {
    domain = domain.trim();
    domain = domain.toLowerCase();
    if (domain.startsWith("http://")) domain = domain.slice(7);
    if (domain.startsWith("https://")) domain = domain.slice(8);
    const domainParts = parseDomain(domain);
    if (domainParts.type === ParseResultType.Listed) return domainParts.hostname;
    throw new Error("Invalid domain");
}

async function AddDomain(session: Session, domain: string): Promise<void> {
    domain = CleanDomain(domain);
    console.log(domain);
}

export {
    CleanDomain,
    AddDomain
}