import { QueryClient } from 'react-query'

export interface requestOptions extends RequestInit {}

const baseUrl = new URL('/api/', location.origin)
const defaultOptions: RequestInit = {
  credentials: 'same-origin',
  headers: {
    'Content-Type': 'application/json',
  },
}

export function request<T>(url: string, options?: RequestInit): Promise<T> {
  const optionsWithDefault = Object.assign({}, defaultOptions, options)
  return fetch(new URL(url, baseUrl), optionsWithDefault).then(res => {
    if (res.ok) {
      return res.json().catch(() => undefined as T)
    }
    if (res.headers.get('content-type')?.includes('application/json')) {
      return res.json().then(e => {
        throw new Error(e.message)
      })
    }
    throw new Error(res.statusText)
  })
}
interface RequestFn<Result, Input> {
  (data?: Input, options?: RequestInit): Promise<Result>
  cacheName: string
}

export function crf<Result = any, Input = any>(
  url: string,
  o?: RequestInit,
): RequestFn<Result, Input> {
  const fn = function (data?: Input, options?: RequestInit) {
    return request<Result>(
      url,
      Object.assign({ method: 'POST', body: data && JSON.stringify(data) }, o, options),
    )
  }
  fn.cacheName = url.toString()
  return fn
}

export const queryClient = new QueryClient()
