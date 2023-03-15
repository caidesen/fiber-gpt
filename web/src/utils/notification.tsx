import React from 'react'
import { SnackbarProvider, useSnackbar } from 'notistack'

export const useNotification = () => {
  const { enqueueSnackbar } = useSnackbar()
  return {
    info: (message: string) => {
      enqueueSnackbar(message, { variant: 'info' })
    },
    success: (message: string) => {
      enqueueSnackbar(message, { variant: 'success' })
    },
    warning: (message: string) => {
      enqueueSnackbar(message, { variant: 'warning' })
    },
    error: (err: any) => {
      enqueueSnackbar(err instanceof Error ? err.message : err, { variant: 'error' })
    },
  }
}

export const NotificationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <SnackbarProvider
    maxSnack={3}
    anchorOrigin={{ horizontal: 'right', vertical: 'top' }}
    autoHideDuration={2000}
  >
    {children}
  </SnackbarProvider>
)
