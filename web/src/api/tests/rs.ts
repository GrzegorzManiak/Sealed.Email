import * as API from "$api/lib";
import Session from "../lib/session/session";
import {HandleRequest} from "$api/lib/api/common";
import {ClientError} from "$api/lib/errors";


(async () => {
	const details = ['test2', 'test2'];
	const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
	await session.DecryptKeys();
	const domain_ = 'sealed.email';

	const domainId = 'bTJeStr43dvJhPeGIOTY8xHecCoYQ86qr599I9g1BR92VVREOfh5M1hPjUDWsIRC';
	const domain = await API.Domain.GetDomain(session, domainId);
	const domainService = await API.DomainService.Decrypt(session, domain);

	// const emails = await API.Email.GetEmailList(session, { domainID: domainService.DomainID });
	// console.log('Emails:', emails);

	const bucketPath = 'C2LMBhwH9jUj1eLyUfZpJiy1g2W8qWjQAbb5ejg4SZUk6GmkWP9y6ZFcaokTc72A:plain:fa0lux/84dv58hxqkja6lnatc0oacbaiupziua4x4qu=@beta.noise.email';
	const email = await API.Email.GetEmail(session, domainService.DomainID, bucketPath);
	const emailData = await API.Email.GetEmailData(session, domainService.DomainID, email);

	console.log('Email:', emailData);
})();
