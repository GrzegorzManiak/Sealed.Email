
import * as API from '../lib';
import {EncodeToBase64, Hash} from "gowl-client-lib";

const username = 'test5515' + Math.random().toString(36).substring(7);
const password = 'test5515';


console.log(`Registering as ${username}`);
const newUser = await API.Register.RegisterUser(username, password);
if (newUser instanceof API.GenericError) {
    console.log(newUser);
}

console.log(`Logging in as ${username}`);
const session = new API.Session(await API.Login.Login(username, password), true);
await session.DecryptKeys();
console.log("Session token:", session.Token);

const randomString = Math.random().toString(36).substring(7) + '.com';
const domain = await API.Domain.AddDomain(session, randomString);
console.log('Domain:', domain);

//
// console.log(`Logging in as ${username}`);
// const session = new API.Session(await API.Login.Login(username, password), true);
// await session.DecryptKeys();
// console.log("Session token:", session.Token);
//
// // const domains = await API.Domain.GetDomainList(session, 0, 10);
// // console.log('Domains:', domains);
//
// //
// // await API.Domain.AddDomain(session, 'beta.noise.email')
//
//
//
// const domainId = 'X0a/eJjvB7vnxLPA3yLnaiiCELkhEW7884e4YG5cTRQ=';
// // const refresh = await API.Domain.RefreshDomainVerification(session, domainId);
// // console.log('Refresh:', refresh);
// const domain = await API.Domain.GetDomain(session, domainId);
// const domainService = await API.DomainService.Decrypt(session, domain);
// // console.log('Domain:', domain);
//
// const emailKey = API.Sym.NewKey();
// const recipientA = API.Asym.GenerateKeyPair();
// const inboxA = await API.EncryptedInbox.Create(
//     'test@beta.noise.email',
//     'Test',
// 	recipientA.pub,
// 	emailKey
// );
//
// const bccRecipient = API.Asym.GenerateKeyPair();
// const inboxBcc = await API.EncryptedInbox.Create(
//     'bcc@beta.noise.email',
//     '',
//     bccRecipient.pub,
//     emailKey
// );
//
// const sender = await domainService.GetSender(emailKey, 'Greg', 'Grzegorz Maniak')
// const email = new API.EncryptedEmail({
// 	domain: domainService,
// 	key: emailKey,
// 	from: sender,
// 	to: inboxA,
//     bcc: [inboxBcc],
// 	subject: 'Hello world SDFGSDFG SDFG SDFG SDFG SDF',
// 	body: 'Hello world SDFG SDFGS DFG '
// });
//
//
// console.log('From email:', sender.ComputedEncryptedInbox.emailHash);
// console.log('To email:', inboxA.ComputedEncryptedInbox.emailHash);
// console.log('Bcc email:', inboxBcc.ComputedEncryptedInbox.emailHash);
//
// console.log(await email.Send(session))
//
// //
// // if (domainSweep) await API.Domain.AddDomain(session, randomString + 'test.grzegorz.ie').then(async(domain) => {
// //
// //     //-- List domains
// //     let domains = await API.Domain.GetDomainList(session, 0, 10);
// //     domains.domains.forEach(domain => console.log('-', domain.domain));
// //
// //     // -- Get domain
// //     console.log('Getting domain');
// //     const domainFull = await API.Domain.GetDomain(session, domain.domainID);
// //     console.log(domainFull);
// //
// //     // -- Refresh domain verification
// //     console.log('Refreshing domain verification');
// //     await API.Domain.RefreshDomainVerification(session, domain.domainID);
// //
// //     // // -- Delete domain
// //     // console.log('Deleting domain');
// //     // await API.Domain.DeleteDomain(session, domain.domainID);
// //
// //     // -- List domains
// //     domains = await API.Domain.GetDomainList(session, 0, 10);
// //     domains.domains.forEach(domain => console.log('-', domain));
// // });
//
//
//
// // const sentEmail = await API.Email.SendPlainEmail(session, {
// //     ...email,
// // 	signature: await domainService.SignData(EncodeToBase64(await Hash(email.body))),
// // 	nonce: EncodeToBase64(API.Sym.NewKey())
// // });
//
// // console.log('Sent email:', sentEmail);
