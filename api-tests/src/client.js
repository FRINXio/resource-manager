import dotenv from "dotenv";
import fetch from "cross-fetch";
import { ApolloClient, InMemoryCache} from '@apollo/client/core/core.cjs';
import { HttpLink } from '@apollo/client/link/http/http.cjs';
import { setContext } from '@apollo/client/link/context/context.cjs';

let isConfigOk = dotenv.config();

if (isConfigOk.error) {
    //reading config run locally from the IDE
    isConfigOk = dotenv.config({ path: process.cwd() + '/../../.env' });
    if (isConfigOk.error) {
        throw isConfigOk.error
    }
}

const config = process.env;

const authLink = setContext((_, { headers }) => {
    return {
        headers: {
            'x-tenant-id': config.X_TENANT_ID,
            'x-auth-user-roles': config.X_AUTH_USER_ROLES,
            'from': config.FROM,
        }
    }
});

const defaultOptions = {
    watchQuery: {
        fetchPolicy: 'no-cache',
        errorPolicy: 'ignore',
    },
    query: {
        fetchPolicy: 'no-cache',
        errorPolicy: 'all',
    },
}

export const client = new ApolloClient({
    link: authLink.concat(new HttpLink({ uri: config.ENDPOINT_URL, fetch })),
    cache: new InMemoryCache(),
    defaultOptions: defaultOptions,
});
