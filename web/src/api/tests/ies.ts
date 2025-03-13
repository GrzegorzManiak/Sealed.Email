import {CurrentCurve, ROOT} from "$api/lib/constants";
import {Decompress, Decrypt} from "$api/lib/symetric";
import {ExtractEphemeralPubKey, Decrypt} from "$api/lib/asymmetric";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "$api/lib/common";
import {BigIntToByteArray, BytesToBigInt, GetCurve, Hash} from "gowl-client-lib";
import { sha256 } from "js-sha256"; // or your preferred SHA-256 implementation

type Output = {
	publicKey: string;
	privateKey: string;
	ephemeralPublicKey: string;
	ephemeralPrivateKey: string;
	ephemeralKeyLength: number;
	sharedX: string;
	sharedKey: string;

	iv: string;
	ciphertext: string;
	encrypted: string;
}

const route = '/dev/encryption';
const method = 'GET';

const encryptionData = await fetch(ROOT + route, { method });
const encryptionDataJson = await encryptionData.json() as Output;
console.log(encryptionDataJson);

// console.log(Decrypt(ExtractEphemeralPubKey(UrlSafeBase64Decode(encryptionDataJson.encrypted)
// ), UrlSafeBase64Decode(encryptionDataJson.ephemeralPrivateKey)));

// const extractedKeys = ExtractEphemeralPubKey(UrlSafeBase64Decode(encryptionDataJson.encrypted));
// const decryptedData = Decrypt(extractedKeys, UrlSafeBase64Decode(encryptionDataJson.ephemeralPrivateKey));
// console.log(await decryptedData);

// async function decryptMessage(encryptionData: Output) {
// 	// Decode Base64 keys
// 	const bobPrivateKey = UrlSafeBase64Decode(encryptionData.ephemeralPrivateKey);
// 	const alicePublicKey = UrlSafeBase64Decode(encryptionData.publicKey);
//
// 	// Convert Alice's Public Key to Curve Point
// 	const alicePubPoint = GetCurve(CurrentCurve).ProjectivePoint.fromHex(alicePublicKey);
//
// 	// Compute Shared Secret (ECDH)
// 	const sharedSecretPoint = alicePubPoint.multiply(BytesToBigInt(bobPrivateKey));
//
// 	// Use only the X coordinate, then hash with SHA-256
// 	const sharedX = sharedSecretPoint.toRawBytes(true).slice(1); // Remove prefix
// 	const sharedKey = sha256(sharedX);
//
// 	// Decode IV & Ciphertext
// 	const iv = UrlSafeBase64Decode(encryptionData.iv);
// 	const ciphertext = UrlSafeBase64Decode(encryptionData.ciphertext);
//
// 	console.log(UrlSafeBase64Decode(encryptionData.sharedX));
// 	console.log(sharedX);
// }

// Usage
// decryptMessage(encryptionDataJson);

// const extractedKeys = ExtractEphemeralPubKey(UrlSafeBase64Decode(encryptionDataJson.encrypted));
// if (extractedKeys.keyLength !== encryptionDataJson.ephemeralKeyLength) throw new Error('Key length mismatch');
// if (UrlSafeBase64Encode(extractedKeys.ephemeralPub) !== encryptionDataJson.ephemeralPublicKey) throw new Error('Public key mismatch');
//
// const decompressedSymData = Decompress(extractedKeys.ciphertext);
// if (UrlSafeBase64Encode(decompressedSymData.iv) !== encryptionDataJson.iv) throw new Error('IV mismatch');
// if (UrlSafeBase64Encode(decompressedSymData.data) !== encryptionDataJson.ciphertext) throw new Error('Ciphertext mismatch');
//
// // load the ephemeral private key
// const ephemeralPrivateKey = BytesToBigInt(UrlSafeBase64Decode(encryptionDataJson.ephemeralPrivateKey));
// const generatedEphemeralPub = MessageCurve.getPublicKey(ephemeralPrivateKey);
// if (UrlSafeBase64Encode(generatedEphemeralPub) !== encryptionDataJson.ephemeralPublicKey) throw new Error('Public key mismatch');
//
// const privateKey = BytesToBigInt(UrlSafeBase64Decode(encryptionDataJson.privateKey));
// const generatedPub = MessageCurve.getPublicKey(privateKey);
// if (UrlSafeBase64Encode(generatedPub) !== encryptionDataJson.publicKey) throw new Error('Public key mismatch');
//
// // You need to swap the parameters in your client code
// // You need to swap the parameters in your client code
// const sharedSecret = MessageCurve.getSharedSecret(
// 	BytesToBigInt(UrlSafeBase64Decode(encryptionDataJson.ephemeralPrivateKey)), // Use the static private key
// 	UrlSafeBase64Decode(encryptionDataJson.publicKey), // Use the ephemeral public key
// 	false
// );
// console.log(UrlSafeBase64Encode(sharedSecret));