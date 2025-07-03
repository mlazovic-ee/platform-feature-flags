<br />
<div align="center">
  <a href="https://github.com/KickEngineering/platform-feature-flags">
    <img src="docs/images/openfeature.svg" alt="Logo" height="150">
  </a>

  <h3 align="center">platform-feature-flags</h3>

  <p align="center">
    Feature Flag management!
    <br />
  </p>
</div>

## About
Feature Flags are cool.


- [About](#about)
- [Links](#links)
- [Architecture](#architecture)
- [Usage](#usage)
  - [Client Side](#client-side)
    - [Mobile](#mobile)
    - [Web](#web)
  - [Server Side](#server-side)
    - [PHP](#php)
    - [Golang](#golang)
    - [NextJS](#nextjs)
- [Structure](#structure)

## Links

| Environment   | URL                                                                      | Configuration                                                                                                   |
| :------------ | :----------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------- |
| Production    | [https://flags.kick.com](https://flags.kick.com)                         | [Production](https://github.com/KickEngineering/platform-feature-flags/tree/main/k8s/flags/overlays/production) |
| Staging       | [https://flags.easygo-drop-kick.com](https://flags.easygo-drop-kick.com) | [Staging](https://github.com/KickEngineering/platform-feature-flags/tree/main/k8s/flags/overlays/staging)       |
| Preproduction | TBA                                                                      | TBA                                                                                                             |
| Development   | TBA                                                                      | TBA                                                                                                             |


## Architecture

![Architecture](./docs/images/diagram.svg)

## Usage
Consumption of feature flags is done through the OpenFeature specification, which allows for a consistent way to manage feature flags across different platforms and avoids vendor lock-in.

### Client Side
For consumers that run in an insecure environment, such as a web browser or mobile app, feature flag evaluation has to be done on a backend that is controlled by us. This is due to flag definitions having the potential to contain sensitive information that should not be exposed to end users. Only one options is available for client-side, which is remote evaluation using the client-side deployment of FlagD that is accessible via a publically exposed endpoint.

#### Mobile
TODO

#### Web

The browser starts with the flag values that were rendered by the Next.js server (see the **NextJS** section below). Immediately after hydration it bootstraps OpenFeature and subscribes to **streaming updates** from FlagD. When a flag definition changes, FlagD sends a [[Server-Sent Event (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)] to the browser and the SDK updates its cached values.

**Architecture Note:** We use a single source of flags (`kick-web`) that is accessible both server-side (via in-process evaluation) and client-side (via the public endpoint). This approach:
- ✅ Enables local development and Vercel preview environments to work without gRPC connections
- ✅ Simplifies flag management with a single source of truth
- ⚠️ Makes all flag keys publicly visible

**High-level flow**

```text
Next.js server ── initial flag values (SSR) ─▶ Browser
Browser ── streaming updates ─▶ FlagD (`kick-web`)
```

**Browser bootstrap (React / OpenFeature)**

```ts
import { OpenFeature } from '@openfeature/react-sdk';
import { FlagdWebProvider } from '@openfeature/flagd-web-provider';

await OpenFeature.setProviderAndWait(
  new FlagdWebProvider({
    url: 'https://flags.kick.com',
    tls: true,
  }),
);

const client = OpenFeature.getClient();

const webTestEnabled = client.getBooleanValue('web-test', false);
```

All flags are defined in `kick-web` and are available for both server-side rendering and client-side real-time updates via SSE.

### Server Side
For consumers that do support it, we recommend inprocess evaluation of feature flags, which is the fastest and most efficient way to consume them. For consumers that do not support inprocess evaluation, we recommend remote evaluation using the server-side deployment of FlagD.

#### PHP
Kubernetes
```yaml
metadata:
  annotations:
    openfeature.dev/enabled: "true"
    openfeature.dev/inprocessconfiguration: platform-open-feature/backend-inprocess-config
```

Code
```
TODO
```

#### Golang

Kubernetes
```yaml
metadata:
  annotations:
    openfeature.dev/enabled: "true"
    openfeature.dev/inprocessconfiguration: platform-open-feature/backend-inprocess-config
```

Code
```golang
	flagDProvider, err := flagd.NewProvider(
		flagd.WithInProcessResolver(),
	)
	if err != nil {
		slog.Error("Failed to initialize FlagD provider", "error", err)
		os.Exit(1)
	}

	err = openfeature.SetProviderAndWait(flagDProvider)
	if err != nil {
		slog.Error("Failed to setup OpenFeature", "error", err)
		os.Exit(1)
	}

	client := openfeature.NewClient("template-demo")

	state := config.State{
		Flags: client,
	}
```

#### NextJS
Kubernetes
```yaml
metadata:
  annotations:
    openfeature.dev/enabled: "true"
    openfeature.dev/inprocessconfiguration: platform-open-feature/backend-inprocess-config
```

Flags in Next.js are evaluated *server-side* in line with [Flags SDK's server-side guidance](https://flags-sdk.dev/principles/server-side-vs-client-side). This avoids layout shift, reduces latency in hydrating components, and prevents leaking flag names to the browser.

**Environment-aware evaluation:**
- **Production/Staging (Kubernetes):** Uses in-process evaluation via the OpenFeature operator for optimal performance
- **Local/Vercel Preview:** Falls back to remote evaluation against the public FlagD endpoint via http when in-process is unavailable

##### 1 – OpenFeature ⇄ flagd adapter
```ts
// libs/shared/src/feature-flags/adapter.ts
import { OpenFeature } from "@openfeature/server-sdk";
import { createOpenFeatureAdapter } from "@flags-sdk/openfeature";
import { FlagdProvider } from "@openfeature/flagd-provider";

export const openFeatureAdapter = createOpenFeatureAdapter(async () => {
  const provider = new FlagdProvider(); // in-process resolver
  await OpenFeature.setProviderAndWait(provider);
  return OpenFeature.getClient();
});
```

##### 2 – Declare typed flags
```ts
// libs/shared/src/feature-flags/flags.ts
import { flag } from "flags/next";
import type { EvaluationContext } from "@openfeature/server-sdk";
import { openFeatureAdapter } from "./adapter";

export const webTest = flag<boolean, EvaluationContext>({
  key: "web-test",
  defaultValue: false,
  adapter: openFeatureAdapter.booleanValue(),
});
```

##### 3 – Use flags in a **Server Component** (App Router)
```tsx
// app/page.tsx – or any React Server Component
after using the app router.
import { webTest } from "@/lib/flags";

export default async function Page() {
  const isEnabled = await webTest();

  return <div>{isEnabled ? "Enabled" : "Disabled"}</div>;
}
```

The adapter pattern above follows the setup in [Flags SDK's flagd provider docs](https://flags-sdk.dev/providers/openfeature/flagd). After the initial render the server does **not** push updates; the browser continues listening to FlagD for changes (see Web section above).

## Structure

```
k8s
├── kustomization.yaml
├── environments.yaml
├── flags
│   ├── base
│   │   ├── kustomization.yaml
│   │   └── ...
│   └── overlays
│       ├── production
│       │   ├── kustomization.yaml
|       |   └── ...
│       └── staging
│           ├── kustomization.yaml
|           └── ...
└── ...
```
