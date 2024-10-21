# Dependencies & Tools

There are a few dependencies & tools that you need if you want to run this project.

## Dependencies

I won't do an 'install guide' for these as everyone's system is different, I am building this on linux, RedHat 40.

### Backend

- [Golang](https://golang.org/) - The main programming language that we use.
- [Etcd](https://etcd.io/) - A distributed key-value store that we use for configuration.
- [PostgreSQL](https://www.postgresql.org/) - The database that we use.

### Frontend

> This can run perfectly fine on node, it's just if you've ever had to set up a node / typescript project, you know how much of a pain it is.
> Bun.js is a drop-in replacement with first-class TypeScript support.

- [TypeScript](https://www.typescriptlang.org/) - The language that we use for the frontend (With SvelteKit).
- [Bun.js](https://bun.sh/) - Typescript runtime (Replacement for Node.js).

## Tools

These are just a small subset of the tools that I found useful while developing this project.

- [DBBeaver](https://dbeaver.io/) - A database management tool that I use to manage the PostgreSQL database.
- [Insomnia](https://insomnia.rest/) - A REST & gRPC client that I use to test the API's, I use it over Postman as Postman got acquired and they have ruined it.
- [Goland](https://www.jetbrains.com/go/) - The IDE that I use for Golang development.
- [WebStorm](https://www.jetbrains.com/webstorm/) - The IDE that I use for TypeScript / SvelteKit development.
- [Etcd Manager](http://etcdmanager.io/) - A GUI for managing Etcd.

## Testing

For **all** testing, I will only use the built-in testing tools that come with Golang & SvelteKit & Bun.js.

- [Golang Testing](https://pkg.go.dev/testing) - The built-in testing tools that come with Golang.
- [SvelteKit Testing](https://kit.svelte.dev/docs#testing) - The built-in testing tools that come with SvelteKit.
- [Bun.js Testing](https://bun.sh/guides/test) - The built-in testing tools that come with Bun.js.