import { Login } from "./auth/login";
import { RegisterUser } from "./auth/register";
import Session from "./session/session";
import {EncodeToBase64} from "gowl-client-lib";
import {AddDomain} from "./api/domain";

const username = 'bob1didbob2';
console.log(`Logging in as ${username}`);
/*
const newUser = await RegisterUser(username, 'Test');
if (newUser instanceof Error) throw newUser;*/

const session = new Session(await Login(username, 'Test'), true);
await session.DecryptKeys();


// console.log("Session token:", session.Token);
await AddDomain(session, 'gfaa.dis-email.grzegorz.ie');