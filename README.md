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
- [How To](#how-to)
  - [Create a new flag](#create-a-new-flag)
  - [Update an existing flag](#update-an-existing-flag)
- [Use evaluation context to target users](#use-evaluation-context-to-target-users)
  - [Examples](#examples)
- [Consumption](#consumption)
  - [Client Side](#client-side)
    - [Mobile](#mobile)
    - [Web](#web)
  - [Server Side](#server-side)
    - [PHP](#php)
    - [Golang](#golang)
    - [NextJS](#nextjs)
      - [1 – OpenFeature ⇄ flagd adapter](#1--openfeature--flagd-adapter)
      - [2 – Declare typed flags](#2--declare-typed-flags)
      - [3 – Use flags in a **Server Component** (App Router)](#3--use-flags-in-a-server-component-app-router)
- [Structure](#structure)

## Links

| Environment   | URL                                                                      | Configuration                                                                                                   | ArgoCD Application                                                                              |
| :------------ | :----------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------- |
| Production    | [https://flags.kick.com](https://flags.kick.com)                         | [Production](https://github.com/KickEngineering/platform-feature-flags/tree/main/k8s/flags/overlays/production) | [platform-feature-flags](https://argo.easygo-kick.com/applications/platform-feature-flags)      |
| Staging       | [https://flags.easygo-drop-kick.com](https://flags.easygo-drop-kick.com) | [Staging](https://github.com/KickEngineering/platform-feature-flags/tree/main/k8s/flags/overlays/staging)       | [platform-feature-flags](https://argo.easygo-drop-kick.com/applications/platform-feature-flags) |
| Preproduction | TBA                                                                      | TBA                                                                                                             | TBA                                                                                             |
| Development   | TBA                                                                      | TBA                                                                                                             | TBA                                                                                             |

## Architecture

![Architecture](./docs/images/diagram.svg)

## How To

### Create a new flag

1. Add a new flag to an existing `FeatureFlag` Manifest located in your desired environment overlay, e.g. `k8s/flags/overlays/staging/kick-web.yaml`:

```yaml
# My flag to do something cool
# disabled: Not cool
# enabled: Cool
my-flag:
  state: ENABLED
  defaultVariant: disabled # What everyone will evaluate if no targetting rules
  variants:
    enabled: true
    disabled: false
```

2. Raise a Pull Request to merge your changes into the `main` branch. This will trigger a deployment of the updated flag definition to your desired environment.

3. Once the Pull Request is merged, ArgoCD will automatically sync your changes to the Kubernetes cluster, and the new flag will be available for consumption.

> [!TIP]
> If you need your flag available quicker than ArgoCD's sync internval, you can manually trigger a sync in the ArgoCD UI.

### Update an existing flag

1. Update your desired flag within a `FeatureFlag` Manifest located in your desired environment overlay, e.g. `k8s/flags/overlays/staging/kick-web.yaml`:

2. Raise a Pull Request to merge your changes into the `main` branch. This will trigger a deployment of the updated flag definition to your desired environment.

3. Once the Pull Request is merged, ArgoCD will automatically sync your changes to the Kubernetes cluster, and the new flag will be available for consumption.

> [!TIP]
> If you need your flag available quicker than ArgoCD's sync internval, you can manually trigger a sync in the ArgoCD UI.

## Use evaluation context to target users

OpenFeature allows you to use an evaluation context to target specific users or groups of users with feature flags. This is useful for A/B testing, gradual rollouts, or targeting specific user segments. OpenFeature achieves this by passing an [evaluation context](https://openfeature.dev/docs/reference/concepts/evaluation-context) where you evaluate the flag, and [targetting rules](https://flagd.dev/reference/flag-definitions/#targeting-rules) on the flag definition.

### Examples

<details>
<summary>Target users that have the `internal` role</summary>

```yaml
web1:
  defaultVariant: disabled
  state: ENABLED
  targeting:
    if:
      - in:
          - internal
          - var: user.roles
      - enabled
      - disabled
  variants:
    disabled: false
    enabled: true
```

</details>

<details>
<summary>Target users that have `a` or `b` role</summary>

```yaml
web1:
  defaultVariant: disabled
  state: ENABLED
  targeting:
    if:
      - or:
          - in:
              - a
              - var: user.roles
          - in:
              - b
              - var: user.roles
      - enabled
      - disabled
  variants:
    disabled: false
    enabled: true
```

</details>

<details>
<summary>Target a single user by email</summary>

```yaml
web2:
  state: ENABLED
  defaultVariant: disabled
  targeting:
    if:
      - in:
          - foo@gmail.com
          - var: user.email
      - enabled
      - disabled
```

</details>

<details>
<summary>Target list of IDs</summary>

```yaml
platform_batch_bans_enabled:
  state: ENABLED
  variants:
    enabled: true
    disabled: false
  defaultVariant: disabled
  targeting:
    if:
      - in:
          - var: ID
          - - "1"
            - "2"
            - "3"
      - enabled
      - disabled
```

</details>

<details>
<summary>Return an object</summary>

```yaml
web3:
  state: ENABLED
  variants:
    object1:
      key: val
    object2:
      key: true
  defaultVariant: object1
```

</details>

<details>
<summary>Return a float</summary>

```yaml
web4:
  state: ENABLED
  variants:
    one: 1.23
    two: 2.34
  defaultVariant: one
```

</details>

<details>
<summary>Email ending in `@easygo.io`</summary>

```yaml
web5:
  state: ENABLED
  defaultVariant: "false"
  variants:
    "true": true
    "false": false
  targeting:
    if:
      - ends_with:
          - var: email
          - "@easygo.io"
      - "true"
      - "false"
```

</details>

<details>
<summary>Fractional Rollout</summary>

```yaml
web6:
  state: ENABLED
  defaultVariant: grey
  variants:
    red: "#b91c1c"
    blue: "#0284c7"
    green: "#16a34a"
    grey: "#4b5563"
  targeting:
    fractional:
      - - red
        - 25
      - - blue
        - 25
      - - green
        - 25
      - - grey
        - 25
```

</details>

<details>
<summary>Progressive Rollout w/ other condition</summary>

This will only work for authenticated users!

```yaml
web7:
  state: ENABLED
  defaultVariant: disabled
  targeting:
    if:
      - in:
          - internal
          - var: user.roles
      - enabled
      - fractional:
          - - enabled
            - 50
          - - disabled
            - 50
  variants:
    enabled: true
    disabled: false
```

</details>

More Examples:

- https://flagd.dev/playground/
- https://yzx396.github.io/flagd-ui/

## Consumption

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
import { OpenFeature } from "@openfeature/react-sdk";
import { FlagdWebProvider } from "@openfeature/flagd-web-provider";

await OpenFeature.setProviderAndWait(
  new FlagdWebProvider({
    url: "https://flags.kick.com",
    tls: true,
  })
);

const client = OpenFeature.getClient();

const webTestEnabled = client.getBooleanValue("web-test", false);
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

Minimal PHP usage example:

```php
use App\Services\FeatureFlag\FeatureFlagServiceInterface;
use Illuminate\Support\Facades\Auth;

class ExampleController extends Controller
{
  private $featureFlags;

  public function __construct(FeatureFlagServiceInterface $featureFlags)
  {
    $this->featureFlags = $featureFlags;
  }

  public function index()
  {
    // Boolean flag with default value
    $isNewFeature = $this->featureFlags->boolean('new_feature_enabled', false);

    // String flag with default value
    $theme = $this->featureFlags->string('ui_theme', 'light');

    // Number flag with default value
    $maxItems = $this->featureFlags->number('max_items_per_page', 10);

    // With evaluation context (targeting)
    $user = Auth::user();
    $showBetaFeatures = $this->featureFlags->boolean('beta_features', false, [
      'targetingKey' => $user->id
    ]);
    $chatProvider = $this->featureFlags->string('chat_provider', 'default', [
      'user_role' => $user->role,
      'subscription_tier' => $user->subscription_tier,
      'region' => $user->region
    ]);

    return null;
  }
}
```

#### Golang

Kubernetes

```yaml
# On Pod manifest
metadata:
  annotations:
    openfeature.dev/enabled: "true"
    openfeature.dev/inprocessconfiguration: platform-open-feature/backend-inprocess-config
```

or

```yaml
# In Helm chart values
featureflagging:
  enabled: true
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

Flags in Next.js are evaluated _server-side_ in line with [Flags SDK's server-side guidance](https://flags-sdk.dev/principles/server-side-vs-client-side). This avoids layout shift, reduces latency in hydrating components, and prevents leaking flag names to the browser.

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
