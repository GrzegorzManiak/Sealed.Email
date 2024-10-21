# User Key Structure

This is a list of the keys that are generated for a user when they sign up for the **NOISE** email service,
when they add custom domains, create inboxes etc.

Each `Private Key` is of course matched with a `Public Key`, but I am only listing the `Private Keys` here, All
keys are encrypted on the client before being sent to the server, and all keys can be traced back to the user's
`Root Key`.

### User Keys

**Private Data:** What the user sees.

- `Password Derived Key` **Symmetric Key** used to encrypt the **Root Key** and only the **Root Key**, 
this enables us to allow the user to change their password without having to re-encrypt all of their data.
- `Private Key` **Asymmetric Key**, High Entropy Random Key, it is used to receive encrypted messages from our servers (Not Email).
- `Root Key` **Symmetric Key**, High Entropy Random Key, it is used to encrypt all the users other keys.

**Public Data:** What the server sees.

- `Public Key`
- `Encrypted Root Key` 
- `Encrypted Private Key`

### Domain Keys

**Private Data:** What the user sees.

- `Domain Root Key` **Symmetric Key**, High Entropy Random Key, this key is encrypted with the user's **Root Key** and is used to encrypt all the other keys for the domain.

**Public Data:** What the server sees.

- `DKIM Public Key`
- `Encrypted Domain Root Key`

### Inbox Keys

> To alleviate some confusion:
> Email address is structured like this **USER**@**DOMAIN**, the **USER** part is what you type before the `@` symbol, and the **DOMAIN** part is what you type after the `@` symbol.
> e.g. `greg@gmail.com` the **USER** part is `greg` and the **DOMAIN** part is `gmail.com`.

**Private Data:** What the user sees.

- `Inbox Private Key` **Asymmetric Key**, High Entropy Random Key, used to decrypt incoming emails.
- `Inbox Root Key` **Symmetric Key**, High Entropy Random Key, this key is encrypted with the user's **Root Key** and is used to encrypt all the other keys for the inbox.
- `Inbox User` **Not a key**, this is the email address that the inbox is associated with.

**Public Data:** What the server sees.

- `Inbox Public Key`
- `H(User + Domain)` **Hash**, used to identify the inbox for incoming emails.
- `Encrypted User Name`
- `Encrypted Inbox Root Key`

### Message Keys

> By `Message` I mean an email, inbound or outbound.

**Private Data:** What the user sees

- `Key` **Symmetric Key**, High Entropy Random Key, used to encrypt / decrypt the email content.
- `Data` **Not a key**, the actual email content.
- `Headers` **Not a key**, the email headers.
- `Sender` **Not a key**, the email sender (Email Address) in plain text.

**Public Data:** What the server sees

- `Encrypted Message Key` A **Symmetric Key** only when the email has just been received, and the `Recipient` has not yet re-encrypted this key with their own `Inbox Root Key`.
- `Encrypted Message Data`
- `Rencryted Flag` **Boolean**, used to check if the `Recipient` has re-encrypted the `Message Key` with their own `Inbox Root Key`.

### How the keys are used in practice

**Registering a User**

1. The user signs up for the **NOISE** email service.
2. The user picks a `Password` and a `Username`.
3. The user's `Password` is hashed and used to derive a `Password Derived Key`.
4. The user generates a `Root Key` and encrypts it with the `Password Derived Key`.
5. The user generates a `Private Key` and encrypts it with the `Root Key`.
6. The user provides us with `Public Key`, `Encrypted Root Key`, `Encrypted Private Key` alongside with a `Proof` that they have access to the `Private Key`.
7. We verify the `Proof` and store the `Public Key`, `Encrypted Root Key`, `Encrypted Private Key` in our database.

**Logging in**

1. The user provides us with their `Username`.
2. We retrieve the user's keys that are necessary for OWL aPAKE.
3. After the user has been authenticated we provide the user with their `Encrypted Root Key` and `Encrypted Private Key`.
4. The user decrypts their `Root Key` and `Private Key` and uses them to decrypt their other keys.

**Adding a Domain**

1. The user provides us with a `Domain` in plain text, there's no way of anything being encrypted here. 
2. The Server generates a `Private DKIM Key` and a `Public DKIM Key`, and encrypts the `Domain Root Key` with the servers `DKIM Key`.
3. We verify the `Proof` and store the `DKIM Public Key`, `Encrypted DKIM Private Key`, `Encrypted Domain Root Key` in our database.
4. We provide the user with the `DKIM Public Key` and a `txt` record that they have to add to their DNS to verify that they own the domain.
5. We perform a DNS check to make sure that the `DKIM Public Key` is correct, and that the `txt` verification record is correct.

**Creating an Inbox**

> This is a bit more complicated as the user can create an inbox for a domain that they do not own, we will
> provide the users with some domains that they can use, e.g. `@noise.email`, but that will also open a window
> for abuse, so we will have to be careful with this, maybe `@me.noise.email` or something like that.

1. The user provides us with a `Domain` (Plain Text) and a `H(User + Domain)` (Hashed full email), `Inbox Public Key`, `Encrypted User Name` and a `Proof` that they have access to the `Inbox Private Key`.
2. We verify the `Proof` and store the `H(User + Domain)`, `Inbox Public Key`, `Encrypted User Name` in our database.

**Sending an Email**

> I know signing the email AFTER encrypting it is not the best idea.
> However, it makes its adaptation easier, as the receiving
> server won't have to care if the email is encrypted or not, it will just have to check the signature.

1. The user sends an email to another user, call them `Recipient`, who resides on `Server B`.
2. **NOISE** _OR_ the `Sender` performs a request to the `Recipient's` server to get the `Recipient's` `Public Key`.
3. The `Sender` generates a new one time `Message Key` and encrypts the email with the `Message Key`. 
4. The `Sender` encrypts the `Message Key` with the `Recipient's` `Public Key`, and signs the encrypted email content with their `Private Key`
5. The `Sender` encrypts the email with the `Message Key`.
6. **NOISE** sends the email to `Server B` and stores the email in the `Recipient`'s inbox.

**Receiving an Email**

1. **NOISE** receives the email from the `Sender`.
2. **NOISE** verifies the signature of the email, SPF, DMARC, and DKIM.
3. **NOISE** stores the email in the `Recipient's` inbox and sends a notification to the `Recipient`.

**Reading an Email**

> You might be wondering why the `Recipient` has to re-encrypt the `Message Key` with their own `Inbox Root Key`,
> this is because Asymmetric encryption is slow, and we want to minimize the amount of time we spend decrypting,
> Whereas Symmetric encryption is fast (Like really fast, [even in the browser fast](https://github.com/mnasyrov/cryptobench-js?tab=readme-ov-file#google-chrome)).

1. The `Recipient` sends a request to **NOISE** to read the email.
2. **NOISE** retrieves the email from the `Recipient's` inbox, and returns the email to the `Recipient`.
3. The `Recipient` decrypts the `Message Key` with their `Private Key`, and decrypts the email with the `Message Key`.
4. The `Recipient` Re-encrypts the `Message Key` with their own `Inbox Root Key` and stores the email back in their inbox.