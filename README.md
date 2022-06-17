# MongoDB Data Masker

This package masks MongoDB data exported to JSON by `mongoexport`.

## Usage

### ES modules

```ts
import { mask } from '@filtered/mongodb-data-masker';
```

### CommonJS

```js
const { mask } = require('@filtered/mongodb-data-masker');
```

## Local development

Install dependencies.

```sh
yarn install
```

### Test

Put test data generated with `mongoexport --jsonArray...` in `export/`.

Tests are using [Jest](https://jestjs.io/) as the testing framework.

```sh
yarn test
```

## Publish

Publish as a private package.

```sh
yarn pub
```
