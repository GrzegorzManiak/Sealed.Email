import Session from "../session/session";
import { ClientError } from "../errors";
import { GenericError } from "../index";

type RequiredOptions = {
    session: Session;
    endpoint: [string, string];
    fallbackError: ClientError;
};

type BodyOptions = {
    body?: Record<string, unknown>;
};

type HeaderOptions = {
    headers?: Headers;
};

type Options = RequiredOptions & BodyOptions & HeaderOptions;

async function HandleRequest<T>(options: Options): Promise<T> {
    const { session, endpoint, fallbackError, body, headers: customHeaders } = options;

    const headers = customHeaders ?? new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const response = await fetch(endpoint[0], {
        method: endpoint[1],
        body: body ? JSON.stringify(body) : undefined,
        headers,
    });

    if (!response.ok) throw GenericError.fromServerString(
        await response.text(),
        fallbackError
    );

    return await response.json();
}

export {
    HandleRequest,

    type RequiredOptions,
    type BodyOptions,
    type HeaderOptions,
    type Options
}