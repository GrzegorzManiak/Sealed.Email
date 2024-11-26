type UserKeys = {
    RootKey: Uint8Array,
    PublicKey: Uint8Array,
    PrivateKey: Uint8Array,
    ContactsKey: Uint8Array,
    EncryptedRootKey: Uint8Array,
    EncryptedPrivateKey: Uint8Array,
    EncryptedContactsKey: Uint8Array,
    IntegrityHash: string,
};

type RefID = {
    RID: string,
};

type ReturnedVerifyData = {
    _headers: Headers,
    integrityHash: string,

    encryptedSymmetricRootKey: string,
    encryptedAsymmetricPrivateKey: string,
    encryptedSymmetricContactsKey: string,

    totalInboundEmails: number,
    totalInboundBytes: number,
    totalOutboundEmails: number,
    totalOutboundBytes: number,
};

export {
    UserKeys,
    RefID,
    ReturnedVerifyData
};