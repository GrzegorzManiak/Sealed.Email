import { Login } from "./auth/login";
import { RegisterUser } from "./auth/register";
import Session from "./session/session";
import {EncodeToBase64} from "gowl-client-lib";
import {AddDomain, DeleteDomain, GetDomainList, RefreshDomainVerification} from "./api/domain";

const username = '1234';
console.log(`Logging in as ${username}`);
//
// const newUser = await RegisterUser(username, 'Test');
// if (newUser instanceof Error) throw newUser;

const session = new Session(await Login(username, 'Test'), true);
await session.DecryptKeys();

//
console.log("Session token:", session.Token);
// const randomString = Math.random().toString(36).substring(2);
// const domain = await AddDomain(session, randomString + '.grzegorz.ie');
// console.log(domain);
//
// // await Bun.sleep(1000);
// // await RefreshDomainVerification(session, domain.domainID);
// // await Bun.sleep(1000);
// // await DeleteDomain(session, domain.domainID);
//
const domains = await GetDomainList(session, 0, 10);
console.log(domains);