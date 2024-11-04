type UserKeys = {
    RootKey: Uint8Array,
    PublicKey: Uint8Array,
    PrivateKey: Uint8Array,
    ContactsKey: Uint8Array,
    EncryptedRootKey: Uint8Array,
    EncryptedPrivateKey: Uint8Array,
    EncryptedContactsKey: Uint8Array,
};

type RefID = {
    RID: string,
};

type ReturnedVerifyData = {
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