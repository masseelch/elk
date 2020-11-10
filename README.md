# elk

This package aims to generated fully functional code for the basic crud operations on a defined set of entities. It relies heavily on [facebook/ent](https://github.com/facebook/ent).

`elk` is meant to be used as drop-in replacement of the provided `entc` command of the [facebook/ent](https://github.com/facebook/ent) package.

> :warning: **This is work in progress**: The API may change without further notice!

### Usage examples
Generate the entity graph (similar to `entc generate`):
```shell script
elk generate
```

Generate crud handlers for the entities:
```shell script
elk generate handler
```

Generate dart models (to use in flutter):
```shell script
elk generate flutter
```