import { Button, Stack, TextField } from '@mui/material'
import React from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from 'react-query'
import { getGptSettings, setGptSettings } from '../../api/settings'
import { useNotification } from '../../utils/notification'
import { Controller, useForm } from 'react-hook-form'
import { LoadingButton } from '@mui/lab'
import _ from 'lodash'

const getMaxToken = _.partial(_.get, _, 'maxToken', 64)
const getTemperature = _.partial(_.get, _, 'temperature', 0.7)
export function GPTOptionsForm() {
  const { error } = useNotification()
  const { data: originData } = useQuery(getGptSettings.cacheName, () => getGptSettings(), {
    staleTime: 5 * 60 * 1000,
    onSuccess(res) {
      setValue('maxToken', getMaxToken(res))
      setValue('temperature', getTemperature(res))
    },
  })
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const { mutate, isLoading: updateLoading } = useMutation(setGptSettings, {
    onSuccess: (data, variables) => {
      queryClient.setQueryData(getGptSettings.cacheName, Object.assign({}, data, variables))
      navigate('/')
    },
    onError: error,
  })
  const { handleSubmit, control, setValue } = useForm({
    defaultValues: {
      maxToken: getMaxToken(originData),
      temperature: getTemperature(originData),
    },
  })
  const onSubmit = handleSubmit(_.flow(_.partial(_.mapValues, _, Number), mutate))
  return (
    <form onSubmit={onSubmit}>
      <Stack direction="column" spacing={2}>
        <Controller
          name="maxToken"
          control={control}
          rules={{
            required: '请输入maxToken',
          }}
          render={({ field, formState }) => (
            <TextField
              {...field}
              error={!!formState.errors.maxToken}
              name="maxToken"
              size="small"
              label="maxToken"
              helperText={String(formState.errors.maxToken?.message) ?? '生成的最大token数'}
              placeholder="请输入"
              fullWidth
            />
          )}
        />
        <Controller
          name="temperature"
          control={control}
          rules={{
            required: '请输入temperature',
            min: { value: 0, message: '不能小于0' },
            max: { value: 1, message: '不能大于1' },
            validate: { isNumber: value => (_.isFinite(value) ? true : '请输入数字') },
          }}
          render={({ field, formState }) => (
            <TextField
              {...field}
              type="number"
              error={!!formState.errors.temperature}
              name="temperature"
              size="small"
              label="temperature"
              helperText={String(formState.errors.temperature?.message ?? '随机性')}
              placeholder="请输入"
              fullWidth
            />
          )}
        />
        <Stack direction="row-reverse">
          <LoadingButton loading={updateLoading} variant="contained" type="submit" size="small">
            完成
          </LoadingButton>
          <Button
            onClick={() => {
              navigate('/setup/1')
            }}
          >
            上一步
          </Button>
        </Stack>
      </Stack>
    </form>
  )
}
