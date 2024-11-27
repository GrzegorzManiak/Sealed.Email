import {Session} from "$api/lib";
import { persisted } from 'svelte-persisted-store'

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

export { user };
export type { User };