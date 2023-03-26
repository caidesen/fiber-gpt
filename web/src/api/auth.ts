import { crf } from '@/utils/request'

export interface User {
  id: string
  username: string
  nickname: string
}
export interface LoginInput {
  username: string
  password: string
}
export const login = crf<User, LoginInput>('auth/login')
export const register = crf<User, LoginInput>('auth/register')
