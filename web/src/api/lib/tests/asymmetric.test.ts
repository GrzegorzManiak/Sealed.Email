import { expect, test } from "bun:test";
import * as Asym from "../asymmetric";
import {BytesToBigInt} from "gowl-client-lib";

test("Key Pair", () => {
	const keyPair = Asym.GenerateKeyPair();
	expect(keyPair.priv).toBeTruthy();
	expect(keyPair.pub).toBeTruthy();
});

test("Sign Data", async () => {
	const keyPair = Asym.GenerateKeyPair();
	const data = "Hello, world!";
	const signature = await Asym.SignData(data, keyPair.priv);
	expect(signature).toBeTruthy();
});

test("Verify Signature", async () => {
	const keyPair = Asym.GenerateKeyPair();
	const data = "Hello, world!";
	const signature = await Asym.SignData(data, keyPair.priv);
	const verified = await Asym.VerifySignature(data, signature, keyPair.pub);
	expect(verified).toBeTruthy();
});

test("Compress Signature", () => {
	const signature = {
		nonce: "nonce",
		signature: "signature"
	};
	const compressed = Asym.CompressSignature(signature);
	expect(compressed).toBe("nonce.signature");
});

test("Decompress Signature", () => {
	const compressed = "nonce.signature";
	const signature = Asym.DecompressSignature(compressed);
	expect(signature.nonce).toBe("nonce");
	expect(signature.signature).toBe("signature");
})

test("Shared Key", async () => {
	const keyPair1 = Asym.GenerateKeyPair();
	const keyPair2 = Asym.GenerateKeyPair();
	const shared1 = await Asym.SharedKey(keyPair1.priv, keyPair2.pub);
	const shared2 = await Asym.SharedKey(keyPair2.priv, keyPair1.pub);
	expect(BytesToBigInt(shared1)).toBe(BytesToBigInt(shared2));
});