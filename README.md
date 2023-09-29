# BlockBuard-Desktop

this reposetory is the Desctop implementation of the project  [`BlockGuard-UI`](https://github.com/Farsight-CDA/BlockGuard-UI) 

## Building

To create a production version of the app:

```bash
git submodule init
git submodule update --remote
```
to update the submodule 

```bash
cd frondend
npm i

to install the packages of the frondend

now you need to past the folowing into
```bash
frontend\src\lib\native-api\native-api.ts
```
sry for the workaround we will try to fix this in the future

```ts
setNativeAPIInitializer(() => {
    if (!Object.hasOwn(window, 'go')) {
        window.location.href = '/';
        throw Error('Missing window.go');
    }

    return Promise.resolve({
        loadFile: (path) => Promise.resolve(localStorage.getItem(path)),
        saveFile: (path, content) =>
            Promise.resolve(localStorage.setItem(path, content)),
        clearFile: (path) => Promise.resolve(localStorage.removeItem(path)),
        mtlsFetch: async <T>(
            method: HttpMethod,
            url: string,
            body: string,
            csr: string,
            privateKey: string
        ) => {
            const res = await MTLSFetch(method, url, body, csr, privateKey);

            if (!res.success) {
                throw Error(`MTLS Fetch failed: StatusCode ${res.statusCode}`);
            }

            try {
                return JSON.parse(res.body) as T;
            } catch (error) {
                return res as T; //T better be string
            }
        },
        getVPNClientStatus: async () =>
            (await SoftEtherStatus()) as VPNClientStatus,
        connectVPN: async (host, username, password) =>
            await ConnectVPN(host, username, password),
        disconnectVPN: async () => await DisconnectVPN(),
        getConnectionStatus: async () =>
            (await GetConnectionStatus()) as VPNConnectionStatus
    } satisfies NativeAPIs);
});
```

now finialy we can build the app run:

```bash
wails build
```
You can preview the production build with `wails dev`.

## Dependencies

this project depends on a modifide version of [`akashjs`](https://github.com/akash-network/akashjs).

ouer version is [`akashjs`](https://github.com/Farsight-CDA/akashjs).

## Roadmap

Operating System Support
The entire app has been built with maximum possible operating system in mind. 
All it takes to onboard another operating system is to implement a handful of native platform APIs.
We are planning to support:
Windows (done)
Linux (more manual installation process)
MacOS
Android
IOS
look related [`Reposetories`](https://github.com/Farsight-CDA/Blockguard) for mor infos.
