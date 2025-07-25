# KT Chat Frontend

A modern, real-time chat application built with Nuxt.js 3, featuring AWS Cognito authentication and WebSocket communication.

## 🚀 Features

- **Real-time Chat**: WebSocket-based messaging with room support
- **User Authentication**: AWS Cognito integration with JWT tokens
- **Modern UI**: Beautiful, responsive design with Tailwind CSS
- **File Upload**: Support for images, documents, and videos
- **User Management**: Profile pages and account settings
- **Protected Routes**: Authentication middleware for secure access

## 🛠️ Technology Stack

- **Framework**: Nuxt.js 3
- **UI Library**: Nuxt UI (based on Tailwind CSS)
- **Authentication**: AWS Cognito + Custom JWT
- **Real-time**: WebSocket with custom composable
- **State Management**: Nuxt's built-in state management
- **TypeScript**: Full TypeScript support

## 📁 Project Structure

```
frontend/
├── app.vue                 # Main app layout
├── nuxt.config.ts          # Nuxt configuration
├── middleware/             # Route middleware
│   └── auth.ts            # Authentication middleware
├── pages/                  # Application pages (file-based routing)
│   ├── index/             # Home page (chat interface)
│   │   └── index.vue      # Main chat page
│   ├── login/             # Authentication pages
│   │   └── index.vue      # Login/register page
│   └── profile/           # User profile pages
│       └── index.vue      # User profile page
├── composables/            # Reusable composables
│   ├── useAuth.ts         # Authentication logic
│   └── useWebSocket.ts    # WebSocket management
├── layouts/               # Page layouts
├── plugins/               # Nuxt plugins
├── assets/                # Static assets
└── public/                # Public files
```

## 🔐 Authentication System

### Protected Routes
All pages except `/login` require authentication. The authentication middleware automatically:
- Redirects unauthenticated users to `/login`
- Redirects authenticated users away from `/login` to `/`
- Handles token persistence and refresh

### Authentication Flow
1. **Login**: Users authenticate via AWS Cognito
2. **Token Storage**: JWT tokens stored in localStorage
3. **Auto-refresh**: Tokens automatically refreshed when needed
4. **Logout**: Clears tokens and redirects to login

## 📄 Page Structure Pattern

The application uses a folder-based structure for pages:

```
pages/
├── page-name/             # Page folder
│   └── index.vue         # Main page component
├── page-name/             # Another page
│   ├── index.vue         # Main page
│   └── [id].vue          # Dynamic route (optional)
└── nested/                # Nested pages
    └── page-name/
        └── index.vue
```

### Benefits of This Structure:
- **Better Organization**: Related files grouped together
- **Scalability**: Easy to add sub-pages and components
- **Maintainability**: Clear separation of concerns
- **Future-proof**: Supports complex routing patterns

## 🚀 Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn
- Backend API running (see backend README)

### Installation
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Variables
Create a `.env` file in the frontend directory:
```env
# API Configuration
API_BASE_URL=http://localhost:8080
WS_URL=ws://localhost:8080

# AWS Cognito (optional for development)
COGNITO_USER_POOL_ID=your-user-pool-id
COGNITO_CLIENT_ID=your-client-id
```

## 🔧 Development

### Adding New Pages
1. Create a new folder in `pages/` with the page name
2. Add `index.vue` inside the folder
3. Apply authentication middleware if needed:
   ```vue
   <script setup>
   definePageMeta({
     middleware: 'auth'
   })
   </script>
   ```

### Example: Creating a Settings Page
```bash
mkdir pages/settings
touch pages/settings/index.vue
```

```vue
<template>
  <div>
    <h1>Settings</h1>
    <!-- Your settings content -->
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})
</script>
```

### Authentication Composable
The `useAuth` composable provides:
- `user`: Current user data
- `isAuthenticated`: Authentication status
- `login()`: User login
- `register()`: User registration
- `logout()`: User logout
- `getAuthHeaders()`: Headers for API requests

### WebSocket Composable
The `useWebSocket` composable provides:
- `isConnected`: Connection status
- `messages`: Real-time messages
- `connect()`: Connect to WebSocket
- `disconnect()`: Disconnect from WebSocket
- `sendMessage()`: Send a message

## 🎨 UI Components

The application uses Nuxt UI components:
- `UButton`: Buttons with various styles
- `UInput`: Form inputs
- `UCard`: Card containers
- `UModal`: Modal dialogs
- `UIcon`: Icons from Heroicons

## 🔒 Security Features

- **Route Protection**: All routes except login require authentication
- **Token Validation**: JWT tokens validated on each request
- **Auto-logout**: Automatic logout on token expiration
- **Secure Storage**: Tokens stored securely in localStorage
- **CORS Protection**: Proper CORS configuration

## 📱 Responsive Design

The application is fully responsive and works on:
- Desktop computers
- Tablets
- Mobile phones
- Various screen sizes

## 🧪 Testing

```bash
# Run tests
npm run test

# Run tests in watch mode
npm run test:watch

# Generate coverage report
npm run test:coverage
```

## 🚀 Deployment

### Production Build
```bash
# Build the application
npm run build

# The built files will be in .output/
```

### Docker Deployment
```bash
# Build Docker image
docker build -t ktchat-frontend .

# Run container
docker run -p 3000:3000 ktchat-frontend
```

## 📚 API Integration

The frontend communicates with the backend API:
- **Base URL**: Configurable via `API_BASE_URL`
- **Authentication**: JWT tokens in Authorization header
- **WebSocket**: Real-time messaging
- **File Upload**: Multipart form data

## 🔄 State Management

The application uses Nuxt's built-in state management:
- `useState()`: Global state
- `useAuth()`: Authentication state
- `useWebSocket()`: WebSocket state
- Local component state with `ref()` and `reactive()`

## 🎯 Best Practices

1. **Component Structure**: Use the folder-based page structure
2. **Authentication**: Always apply middleware to protected pages
3. **Error Handling**: Handle API errors gracefully
4. **Loading States**: Show loading indicators for async operations
5. **Responsive Design**: Ensure mobile compatibility
6. **TypeScript**: Use proper typing for better development experience

## 🐛 Troubleshooting

### Common Issues

1. **Authentication not working**
   - Check if backend is running
   - Verify AWS Cognito configuration
   - Check browser console for errors

2. **WebSocket connection failed**
   - Ensure backend WebSocket endpoint is available
   - Check network connectivity
   - Verify authentication token

3. **Build errors**
   - Clear node_modules and reinstall
   - Check TypeScript errors
   - Verify Nuxt configuration

## 📞 Support

For questions and support:
- Check the documentation in each component
- Review the API documentation
- Check the browser console for errors
- Run the test suite to verify functionality

## 🎉 Next Steps

After setting up the frontend:

1. **Customize UI**: Modify colors, fonts, and layout
2. **Add Features**: Implement additional chat features
3. **Testing**: Add comprehensive tests
4. **Performance**: Optimize for production
5. **Deployment**: Deploy to your hosting platform

Happy coding! 🚀
