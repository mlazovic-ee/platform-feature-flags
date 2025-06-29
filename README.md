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
TODO

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

Code
```
TODO
```

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
