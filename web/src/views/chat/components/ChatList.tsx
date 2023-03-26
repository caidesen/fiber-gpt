import { AppBar, Fab, IconButton, Paper, Stack, Toolbar } from '@mui/material'
import Add from '@mui/icons-material/Add'
import Settings from '@mui/icons-material/Settings'
import { useQuery } from 'react-query'
import { getChats } from '@/api/chat'
import React from 'react'
import { ChatCard } from './ChatCard'
import { useNavigate, useParams } from 'react-router-dom'
import _ from '@/utils/lodash'
import { useNotification } from '@/utils/notification'
import { useMd } from '@/hooks/useMd'

export function ChatList() {
  const navigate = useNavigate()
  const notification = useNotification()
  const upMd = useMd()
  const params = useParams<{ chatId?: string }>()
  const { data } = useQuery(getChats.cacheName, () => getChats(), {
    staleTime: 1000 * 60 * 5,
    onError: notification.error,
    onSuccess: res => {
      if ((!params.chatId || !res?.find(chat => chat.id === params.chatId)) && upMd) {
        const id = _.head(res)?.id
        if (id) navigate(`/chat/${id}`)
      }
    },
  })

  return (
    <Paper
      square
      sx={{ height: '100vh', position: 'relative', display: 'flex', flexDirection: 'column' }}
    >
      <AppBar position="static">
        <Toolbar variant="dense">
          <IconButton
            size="small"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
            onClick={() => navigate('/setup')}
          >
            <Settings />
          </IconButton>
        </Toolbar>
      </AppBar>
      <Stack spacing={1} sx={{ p: 1, flex: 1, overflow: 'auto' }}>
        {params.chatId === 'new' && (
          <ChatCard
            chat={{ id: 'new', createdAt: new Date().toString(), updatedAt: '', questions: [] }}
            key="new"
          />
        )}
        {data?.map(chat => (
          <ChatCard chat={chat} key={chat.id} />
        ))}
      </Stack>
      <Fab
        size="small"
        color="primary"
        aria-label="add"
        sx={{ position: 'absolute', bottom: 10, right: 10 }}
        onClick={() => navigate('/chat/new')}
      >
        <Add />
      </Fab>
    </Paper>
  )
}
