FROM oven/bun:latest

WORKDIR /app

# Kopiraj package.json in bun.lockb
COPY package.json bun.lockb ./

# Namesti odvisnosti
RUN bun install --frozen-lockfile

# Kopiraj izvorno kodo
COPY . .

# Izpostavi port (opcijsko, če bi v prihodnosti dodali HTTP endpoint)
EXPOSE 3000

# Zaženi aplikacijo
CMD ["bun", "run", "src/index.ts"] 
