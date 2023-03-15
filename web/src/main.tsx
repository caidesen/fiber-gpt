import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'
import { CssBaseline } from '@mui/material'
import { BrowserRouter } from 'react-router-dom'
import { ThemeContextProvider } from './theme'
import { QueryClientProvider } from 'react-query'
import { queryClient } from './utils/request'
import { NotificationProvider } from './utils/notification'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <BrowserRouter>
      <ThemeContextProvider>
        <CssBaseline />
        <NotificationProvider>
          <QueryClientProvider client={queryClient}>
            <App />
          </QueryClientProvider>
        </NotificationProvider>
      </ThemeContextProvider>
    </BrowserRouter>
  </React.StrictMode>,
)
