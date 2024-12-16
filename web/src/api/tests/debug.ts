import * as API from '../lib';

const username = '1234';
const password = 'Test';
const domain = 'hello';
const randomString = Math.random().toString(36).substring(2);

const register = false;
const domainSweep = false;
const inboxSweep = true;

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

if (domainSweep) API.Domain.AddDomain(session, randomString + 'test.grzegorz.ie').then(async(domain) => {
    // -- List domains
    let domains = await API.Domain.GetDomainList(session, 0, 10);
    domains.domains.forEach(domain => console.log('-', domain.domain));

    // -- Get domain
    console.log('Getting domain');
    const domainFull = await API.Domain.GetDomain(session, domain.domainID);
    console.log(domainFull);

    // -- Refresh domain verification
    console.log('Refreshing domain verification');
    await API.Domain.RefreshDomainVerification(session, domain.domainID);

    // -- Delete domain
    console.log('Deleting domain');
    await API.Domain.DeleteDomain(session, domain.domainID);

    // -- List domains
    domains = await API.Domain.GetDomainList(session, 0, 10);
    domains.domains.forEach(domain => console.log('-', domain.domain));
});

if (inboxSweep) {
    let domains = await API.Domain.GetDomainList(session, 0, 10);
    const firstDomain = domains.domains[0];

    console.log('PID', firstDomain.domainID)

    // -- Build Domain service
    const domainService = await API.DomainService.Decrypt(
        session,
        await API.Domain.GetDomain(session, firstDomain.domainID)
    );

    // -- Create inbox
    const userInbox = "hello" + randomString;
    const inbox = await API.Inbox.AddInbox(session, domainService, userInbox);
    console.log(inbox.InboxName);

    // -- List inboxes
    const inboxes = await API.Inbox.ListInboxes(session, firstDomain.domainID, 0, 10);
    inboxes.inboxes.forEach(inbox => console.log('-', inbox.inboxName));
}