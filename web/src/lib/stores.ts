import { persisted } from 'svelte-persisted-store'
import type {DomainBrief, DomainRefID, DomainFull} from "$api/lib/api/domain";
import * as API from "$api/lib";
import {get} from "svelte/store";
import type {SerializedSession} from "../api/lib/session/session";
import {Login} from "../api/lib/auth/login";

/**
 * User store
 */
type User = {
    isLoggedIn: true;
    session: string;
    username: string;
} | {
    isLoggedIn: false;
    session: null;
    username: null;
};

const user = persisted<User>('user', {
    isLoggedIn: false,
    session: null,
    username: null
});

/**
 * Notification Count store
 */
type NotificationCount = Map<string, Map<string, number>>;

const notificationCount = persisted<NotificationCount>('notificationCount', new Map());

/**
 * Domain store
 */

const domains = persisted<Record<DomainRefID, {
    brief: DomainBrief,
    full: DomainFull,
    service: API.DomainService
}>>('domains', {});

const selectedDomain = persisted<{
    domainID: DomainRefID;
    domainName: string;
} | null>('selectedDomain', null);

async function setAllDomains() {
    if (!user) {
        console.info('No user');
        return
    }

    const fetchedUser: User = get(user);
    if (!fetchedUser || !fetchedUser.session || !fetchedUser.isLoggedIn) {
        console.info('No session');
        return;
    }

    const session = API.Session.Deserialize(fetchedUser.session);
    if (!session || session instanceof Error) {
        console.info('No session');
        return;
    }

    const fetchedDomains = await API.Domain.GetDomainList(session, 0, 15);
    const domainMap: Record<DomainRefID, { brief: DomainBrief, service: API.DomainService }> = {};

    for (const domain of fetchedDomains.domains) {
        const domainFull = await API.Domain.GetDomain(session, domain.domainID);
        const domainService = await API.DomainService.Decrypt(session, domainFull);
        domainMap[domain.domainID] = { brief: domain, service: domainService, full: domainFull };
    }

    if (get(selectedDomain) === null && fetchedDomains.domains.length > 0) selectedDomain.set({
        domainID: fetchedDomains.domains[0].domainID,
        domainName: fetchedDomains.domains[0].domain
    });

    domains.set(domainMap);
}


export { user, notificationCount, setAllDomains, domains, selectedDomain };
export type { User, NotificationCount, DomainBrief, DomainRefID };