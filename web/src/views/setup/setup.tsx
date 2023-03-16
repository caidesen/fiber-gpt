import { Container, Step, StepLabel, Stepper } from '@mui/material'
import React, { useMemo } from 'react'
import { useParams } from 'react-router-dom'
import { ApiKeyForm } from './ApiKey'
import { GPTOptionsForm } from './OtherSettings'

export default function Setup() {
  const params = useParams<{ step: string }>()
  const step = useMemo(() => +(params.step ?? 1) - 1, [params])
  const paperForm = useMemo(() => {
    switch (step) {
      case 0:
        return <ApiKeyForm />
      case 1:
        return <GPTOptionsForm />
    }
  }, [step])
  return (
    <Container maxWidth="sm" fixed sx={{ pt: '20vh' }}>
      <Stepper activeStep={step} alternativeLabel sx={{ mb: 5 }}>
        <Step>
          <StepLabel>设置apiKey</StepLabel>
        </Step>
        <Step>
          <StepLabel>设置GPT参数</StepLabel>
        </Step>
      </Stepper>
      {paperForm}
    </Container>
  )
}
