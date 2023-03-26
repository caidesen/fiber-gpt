import { Avatar, Box, Paper, Stack, SxProps, Theme, Typography } from '@mui/material'
import SentimentSatisfiedAlt from '@mui/icons-material/SentimentSatisfiedAlt'
import { grey } from '@mui/material/colors'
import OpenAiIcon from '@/assets/openai.svg'
import React, { useEffect, useRef } from 'react'
import { getQuestions, Question, readAnswer } from '@/api/chat'
import './MessageRow.css'
import { useParams } from 'react-router-dom'
import { useQueryClient } from 'react-query'
export function MessageRow({ message }: { message: Question }) {
  const bubbleSx: SxProps<Theme> = {
    p: 1,
  }
  const spanRef = useRef<HTMLSpanElement>(null)
  const params = useParams<{ chatId: string }>()
  const queryClient = useQueryClient()
  useEffect(() => {
    if (message?.status !== 'pending') return
    const abortController = new AbortController()
    readAnswer({ qid: message.id }, { signal: abortController.signal })
      .then(res => {
        const body = (res as unknown as Response).body
        if (!body) return Promise.reject(new Error('no body'))
        return body
      })
      .then(stream =>
        stream.pipeTo(
          new WritableStream({
            write(chunk) {
              let str = new TextDecoder().decode(chunk)
              if (spanRef.current) {
                if (spanRef.current.children.length === 0) str = str.replace(/^\n\n/, '')
                const newWord = document.createElement('span')
                newWord.className = 'MessageRow__word'
                newWord.innerText = str
                spanRef.current.appendChild(newWord)
                newWord.scrollIntoView({ block: 'start' })
              }
            },
          }),
        ),
      )
      .then(() => {
        queryClient.setQueryData(
          [getQuestions.cacheName, params.chatId],
          (oldData?: Question[]) =>
            oldData?.map(it => {
              if (it.id !== message.id) return it
              return { ...it, a: spanRef.current?.innerText ?? '', status: 'success' } as Question
            }) ?? [],
        )
      })
    return abortController.abort.bind(abortController)
  }, [message?.status])
  return (
    <>
      <Stack direction="row-reverse" alignItems="end" spacing={1}>
        <Avatar>
          <SentimentSatisfiedAlt />
        </Avatar>
        <Paper sx={bubbleSx}>{message.q}</Paper>
        <Box sx={{ width: 120 }}></Box>
      </Stack>
      <Stack direction="row" alignItems="end" spacing={1}>
        <Avatar sx={{ p: 0.5, bgcolor: grey[300] }}>
          <img src={OpenAiIcon} alt="" />
        </Avatar>
        <Paper sx={bubbleSx}>
          <Typography whiteSpace="pre-line" sx={{ wordBreak: 'break-all' }} ref={spanRef}>
            {message.a.replace(/^\n\n/, '')}
          </Typography>
        </Paper>
        <Box sx={{ width: 120 }}></Box>
      </Stack>
    </>
  )
}
