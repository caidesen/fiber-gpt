import Send from '@mui/icons-material/Send'
import ArrowBack from '@mui/icons-material/ArrowBack'
import { AppBar, Container, Fab, IconButton, InputBase, Paper, Stack, Toolbar } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { useMutation, useQuery, useQueryClient } from 'react-query'
import React, { useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { MessageRow } from './MessageRow'
import { ChatBackground } from './ChatBackground'
import _ from '@/utils/lodash'
import { Chat, createChat, createQuestion, getChats, getQuestions, Question } from '@/api/chat'
import { useMd } from '@/hooks/useMd'
import { useNotification } from '@/utils/notification'

export function ChatWindow() {
  const params = useParams<{ chatId?: string }>()
  const upMd = useMd()
  const paddingSize = upMd ? 0 : 1
  const navigate = useNavigate()
  const notification = useNotification()
  useEffect(() => {
    scrollBottom()
  }, [params.chatId])
  const scrollViewRef = useRef<HTMLDivElement>(null)
  const scrollBottom = () => {
    if (!scrollViewRef.current) return
    scrollViewRef.current.scrollTop = scrollViewRef.current.scrollHeight
  }
  const { data } = useQuery(
    [getQuestions.cacheName, params.chatId],
    () => getQuestions({ chatId: params.chatId! }),
    {
      enabled: !!params.chatId && params.chatId !== 'new',
      staleTime: 1000 * 60,
      cacheTime: 1000 * 60,
      onSuccess: () => setTimeout(scrollBottom),
    },
  )
  const latestMessage = _.last(data)
  const queryClient = useQueryClient()
  const { mutateAsync: createNewChat, isLoading: createNewChatLoading } = useMutation(
    () => createChat(),
    {
      onSuccess(res) {
        queryClient.setQueryData(getChats.cacheName, (oldData?: Chat[]) => {
          if (!oldData) return [res]
          return [res, ...oldData]
        })
      },
      onError: notification.error,
    },
  )
  const { mutateAsync: send, isLoading: sendLoading } = useMutation(createQuestion, {
    onError: notification.error,
    onSuccess(res, variables) {
      queryClient.setQueryData(
        [getQuestions.cacheName, variables!.chatId],
        (oldData?: Question[]) => {
          if (!oldData) return [res]
          return [...oldData, res]
        },
      )
      if (params.chatId === 'new') {
        setTimeout(() => navigate(`/chat/${variables!.chatId}`), 10)
      }
      setTimeout(scrollBottom, 100)
    },
  })
  const mainLoading = useMemo(() => {
    if (sendLoading || createNewChatLoading) return true
    return latestMessage?.status === 'pending'
  }, [latestMessage])
  const [text, setText] = useState('')
  const onSend = async () => {
    setText('')
    scrollBottom()
    let chatId = params.chatId!
    if (chatId === 'new') {
      const chat = await createNewChat()
      chatId = chat.id
    }
    await send({ chatId: chatId, q: text })
  }
  return (
    <Container
      disableGutters
      sx={{ display: 'flex', flexDirection: 'column', height: '100vh', position: 'relative' }}
    >
      <ChatBackground
        sx={{
          position: 'absolute',
          zIndex: 0,
          top: 0,
          bottom: 0,
          left: 0,
          right: 0,
        }}
      />
      <AppBar
        position="static"
        sx={{ backgroundColor: 'background.paper', color: 'text.secondary', position: 'relative' }}
      >
        <Toolbar variant="dense">
          <Stack direction="row" alignItems="center" spacing={1}>
            {!upMd && (
              <IconButton
                size="small"
                edge="start"
                color="inherit"
                aria-label="menu"
                sx={{ mr: 2 }}
                onClick={_.partial(navigate, -1)}
              >
                <ArrowBack />
              </IconButton>
            )}
          </Stack>
        </Toolbar>
      </AppBar>
      <Grid
        container
        justifyContent="center"
        ref={scrollViewRef}
        sx={{ height: 0, flex: 1, overflow: 'auto', pb: 1, pt: 1, position: 'relative' }}
      >
        <Grid xs={12} lg={8}>
          <Stack spacing={2} sx={{ pl: paddingSize, pr: paddingSize }}>
            {data?.map(message => (
              <MessageRow message={message} key={message.id} />
            ))}
          </Stack>
        </Grid>
      </Grid>
      <Grid
        container
        justifyContent="center"
        sx={{ pt: 1, pb: upMd ? 2 : 1, position: 'relative' }}
      >
        <Grid xs={12} lg={8}>
          <Stack
            direction="row"
            spacing={1}
            sx={{ width: '100%', pl: paddingSize, pr: paddingSize }}
          >
            <Paper component="label" sx={{ display: 'flex', alignItems: 'center', flex: 1 }}>
              <InputBase
                placeholder="question"
                value={text}
                sx={{ width: '100%', height: '100%', ml: 1 }}
                onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
                  setText(event.target.value)
                }
              />
            </Paper>
            <Fab
              size="small"
              color="primary"
              aria-label="send"
              disabled={mainLoading}
              onClick={onSend}
            >
              <Send />
            </Fab>
          </Stack>
        </Grid>
      </Grid>
    </Container>
  )
}
