import { Stack, TextField } from '@mui/material'
import React from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from 'react-query'
import { getGptSettings, setGptSettings } from '../../api/settings'
import { Controller, useForm } from 'react-hook-form'
import { LoadingButton } from '@mui/lab'
import { useNotification } from '../../utils/notification'
import _ from 'lodash'

const geApiKey = _.partial(_.get, _, 'apiKey', '')
export function ApiKeyForm() {
  const { error } = useNotification()
  const { data: originData } = useQuery(getGptSettings.cacheName, () => getGptSettings(), {
    staleTime: 5 * 60 * 1000,
    onSuccess(res) {
      setValue('apiKey', geApiKey(res))
    },
  })
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { mutate, isLoading: updateLoading } = useMutation(setGptSettings, {
    onSuccess: (data, variables) => {
      queryClient.setQueryData(getGptSettings.cacheName, Object.assign({}, originData, variables))
      navigate('/setup/2')
    },
    onError: error,
  })
  const { handleSubmit, control, setValue } = useForm({
    defaultValues: {
      apiKey: geApiKey(originData),
    },
  })
  return (
    <form onSubmit={handleSubmit(_.unary(mutate))}>
      <Stack direction="column" spacing={2}>
        <Controller
          name="apiKey"
          control={control}
          rules={{
            required: '请输入apiKey',
          }}
          render={({ field, formState }) => (
            <TextField
              {...field}
              error={!!formState.errors.apiKey}
              name="apiKey"
              size="small"
              label="apiKey"
              helperText={String(
                formState.errors.apiKey?.message ?? '请前往 https://openai.com/developers/ 申请',
              )}
              placeholder="请输入"
              fullWidth
              autoFocus
            />
          )}
        />
        <Stack direction="row-reverse">
          <LoadingButton loading={updateLoading} variant="contained" type="submit" size="small">
            下一步
          </LoadingButton>
        </Stack>
      </Stack>
    </form>
  )
}
