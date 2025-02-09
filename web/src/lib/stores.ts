import {Session} from "$api/lib";
import { persisted } from 'svelte-persisted-store'
import type {DomainBrief, DomainRefID} from "$api/lib/api/domain";

/**
 * User store
 */
type User = {
    isLoggedIn: true;
    session: Session;
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

const domains = persisted<Map<DomainRefID, DomainBrief>>('domains', new Map());

const selectedDomain = persisted<DomainRefID | null>('selectedDomain', null);


export { user, notificationCount };
export type { User, NotificationCount, DomainBrief, DomainRefID };