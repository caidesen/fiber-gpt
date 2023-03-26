import { Box, Container, Drawer } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { ChatList } from './components/ChatList'
import { ChatWindow } from './components/ChatWindow'
import { useQuery } from 'react-query'
import { getChats } from '@/api/chat'
import _ from '@/utils/lodash'
import { useNavigate, useParams } from 'react-router-dom'
import { useMd } from '@/hooks/useMd'
import { getGptSettings } from '@/api/settings'

export default function Chat() {
  const params = useParams<{ chatId?: string }>()
  const navigate = useNavigate()
  const upMd = useMd()
  useQuery(getGptSettings.cacheName, () => getGptSettings(), {
    staleTime: 5 * 60 * 1000,
    onSuccess(res) {
      if (!res.apiKey) {
        navigate('/setup')
      }
    },
  })
  useQuery(getChats.cacheName, () => getChats(), {
    staleTime: 1000 * 60 * 5,
    onSuccess: res => {
      if ((!params.chatId || !res?.find(chat => chat.id === params.chatId)) && upMd) {
        const id = _.head(res)?.id
        if (id) navigate(`/chat/${id}`)
      }
    },
  })
  if (!upMd) {
    return (
      <Container maxWidth="md" disableGutters>
        <ChatList />
        <Drawer anchor="right" open={!!params.chatId}>
          <Box sx={{ width: '100vw' }}>
            <ChatWindow />
          </Box>
        </Drawer>
      </Container>
    )
  }
  return (
    <Container maxWidth="xl" disableGutters>
      <Grid container>
        <Grid xs={12} md={3}>
          <ChatList />
        </Grid>
        <Grid sx={{ flex: 1 }}>
          <ChatWindow />
        </Grid>
      </Grid>
    </Container>
  )
}
