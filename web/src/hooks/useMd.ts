import { useTheme } from '@mui/material'
import useMediaQuery from '@mui/material/useMediaQuery'

export function useMd() {
  const theme = useTheme()
  return useMediaQuery(theme.breakpoints.up('md'))
}
