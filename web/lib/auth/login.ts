import { BigIntToByteArray, Client, Hash, EncodeToBase64, ServerAuthInit, ServerAuthVerify } from "gowl-client-lib";
import { ProcessDetails } from "../common";
import { CurrentCurve, Endpoints, ServerName } from "../constants";
import { ClientError, GenericError } from "../errors";
import { RefID, ReturnedKeys } from "./types";

async function LoginInit(client: Client): Promise<ServerAuthInit & RefID> {
    const payload = await client.AuthInit();
    if (payload instanceof Error) throw payload;

    const response = await fetch(Endpoints.LOGIN_INIT[0], {
        method: Endpoints.LOGIN_INIT[1],
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    });

    if (!response.ok) {
        const errorText = await response.text();
        console.error('Login verification failed:', errorText);
        throw new ClientError(
            'Failed to login', 
            'Sorry, we were unable to log you into your account', 
            'LOGIN-REQ-FAIL-1'
        );
    }

    return await response.json();
};

async function LoginVerify(client: Client, data: ServerAuthInit & RefID): Promise<ServerAuthVerify & ReturnedKeys> {
    const payload = await client.AuthVerify(data);
    if (payload instanceof Error) throw payload;

    const response = await fetch(Endpoints.LOGIN_VERIFY[0], {
        method: Endpoints.LOGIN_VERIFY[1],
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            ...payload,
            RID: data.RID
        })
    });

    if (!response.ok) {
        const errorText = await response.text();
        console.error('Login verification failed:', errorText);
        throw new ClientError(
            'Failed to login', 
            'Sorry, we were unable to log you into your account', 
            'LOGIN-REQ-FAIL-2'
        );
    }

    return await response.json();
};

async function Login(username: string, password: string): Promise<Client | ClientError> {
    try {
        const { usernameHash, passwordHash } = await ProcessDetails(username, password);
        const client = new Client(usernameHash, passwordHash.encoded, ServerName, CurrentCurve);

        const init = await LoginInit(client);
        if (init instanceof Error) throw init;
    
        const verify = await LoginVerify(client, init);
        if (verify instanceof Error) throw verify;

        if (client.ValidateServer(verify) instanceof Error) throw new ClientError(
            'Possible MITM Attack - Server Verification Failed', 
            'Sorry, the server did not respond correctly!', 
            'LOGIN-FAIL-SV'
        );

        return client;
    }   

    catch (UnknownError) {
        return ClientError.from_unknown(UnknownError, new ClientError(
            'Failed to login', 
            'Sorry, we were unable to login to your account',
            'LOGIN-CATCH'
        ));
    }
};

export {
    ProcessDetails,
    LoginInit,
    LoginVerify,
    Login
};