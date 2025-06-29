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
  - [Web](#web)
    - [Client side](#client-side)
    - [Server side](#server-side)
  - [Mobile](#mobile)
  - [Backend](#backend)
    - [PHP](#php)
    - [Golang](#golang)
- [Structure](#structure)
- [Contributing](#contributing)
  - [Creating a new feature flag](#creating-a-new-feature-flag)
  - [Updating an existing feature flag](#updating-an-existing-feature-flag)

## Links

| Environment   | URL                                                                      |
| :------------ | :----------------------------------------------------------------------- |
| Production    | [https://flags.kick.com](https://flags.kick.com)                         |
| Staging       | [https://flags.easygo-drop-kick.com](https://flags.easygo-drop-kick.com) |
| Preproduction | TBA                                                                      |
| Development   | TBA                                                                      |


## Architecture

![Architecture](./docs/images/diagram.svg)

## Usage
TODO

### Web
TODO

#### Client side
TODO

#### Server side
TODO

### Mobile
TODO

### Backend
TODO

#### PHP
TODO

#### Golang
TODO


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


## Contributing
All contributions welcome!

### Creating a new feature flag
TODO

### Updating an existing feature flag
TODO
