# _"Contentive"_ 📖 Headless CMS

> Currently, this project has not been completed and we are working hard on it...

<img src="./documents/public/contentive_post.png" alt="Contentive Logo">
<div align="center">
  <h2>A Modern Headless CMS for Content Management</h3>
</div>

## 🌟 Features

- 🚀 **High Performance**: Built with Go and PostgreSQL
- 🔐 **Secure Authentication**: JWT-based token system
- 📝 **Flexible Content Types**: Create custom content schemas
- 🗄️ **Content Versioning**: Track and manage content changes
- 📦 **Media Management**: Support multiple storage providers
- 🔑 **Role-Based Access**: Granular permission control
- 🌐 **RESTful API**: Well-documented API endpoints
- 💻 **Modern Admin UI**: Built with Next.js and Tailwind CSS

## 🛠️ Tech Stack

- **Backend**
  - Go
  - PostgreSQL
  - Fiber (Web Framework)
  - GORM (ORM)
- **Frontend**
  - Next.js
  - Tailwind CSS
  - TypeScript
  - Shadcn UI
- **Documentation**
  - Nextra

## 📦 Installation

### Prerequisites

- Go 1.23.1 or higher
- PostgreSQL 12 or higher
- Node.js 18 or higher
- npm or yarn

### Quick Start

1. Clone the repository:

```bash
git clone https://github.com/inannan423/contentive.git
cd contentive
```

2. Set up environment variables:

```bash
cp .env.example .env
```

3. Install dependencies:

```bash
# Backend dependencies
go mod download

# Admin UI dependencies
cd admin && npm install
cd ../documents && npm install
```

4. Start the development server:

```bash
# Start backend server
go run app.go

# Start admin UI (in another terminal)
cd admin && npm run dev

# Start documentation (in another terminal)
cd documents && npm run dev
```

## 🔧 Configuration

Contentive uses environment variables for configuration. Key configurations include:

- Database settings
- JWT authentication
- Media storage (Local/Aliyun OSS)
- Super admin account

See the [Environment Configuration](https://contentive-docs.vercel.app/getting-started/environments) guide for details.

## 📚 Documentation

- [Getting Started](https://contentive-docs.vercel.app/getting-started)
- [API Reference](https://contentive-docs.vercel.app/api)
- [Admin Guide](https://contentive-docs.vercel.app/admin)

## 🔒 Security

Contentive implements several security measures:

- JWT-based authentication
- Role-based access control
- API token management
- Secure password hashing
- File upload validation

## 🤝 Contributing

We welcome contributions!

## 🌟 Support

If you find Contentive helpful, please consider giving it a star ⭐️
