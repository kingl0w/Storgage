##Storage Webapp

A self storage application I built to share files between my work and personal computer

It utilizes Azure Blob storage, Go, Svelte, and an auth invitation functionality

##Setup

Pretty much ready for deployment you will need to create an .env file in the backend with:

- the server port (8080)
- JWT secret
- Azure storage account name
- Azure storage access key
- Azure storage container
- Admin username
- Admin password

And an .env file in the frontend with the VITE_API_URL (http://localhost:8080/api)

Just clone and run docker compose 
