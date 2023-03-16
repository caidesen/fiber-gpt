import { crf } from '../utils/request'

export interface Settings {
  apiKey?: string
  maxToken?: number
  temperature?: number
}
export const getGptSettings = crf<Settings>('gptsetting/get')
export const setGptSettings = crf<undefined, Settings>('gptsetting/update')
// export const getGPTModels = crf<string[]>('config/gpt-models')
