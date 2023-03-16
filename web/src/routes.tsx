import { RouteObject } from 'react-router-dom'
import { Box } from '@mui/material'
import Login from './views/login'
import Setup from './views/setup/setup'

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
  {
    path: '/setup/:step',
    name: '设置',
    element: <Setup />,
  },
]
