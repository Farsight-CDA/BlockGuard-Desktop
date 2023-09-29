# BlockBuard-Desktop

This repository is native Desktopa application using wails.
The frontend implementation hosted through wails can be found at [`BlockGuard-UI`](https://github.com/Farsight-CDA/BlockGuard-UI) 

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
```
to install the packages of the frontend

Now you will need to manually connect the web project to the native APIs implemented in go.
To do that you will have to paste some code into `frontend\src\lib\native-api\native-api.ts`

Its a messy workaround, but hey: it works!

```ts
import {
	MTLSFetch,
	SoftEtherStatus,
	ConnectVPN,
	DisconnectVPN,
	GetConnectionStatus
} from '$lib/wailsjs/go/main/App';
```
Exising file here
[...]
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

now finally we can build the app

```bash
wails build
```

You can preview the production build with `wails dev`.

## Dependencies

this project depends on a modified version of [`akashjs`](https://github.com/akash-network/akashjs).
Our version is [`akashjs`](https://github.com/Farsight-CDA/akashjs).

## Roadmap

Operating System Support
The entire app has been built with maximum possible operating system support in mind. 
All it takes to onboard another operating system is to implement a handful of native platform APIs.
We are planning to support:
- Windows (done)
- Linux (more manual installation process)
- MacOS
- Android
- IOS
