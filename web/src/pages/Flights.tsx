import { createSignal, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { createQuery } from '@tanstack/solid-query'
import { Layout } from '../components/Layout'
import { Table, type Column } from '../components/Table'
import { Button } from '../components/Button'
import { Pagination } from '../components/Pagination'
import { getFlights } from '../api/flights'
import type { Flight, FlightStatus } from '../api/types'
import { formatDate, formatCoordinates, getFlightStatusLabel } from '../lib/utils'

export function Flights() {
  const navigate = useNavigate()
  const [page, setPage] = createSignal(1)

  const flightsQuery = createQuery(() => ({
    queryKey: ['flights', page()],
    queryFn: () => getFlights({ page: page(), limit: 15 }),
  }))

  const columns: Column<Flight>[] = [
    { key: 'id', label: 'ID' },
    {
      key: 'source',
      label: 'Source',
      render: (flight) => formatCoordinates(flight.source_lat, flight.source_lon),
    },
    {
      key: 'destination',
      label: 'Destination',
      render: (flight) => formatCoordinates(flight.destination_lat, flight.destination_lon),
    },
    {
      key: 'start_time',
      label: 'Start Time',
      render: (flight) => formatDate(flight.start_time),
    },
    {
      key: 'end_time',
      label: 'End Time',
      render: (flight) => formatDate(flight.end_time),
    },
    {
      key: 'status',
      label: 'Status',
      render: (flight) => <span class={getStatusColor(flight.status)}>{getFlightStatusLabel(flight.status)}</span>,
    },
  ]

  return (
    <Layout>
      <div class='space-y-4'>
        <div class='flex justify-between items-center'>
          <h1 class='text-2xl font-bold text-gray-900'>Flights</h1>
          <Button onClick={() => navigate('/new-flight')}>Schedule New Flight</Button>
        </div>

        <Table
          columns={columns}
          data={flightsQuery.data?.data || []}
          isLoading={flightsQuery.isLoading}
          onRowClick={(flight) => navigate(`/flights/${flight.id}`)}
          emptyMessage='No flights found'
        />

        <Show when={flightsQuery.data}>
          <Pagination pagination={flightsQuery.data!} onPageChange={setPage} />
        </Show>
      </div>
    </Layout>
  )
}

function getStatusColor(status: FlightStatus) {
  const colors = [
    'text-blue-600 bg-blue-50 px-2 py-1 rounded-full text-xs font-medium', // Scheduled
    'text-purple-600 bg-purple-50 px-2 py-1 rounded-full text-xs font-medium', // InFlight
    'text-green-600 bg-green-50 px-2 py-1 rounded-full text-xs font-medium', // Completed
    'text-red-600 bg-red-50 px-2 py-1 rounded-full text-xs font-medium', // Failed
  ]
  return colors[status] || 'text-gray-600 bg-gray-50 px-2 py-1 rounded-full text-xs font-medium'
}
