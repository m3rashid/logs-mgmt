import { Router, Route } from '@solidjs/router'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import { Flights } from './pages/Flights'
import { Aircrafts } from './pages/Aircrafts'
import { NewAircraft } from './pages/NewAircraft'
import { NewFlight } from './pages/NewFlight'
import { FlightDetail } from './pages/FlightDetail'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 30000, // 30 seconds
    },
  },
})

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Route path='/' component={Flights} />
        <Route path='/aircrafts' component={Aircrafts} />
        <Route path='/new-aircraft' component={NewAircraft} />
        <Route path='/new-flight' component={NewFlight} />
        <Route path='/flights/:id' component={FlightDetail} />
      </Router>
    </QueryClientProvider>
  )
}
