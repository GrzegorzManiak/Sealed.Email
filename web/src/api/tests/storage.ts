import * as API from "$api/lib";
import Session from "../lib/session/session";
import {HandleRequest} from "$api/lib/api/common";
import {ClientError} from "$api/lib/errors";


const details = ['test2', 'test2'];
const createAccount = false;
const addDomain = false;
const domain_ = 'sealed.email.';

const session = new API.Session(await API.Login.Login(details[0], details[1]), true);
await session.DecryptKeys();

const dummyStorageService = new API.StorageServices.DummyStorageService();
const emailService = new API.EmailStorage(dummyStorageService, session);
await emailService.init();
const emailProvider = new API.EmailProvider(emailService, session);


const domains = await API.Domain.GetDomainList(session, 0, 1);
console.log('Domains:', domains);
const domainId = domains.domains[0].domainID;
const domain = await API.Domain.GetDomain(session, domainId);
const domainService = await API.DomainService.Decrypt(session, domain);
console.log('Domain:', domainService.Domain, domain);

//Email: {
//   domainID: "hZBt3ARb0YDuVDeuwSRJoUat46brUnCXVOkEp8b4Be2QV2tYz695MA3F9YN0d18V",
//   subject: "DHiShl0tUA05NXhQEKbg57ejVnxDv7MtTSwUiQAgF2q0G2AzRrIX5cTpaww2BPbcDyXuQfKokriO3S0t4IgC3YGyXfB5UMwSdBqvTk-lPJ5e",
//   body: "DGwlg2W6RoG8VaFkhDeqZTI6-u7ifrjae2tVObilZ-ojaYWA60MQ7rouMF0wcN8n43WfsQaiSFGZPIO-9_2IuzM",
//   from: {
//     displayName: "DGIK0M3NIqebTSoTNvC3pUZJBmp3LrueFwlWG1XLi5rlr7L3HvhPZJUuJEM",
//     emailHash: "EUJVC9vwBogvilCC4yL7YAcnk82g8xYus4TaUuB0FhE@beta.grzegorz.ie.",
//     publicKey: "AlhgO5lRh7IX57uBm0Y6oM-M5EGpghovAs1yuI6x7Y9c",
//     encryptedEmailKey: "DBAj1d1ERvsJLSoRCDnosJXbz3fC74UTNzThfv2Eyexs0HlmGT2_huwfvhpWTlArAGL2NcuIN4e-fr9-L6qdQQa_NrQrXRr2",
//   },
//   to: {
//     displayName: "DIEMcIby661GamBQeuEr5Fal1bEouoc7AaL_85CkiZPK",
//     emailHash: "J84CaH1v2UtzmPMCn8LvlRjhChD6WW3R8SdEVFLFPSc@sealed.email",
//     publicKey: "Akp_7vMAEkKFwdy6f0wCZilniMlaees-qUBJO61_UJEK",
//     encryptedEmailKey: "DE9QcbTdWhp_K-58Oxtxdtgj4SOJaPdv1F_y_O58ksecIP9z4dZF9VzX-h3RcwBpDazebZ8d4tw-Z6Tb2gy5O9MKQveDIl17",
//   },
//   cc: [],
//   bcc: [],
//   inReplyTo: "",
//   references: [],
// }

console.log(await emailProvider.getEmails(domainService, {domainID: domainService.DomainID, order: 'asc', perPage: 1}));