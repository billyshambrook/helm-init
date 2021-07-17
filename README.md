# Helm Init

This is a Helm plugin that initializes a new chart. This functions similar to helm create with the additional features:

- Support starter scaffolds from VCS

## Install

```
helm plugin install https://github.com/billyshambrook/helm-init
```

## Usage

### Use starter scaffold from github repository

```
helm init mychart -p https://github.com/billyshambrook/helm-starter-scaffold
```

By default, it expects the repository to contain a `scaffold/` directory with the starter scaffold. This can be changed using `--directory` flag.
