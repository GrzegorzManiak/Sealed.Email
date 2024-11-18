\# Points Of Interest

Just some random tidbits of information that I found interesting / Ideas that I want
to explore further.

- `BIP39` / Similar concepts to verify public keys Out of Band
- Holds system for:
    - Workspaces, e.g., admin blocking read/write access to a user
    - Us, so we can block a user from sending / receiving emails if they go over their quota
    - Same with workspaces if they go over their quota
- Go routines for:
    - Verifying DNS records
- Deploying a _identify.{domain} for 
- Encourage users to change username and password as often as possible as the username is a salt for the user's password, and if they don't change it, it could be a security risk.
- SPLIT `Workspaces` and `Users` into separate services, independent of each other, only shared table should be the `domain` table.
- Add required api roles to `GroupFilter` in `dto` when we get to implementing API's
- USe the session key to encrypt initial login covo, make a encrypt struct tag
- 