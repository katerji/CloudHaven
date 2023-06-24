# CloudHaven

CloudHaven is a cloud storage API that integrates with Google Cloud Storage.

## Description

CloudHaven is a cloud storage API that allows users to interact with Google Cloud Storage. It provides functionality for uploading files, deleting files, sharing files, and retrieving analytics data on shared files. The API is designed to be secure, scalable, and user-friendly.

## Capabilities

CloudHaven offers the following capabilities:

1. **Upload File**: Users can upload files to the cloud storage. The API handles file storage and management seamlessly.

2. **Delete File**: Users can delete files from the cloud storage. This ensures that unnecessary files can be removed easily.

3. **Share File**: Users can generate shareable links for files. These links can be shared with others, allowing them to access the file directly.

4. **Get Analytics on Shared File**: CloudHaven tracks analytics data on shared files, such as the open rate. This data can be retrieved to gain insights into file popularity and usage.

## API Endpoints

The CloudHaven API provides the following endpoints:

- `GET /api/landing`: Retrieve the landing page information.
- `POST /api/auth/register`: Register a new user.
- `POST /api/auth/login`: Log in with user credentials.
- `POST /api/auth/refresh`: Refresh the authentication token.
- `GET /api/files`: Get a list of files owned by the authenticated user.
- `POST /api/file`: Upload a file to the cloud storage.
- `DELETE /api/file`: Delete a file from the cloud storage.
- `POST /api/file/share`: Generate a shareable link for a file.
- `GET /api/file/share/info/:file_id`: Get analytics data on a shared file.

## Dependencies

CloudHaven relies on the following dependencies:

- `gin`: HTTP web framework for routing and middleware.
- `katerji/UserAuthKit`: User authentication and authorization package.
- `gcp`: Package for integrating with Google Cloud Storage.
- `redis`: Package for integerating with Redis. Used for shared link analytics
  
## Getting Started

To get started with CloudHaven, follow these steps:

1. Clone the repository: `git clone https://github.com/katerji/cloudhaven.git`
2. Install the dependencies: `go mod download`
3. Configure the necessary environment variables (check .env.example).
4. You can use the docker-compose.yml file in local-env directory for your local db.
5. Build and run the application: `go run main.go`

## Deloyment

CloudHaven is deployed using AWS ECS service.

API: [url](http://18.202.249.247/api)
