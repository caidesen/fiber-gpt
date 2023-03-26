import { RouteObject } from 'react-router-dom'
import { Box } from '@mui/material'
import Login from './views/login'
import Setup from './views/setup/setup'
import Chat from './views/chat'
import Register from '@/views/register'

export const routes: Array<RouteObject & { name: string }> = [
  {
    path: '/',
    name: '首页',
    element: <Chat />,
  },
  {
    path: '/chat',
    name: '聊天',
    element: <Chat />,
  },
  {
    path: '/chat/:chatId',
    name: '聊天',
    element: <Chat />,
  },
  {
    path: '/login',
    name: '登录',
    element: <Login />,
  },
  {
    path: '/register',
    name: '注册',
    element: <Register />,
  },
  {
    path: '/setup',
    name: '设置',
    element: <Setup />,
  },
  {
    path: '/setup/:step',
    name: '设置',
    element: <Setup />,
  },
  {
    path: '*',
    name: '404',
    element: <Box sx={{ width: '100%', height: '100%' }}>404</Box>,
  },
]
