import * as API from "$api/lib";
import { UrlSafeBase64Encode } from "$api/lib/common";
import {Hash} from "gowl-client-lib";

const username = 'test';
const password = 'test';
const session = new API.Session(await API.Login.Login(username, password), true);
await session.DecryptKeys();

const domainId = 'X0a/eJjvB7vnxLPA3yLnaiiCELkhEW7884e4YG5cTRQ=';
const domainService = await API.DomainService.Decrypt(session, await API.Domain.GetDomain(session, domainId));

const email = {
	domainID: domainService.DomainID,

	inReplyTo: '',


	from:   { displayName: 'Greg', email: 'hello@beta.noise.email' },
	to:     { displayName: '', email: 'x00189661@mytudublin.ie' },
	// cc:     [{ displayName: 'Greg Maniak', email: 'x00189661@mytudublin.ie' }],
	// bcc:    [{ displayName: '', email: 'ap3xdigital@gmail.com' }],
	bcc:    [],
	cc:    [],

	subject: 'Wtf lol it aint sending',
	body: 'Bruhh wtf lol',

}

const sentEmail = await API.Email.SendPlainEmail(session, {
    ...email,
	signature: await domainService.SignData(UrlSafeBase64Encode(await Hash(email.body))),
	nonce: UrlSafeBase64Encode(API.Sym.NewKey())
});
