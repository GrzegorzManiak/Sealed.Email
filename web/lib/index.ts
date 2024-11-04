import { Login } from "./auth/login";
import { RegisterUser } from "./auth/register";
import Session from "./session/session";
import {EncodeToBase64} from "gowl-client-lib";

const username = 'bob1' + Math.random() * 1000;
console.log(username);

const newUser = await RegisterUser(username, 'Test');
if (newUser instanceof Error) throw newUser;
console.log('REGISTER');

console.log(EncodeToBase64(newUser.RootKey));
console.log(EncodeToBase64(newUser.PrivateKey))
console.log(EncodeToBase64(newUser.ContactsKey))

const session = new Session(await Login(username, 'Test'));
await session.DecryptKeys();
console.log(session.SessionKey)
console.log('DONE');