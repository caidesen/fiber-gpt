import { RouteObject } from 'react-router-dom'
import { Box } from '@mui/material'
import Login from './views/login'

export const routes: Array<RouteObject & { name: string }> = [
  {
    path: '/',
    name: '首页',
    element: <Box sx={{ width: '100%', height: '100%' }} />,
  },
  {
    path: '/login',
    name: '登录',
    element: <Login />,
  },
]
