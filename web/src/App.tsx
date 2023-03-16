import './App.css'
import React from 'react'
import { routes } from './routes'
import { useRoutes } from 'react-router-dom'
import { useQuery } from 'react-query'
import { getGptSettings } from './api/settings'

function App() {
  const Routes = useRoutes(routes.map(it => it))
  return Routes
}

export default App
