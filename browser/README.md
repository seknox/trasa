# TRASA browser extension

TRASA browser extension is used to register user device and manage authenticated session for TRASA protected http service.

For full project reference, refer to [github repo](https://github.com/seknox/trasa) or [website](https://www.trasa.io)

## Project structure

- _/src_ : This directory contains popup ui and login/registration, sync functions for extension. Files inside this directory are of react codes and is managed by webpack.
-
- _/extension_ : This directory contains extension file's (background and content scripts).

## Building from source

1. `npm install` in root directory of this project
2. `npm run build` to transpile react code
3. In new terminal window (root directory of this project), enter `npm run firefox`. This will invoke `web-ext` command to start firefox browser with the addon.
