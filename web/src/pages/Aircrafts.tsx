import { createSignal, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { Layout } from '../components/Layout'
import { Table, type Column } from '../components/Table'
import { Button } from '../components/Button'
import { Pagination } from '../components/Pagination'
import { useAircraftsQuery } from '../api/aircrafts'
import type { Aircraft, AircraftStatus } from '../api/types'
import { formatDate, getAircraftStatusLabel } from '../lib/utils'

export function Aircrafts() {
  const navigate = useNavigate()
  const [page, setPage] = createSignal(1)

  const aircraftsQuery = useAircraftsQuery(() => ({ page: page(), limit: 15 }))

  const columns: Column<Aircraft>[] = [
    { key: 'id', label: 'ID' },
    { key: 'name', label: 'Name' },
    { key: 'model', label: 'Model' },
    { key: 'serial_number', label: 'Serial #' },
    {
      key: 'battery_health',
      label: 'Battery',
      render: (aircraft) => `${aircraft.battery_health}%`,
    },
    {
      key: 'status',
      label: 'Status',
      render: (aircraft) => (
        <span class={getStatusColor(aircraft.status)}>{getAircraftStatusLabel(aircraft.status)}</span>
      ),
    },
    {
      key: 'last_maintenance_date',
      label: 'Last Maintenance',
      render: (aircraft) => formatDate(aircraft.last_maintenance_date),
    },
  ]

  return (
    <Layout>
      <div class='space-y-4'>
        <div class='flex justify-between items-center'>
          <h1 class='text-2xl font-bold text-gray-900'>Aircrafts</h1>
          <Button onClick={() => navigate('/new-aircraft')}>Add New Aircraft</Button>
        </div>

        <Table
          columns={columns}
          data={aircraftsQuery.data?.docs || []}
          isLoading={aircraftsQuery.isLoading}
          emptyMessage='No aircrafts found'
        />

        <Show when={aircraftsQuery.data}>
          <Pagination pagination={aircraftsQuery.data!} onPageChange={setPage} />
        </Show>
      </div>
    </Layout>
  )
}

function getStatusColor(status: AircraftStatus) {
  const colors = [
    'text-gray-600 bg-gray-50 px-2 py-1 rounded-full text-xs font-medium', // Inactive
    'text-green-600 bg-green-50 px-2 py-1 rounded-full text-xs font-medium', // Active
    'text-yellow-600 bg-yellow-50 px-2 py-1 rounded-full text-xs font-medium', // Maintenance
  ]
  return colors[status] || 'text-gray-600 bg-gray-50 px-2 py-1 rounded-full text-xs font-medium'
}
