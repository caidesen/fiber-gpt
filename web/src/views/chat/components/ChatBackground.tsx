import { Box, SxProps, Theme } from '@mui/material'
import { styled } from '@mui/material/styles'
import { useEffect, useRef, useState } from 'react'
import patternSvg from '@/assets/pattern.svg'
import _ from '@/utils/lodash'
const Canvas = styled('canvas')`
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  z-index: 0;
`
const CanvasSingle = styled(Canvas)`
  opacity: 0.1;
`
const loadImage = (url: string) =>
  new Promise<HTMLImageElement>((resolve, reject) => {
    const img = new Image()
    img.src = url
    img.onload = () => resolve(img)
    img.onerror = reject
  })
export function ChatBackground(props: { sx?: SxProps<Theme> }) {
  const canvas = useRef<HTMLCanvasElement>(null)
  const [size, setSize] = useState({})
  useEffect(() => {
    const onReset = _.throttle(() => setSize({}), 500, { trailing: true })
    window.addEventListener('resize', onReset)
    return () => window.removeEventListener('resize', onReset)
  }, [])
  useEffect(() => {
    if (!canvas.current) return
    const ctx = canvas.current.getContext('2d')
    if (!ctx) return
    const rect = canvas.current.getBoundingClientRect()
    const maxSize = Math.max(rect.width, rect.height)
    canvas.current.width = maxSize * 2
    canvas.current.height = maxSize * 2
    ctx.fillStyle = '#c5cc89'
    ctx.fillRect(0, 0, canvas.current.width, canvas.current.height)
    const rate = 1.8
    let gradient = ctx.createRadialGradient(
      0,
      0,
      canvas.current.width / rate / 2,
      0,
      0,
      canvas.current.width / rate,
    )
    gradient.addColorStop(0, '#7ca885')
    gradient.addColorStop(1, '#c5cc89')
    ctx.fillStyle = gradient
    ctx.fillRect(0, 0, canvas.current.width / rate, canvas.current.width / rate)
    gradient = ctx.createRadialGradient(
      canvas.current.width,
      canvas.current.height,
      canvas.current.width / rate / 2,
      canvas.current.width,
      canvas.current.height,
      canvas.current.width / rate,
    )
    gradient.addColorStop(0, '#7ca885')
    gradient.addColorStop(1, 'rgba(124,168,133,0.01)')
    ctx.fillStyle = gradient
    ctx.fillRect(
      canvas.current.width - canvas.current.width / rate,
      canvas.current.height - canvas.current.width / rate,
      canvas.current.width / rate,
      canvas.current.width / rate,
    )
  }, [])
  const signCanvas = useRef<HTMLCanvasElement>(null)
  useEffect(() => {
    if (!signCanvas.current) return
    const ctx = signCanvas.current.getContext('2d')
    const rect = signCanvas.current.getBoundingClientRect()
    signCanvas.current.height = rect.height * 2
    signCanvas.current.width = rect.width * 2
    const canvasHeight = signCanvas.current.height
    const canvasWidth = signCanvas.current.width
    if (!ctx) return
    loadImage(patternSvg).then(img => {
      const pattern = ctx.createPattern(img, 'repeat-x')
      if (!pattern) return
      ctx.fillStyle = pattern
      ctx.fillRect(0, 0, canvasWidth, canvasHeight)
    })
  }, [size])
  return (
    <Box sx={{ position: 'relative', ...props.sx }}>
      <Canvas ref={canvas}></Canvas>
      <CanvasSingle ref={signCanvas}></CanvasSingle>
    </Box>
  )
}
