import { Container, TextField, Typography } from '@mui/material'
import { useMutation } from 'react-query'
import { register } from '@/api/auth'
import { useNotification } from '@/utils/notification'
import { useNavigate } from 'react-router-dom'
import { Controller, useForm } from 'react-hook-form'
import { LoadingButton } from '@mui/lab'
import _ from '@/utils/lodash'

export default function Register() {
  const navigate = useNavigate()
  const notification = useNotification()
  const { mutate, isLoading } = useMutation(register, {
    onSuccess() {
      notification.success('注册成功')
      navigate('/login', { replace: true })
    },
    onError: notification.error,
  })
  const { handleSubmit, control } = useForm({
    defaultValues: {
      username: '',
      password: '',
      passwordAgain: '',
    },
  })
  return (
    <Container maxWidth="xs" sx={{ pt: '20vh' }}>
      <form
        onSubmit={handleSubmit(
          data => {
            mutate(data)
          },
          errors => {
            notification.error(_.head(Object.values(errors))?.message)
          },
        )}
      >
        <Typography variant="h4" component="h1" sx={{ mb: '10vh' }} textAlign="center">
          Fiber GPT
        </Typography>
        <Controller
          name="username"
          rules={{ required: '用户名不能为空' }}
          render={({ field }) => (
            <TextField {...field} label="用户名" sx={{ mb: 4 }} variant="standard" fullWidth />
          )}
          control={control}
        />
        <Controller
          name="password"
          rules={{ required: '密码不能为空' }}
          render={({ field }) => (
            <TextField
              {...field}
              type="password"
              variant="standard"
              fullWidth
              label="密码"
              sx={{ mb: 4 }}
            />
          )}
          control={control}
        />
        <Controller
          name="passwordAgain"
          rules={{
            required: '密码不能为空',
            validate(val, formState) {
              if (val !== formState.password) {
                return '两次输入的密码不一致'
              }
              return true
            },
          }}
          render={({ field }) => (
            <TextField
              {...field}
              type="password"
              variant="standard"
              fullWidth
              label="再次输入密码"
              sx={{ mb: 4 }}
            />
          )}
          control={control}
        />
        <LoadingButton
          loading={isLoading}
          type="submit"
          variant="contained"
          color="primary"
          fullWidth
          size="large"
        >
          注册
        </LoadingButton>
      </form>
    </Container>
  )
}
