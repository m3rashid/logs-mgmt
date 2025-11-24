import { createSignal, onMount, onCleanup, Show, Index } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { Layout } from '../components/Layout'
import { Button } from '../components/Button'
import { FlightDetailModal } from '../components/FlightDetailModal'
import { getFlightRealtime } from '../api/flights'
import type { FlightEvent } from '../api/types'
import { formatCoordinates } from '../lib/utils'

export function Flights() {
  const navigate = useNavigate()
  const [realtimeFlights, setRealtimeFlights] = createSignal<FlightEvent[]>([])
  const [isConnected, setIsConnected] = createSignal(false)
  const [selectedFlightId, setSelectedFlightId] = createSignal<number | null>(null)

  onMount(() => {
    const eventSource = getFlightRealtime()

    eventSource.onopen = () => {
      setIsConnected(true)
    }

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as FlightEvent[]
        setRealtimeFlights(data)
      } catch (err) {
        console.error('Failed to parse realtime data:', err)
      }
    }

    eventSource.onerror = () => {
      setIsConnected(false)
    }

    onCleanup(() => {
      eventSource.close()
    })
  })

  return (
    <Layout>
      <div class='space-y-6'>
        {/* Header */}
        <div class='flex justify-between items-center'>
          <div>
            <h1 class='text-3xl font-bold text-gray-900'>Active Flights Dashboard</h1>
            <p class='text-gray-600 mt-1'>Real-time monitoring of flights in progress</p>
          </div>
          <Button onClick={() => navigate('/new-flight')}>Schedule New Flight</Button>
        </div>

        {/* Connection Status */}
        <div class='flex items-center gap-3 bg-white rounded-lg border border-gray-200 px-4 py-3'>
          <div class={`w-3 h-3 rounded-full ${isConnected() ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`} />
          <span class='text-sm font-medium text-gray-700'>
            {isConnected() ? 'Connected to Live Stream' : 'Disconnected'}
          </span>
          <span class='text-sm text-gray-500 ml-auto'>
            {realtimeFlights().length} {realtimeFlights().length === 1 ? 'flight' : 'flights'} active
          </span>
        </div>

        {/* Flights Grid */}
        <Show
          when={realtimeFlights().length > 0}
          fallback={
            <div class='bg-white rounded-lg border border-gray-200 p-12 text-center'>
              <div class='text-gray-400 mb-2'>
                <svg class='w-16 h-16 mx-auto' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
                  <path
                    stroke-linecap='round'
                    stroke-linejoin='round'
                    stroke-width='2'
                    d='M12 19l9 2-9-18-9 18 9-2zm0 0v-8'
                  />
                </svg>
              </div>
              <h3 class='text-lg font-medium text-gray-900'>No Active Flights</h3>
              <p class='text-gray-600 mt-1'>There are currently no flights in progress</p>
            </div>
          }
        >
          <div class='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
            <Index each={realtimeFlights()}>
              {(flight) => (
                <div
                  onClick={() => setSelectedFlightId(flight().flight_id)}
                  class='bg-white rounded-lg border border-gray-200 p-5 hover:shadow-lg transition-shadow cursor-pointer'
                >
                  {/* Flight Header */}
                  <div class='flex items-center justify-between mb-4'>
                    <div class='flex items-baseline gap-2'>
                      <h3 class='text-lg font-bold text-gray-900'>Flight #{flight().flight_id}</h3>
                      <span class='text-xs text-gray-500'>· Aircraft #{flight().aircraft_id}</span>
                    </div>
                    <span class='text-green-600 bg-green-50 px-2.5 py-1 rounded-full text-xs font-medium'>
                      In Flight
                    </span>
                  </div>

                  {/* Primary Metrics - 2x2 Grid */}
                  <div class='grid grid-cols-2 gap-3 mb-4'>
                    <MetricBox icon='↑' label='Altitude' value={`${flight().data.alt.toFixed(1)} m`} />
                    <MetricBox icon='→' label='Speed' value={`${flight().data.speed.toFixed(1)} m/s`} />
                    <MetricBox
                      icon='🔋'
                      label='Battery'
                      value={`${flight().sys_info.battery.toFixed(0)}%`}
                      status={getBatteryStatus(flight().sys_info.battery)}
                    />
                    <MetricBox icon='💻' label='CPU' value={`${flight().sys_info.cpu.toFixed(0)}%`} />
                  </div>

                  {/* Position - Full Width */}
                  <div class='bg-gray-50 rounded-lg px-3 py-2 mb-3'>
                    <div class='flex items-center gap-2'>
                      <span class='text-base'>📍</span>
                      <div class='flex-1 min-w-0'>
                        <div class='text-xs text-gray-500'>Position</div>
                        <div class='text-xs font-medium text-gray-900 truncate'>
                          {formatCoordinates(flight().data.lat, flight().data.lon)}
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* System Stats - Inline */}
                  <div class='flex items-center justify-between text-xs pt-3 border-t border-gray-200'>
                    <div class='flex items-center gap-3'>
                      <span class='text-gray-500'>
                        Mem: <span class='font-medium text-gray-900'>{flight().sys_info.memory.toFixed(0)}%</span>
                      </span>
                      <span class='text-gray-500'>
                        Disk: <span class='font-medium text-gray-900'>{flight().sys_info.disk.toFixed(0)}%</span>
                      </span>
                    </div>
                    <button class='text-blue-600 hover:text-blue-700 font-medium'>Details →</button>
                  </div>
                </div>
              )}
            </Index>
          </div>
        </Show>
      </div>

      {/* Flight Detail Modal */}
      <Show when={selectedFlightId() !== null}>
        <FlightDetailModal
          isOpen={selectedFlightId() !== null}
          onClose={() => setSelectedFlightId(null)}
          flightId={selectedFlightId()!}
        />
      </Show>
    </Layout>
  )
}

function MetricBox(props: { icon: string; label: string; value: string; status?: 'good' | 'warning' | 'critical' }) {
  const getBgColor = () => {
    switch (props.status) {
      case 'warning':
        return 'bg-yellow-50'
      case 'critical':
        return 'bg-red-50'
      default:
        return 'bg-gray-50'
    }
  }

  return (
    <div class={`${getBgColor()} rounded-lg p-3`}>
      <div class='flex items-center gap-1.5 mb-1'>
        <span class='text-base'>{props.icon}</span>
        <span class='text-xs text-gray-500'>{props.label}</span>
      </div>
      <div class='text-base font-bold text-gray-900'>{props.value}</div>
    </div>
  )
}

function getBatteryStatus(battery: number): 'good' | 'warning' | 'critical' {
  if (battery < 20) return 'critical'
  if (battery < 40) return 'warning'
  return 'good'
}
