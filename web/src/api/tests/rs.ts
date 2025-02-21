import * as API from "$api/lib";
import Session from "../lib/session/session";
import {HandleRequest} from "$api/lib/api/common";
import {ClientError} from "$api/lib/errors";


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
	console.log('Domains:', domains);
	const domainId = domains.domains[0].domainID;
	const domain = await API.Domain.GetDomain(session, domainId);
	const domainService = await API.DomainService.Decrypt(session, domain);

	const folder = await API.Folder.FolderCreate(session, domainService, 'test');
	console.log('Folder:', folder);
	console.log('folders', (await API.Folder.FolderList(session, domainId)).total);
	await API.Folder.FolderDelete(session, domainService.DomainID, folder.folderID);
	console.log('folders', (await API.Folder.FolderList(session, domainId)).total);


	//
	// if (domainService.Domain !== domain_) {
	// 	console.error('Domain mismatch');
	// 	return;
	// }
	//
	// const emails = await API.Email.GetEmailList(session, { domainID: domainService.DomainID, order: 'desc', folders: ['test', 'boss'] });
	// console.log('Emails:', emails);
	//
	// if (emails.emails.length > 0) {
	// 	const email = await API.Email.GetEmail(session, domainService.DomainID, emails.emails[0].bucketPath);
	// 	const emailData = await API.Email.GetEmailData(session, domainService.DomainID, email);
	//
	// 	console.log('Email:', emailData);
	// }
	//
	// // const bucketPath = 'C2LMBhwH9jUj1eLyUfZpJiy1g2W8qWjQAbb5ejg4SZUk6GmkWP9y6ZFcaokTc72A:plain:fa0lux/84dv58hxqkja6lnatc0oacbaiupziua4x4qu=@beta.noise.email';
	// // const email = await API.Email.GetEmail(session, domainService.DomainID, bucketPath);
	// // const emailData = await API.Email.GetEmailData(session, domainService.DomainID, email);
	// //
	// // console.log('Email:', emailData);
})();
