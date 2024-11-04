type UserKeys = {
    RootKey: Uint8Array,
    PublicKey: Uint8Array,
    PrivateKey: Uint8Array,
    EncryptedRootKey: Uint8Array,
    EncryptedPrivateKey: Uint8Array,
    EncryptedContactKey: Uint8Array,
};

type RefID = {
    RID: string,
};

type ReturnedKeys = {
    EncryptedRootKey: Uint8Array,
    EncryptedPrivateKey: Uint8Array,
};

export {
    UserKeys,
    RefID,
    ReturnedKeys
};