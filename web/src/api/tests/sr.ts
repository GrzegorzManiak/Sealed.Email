import * as API from '../lib';
import {sleep} from "bun";
import {UrlSafeBase64Encode} from "$api/lib/common";

const details = ['test', 'test'];
const createAccount = false;
const addDomain = false;
const domain_ = 'beta.grzegorz.ie.';

if (createAccount) {
	await (async () => {
		await API.Register.RegisterUser(details[0], details[1]);
		const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
		await session.DecryptKeys();
	})();
}

if (addDomain) {
	await (async () => {
		const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
		await session.DecryptKeys();
		const domain = await API.Domain.AddDomain(session, domain_);
		console.log('Domain:', domain);
	})();
}


(async () => {
	const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
	await session.DecryptKeys();
	const send = false;


	const domains = await API.Domain.GetDomainList(session, 0, 10);
	console.log('Domains:', domains);
	const domainId = domains.domains[0].domainID;
	const domain = await API.Domain.GetDomain(session, domainId);
	const domainService = await API.DomainService.Decrypt(session, domain);

	if (domainService.Domain !== domain_) {
		console.error('Domain mismatch');
		return;
	}




	if (send) {
		const emailKey = API.Sym.NewKey();
		const recipientAKeys = API.Asym.GenerateKeyPair();
		const recipientAInbox = await API.EncryptedInbox.Create(
			'test@test.com',
			'Test',
			recipientAKeys.pub,
			emailKey
		);

		console.log('Recipient A:', recipientAInbox);

		const sender = await domainService.GetSender(emailKey, 'Greg', 'Grzegorz Maniak')
		const email = new API.EncryptedEmail({
			domain: domainService,
			key: emailKey,
			from: sender,
			to: recipientAInbox,
			subject: 'Hello world SDFGSDFG SDFG SDFG SDFG SDF',
			body: 'Hello world SDFG SDFGS DFG '
		});

		console.log(await email.Send(session))
	}

	await API.Email.SendPlainEmail(session, {
		domainID: domainService.DomainID,
		inReplyTo: '',
		from: { displayName: 'Greg', email: `test@${domain_}` },
		to: { displayName: '', email: 'test@sealed.email' },
		bcc: [],
		cc: [],
		subject: 'Hello world SDFGSDFG SDFG SDFG SDFG SDF',
		body: 'Hello world SDFG SDFGS DFG ',
		nonce: UrlSafeBase64Encode(API.Sym.NewKey()),
		signature: UrlSafeBase64Encode(API.Sym.NewKey()),
	})

	await sleep(5000);

	console.log('Plain DOmain', domainService.Domain);
	const emails = await API.Email.GetEmailList(session, { domainID: domainService.DomainID, order: 'desc', });
	if (emails.emails.length > 0) {
		const email = await API.Email.GetEmail(session, domainService.DomainID, emails.emails[0].bucketPath);
		const emailData = await API.Email.GetEmailData(session, domainService.DomainID, email);
		console.log('Email:', emails.emails[0]);
		console.log('Email:', emailData);
	}

})();


// await API.Email.SendPlainEmail(session, {
// 	domainID: domainService.DomainID,
// 	inReplyTo: '',
// 	from: { displayName: 'Greg', email: 'test@beta.grzegorz.ie' },
// 	to: { displayName: '', email: 'test@sealed.email' },
// 	bcc: [],
// 	cc: [],
// 	subject: 'Hello world SDFGSDFG SDFG SDFG SDFG SDF',
// 	body: 'Hello world SDFG SDFGS DFG ',
// 	nonce: UrlSafeBase64Encode(API.Sym.NewKey()),
// 	signature: UrlSafeBase64Encode(API.Sym.NewKey()),
// })