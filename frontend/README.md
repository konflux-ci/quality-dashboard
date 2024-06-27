# Quality Dashboard Frontend

Quality Dashboard Frontend is based on Patternfly Seed (an open source build scaffolding utility for web apps). 

<img width="1058" alt="Out of box dashboard view of patternfly seed" src="https://raw.githubusercontent.com/konflux-ci/quality-dashboard/main/.github/images/frontend-screenshot.png">

## Quick-start

Prerequisites:
* Node.js 16+ & NPM 9+ (you can use [nvm](https://github.com/nvm-sh/nvm))
* Yarn (`npm install -g yarn`)

In case you are running on a macOS, it is recommended to have Python 3.11+ installed (you can use [Homebrew](https://brew.sh) and pyenv).


```bash
git clone git@github.com:konflux-ci/quality-dashboard.git
cd quality-dashboard/frontend
yarn install && yarn run start:dev
```

## Development scripts
```sh
# Install development/build dependencies
yarn install

# Start the development server
yarn run start:dev

# Run a production build (outputs to "dist" dir)
yarn run build

# Run the test suite
yarn run test

# Run the test suite with coverage
yarn run test:coverage

# Run the linter
yarn run lint

# Run the code formatter
yarn run format

# Launch a tool to inspect the bundle size
yarn run bundle-profile:analyze

# Start the express server (run a production build first)
yarn run start
```

## Configurations
* [TypeScript Config](./tsconfig.json)
* [Webpack Config](./webpack.common.js)
* [Jest Config](./jest.config.js)
* [Editor Config](./.editorconfig)

## Multi environment configuration
This project uses [dotenv-webpack](https://www.npmjs.com/package/dotenv-webpack) for exposing environment variables to your code. Either export them at the system level like `export MY_ENV_VAR=http://dev.myendpoint.com && npm run start:dev` or simply drop a `.env` file in the root that contains your key-value pairs like below:

```sh
ENV_1=http://1.myendpoint.com
ENV_2=http://2.myendpoint.com
```

With that in place, you can use the values in your code like `console.log(process.env.ENV_1);`
