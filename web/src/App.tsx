import './App.css'
import React from 'react'
import { routes } from './routes'
import { useRoutes } from 'react-router-dom'

function App() {
  const Routes = useRoutes(routes.map(it => it))
  return Routes
}

export default App
