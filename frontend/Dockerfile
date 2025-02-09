# Build Stage
FROM node:18 as builder

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy source files
COPY . .

# Build the Svelte app
RUN npm run build

# Serve with a lightweight web server
FROM nginx:alpine

# Copy built files from the builder
COPY --from=builder /app/build /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf  

# Expose port
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
