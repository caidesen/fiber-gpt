import { crf } from '@/utils/request'

export interface Question {
  id: string
  q: string
  a: string
  createdAt: string
  updatedAt: string
  status: 'pending' | 'success' | 'failed'
}
export interface Chat {
  id: string
  questions: Question[]
  createdAt: string
  updatedAt: string
}
export interface CreateQuestionInput {
  chatId: string
  q: string
}
export const getChats = crf<Chat[], undefined>('chat/getChats')
export const createChat = crf<Chat, undefined>('chat/createChat')
export const createQuestion = crf<Question, CreateQuestionInput>('chat/createQuestion')
export const getQuestions = crf<Question[], { chatId: string }>('chat/getQuestions')
export const readAnswer = crf<string, { qid: string }>('chat/readAnswer')
