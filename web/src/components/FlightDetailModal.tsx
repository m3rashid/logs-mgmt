import { createSignal, onMount, onCleanup, Show } from 'solid-js'
import { Modal } from './Modal'
import { useFlightQuery, getFlightRealtime } from '../api/flights'
import type { FlightEvent } from '../api/types'
import { formatDate, formatCoordinates, getFlightStatusLabel, getAircraftStatusLabel } from '../lib/utils'

type FlightDetailModalProps = {
  isOpen: boolean
  onClose: () => void
  flightId: number
}

export function FlightDetailModal(props: FlightDetailModalProps) {
  const [realtimeData, setRealtimeData] = createSignal<FlightEvent | null>(null)
  const [isConnected, setIsConnected] = createSignal(false)

  const flightQuery = useFlightQuery(() => props.flightId)

  onMount(() => {
    if (!props.isOpen) return

    const eventSource = getFlightRealtime()

    eventSource.onopen = () => {
      setIsConnected(true)
    }

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as FlightEvent[]
        const flightEvent = data.find((e) => e.flight_id === props.flightId)
        if (flightEvent) setRealtimeData(flightEvent)
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
    <Modal isOpen={props.isOpen} onClose={props.onClose} title={`Flight #${props.flightId}`} size='xl'>
      <Show when={!flightQuery.isLoading && flightQuery.data}>
        {(flight) => (
          <div class='space-y-6'>
            <div class='flex items-center justify-between'>
              <p class='text-gray-600'>
                {formatCoordinates(flight().source_lat, flight().source_lon)} →{' '}
                {formatCoordinates(flight().destination_lat, flight().destination_lon)}
              </p>
              <div class='flex items-center gap-2'>
                <div class={`w-2 h-2 rounded-full ${isConnected() ? 'bg-green-500' : 'bg-red-500'}`} />
                <span class='text-sm text-gray-600'>{isConnected() ? 'Live Updates' : 'Disconnected'}</span>
              </div>
            </div>

            {/* Flight Info Card */}
            <div class='bg-white rounded-lg border border-gray-200 p-6'>
              <h3 class='text-lg font-semibold text-gray-900 mb-4'>Flight Information</h3>
              <div class='grid grid-cols-2 md:grid-cols-3 gap-4'>
                <InfoItem label='Status' value={getFlightStatusLabel(flight().flight_status)} />
                <InfoItem label='Start Time' value={formatDate(flight().start_time)} />
                <InfoItem label='End Time' value={formatDate(flight().end_time)} />
                <InfoItem label='Source' value={formatCoordinates(flight().source_lat, flight().source_lon)} />
                <InfoItem
                  label='Destination'
                  value={formatCoordinates(flight().destination_lat, flight().destination_lon)}
                />
              </div>
            </div>

            {/* Aircraft Info Card */}
            <Show when={flight().aircraft}>
              {(aircraft) => (
                <div class='bg-white rounded-lg border border-gray-200 p-6'>
                  <h3 class='text-lg font-semibold text-gray-900 mb-4'>Aircraft Information</h3>
                  <div class='grid grid-cols-2 md:grid-cols-3 gap-4'>
                    <InfoItem label='Name' value={aircraft().name} />
                    <InfoItem label='Model' value={aircraft().model} />
                    <InfoItem label='Serial Number' value={aircraft().serial_number} />
                    <InfoItem label='Battery Health' value={`${aircraft().battery_health}%`} />
                    <InfoItem label='Status' value={getAircraftStatusLabel(aircraft().aircraft_status)} />
                    <InfoItem label='Last Maintenance' value={formatDate(aircraft().last_maintenance_date)} />
                  </div>
                </div>
              )}
            </Show>

            {/* Realtime Data Card */}
            <Show when={realtimeData()}>
              {(data) => (
                <div class='bg-white rounded-lg border border-gray-200 p-6'>
                  <h3 class='text-lg font-semibold text-gray-900 mb-4'>Real-Time Telemetry</h3>
                  <div class='grid grid-cols-2 md:grid-cols-4 gap-4'>
                    <RealtimeMetric label='Altitude' value={`${data().data.alt.toFixed(1)} m`} icon='↑' />
                    <RealtimeMetric label='Speed' value={`${data().data.speed.toFixed(1)} m/s`} icon='→' />
                    <RealtimeMetric
                      label='Position'
                      value={formatCoordinates(data().data.lat, data().data.lon)}
                      icon='📍'
                    />
                    <RealtimeMetric label='Battery' value={`${data().sys_info.battery.toFixed(0)}%`} icon='🔋' />
                  </div>

                  <div class='mt-4 pt-4 border-t border-gray-200'>
                    <h4 class='text-sm font-semibold text-gray-700 mb-3'>System Information</h4>
                    <div class='grid grid-cols-2 md:grid-cols-3 gap-3'>
                      <InfoItem label='CPU' value={`${data().sys_info.cpu.toFixed(1)}%`} />
                      <InfoItem label='Memory' value={`${data().sys_info.memory.toFixed(1)}%`} />
                      <InfoItem label='Disk' value={`${data().sys_info.disk.toFixed(1)}%`} />
                      <InfoItem label='Network' value={`${data().sys_info.network.toFixed(1)} KB/s`} />
                      <InfoItem label='Uptime' value={`${(data().sys_info.uptime / 3600).toFixed(1)} hrs`} />
                      <InfoItem label='Last Updated' value={formatDate(data().time)} />
                    </div>
                  </div>
                </div>
              )}
            </Show>

            <Show when={!realtimeData() && isConnected()}>
              <div class='bg-blue-50 border border-blue-200 rounded-lg p-4 text-center text-blue-700'>
                Waiting for real-time telemetry data...
              </div>
            </Show>
          </div>
        )}
      </Show>

      <Show when={flightQuery.isLoading}>
        <div class='text-center py-12'>
          <div class='text-gray-600'>Loading flight details...</div>
        </div>
      </Show>

      <Show when={flightQuery.isError}>
        <div class='bg-red-50 border border-red-200 rounded-lg p-4 text-center text-red-700'>
          Failed to load flight details
        </div>
      </Show>
    </Modal>
  )
}

function InfoItem(props: { label: string; value: string }) {
  return (
    <div>
      <dt class='text-sm font-medium text-gray-500'>{props.label}</dt>
      <dd class='mt-1 text-sm text-gray-900'>{props.value}</dd>
    </div>
  )
}

function RealtimeMetric(props: { label: string; value: string; icon: string }) {
  return (
    <div class='bg-gray-50 rounded-lg p-4'>
      <div class='flex items-center gap-2 mb-1'>
        <span class='text-xl'>{props.icon}</span>
        <span class='text-sm font-medium text-gray-500'>{props.label}</span>
      </div>
      <div class='text-xl font-bold text-gray-900'>{props.value}</div>
    </div>
  )
}
