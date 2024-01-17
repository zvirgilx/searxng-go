# Nuxt 3 Minimal Starter

Look at the [Nuxt 3 documentation](https://nuxt.com/docs/getting-started/introduction) to learn more.

## Setup Base

1. need node 16.17+
2. node install pnpm@8 -g

## Install

Make sure to install the dependencies:

```bash
# pnpm
cd web
pnpm install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
# 1
go run main.go -l debug api -m debug -a :9999

# 2
cd web
pnpm run dev

```

## Production

Build the application for production:

```bash
# pnpm
# need ng
# You can now deploy .output/public to any static hosting! 
pnpm run generate
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.
