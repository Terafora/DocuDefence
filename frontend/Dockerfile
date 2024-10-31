# Use the official Node.js image as a base
FROM node:18-alpine

# Set working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json (if present)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React app for production
RUN npm run build

# Use a lightweight web server to serve the built app
RUN npm install -g serve
CMD ["serve", "-s", "build", "-l", "3000"]

# Expose the port that serve runs on
EXPOSE 3000
