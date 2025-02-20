import * as API from '../lib';
import {EncodeToBase64} from "gowl-client-lib";
import Session from "../lib/session/session";
import {HandleRequest} from "$api/lib/api/common";
import {Endpoints} from "$api/lib/constants";
import {ClientError} from "$api/lib/errors";

(async () => {
	const details = ['test', 'test'];
	const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
	await session.DecryptKeys();
	const domain_ = 'beta.grzegorz.ie';
	const send = false;


	const domainId = 'NZ1lQfo8t3HV47t2E99KmEIATbIIQmLTGFRFeqO9CD9fh4hrhsLhbxMlrali9ohr';
	const domain = await API.Domain.GetDomain(session, domainId);
	const domainService = await API.DomainService.Decrypt(session, domain);

	console.log('Plain DOmain', domainService.Domain);
	// const emails = await API.Email.GetEmailList(session, { domainID: domainService.DomainID });
	// console.log('Emails:', emails);

	if (send) {
		const emailKey = API.Sym.NewKey();
		const recipientAKeys = API.Asym.GenerateKeyPair();
		const recipientAInbox = await API.EncryptedInbox.Create(
			'test@sealed.email',
			'Test',
			recipientAKeys.pub,
			emailKey
		);

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
		nonce: EncodeToBase64(API.Sym.NewKey()),
		signature: EncodeToBase64(API.Sym.NewKey()),
	})
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
// 	nonce: EncodeToBase64(API.Sym.NewKey()),
// 	signature: EncodeToBase64(API.Sym.NewKey()),
// })