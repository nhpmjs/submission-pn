import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom'
import React from 'react'

import { LoadingSpinner } from './components/LoadingSpinner'

const Home = React.lazy(() => import('./pages/Home'))
const GamePlay = React.lazy(() => import('./pages/GamePlay'))

const queryClient = new QueryClient()

export default function App() {
  return (
    <React.Suspense fallback={<LoadingSpinner />}>
      <QueryClientProvider client={queryClient}>
        <Router>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/game/:gameId" element={<GamePlay />} />
          </Routes>
        </Router>
      </QueryClientProvider>
    </React.Suspense>
  )
}
