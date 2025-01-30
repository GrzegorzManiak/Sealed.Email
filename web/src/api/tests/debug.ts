import * as API from '../lib';

const username = '1234';
const password = 'Test';
const domain = 'hello';
const randomString = Math.random().toString(36).substring(2);

const register = false;
const domainSweep = false;

if (register) {
    console.log(`Registering as ${username}`);
    const newUser = await API.Register.RegisterUser(username, password);
    if (newUser instanceof API.GenericError) {
        console.log(newUser);
    }
}

console.log(`Logging in as ${username}`);
const session = new API.Session(await API.Login.Login(username, password), true);
await session.DecryptKeys();
console.log("Session token:", session.Token);

if (domainSweep) await API.Domain.AddDomain(session, randomString + 'test.grzegorz.ie').then(async(domain) => {

    //-- List domains
    let domains = await API.Domain.GetDomainList(session, 0, 10);
    domains.domains.forEach(domain => console.log('-', domain.domain));

    // -- Get domain
    console.log('Getting domain');
    const domainFull = await API.Domain.GetDomain(session, domain.domainID);
    console.log(domainFull);

    // -- Refresh domain verification
    console.log('Refreshing domain verification');
    await API.Domain.RefreshDomainVerification(session, domain.domainID);

    // // -- Delete domain
    // console.log('Deleting domain');
    // await API.Domain.DeleteDomain(session, domain.domainID);

    // -- List domains
    domains = await API.Domain.GetDomainList(session, 0, 10);
    domains.domains.forEach(domain => console.log('-', domain));
});


const domainId = 'TWXjSnVnc6+HQ/WUZaJA6vl3DdkyvNHuowp3TcevrbM=';
const domainService = await API.DomainService.Decrypt(session, await API.Domain.GetDomain(session, domainId));

const email = {
    domainID: domainService.DomainID,
    signature: '',
    nonce: '',

    from:   { displayName: 'Grzegorz Maniak <ADMIN@beta.noise.email>, ADMIN@beta.noise.email', email: 'hello@beta.noise.email' },
    to:     { displayName: 'Grzegorz Maniak', email: 'gregamaniak@gmail.com' },
    // cc:     [{ displayName: '', email: 'x00189661@mytudublin.ie' }],
    // bcc:    [{ displayName: '', email: 'ap3xdigital@gmail.com' }],
    cc:     [],
    bcc:    [],

    subject: 'I am a dwarf and i am digging a hole!!',
    body: 'Diggy diggy hole!!',
}


const sentEmail = await API.Email.SendPlainEmail(session, {
    ...email,
    ...await domainService.SignEmail(email)
});

console.log('Sent email:', sentEmail);
