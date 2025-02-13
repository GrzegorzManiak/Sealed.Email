import * as API from '../lib';
import {EncodeToBase64} from "gowl-client-lib";

(async () => {
	const details = ['test', 'test'];
	const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
	await session.DecryptKeys();



	const domainId = '5WyfNlZ1VUIUOYNeyvoFHeQe7O+W+bJwy1ThXKZPhOo=';
	const domain = await API.Domain.GetDomain(session, domainId);
	const domainService = await API.DomainService.Decrypt(session, domain);


	const emailKey = API.Sym.NewKey();
	const recipientA = API.Asym.GenerateKeyPair();
	const inboxA = await API.EncryptedInbox.Create(
		'test@sealed.email',
		'Test',
		recipientA.pub,
		emailKey
	);


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
	// console.log(await email.Send(session))

	await API.Email.SendPlainEmail(session, {
		domainID: domainService.DomainID,
		inReplyTo: '',
		from: { displayName: 'Greg', email: 'test@beta.grzegorz.ie' },
		to: { displayName: '', email: 'test@sealed.email' },
		bcc: [],
		cc: [],
		subject: 'Hello world SDFGSDFG SDFG SDFG SDFG SDF',
		body: 'Hello world SDFG SDFGS DFG ',
		nonce: EncodeToBase64(API.Sym.NewKey()),
		signature: EncodeToBase64(API.Sym.NewKey()),
	})
})();


// (async () => {
// 	const details = ['test2', 'test2'];
// 	const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
// 	await session.DecryptKeys();
//
//
// 	const domainId = '49peyLQOSHIylIeftxVzGNFivTtqbwDzK5nN7VjJFkg=';
// 	const domain = await API.Domain.GetDomain(session, domainId);
// 	const domainService = await API.DomainService.Decrypt(session, domain);
//
// 	// console.log(domainService);
// })();