
import * as API from '../lib';
import {EncodeToBase64, Hash} from "gowl-client-lib";

const username = 'test';
const password = 'test';


// console.log(`Registering as ${username}`);
// const newUser = await API.Register.RegisterUser(username, password);
// if (newUser instanceof API.GenericError) {
//     console.log(newUser);
// }


// console.log(`Logging in as ${username}`);
// const session = new API.Session(await API.Login.Login(username, password), true);
// await session.DecryptKeys();
// console.log("Session token:", session.Token);

// const domains = await API.Domain.GetDomainList(session, 0, 10);
// console.log('Domains:', domains);


// await API.Domain.AddDomain(session, 'beta.noise.email')


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
// const sender = await domainService.GetSender(emailKey, 'Greg', 'Grzegorz Maniak')
// const email = new API.EncryptedEmail({
// 	domain: domainService,
// 	key: emailKey,
// 	from: sender,
// 	to: inboxA,
// 	subject: 'Hello world SDFGSDFG SDFG SDFG SDFG SDF',
// 	body: 'Hello world SDFG SDFGS DFG '
// });
//
// console.log('Email:', await email.Encrypt());
// console.log(await email.Send(session))

//
// if (domainSweep) await API.Domain.AddDomain(session, randomString + 'test.grzegorz.ie').then(async(domain) => {
//
//     //-- List domains
//     let domains = await API.Domain.GetDomainList(session, 0, 10);
//     domains.domains.forEach(domain => console.log('-', domain.domain));
//
//     // -- Get domain
//     console.log('Getting domain');
//     const domainFull = await API.Domain.GetDomain(session, domain.domainID);
//     console.log(domainFull);
//
//     // -- Refresh domain verification
//     console.log('Refreshing domain verification');
//     await API.Domain.RefreshDomainVerification(session, domain.domainID);
//
//     // // -- Delete domain
//     // console.log('Deleting domain');
//     // await API.Domain.DeleteDomain(session, domain.domainID);
//
//     // -- List domains
//     domains = await API.Domain.GetDomainList(session, 0, 10);
//     domains.domains.forEach(domain => console.log('-', domain));
// });


// const domainId = 'TWXjSnVnc6+HQ/WUZaJA6vl3DdkyvNHuowp3TcevrbM=';
// const domainService = await API.DomainService.Decrypt(session, await API.Domain.GetDomain(session, domainId));
// console.log('Domain service:', domainService);
// const email = {
//     domainID: domainService.DomainID,
//     signature: '',
//     nonce: '',
//     inReplyTo: '',
//
//
//     from:   { displayName: 'Greg your beloved', email: 'hello@beta.noise.email' },
//     to:     { displayName: '', email: 'example@beta.noise.email' },
//     // cc:     [{ displayName: 'Greg Maniak', email: 'x00189661@mytudublin.ie' }],
//     // bcc:    [{ displayName: '', email: 'ap3xdigital@gmail.com' }],
//     bcc:    [],
//     cc:    [],
//
//     subject: 'Re: Your refund for "Bague ambre et argent" is on its way',
//     body: 'Bruhh wtf lol',
//
// }


// const sentEmail = await API.Email.SendPlainEmail(session, {
//     ...email,
// 	signature: await domainService.SignData(EncodeToBase64(await Hash(email.body))),
// 	nonce: EncodeToBase64(API.Sym.NewKey())
// });

// console.log('Sent email:', sentEmail);

const emailKey = API.Sym.NewKey();

const recipientA = API.Asym.GenerateKeyPair();
const recipientB = API.Asym.GenerateKeyPair();
const recipientC = API.Asym.GenerateKeyPair();

const Sender = API.Asym.GenerateKeyPair();

const [PubA, PubB, PubC] = [recipientA.pub, recipientB.pub, recipientC.pub];

const [
    SharedA, SharedB, SharedC
] = [
    await API.Asym.SharedKey(Sender.priv, PubA),
    await API.Asym.SharedKey(Sender.priv, PubB),
    await API.Asym.SharedKey(Sender.priv, PubC),
];

const [
    EncryptedA, EncryptedB, EncryptedC
] = [
    await API.Asym.Encrypt(emailKey, SharedA),
    await API.Asym.Encrypt(emailKey, SharedB),
    await API.Asym.Encrypt(emailKey, SharedC),
];



const [PrivA, PrivB, PrivC] = [recipientA.priv, recipientB.priv, recipientC.priv];

const [
    DecryptedA, DecryptedB, DecryptedC
] = [
    await API.Asym.Decrypt(EncryptedA, await API.Asym.SharedKey(PrivA, Sender.pub)),
    await API.Asym.Decrypt(EncryptedB, await API.Asym.SharedKey(PrivB, Sender.pub)),
    await API.Asym.Decrypt(EncryptedC, await API.Asym.SharedKey(PrivC, Sender.pub)),
];

console.log('Email key:', emailKey);
console.log('Encrypted A:', DecryptedA);
console.log('Encrypted B:', DecryptedB);
console.log('Encrypted C:', DecryptedC);