import { expect, test } from "bun:test";
import * as Asym from "../asymmetric";
import {BytesToBigInt} from "gowl-client-lib";

// test("Key Pair", () => {
// 	const keyPair = Asym.GenerateKeyPair();
// 	expect(keyPair.priv).toBeTruthy();
// 	expect(keyPair.pub).toBeTruthy();
// });
//
// test("Sign Data", async () => {
// 	const keyPair = Asym.GenerateKeyPair();
// 	const data = "Hello, world!";
// 	const signature = await Asym.SignData(data, keyPair.priv);
// 	expect(signature).toBeTruthy();
// });
//
// test("Verify Signature", async () => {
// 	const keyPair = Asym.GenerateKeyPair();
// 	const data = "Hello, world!";
// 	const signature = await Asym.SignData(data, keyPair.priv);
// 	const verified = await Asym.VerifySignature(data, signature, keyPair.pub);
// 	expect(verified).toBeTruthy();
// });
//
// test("Shared Key", async () => {
// 	const keyPair1 = Asym.GenerateKeyPair();
// 	const keyPair2 = Asym.GenerateKeyPair();
// 	const shared1 = await Asym.SharedKey(keyPair1.priv, keyPair2.pub);
// 	const shared2 = await Asym.SharedKey(keyPair2.priv, keyPair1.pub);
// 	expect(BytesToBigInt(shared1)).toBe(BytesToBigInt(shared2));
// });
//
// test("Encrypt / Decrypt", async () => {
// 	const keyPair = Asym.GenerateKeyPair();
// 	const data = "Hello, world!";
//
//
// 	const shared = await Asym.SharedKey(keyPair.priv, keyPair.pub);
//
// 	const encrypted = await Asym.Encrypt(data, shared);
// 	const decrypted = await Asym.Decrypt(encrypted, shared);
//
// 	expect(decrypted).toBe(data);
// });
//
// test("Encrypt / Decrypt, Alice / Bob", async () => {
// 	const alice = Asym.GenerateKeyPair();
// 	const bob = Asym.GenerateKeyPair();
//
// 	// -- alice sends her public key to bob
// 	const sharedAlice = await Asym.SharedKey(alice.priv, bob.pub);
//
// 	// -- bob sends his public key to alice
// 	const sharedBob = await Asym.SharedKey(bob.priv, alice.pub);
//
// 	// -- alice encrypts data with shared key
// 	const data = "Hello, world!";
// 	const encrypted = await Asym.Encrypt(data, sharedAlice);
//
// 	// -- bob decrypts data with shared key
// 	const decrypted = await Asym.Decrypt(encrypted, sharedBob);
//
// 	expect(decrypted).toBe(data);
// });
//
// test("Same Key, Same Data, Different Results", async () => {
// 	const keyPair = Asym.GenerateKeyPair();
// 	const data = "Hello, world!";
// 	const shared = await Asym.SharedKey(keyPair.priv, keyPair.pub);
//
// 	const encrypted1 = await Asym.Encrypt(data, shared);
// 	const encrypted2 = await Asym.Encrypt(data, shared);
//
// 	expect(encrypted1).not.toBe(encrypted2);
//
// 	const decrypted1 = await Asym.Decrypt(encrypted1, shared);
// 	const decrypted2 = await Asym.Decrypt(encrypted2, shared);
//
// 	expect(decrypted1).toBe(decrypted2);
// });

