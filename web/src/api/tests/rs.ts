import * as API from "$api/lib";
import Session from "../lib/session/session";
import {HandleRequest} from "$api/lib/api/common";
import {ClientError} from "$api/lib/errors";
import {sleep} from "bun";
import {UrlSafeBase64Encode} from "$api/lib/common";


const details = ['test2', 'test2'];
const createAccount = false;
const addDomain = false;
const domain_ = 'sealed.email.';

if (createAccount) {
	await (async () => {
		await API.Register.RegisterUser(details[0], details[1]);
		const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
		await session.DecryptKeys();
		const domain = await API.Domain.AddDomain(session, domain_);
		console.log('Domain:', domain);
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

	const domains = await API.Domain.GetDomainList(session, 0, 10);
	const domainId = domains.domains[0].domainID;
	const domain = await API.Domain.GetDomain(session, domainId);
	const domainService = await API.DomainService.Decrypt(session, domain);

	const randomString = Math.random().toString(36).substring(7);
	console.log(`${randomString}`);

	// await API.Email.SendPlainEmail(session, {
	// 	domainID: domainService.DomainID,
	// 	inReplyTo: '',
	// 	from: { displayName: 'Greg', email: `test@${domain_}` },
	// 	to: { displayName: '', email: randomString + '@beta.grzegorz.ie' },
	// 	bcc: [],
	// 	cc: [],
	// 	subject: randomString + ' - Test 1234',
	// 	body: 'When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example: When updating a single column with Update, it needs to have any conditions or it will raise error ErrMissingWhereClause, checkout Block Global Updates for details. When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example:',
	// 	nonce: UrlSafeBase64Encode(API.Sym.NewKey()),
	// 	signature: UrlSafeBase64Encode(API.Sym.NewKey()),
	// })

	const dummyStorageService = new API.StorageServices.DummyStorageService();
	const emailService = new API.EmailStorage(dummyStorageService, session);
	await emailService.init();
	const emailProvider = new API.EmailProvider(emailService, session);
	console.log(await emailProvider.getEmails(domainService, {domainID: domainService.DomainID, order: 'desc', perPage: 10, sent: 'out'}));

	// await sleep(5000);
	//
	// const emails = await API.Email.GetEmailList(session, { domainID: domainService.DomainID, order: 'desc' });
	// if (emails.emails.length > 0) {
	// 	console.log('Emails:', emails.emails[0].bucketPath);
	// 	const email = await API.Email.GetEmail(session, domainService.DomainID, emails.emails[0].bucketPath);
	// 	console.log('Email:', email);
	// 	//
	// 	// const emailData = await API.Email.GetEmailData(session, domainService.DomainID, email);
	// 	// console.log('Email:', emails.emails[0]);
	// 	// console.log('Email:', emailData);
	// }
})();
