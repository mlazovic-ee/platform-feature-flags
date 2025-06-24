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
TODO

## Usage
TODO


## Structure

```
k8s
├── kustomization.yaml
├── environments.yaml
├── production
│   └── featureflags.yaml
├── staging
│   └── featureflags.yaml
├── preproduction
│   └── featureflags.yaml
├── development
│   └── featureflags.yaml
└── ...
```


## Contributing
All contributions welcome!

### Creating a new feature flag
TODO

### Updating an existing feature flag
TODO
