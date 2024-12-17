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

if (domainSweep) await API.Domain.AddDomain(session, randomString + 'test.grzegorz.ie').then(async(domain) => {
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
    // console.log('Deleting domain');
    // await API.Domain.DeleteDomain(session, domain.domainID);

    // -- List domains
    domains = await API.Domain.GetDomainList(session, 0, 10);
    domains.domains.forEach(domain => console.log('-', domain.domain));
});

if (inboxSweep) {
    const domains = await API.Domain.GetDomainList(session, 1, 10);
    const lastDomain = domains.domains[domains.domains.length - 1];
    console.log('PID', lastDomain.domainID)

    // -- Build Domain service
    const domainService = await API.DomainService.Decrypt(session, await API.Domain.GetDomain(session, lastDomain.domainID));

    // -- Create inbox
    const inbox = await API.Inbox.AddInbox(session, domainService, randomString);
    console.log("New inbox:", inbox.InboxName);

    // -- List inboxes
    const inboxes = await API.Inbox.ListInboxes(session, lastDomain.domainID, 0, 10);
    inboxes.inboxes.forEach(inbox => console.log('-', inbox.emailHash));

    // -- Get first inbox
    const lastInboxID = inboxes.inboxes[inboxes.inboxes.length - 1].inboxID;
    const anInbox = await API.Inbox.GetInbox(session, domainService, lastInboxID);
    console.log("Got inbox:", anInbox.InboxName);
}