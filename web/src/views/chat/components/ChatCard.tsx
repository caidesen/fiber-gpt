import { Avatar, Box, Stack } from '@mui/material'
import React, { useMemo } from 'react'
import { Chat } from '@/api/chat'
import dayjs from 'dayjs'
import { useNavigate, useParams } from 'react-router-dom'
import { styled } from '@mui/material/styles'
import OpenAiIcon from '@/assets/openai.svg'
import { grey } from '@mui/material/colors'
import _ from '@/utils/lodash'

const CardWrapper = styled(Box)`
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
  user-select: none;
`

export function ChatCard({ chat }: { chat: Chat }) {
  const params = useParams<{ chatId?: string }>()
  const navigate = useNavigate()
  const lastMessage = _.last(chat.questions)
  const lastMessageTime = useMemo(() => {
    const updateAt = dayjs(lastMessage?.updatedAt ?? chat.createdAt)
    const now = dayjs()
    if (updateAt.isSame(now, 'minute')) {
      return '刚刚'
    }
    if (updateAt.isSame(now, 'day')) {
      return updateAt.format('HH:mm')
    }
    if (updateAt.isSame(now, 'year')) {
      return updateAt.format('MM-DD')
    }
    return updateAt.format('YYYY')
  }, [chat, lastMessage])
  const active = params.chatId === chat.id
  return (
    <CardWrapper
      sx={{
        p: 1,
        ':hover': {
          backgroundColor: 'action.hover',
        },
        ...(active && {
          backgroundColor: 'primary.main',
          ':hover': {},
          color: 'primary.contrastText',
        }),
      }}
      onClick={() => navigate(`/chat/${chat.id}`)}
    >
      <Stack direction="row" spacing={2}>
        <Avatar sx={{ p: 0.5, bgcolor: grey[300] }}>
          <img src={OpenAiIcon} alt="" />
        </Avatar>
        {/*超出两行剩余*/}
        <Box
          sx={{
            flex: 1,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
          }}
        >
          {lastMessage?.a || lastMessage?.q || '新建聊天'}
        </Box>
        <Box>{lastMessageTime}</Box>
      </Stack>
    </CardWrapper>
  )
}
