import { Login } from "./auth/login";
import { RegisterUser } from "./auth/register";

const username = 'usfername' + Math.random() * 1000;
console.log(username);

const newUser = await RegisterUser(username, 'Test');
if (newUser instanceof Error) throw newUser;
console.log('REGISTER');

const loggedIn = await Login(username, 'Test');
if (loggedIn instanceof Error) throw loggedIn;
console.log('LOGGED IN');

console.log('DONE');