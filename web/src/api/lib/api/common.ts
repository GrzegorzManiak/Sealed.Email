import Session from "../session/session";
import { ClientError } from "../errors";
import { GenericError } from "../index";

type RequiredOptions = {
    session: Session;
    parse?: boolean;
    endpoint: [string, string];
    fallbackError: ClientError;
};

type BodyOptions = {
    body?: Record<string, unknown>;
};

type HeaderOptions = {
    headers?: Headers;
};

type QueryOptions = {
    query?: Record<string, string | number | boolean | undefined | null | Array<string | number | boolean | undefined | null>>;
};

type Options = RequiredOptions & BodyOptions & HeaderOptions & QueryOptions;

function StringifyArray(key: string, arr: Array<string | number | boolean | undefined | null>): string {
    return arr.map((value) => `${key}=${value}`).join('&');
}

async function HandleRequest<T>(options: Options): Promise<T> {
    const { session, endpoint, fallbackError, query, body, headers: customHeaders } = options;
    if (options.parse === null || options.parse === undefined) options.parse = true;

    const headers = customHeaders ?? new Headers();
    if (session.IsTokenAuthenticated) headers.set("cookie", session.CookieToken);

    const parsedQuery = new URLSearchParams();
    const arrays: Array<string> = [];
    for (const [key, value] of Object.entries(query ?? {})) {
        if (Array.isArray(value)) arrays.push(StringifyArray(key, value));
        else parsedQuery.append(key, (value ?? '').toString());
    }

    let endpointWithQuery = query ? `${endpoint[0]}?${parsedQuery}` : endpoint[0];
    if (arrays.length > 0) endpointWithQuery += parsedQuery.size > 0 ? '&' : '?' + `${arrays.join('&')}`;

    const response = await fetch(endpointWithQuery, {
        method: endpoint[1],
        body: body ? JSON.stringify(body) : undefined,
        headers,
    });

    if (!response.ok) throw GenericError.fromServerString(
        await response.text(),
        fallbackError
    );

    if (options.parse) return await response.json();
    return response as unknown as T;
}

export {
    HandleRequest,

    type RequiredOptions,
    type BodyOptions,
    type HeaderOptions,
    type Options
}