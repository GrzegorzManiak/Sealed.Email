import * as API from '../lib';
//
// const username = '1234';
// const password = 'Test';
// const domain = 'hello';
// const randomString = Math.random().toString(36).substring(2);
//
// const register = false;
// const domainSweep = false;
//
// if (register) {
//     console.log(`Registering as ${username}`);
//     const newUser = await API.Register.RegisterUser(username, password);
//     if (newUser instanceof API.GenericError) {
//         console.log(newUser);
//     }
// }
//
// console.log(`Logging in as ${username}`);
// const session = new API.Session(await API.Login.Login(username, password), true);
// await session.DecryptKeys();
// console.log("Session token:", session.Token);
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
//     to:     { displayName: '', email: 'test-e9u8xwz7m@srv1.mail-tester.com' },
//     // cc:     [{ displayName: 'Greg Maniak', email: 'x00189661@mytudublin.ie' }],
//     // bcc:    [{ displayName: '', email: 'ap3xdigital@gmail.com' }],
//     bcc:    [],
//     cc:    [],
//
//     subject: 'Re: Your refund for "Bague ambre et argent" is on its way',
//     body: 'Bruhh wtf lol',
// }
//
//
// const sentEmail = await API.Email.SendPlainEmail(session, {
//     ...email,
//     ...await domainService.SignEmail(email)
// });
//
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